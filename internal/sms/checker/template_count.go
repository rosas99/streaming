package checker

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/rosas99/streaming/internal/pkg/errno"
	"github.com/rosas99/streaming/internal/sms/store"
	"github.com/rosas99/streaming/internal/sms/types"
	"github.com/rosas99/streaming/pkg/log"
	"time"
)

type MessageCountForTemplateRule struct {
	DS  store.IStore
	RDS *redis.Client
}

func NewMessageCountForTemplateRule(DS store.IStore, RDS *redis.Client) *MessageCountForTemplateRule {
	return &MessageCountForTemplateRule{DS: DS, RDS: RDS}
}

var _ Rule = (*MessageCountForTemplateRule)(nil)

// isValid verifies if the message count for a given template code is within the allowed limit.
func (m *MessageCountForTemplateRule) isValid(ctx context.Context, rq *Request) error {

	start := time.Now().Unix()
	key := types.WrapperTemplateCount(rq.TemplateCode)
	sentCount, err := m.RDS.Incr(ctx, key).Result()
	if err != nil {
		log.Errorf("Failed to increment count for key: %s, error: %v", key, err)
		return err
	}

	if sentCount == 1 {
		err = m.RDS.Expire(ctx, key, types.LimitLeftTime).Err()
		if err != nil {
			// 处理错误
			log.Errorf("Error setting expiration for key: %v", err)
		}
	}
	log.C(ctx).Infof("Template times checker took %d seconds", time.Now().Unix()-start)

	isValid := sentCount <= rq.LimitValue
	if !isValid {
		log.Errorf("Exceed limit for this template: %v", rq.TemplateCode)
		return errno.ErrTemplateCount
	}
	return nil
}
