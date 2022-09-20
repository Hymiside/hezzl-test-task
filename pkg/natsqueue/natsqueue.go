package natsqueue

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
)

type Nats struct {
	nc *nats.Conn
}

type ConfigNats struct {
	Host string
	Port string
}

func NewNats(ctx context.Context, cfg ConfigNats) (*Nats, error) {
	nc, err := nats.Connect(fmt.Sprintf("nats://%s:%s", cfg.Host, cfg.Port))
	if err != nil {
		return nil, fmt.Errorf("failed to connect nats: %w", err)
	}

	go func(ctx context.Context) {
		<-ctx.Done()
		nc.Close()
	}(ctx)

	return &Nats{nc: nc}, nil
}

func (n *Nats) LogPublish(ctx, log string) error {
	return nil
}

func (n *Nats) LogSubscribe(ctx, log string) error {
	_, err := n.nc.Subscribe("logs", func(m *nats.Msg) {
		fmt.Printf("Received a message: %s\n", string(m.Data))
	})
	if err != nil {
		return err
	}
	return nil
}
