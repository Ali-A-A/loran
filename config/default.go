package config

import (
	"time"

	"github.com/nats-io/nats.go"
)

// nolint:gomnd,funlen
func Default() Config {
	return Config{
		NATS: NATS{
			URL:            "127.0.0.1",
			ReconnectWait:  1 * time.Second,
			MaxReconnect:   120,
			PublishEnabled: false,
			JetStream: JetStream{
				Enable:   true,
				MaxWait:  500 * time.Millisecond,
				Replicas: 1,
				MaxAge:   1 * time.Minute,
				Storage:  nats.MemoryStorage,
				Consumer: Consumer{
					Durable: "durable",
					Stream:  "stream",
					Subject: "stream.subject",
				},
			},
		},
		Redis: Redis{
			Master: RedisConfig{
				Address: "127.0.0.1:6379",
			},
		},
	}
}
