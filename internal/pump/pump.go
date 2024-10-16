//go:generate wire .
package pump

import (
	"github.com/rosas99/streaming/pkg/log"
	"github.com/rosas99/streaming/pkg/streams/flow"
	"time"

	"github.com/segmentio/kafka-go"
	"k8s.io/apimachinery/pkg/util/wait"

	"go.mongodb.org/mongo-driver/mongo"

	genericoptions "github.com/rosas99/streaming/pkg/options"
	kafkaconnector "github.com/rosas99/streaming/pkg/streams/connector/kafka"
	mongoconnector "github.com/rosas99/streaming/pkg/streams/connector/mongo"
)

// Config defines the config for the apiserver.
type Config struct {
	KafkaOptions *genericoptions.KafkaOptions
	MongoOptions *genericoptions.MongoOptions
	RedisOptions *genericoptions.RedisOptions
	MySQLOptions *genericoptions.MySQLOptions
}

// Server contains state for a Kubernetes cluster master/api server.
type Server struct {
	kafkaReader kafka.ReaderConfig
	colName     string
	db          *mongo.Database
}

type completedConfig struct {
	*Config
}

// addUTC appends a UTC timestamp to the beginning of the message value.
var addUTC = func(msg kafka.Message) kafka.Message {
	timestamp := time.Now().Format(time.DateTime)

	// Concatenate the UTC timestamp with msg.Value
	msg.Value = []byte(timestamp + " " + string(msg.Value))
	return msg
}

// Complete fills in any fields not set that are required to have valid data. It's mutating the receiver.
func (cfg *Config) Complete() completedConfig {
	return completedConfig{cfg}
}

// New returns a new instance of Server from the given config.
// Certain config fields will be set to a default value if unset.
func (c completedConfig) New() (*Server, error) {
	client, err := c.MongoOptions.NewClient()
	if err != nil {
		return nil, err
	}

	server := &Server{
		kafkaReader: kafka.ReaderConfig{
			Brokers:           c.KafkaOptions.Brokers,
			Topic:             c.KafkaOptions.Topic,
			GroupID:           c.KafkaOptions.ReaderOptions.GroupID,
			QueueCapacity:     c.KafkaOptions.ReaderOptions.QueueCapacity,
			MinBytes:          c.KafkaOptions.ReaderOptions.MinBytes,
			MaxBytes:          c.KafkaOptions.ReaderOptions.MaxBytes,
			MaxWait:           c.KafkaOptions.ReaderOptions.MaxWait,
			ReadBatchTimeout:  c.KafkaOptions.ReaderOptions.ReadBatchTimeout,
			HeartbeatInterval: c.KafkaOptions.ReaderOptions.HeartbeatInterval,
			CommitInterval:    c.KafkaOptions.ReaderOptions.CommitInterval,
			RebalanceTimeout:  c.KafkaOptions.ReaderOptions.RebalanceTimeout,
			StartOffset:       c.KafkaOptions.ReaderOptions.StartOffset,
			MaxAttempts:       c.KafkaOptions.ReaderOptions.MaxAttempts,
		},
		colName: c.MongoOptions.Collection,
		db:      client.Database(c.MongoOptions.Database),
	}

	return server, nil
}

type PreparedServer struct {
	*Server
}

func (s *Server) PrepareRun() PreparedServer {
	return PreparedServer{s}
}

func (s PreparedServer) Run(stopCh <-chan struct{}) error {
	ctx := wait.ContextForChannel(stopCh)

	source, err := kafkaconnector.NewKafkaSource(ctx, s.kafkaReader)
	if err != nil {
		return err
	}

	filter := flow.NewMap(addUTC, 1)

	sink, err := mongoconnector.NewMongoSink(ctx, s.db, mongoconnector.SinkConfig{
		CollectionName:            s.colName,
		CollectionCapMaxDocuments: 2000,
		CollectionCapMaxSizeBytes: 5 * genericoptions.GiB,
		CollectionCapEnable:       true,
	})
	if err != nil {
		return err
	}
	log.Infof("Successfully start pump server")
	via := source.Via(filter)
	outs := flow.FanOut(via, 2)
	for _, out := range outs {
		out.Via(filter).To(sink) // loyalty mq
		out.Via(filter).To(sink) // cdp mq
	}

	// sink 完成后会结束服务

	// 每次创建定时任务时，启动一个pump
	return err
}
