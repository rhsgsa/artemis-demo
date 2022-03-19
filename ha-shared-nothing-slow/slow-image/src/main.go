package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	_ "time/tzdata"

	"github.com/kwkoo/configparser"
)

var config proxyConfig

var listener = struct {
	wg          sync.WaitGroup
	mux         sync.RWMutex
	targetAlive bool
	ln          net.Listener
	keepRunning bool
}{}
var goRoutinesWaitGroup sync.WaitGroup

// this channel closes after we receive a signal
var signalShutdown chan struct{}

func main() {
	if err := configparser.Parse(&config); err != nil {
		log.Fatal(err)
	}
	if config.ListenAddress == "0.0.0.0" {
		config.ListenAddress = ""
	}

	log.Print(&config)

	// Setup signal handling.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	signalShutdown = make(chan struct{})

	// start web server
	server := &http.Server{
		Addr: fmt.Sprintf(":%d", config.Port),
	}
	goRoutinesWaitGroup.Add(1)
	go initHttpServer(server)

	// start listener goroutines here
	listener.keepRunning = true
	listener.wg.Add(2)
	go targetChecker()
	go startListener()

	<-shutdown
	log.Print("interrupt signal received")
	signal.Reset(os.Interrupt, syscall.SIGTERM)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	server.Shutdown(ctx)

	killListener()
	close(signalShutdown)

	log.Print("waiting for all go routines to terminate...")
	goRoutinesWaitGroup.Wait()
	log.Print("go routines terminated")

	log.Print("waiting for listener to terminate...")
	listener.wg.Wait()
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
	listener.wg.Done()
}

func startListener() {
	defer notifyListenerTermination()
	for listenerRunning() {
		if targetIsAlive() {
			conn := acceptIncomingConnection()
			if conn == nil {
				continue
			}
			addGoRoutine()
			go handleConnection(conn)
		} else {
			time.Sleep(time.Second)
		}
	}
}

func handleConnection(upstreamConn net.Conn) {
	defer notifyTermination()
	connDelay := config.getConnDelay()
	if connDelay > 0 {
		time.Sleep(time.Duration(connDelay) * time.Millisecond)
	}

	downstreamConn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", config.TargetAddress, config.TargetPort))
	if err != nil {
		log.Printf("error connecting to target: %v", err)
		upstreamConn.Close()
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
			log.Printf("reader got error: %v", err)
			break
		}
		if n == 0 {
			continue
		}
		ch <- p[0]
		streamDelay := config.getStreamDelay()
		if streamDelay > 0 {
			time.Sleep(time.Duration(streamDelay) * time.Millisecond)
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

func targetChecker() {
	defer notifyListenerTermination()
	for listenerRunning() {
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", config.TargetAddress, config.TargetPort))
		if err != nil {
			setTargetState(false)
		} else {
			conn.Close()
			setTargetState(true)
		}
		time.Sleep(5 * time.Second)
	}
}

func setTargetState(alive bool) {
	listener.mux.Lock()
	if listener.targetAlive && !alive && listener.ln != nil {
		listener.ln.Close()
		listener.ln = nil
	} else if !listener.targetAlive && alive && listener.ln == nil {
		l := fmt.Sprintf("%s:%d", config.ListenAddress, config.ListenPort)
		ln, err := net.Listen("tcp", l)
		if err != nil {
			log.Fatalf("error starting listener: %v", err)
			return
		}
		listener.ln = ln
	}
	listener.targetAlive = alive
	listener.mux.Unlock()
}

func targetIsAlive() bool {
	listener.mux.RLock()
	defer listener.mux.RUnlock()
	return listener.targetAlive
}

func acceptIncomingConnection() net.Conn {
	listener.mux.RLock()
	if listener.ln == nil {
		listener.mux.RUnlock()
		return nil
	}
	ln := listener.ln
	listener.mux.RUnlock()
	conn, err := ln.Accept()
	if err != nil {
		log.Printf("error accepting connection: %v", err)
		return nil
	}
	return conn
}

func killListener() {
	listener.mux.Lock()
	listener.keepRunning = false
	if listener.ln != nil {
		listener.ln.Close()
		listener.ln = nil
	}
	listener.mux.Unlock()
}

func listenerRunning() bool {
	listener.mux.RLock()
	defer listener.mux.RUnlock()
	return listener.keepRunning
}
