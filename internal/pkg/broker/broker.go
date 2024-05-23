package broker

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	"time"
)

type Config struct {
	Host    string `yaml:"host"`
	Port    string `yaml:"port"`
	Subject string `yaml:"subject"`
	Durable string `yaml:"durable"`
}

type IBroker interface {
	Subscribe(subject string, handler func(msg *nats.Msg)) (*nats.Subscription, error)
}

type NATSBroker struct {
	nc   *nats.Conn
	js   nats.JetStreamContext
	cfg  *Config
	opts *nats.ConsumerConfig
}

func (b *NATSBroker) Subscribe(subject string, handler func(msg *nats.Msg)) (*nats.Subscription, error) {
	sub, err := b.js.PullSubscribe(subject, b.cfg.Durable)
	if err != nil {
		return nil, err
	}
	go func() {
		for {
			msgs, err := sub.Fetch(1, nats.MaxWait(time.Second*10))
			if err != nil {
				log.Println("Error fetching messages:", err)
				continue
			}

			for _, msg := range msgs {
				handler(msg)
				msg.Ack()
			}
		}
	}()

	return sub, nil
}

func NewBroker(cfg *Config) (IBroker, error) {
	nc, err := nats.Connect(fmt.Sprintf("nats://%s:%s", cfg.Host, cfg.Port))
	if err != nil {
		return nil, err
	}

	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		return nil, err
	}

	return &NATSBroker{nc: nc, js: js, cfg: cfg}, nil
}
