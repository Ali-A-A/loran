package cmq

import (
	"github.com/ali-a-a/loran/config"

	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

// Conn represents nats jet stream connections.
type Conn struct {
	NC *nats.Conn
	JS nats.JetStreamContext
}

// CreateConnection returns new nats connection with some options that
// they are specified in config.NATS config.
func CreateConnection(cfg config.NATS) (*nats.Conn, error) {
	opts := connectionOpts(cfg)

	conn, err := nats.Connect(cfg.URL, opts...)
	if err != nil {
		logrus.Errorf("could not connect to nats server %s: %s", cfg.URL, err)
		return nil, err
	}

	return conn, nil
}

// CreateJetStreamConnection returns new Conn using CreateConnection function.
// For more information, see CreateConnection.
func CreateJetStreamConnection(cfg config.NATS) (*Conn, error) {
	nc, err := CreateConnection(cfg)
	if err != nil {
		return nil, err
	}

	js, err := nc.JetStream(nats.MaxWait(cfg.JetStream.MaxWait))
	if err != nil {
		logrus.Errorf("could not connect to jetstream %s: %s", cfg.URL, err)
		return nil, err
	}

	return &Conn{
		NC: nc,
		JS: js,
	}, nil
}

// connectionOpts returns some nats connection options.
func connectionOpts(cfg config.NATS) []nats.Option {
	var opts []nats.Option

	opts = append(opts, nats.ReconnectWait(cfg.ReconnectWait))

	opts = append(opts, nats.MaxReconnects(cfg.MaxReconnect))

	opts = append(opts, nats.ErrorHandler(func(nc *nats.Conn, sub *nats.Subscription, err error) {
		logrus.Errorf("nats error handler: url: %s subject: %s error: %s", nc.ConnectedUrl(), sub.Subject, err)
	}))

	opts = append(opts, nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
		logrus.Errorf("nats disconnected error handler: url: %s error: %s", nc.ConnectedUrl(), err)
	}))

	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
		logrus.Infof("nats reconnect handler: [%s]", nc.ConnectedUrl())
	}))

	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
		logrus.Warnf("nats close handler: %v", nc.LastError())
	}))

	return opts
}
