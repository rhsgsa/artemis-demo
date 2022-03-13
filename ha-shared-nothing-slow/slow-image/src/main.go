package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	_ "time/tzdata"

	"github.com/kwkoo/configparser"
)

var config = struct {
	ListenAddress string `default:"" usage:"Listener address"`
	ListenPort    int    `default:"8080" usage:"Listener port"`
	TargetAddress string `default:"localhost" usage:"Target address"`
	TargetPort    int    `default:"8081" usage:"Target port"`
	ConnDelay     int    `default:"0" usage:"Connection delay in milliseconds"`
	StreamDelay   int    `default:"0" usage:"Delay between bytes in milliseconds"`
}{}

var goRoutinesWaitGroup sync.WaitGroup
var listenerWaitGroup sync.WaitGroup

// this channel closes after we receive a signal
var signalShutdown chan struct{}

func main() {
	if err := configparser.Parse(&config); err != nil {
		log.Fatal(err)
	}
	if config.ListenAddress == "0.0.0.0" {
		config.ListenAddress = ""
	}

	log.Printf("ListenAddress=%s ListenPort=%d TargetAddress=%s TargetPort=%d ConnDelay=%d StreamDelay=%d", config.ListenAddress, config.ListenPort, config.TargetAddress, config.TargetPort, config.ConnDelay, config.StreamDelay)

	// Setup signal handling.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	signalShutdown = make(chan struct{})

	// start listener here
	l := fmt.Sprintf("%s:%d", config.ListenAddress, config.ListenPort)
	ln, err := net.Listen("tcp", l)
	if err != nil {
		log.Fatalf("error starting listener: %v", err)
		return
	}
	log.Printf("listening on %s...", l)
	listenerWaitGroup.Add(1)
	go startListener(ln)

	<-shutdown
	log.Print("interrupt signal received")
	signal.Reset(os.Interrupt, syscall.SIGTERM)

	ln.Close()
	close(signalShutdown)

	log.Print("waiting for all go routines to terminate...")
	goRoutinesWaitGroup.Wait()
	log.Print("go routines terminated")

	log.Print("waiting for listener to terminate...")
	listenerWaitGroup.Wait()
	log.Print("listener terminated")

	log.Print("shutdown successful")
}

func notifyTermination() {
	goRoutinesWaitGroup.Done()
}

func addGoRoutine() {
	goRoutinesWaitGroup.Add(1)
}

func notifyListenerTermination() {
	listenerWaitGroup.Done()
}

func startListener(ln net.Listener) {
	defer notifyListenerTermination()
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("error accepting connection: %v", err)
			return
		}

		addGoRoutine()
		go handleConnection(conn)
	}
}

func handleConnection(upstreamConn net.Conn) {
	defer notifyTermination()
	if config.ConnDelay > 0 {
		time.Sleep(time.Duration(config.ConnDelay) * time.Millisecond)
	}

	downstreamConn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", config.TargetAddress, config.TargetPort))
	if err != nil {
		log.Printf("error connecting to target: %v", err)
		return
	}

	proxyConn := proxyConnection{
		upstreamConn:   upstreamConn,
		downstreamConn: downstreamConn,
		fromUpstream:   make(chan byte),
		fromDownstream: make(chan byte),
		stopReaders:    make(chan struct{}),
	}

	addGoRoutine()
	go proxyConn.readerLoop(upstreamConn, proxyConn.fromUpstream)
	addGoRoutine()
	go proxyConn.writerLoop(upstreamConn, downstreamConn, proxyConn.fromUpstream)

	addGoRoutine()
	go proxyConn.readerLoop(downstreamConn, proxyConn.fromDownstream)
	addGoRoutine()
	go proxyConn.writerLoop(downstreamConn, upstreamConn, proxyConn.fromDownstream)
}

type proxyConnection struct {
	upstreamConn   net.Conn
	downstreamConn net.Conn
	fromUpstream   chan byte
	fromDownstream chan byte
	stopReaders    chan struct{}
	mux            sync.Mutex
}

func (conn *proxyConnection) shutdown() {
	conn.mux.Lock()
	if conn.stopReaders != nil {
		close(conn.stopReaders)
		conn.stopReaders = nil
	}
	conn.mux.Unlock()
}

func (conn *proxyConnection) readerLoop(reader io.ReadCloser, ch chan byte) {
	p := make([]byte, 1)
	for {
		n, err := reader.Read(p)
		if err != nil {
			//log.Printf("reader got error: %v", err)
			break
		}
		if n == 0 {
			continue
		}
		ch <- p[0]
		if config.StreamDelay > 0 {
			time.Sleep(time.Duration(config.StreamDelay) * time.Millisecond)
		}
	}

	reader.Close()
	close(ch)
	conn.shutdown()
	notifyTermination()
	log.Print("reader terminated")
}

func (conn *proxyConnection) writerLoop(incomingConn net.Conn, writer io.Writer, ch chan byte) {
	p := make([]byte, 1)
	keepRunning := true
	for keepRunning {
		select {
		case b, ok := <-ch:
			if !ok {
				keepRunning = false
				break
			}
			p[0] = b
			if _, err := writer.Write(p); err != nil {
				//log.Printf("error writing to stream: %v", err)
				keepRunning = false
			}
		case <-conn.stopReaders:
			keepRunning = false
		case <-signalShutdown:
			keepRunning = false
		}
	}
	incomingConn.Close()
	conn.shutdown()
	notifyTermination()
	log.Print("writer terminated")
}
