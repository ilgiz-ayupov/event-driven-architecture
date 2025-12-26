package sse_broker

import "sync"

type connection struct {
	ch   chan []byte
	once sync.Once
}

func newConnection(ch chan []byte) *connection {
	return &connection{ch: ch}
}
