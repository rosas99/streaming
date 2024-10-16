package monitor

import (
	"github.com/rosas99/streaming/pkg/log"
	genericoptions "github.com/rosas99/streaming/pkg/options"
	"github.com/segmentio/kafka-go"
	"sync"
)

var (
	once sync.Once
	cli  *monitor
)

type monitor struct {
	writer *kafka.Writer
}

// NewMonitor creates a new kafkaLogger instance.
func NewMonitor(monitorOpts *genericoptions.KafkaOptions) (*monitor, error) {
	writer, err := monitorOpts.Writer()
	if err != nil {
		log.Errorf("Failed to create Kafka writer error: %s", err)
		return nil, err
	}

	once.Do(func() {
		cli = &monitor{writer: writer}
	})
	return cli, nil
}

// GetMonitor returns the globally initialized client.
func GetMonitor() *monitor {
	return cli
}
