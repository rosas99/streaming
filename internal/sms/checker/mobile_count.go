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

type MessageCountForMobileRule struct {
	DS  store.IStore
	RDS *redis.Client
}

func NewMessageCountForMobileRule(DS store.IStore, RDS *redis.Client) *MessageCountForMobileRule {
	return &MessageCountForMobileRule{DS: DS, RDS: RDS}
}

var _ Rule = (*MessageCountForMobileRule)(nil)

// isValid validates if the message sending count for a specific mobile and template is within limits.
func (m *MessageCountForMobileRule) isValid(ctx context.Context, rq *Request) error {
	start := time.Now().Unix()
	key := types.WrapperMobileCount(rq.Mobile, rq.TemplateCode)

	sentCount, err := m.RDS.Incr(ctx, key).Result()
	if err != nil {
		log.Errorf("Failed to increment count for key: %s, error: %v", key, err)
		return err
	}

	if sentCount == 1 {
		err = m.RDS.Expire(ctx, key, types.LimitLeftTime).Err()
		if err != nil {
			log.Fatalf("Error setting expiration for key: %v", err)
		}
	}

	log.Infof("Mobile times checker took %d seconds", time.Now().Unix()-start)

	isValid := sentCount <= rq.LimitValue
	if !isValid {
		log.Errorf("Exceed limit for this phone: %v", rq.Mobile)
		return errno.ErrMobileCount

	}
	return nil
}
