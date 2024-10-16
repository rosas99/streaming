package mqs

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/rosas99/streaming/internal/pkg/idempotent"
	"github.com/rosas99/streaming/internal/sms/model"
	factory "github.com/rosas99/streaming/internal/sms/provider"
	"github.com/rosas99/streaming/internal/sms/types"
	"github.com/rosas99/streaming/internal/sms/writer"
	"github.com/rosas99/streaming/pkg/log"
	"github.com/segmentio/kafka-go"
	"time"
)

type CommonMessageConsumer struct {
	ctx       context.Context
	idt       *idempotent.Idempotent
	logger    *writer.Writer
	providers *factory.ProviderFactory
}

func NewCommonMessageConsumer(ctx context.Context, providers *factory.ProviderFactory, idt *idempotent.Idempotent, logger *writer.Writer) *CommonMessageConsumer {

	return &CommonMessageConsumer{
		ctx:       ctx,
		idt:       idt,
		logger:    logger,
		providers: providers,
	}
}

func (l *CommonMessageConsumer) Consume(elem any) error {
	val := elem.(kafka.Message)
	var msg *types.TemplateMsgRequest
	err := json.Unmarshal(val.Value, &msg)
	if err != nil {
		log.C(l.ctx).Errorf("Failed to unmarshal message value: %v, error: %v", val.Value, err)
		return err
	}
	log.C(l.ctx).Infof("Successfully unmarshalled message: %v", msg)

	if err := l.handleSmsRequest(l.ctx, msg); err != nil {
		log.C(l.ctx).Errorf("Error handling SMS request: %v", err)
		return err
	}

	log.C(l.ctx).Infof("SMS request handled successfully")

	return nil

}

func (l *CommonMessageConsumer) handleSmsRequest(ctx context.Context, msg *types.TemplateMsgRequest) error {

	if !l.idt.Check(ctx, msg.RequestId) {
		log.C(ctx).Errorf("Idempotent token check failed: %v", errors.New("idempotent token is invalid"))
		return errors.New("idempotent token is invalid")
	}

	historyM := model.HistoryM{
		Mobile:            msg.PhoneNumber,
		SendTime:          time.Now(),
		Content:           msg.Content,
		MessageTemplateID: msg.TemplateCode,
	}

	successful := false

	for _, provider := range msg.Providers {
		log.C(ctx).Infof("Attempting to use provider: %s", provider)
		providerIns, err := l.providers.GetSMSTemplateProvider(types.ProviderType(provider))
		if err != nil {
			continue
		}
		ret, err := providerIns.Send(ctx, msg)

		if err != nil {
			log.C(ctx).Errorf("Failed to send SMS: %v", err)
			historyM.Status = "Failed"
			historyM.Message = err.Error()
		} else {
			log.C(ctx).Infof("SMS sent successfully: bizId=%v", ret.BizId)
			historyM.Status = "Success"
			historyM.MessageID = ret.BizId
			historyM.Code = ret.Code
			historyM.Message = ret.Message
			successful = true
			break
		}
	}

	if successful {
		historyM.Status = "Success"
	} else {
		historyM.Status = "Failed"
	}

	l.logger.WriterHistory(&historyM)

	return nil
}
