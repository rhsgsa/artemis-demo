package main

import "sync"

type connectionRegistry struct {
	mux     sync.RWMutex
	allConn map[string]*proxyConnection
}

var registry connectionRegistry

func initRegistry() {
	registry = connectionRegistry{
		allConn: make(map[string]*proxyConnection),
	}
}

func (r *connectionRegistry) addConnection(p *proxyConnection) {
	r.mux.Lock()
	r.allConn[p.connId] = p
	r.mux.Unlock()
}

func (r *connectionRegistry) removeConnection(id string) {
	r.mux.RLock()
	p, ok := r.allConn[id]
	r.mux.RUnlock()
	if !ok {
		return
	}

	r.mux.Lock()
	delete(r.allConn, id)
	r.mux.Unlock()

	p.shutdown()
}

/*
func (r *connectionRegistry) closeAllConnections() {
	r.mux.Lock()
	existingConn := r.allConn
	r.allConn = make(map[string]*proxyConnection)
	r.mux.Unlock()

	for _, p := range existingConn {
		p.shutdown()
	}
}
*/

func (r *connectionRegistry) size() int {
	r.mux.RLock()
	defer r.mux.RUnlock()
	return len(r.allConn)
}
