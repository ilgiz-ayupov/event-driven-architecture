package infrastructure

import "github.com/nats-io/nats.go"

func NewNATS(url string) (*nats.Conn, error) {
	return nats.Connect(url)
}
