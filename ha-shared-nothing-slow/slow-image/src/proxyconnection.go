package main

import (
	"io"
	"log"
	"net"
	"sync"
	"time"
)

type proxyConnection struct {
	connId         string
	upstreamConn   net.Conn
	downstreamConn net.Conn
	fromUpstream   chan []byte
	fromDownstream chan []byte
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

func (conn *proxyConnection) readerLoop(reader io.ReadCloser, ch chan []byte) {
	bufsize := config.getBufferSize()
	var newbufsize int
	p := make([]byte, bufsize)
	for {
		n, err := reader.Read(p)
		if err != nil {
			log.Printf("reader got error: %v", err)
			break
		}
		if n == 0 {
			continue
		}
		chunk := make([]byte, n)
		copy(chunk, p[:n])
		ch <- chunk
		streamDelay := config.getStreamDelay()
		if streamDelay > 0 {
			time.Sleep(time.Duration(streamDelay) * time.Millisecond)
		}
		newbufsize = config.getBufferSize()
		if newbufsize != bufsize {
			bufsize = newbufsize
			log.Printf("creating new read buffer of size %d bytes", bufsize)
			p = make([]byte, bufsize)
		}
	}

	reader.Close()
	close(ch)
	registry.removeConnection(conn.connId)
	notifyTermination()
	log.Print("reader terminated")
}

func (conn *proxyConnection) writerLoop(incomingConn net.Conn, writer io.Writer, ch chan []byte) {
	keepRunning := true
	for keepRunning {
		select {
		case b, ok := <-ch:
			if !ok {
				keepRunning = false
				break
			}
			if _, err := writer.Write(b); err != nil {
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
	registry.removeConnection(conn.connId)
	notifyTermination()
	log.Print("writer terminated")
}
