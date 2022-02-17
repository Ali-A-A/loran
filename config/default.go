package config

// nolint:gomnd,funlen
func Default() Config {
	return Config{
		NATS: NATS{
			URL:            "127.0.0.1:4222",
			ReconnectWait:  0,
			MaxReconnect:   0,
			PublishEnabled: false,
			JetStream:      JetStream{},
		},
	}
}
