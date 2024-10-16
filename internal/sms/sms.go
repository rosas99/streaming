package sms

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/rosas99/streaming/internal/pkg/client/usercenter"
	"github.com/rosas99/streaming/internal/pkg/idempotent"
	"github.com/rosas99/streaming/internal/pkg/middleware/header"
	"github.com/rosas99/streaming/internal/pkg/middleware/trace"
	"github.com/rosas99/streaming/internal/sms/biz"
	"github.com/rosas99/streaming/internal/sms/checker"
	"github.com/rosas99/streaming/internal/sms/middleware/validate"
	"github.com/rosas99/streaming/internal/sms/monitor"
	"github.com/rosas99/streaming/internal/sms/mqs"
	providerFactory "github.com/rosas99/streaming/internal/sms/provider"
	"github.com/rosas99/streaming/internal/sms/service"
	"github.com/rosas99/streaming/internal/sms/store"
	"github.com/rosas99/streaming/internal/sms/types"
	"github.com/rosas99/streaming/internal/sms/writer"
	"github.com/rosas99/streaming/pkg/db"
	"github.com/rosas99/streaming/pkg/log"
	genericoptions "github.com/rosas99/streaming/pkg/options"
	ailiyunoptions "github.com/rosas99/streaming/pkg/sdk/ailiyun"
)

// Config represents the configuration of the service.
type Config struct {
	GRPCOptions            *genericoptions.GRPCOptions
	HTTPOptions            *genericoptions.HTTPOptions
	TLSOptions             *genericoptions.TLSOptions
	MySQLOptions           *genericoptions.MySQLOptions
	RedisOptions           *genericoptions.RedisOptions
	CommonKafkaOptions     *genericoptions.KafkaOptions
	VerifyKafkaOptions     *genericoptions.KafkaOptions
	UplinkMessageKqOptions *genericoptions.KafkaOptions
	MonitorKafkaOptions    *genericoptions.KafkaOptions
	Address                string
	Accounts               map[string]string
	AiliyunSmsOptions      *ailiyunoptions.SmsOptions
}

// Complete fills in any fields not set that are required to have valid data. It's mutating the receiver.
func (cfg *Config) Complete() completedConfig {
	return completedConfig{cfg}
}

type completedConfig struct {
	*Config
}

// SmsServer represents the fake server.
type SmsServer struct {
	httpsrv Server
	mqsrv   MqServer
	mqsrv2  MqServer
	mqsrv3  MqServer
	config  completedConfig
}

// New returns a new instance of SmsServer from the given config.
func (c completedConfig) New() (*SmsServer, error) {
	var ds store.IStore

	var dbOptions db.MySQLOptions
	_ = copier.Copy(&dbOptions, c.MySQLOptions)
	ins, err := db.NewMySQL(&dbOptions)
	if err != nil {
		return nil, err
	}
	ds = store.NewStore(ins)

	var redisOptions db.RedisOptions
	value := &c.Config.RedisOptions
	_ = copier.Copy(&redisOptions, value)
	rds, err := db.NewRedis(&redisOptions)
	if err != nil {
		return nil, err
	}

	// registers message check rules
	factory := checker.NewRuleFactory()
	factory.RegisterRule(types.MessageCountForTemplatePerDay, checker.NewMessageCountForTemplateRule(ds, rds))
	factory.RegisterRule(types.MessageCountForMobilePerDay, checker.NewMessageCountForMobileRule(ds, rds))
	factory.RegisterRule(types.TimeIntervalForMobilePerDay, checker.NewTimeIntervalForMobileRule(ds, rds))

	// creates a logger instance
	l, err := writer.NewWriter(c.CommonKafkaOptions, c.VerifyKafkaOptions,
		c.UplinkMessageKqOptions, ds.Histories())
	if err != nil {
		return nil, err
	}

	// creates a monitor instance
	_, err = monitor.NewMonitor(c.MonitorKafkaOptions)
	if err != nil {
		return nil, err
	}

	// creates an idempotent instance
	idt, err := idempotent.NewIdempotent(rds)
	if err != nil {
		return nil, err
	}

	bizIns := biz.NewBiz(ds, rds, idt, l)
	srv := service.NewSmsServerService(bizIns)

	// Sets the running mode for the Gin
	gin.SetMode(gin.ReleaseMode)
	// create a gin engine
	g := gin.New()

	usercenter.NewUserCenterServer()

	installRouters(g, srv)
	mws := []gin.HandlerFunc{
		gin.Recovery(), header.NoCache, header.Cors, header.Secure,
		trace.TraceID(), nil, validate.Validation(ds),
	}
	// add gin middlewares
	g.Use(mws...)

	httpsrv, err := NewHTTPServer(c.HTTPOptions, c.TLSOptions, g)
	if err != nil {
		return nil, err
	}
	providers := providerFactory.NewProviderFactory()
	providers.RegisterProvider(types.ProviderTypeAliyun,
		providerFactory.NewAILIYUNProvider(types.ProviderTypeAliyun, rds, l, c.AiliyunSmsOptions))
	providers.RegisterProvider(types.ProviderTypeDummy,
		providerFactory.NewDummyProvider(types.ProviderTypeDummy))

	handler1 := mqs.NewCommonMessageConsumer(context.Background(), providers, idt, l)
	mqsrv, err := NewMqServer(c.CommonKafkaOptions, handler1)
	if err != nil {
		return nil, err
	}

	handler2 := mqs.NewVerifyMessageConsumer(context.Background(), providers, idt, l)
	mqsrv2, err := NewMqServer(c.VerifyKafkaOptions, handler2)
	if err != nil {
		return nil, err
	}

	handler3 := mqs.NewUplinkMessageConsumer(context.Background(), ds, idt, l)
	mqsrv3, err := NewMqServer(c.UplinkMessageKqOptions, handler3)
	if err != nil {
		return nil, err
	}

	// Need start grpc server first. http server depends on grpc sever.
	return &SmsServer{httpsrv: httpsrv, mqsrv: mqsrv, mqsrv2: mqsrv2, mqsrv3: mqsrv3, config: c}, nil
}

// Run is a method of the SmsServer struct that starts the server.
func (s *SmsServer) Run(stopCh <-chan struct{}) error {

	log.Infof("Successfully start sms server")
	go s.httpsrv.RunOrDie()
	go s.mqsrv.RunOrDie()
	go s.mqsrv2.RunOrDie()
	go s.mqsrv3.RunOrDie()
	<-stopCh

	log.Infof("Gracefully shutting down sms server ...")

	// The most gracefully way is to shut down the dependent service first,
	// and then shutdown the depended on service.
	s.httpsrv.GracefulStop()
	s.mqsrv.GracefulStop()
	s.mqsrv2.GracefulStop()
	s.mqsrv3.GracefulStop()

	return nil
}
