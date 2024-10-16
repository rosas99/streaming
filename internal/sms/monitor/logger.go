package monitor

import (
	"context"
	"encoding/json"
	"github.com/rosas99/streaming/internal/pkg/meta"
	"github.com/rosas99/streaming/pkg/log"
	"github.com/segmentio/kafka-go"
)

const (
	// AppName is the appName name when starting streaming-sms server.
	AppName = "streaming-sms"
)

// LogKpi writes a log message for the api request.
func (i *monitor) LogKpi(kpiName, traceId, templateCode string, status bool, costTime int64) {
	extra := map[string]any{"template_code": templateCode}

	kpi := meta.NewKpiOptions(meta.WithAppName(AppName), meta.WithKpiName(kpiName), meta.WithTraceId(traceId),
		meta.WithStatus(status), meta.WithCostTime(costTime), meta.WithExtra(extra)).Kpi

	out, _ := json.Marshal(kpi)
	if err := i.writer.WriteMessages(context.Background(), kafka.Message{Value: out}); err != nil {
		log.Errorw(err, "Failed to write kafka messages")
	}
}

func (i *monitor) LogTemplateKpi(kpiName, traceId string, status bool, costTime int64) {
	kpi := meta.NewKpiOptions(meta.WithAppName(AppName), meta.WithKpiName(kpiName), meta.WithTraceId(traceId),
		meta.WithStatus(status), meta.WithCostTime(costTime)).Kpi

	out, _ := json.Marshal(kpi)
	if err := i.writer.WriteMessages(context.Background(), kafka.Message{Value: out}); err != nil {
		log.Errorw(err, "Failed to write kafka messages")
	}
}
