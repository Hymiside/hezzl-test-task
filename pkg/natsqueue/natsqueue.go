package natsqueue

import (
	"context"
	"fmt"

	"github.com/Hymiside/hezzl-test-task/pkg/repository/postgres"
	"github.com/nats-io/nats.go"
)

type Nats struct {
	nc *nats.Conn
	r  *postgres.Repository
}

type ConfigNats struct {
	Host string
	Port string
}

var Logs [][]byte

func NewNats(ctx context.Context, cfg ConfigNats, r *postgres.Repository) (*Nats, error) {
	nc, err := nats.Connect(fmt.Sprintf("nats://%s:%s", cfg.Host, cfg.Port))
	if err != nil {
		return nil, fmt.Errorf("failed to connect nats: %w", err)
	}

	go func(ctx context.Context) {
		<-ctx.Done()
		nc.Close()
	}(ctx)

	return &Nats{nc: nc, r: r}, nil
}

func (n *Nats) Pub(data []byte) error {
	return n.nc.Publish("hezzl", data)
}

func (n *Nats) Sub() error {
	_, err := n.nc.Subscribe("hezzl", func(msg *nats.Msg) {
		if len(Logs) < 24 {
			Logs = append(Logs, msg.Data)
		} else {
			if err := n.r.CreateLog(context.Background(), Logs); err != nil {
				fmt.Println(err.Error())
			}
			Logs = nil
		}
	})
	if err != nil {
		return err
	}
	return nil
}
