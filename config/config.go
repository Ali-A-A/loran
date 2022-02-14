package config

import (
	"github.com/nats-io/nats.go"
	"strings"
	"time"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"

	"github.com/knadh/koanf"
	"github.com/sirupsen/logrus"
)

const Prefix = "LORAN_"

type (
	// Config represents application configuration struct.
	Config struct {
		Logger Logger `koanf:"logger"`
		NATS   NATS   `koanf:"nats"`
	}

	// Logger represents logger configuration struct.
	Logger struct {
		Level string `koanf:"level"`
	}

	// NATS represents nats configuration struct.
	// Its dependency is JetStream struct.
	// For more information, see JetStream.
	NATS struct {
		URL            string        `koanf:"url"`
		ReconnectWait  time.Duration `koanf:"reconnect-wait"`
		MaxReconnect   int           `koanf:"max-reconnect"`
		PublishEnabled bool          `koanf:"publish-enabled"`
		JetStream      JetStream     `koanf:"jet-stream"`
	}

	// JetStream represents jet stream configuration struct.
	// It has just some configurations of nats jet stream
	// that needed in this application.
	// Its dependency is Consumer struct.
	JetStream struct {
		Enable    bool             `koanf:"enable"`
		Consumers []Consumer       `koanf:"consumers"`
		MaxWait   time.Duration    `koanf:"max-wait"`
		Replicas  int              `koanf:"replicas"`
		MaxAge    time.Duration    `koanf:"max-age"`
		Storage   nats.StorageType `koanf:"storage"`
	}

	// Consumer represents consumer configuration struct.
	Consumer struct {
		Durable string `koanf:"durable"`
		Stream  string `koanf:"stream"`
		Subject string `koanf:"subject"`
	}
)

func Init() Config {
	var cfg Config

	k := koanf.New(".")

	if err := k.Load(structs.Provider(Default(), "koanf"), nil); err != nil {
		logrus.Fatalf("error loading default: %s", err)
	}

	if err := k.Load(file.Provider("config.yaml"), yaml.Parser()); err != nil {
		logrus.Errorf("error loading config.yml: %s", err)
	}

	if err := k.Load(env.Provider(Prefix, ".", func(s string) string {
		parsedEnv := strings.Replace(strings.ToLower(strings.TrimPrefix(s, Prefix)), "__", "-", -1)
		return strings.Replace(parsedEnv, "_", ".", -1)
	}), nil); err != nil {
		logrus.Errorf("error loading environment variables: %s", err)
	}

	if err := k.Unmarshal("", &cfg); err != nil {
		logrus.Fatalf("error unmarshalling config: %s", err)
	}

	return cfg
}
