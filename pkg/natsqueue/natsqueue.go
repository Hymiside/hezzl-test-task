package natsqueue

import (
	"fmt"

	"github.com/nats-io/nats.go"
)

type Nats struct {
	n *nats.Conn
}

type ConfigNats struct {
	Host string
	Port string
}

func NewNats(c ConfigNats) (*Nats, error) {
	nc, err := nats.Connect(fmt.Sprintf("nats://%s:%s", c.Host, c.Port))
	if err != nil {
		return nil, fmt.Errorf(" i'm cry, connection died %s", err)
	}
	return &Nats{n: nc}, nil
}

func (nc *Nats) CloseNats() {
	nc.n.Close()
}
