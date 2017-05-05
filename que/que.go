package que

import (
	"flag"
	"time"

	nats "github.com/nats-io/go-nats"
)

var (
	NatsURL string
)

type Que struct {
	NatsClient *nats.Conn
	Channel    string
	Timeout    time.Duration
}

func Flags() {
	flag.StringVar(&NatsURL, "nats-url", nats.DefaultURL, "Url for Nats")
}

func New(ch string) (*Que, error) {
	q := Que{Channel: ch, Timeout: time.Second}
	var err error
	q.NatsClient, err = nats.Connect(NatsURL)
	return &q, err
}

func (q *Que) Sync(b []byte) ([]byte, error) {
	m, err := q.NatsClient.Request(q.Channel, b, q.Timeout)
	return m.Data, err
}

func (q *Que) Async(b []byte) error {
	return q.NatsClient.Publish(q.Channel, b)
}

func (q *Que) Close() {
	q.NatsClient.Close()
}
