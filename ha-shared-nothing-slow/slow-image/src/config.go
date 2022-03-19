package main

import (
	"errors"
	"fmt"
	"sync"
)

type proxyConfig struct {
	Port          int    `default:"8000" usage:"HTTP listener port"`
	ListenAddress string `default:"" usage:"Listener address"`
	ListenPort    int    `default:"8080" usage:"Listener port"`
	TargetAddress string `default:"localhost" usage:"Target address"`
	TargetPort    int    `default:"8081" usage:"Target port"`
	ConnDelay     int    `default:"0" usage:"Connection delay in milliseconds"`
	StreamDelay   int    `default:"0" usage:"Delay between bytes in milliseconds"`
	mux           sync.RWMutex
}

func (p *proxyConfig) String() string {
	return fmt.Sprintf("Port=%d ListenAddress=%s ListenPort=%d TargetAddress=%s TargetPort=%d ConnDelay=%d StreamDelay=%d", p.Port, p.ListenAddress, p.ListenPort, p.TargetAddress, p.TargetPort, p.ConnDelay, p.StreamDelay)
}

func (p *proxyConfig) getStreamDelay() int {
	p.mux.RLock()
	defer p.mux.RUnlock()
	return p.StreamDelay
}

func (p *proxyConfig) setStreamDelay(d int) error {
	if d < 0 {
		return errors.New("stream delay cannot be below 0")
	}
	p.mux.Lock()
	p.StreamDelay = d
	p.mux.Unlock()

	return nil
}

func (p *proxyConfig) getConnDelay() int {
	p.mux.RLock()
	defer p.mux.RUnlock()
	return p.ConnDelay
}

func (p *proxyConfig) setConnDelay(d int) error {
	if d < 0 {
		return errors.New("connection delay cannot be below 0")
	}
	p.mux.Lock()
	p.ConnDelay = d
	p.mux.Unlock()

	return nil
}
