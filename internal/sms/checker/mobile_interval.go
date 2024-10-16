package checker

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"github.com/rosas99/streaming/internal/pkg/errno"
	"github.com/rosas99/streaming/internal/sms/store"
	"github.com/rosas99/streaming/internal/sms/types"
	"github.com/rosas99/streaming/pkg/log"
	"strconv"
	"time"
)

type TimeIntervalForMobileRule struct {
	DS  store.IStore
	RDS *redis.Client
}

func NewTimeIntervalForMobileRule(DS store.IStore, RDS *redis.Client) *TimeIntervalForMobileRule {
	return &TimeIntervalForMobileRule{DS: DS, RDS: RDS}
}

var _ Rule = (*TimeIntervalForMobileRule)(nil)

// isValid checks if the mobile device request exceeds the specified time interval limit.
func (m *TimeIntervalForMobileRule) isValid(ctx context.Context, rq *Request) error {
	start := time.Now().UnixMilli()
	key := types.WrapperTimeInterval(rq.Mobile, rq.TemplateCode)

	timeStampStr, err := m.RDS.Get(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		log.Errorf("Failed to get timestamp from Redis for key: %s, error: %v", key, err)
		return err
	}

	if timeStampStr == "" {
		_, err2 := m.RDS.Set(ctx, key, time.Now().UnixMilli(), 24*time.Hour).Result()
		if err2 != nil {
			return err
		}
		return nil
	}

	remainingTTL, err2 := m.RDS.TTL(ctx, key).Result()
	if err2 != nil {
		return err
	}
	_, err2 = m.RDS.Set(ctx, key, time.Now().UnixMilli(), remainingTTL).Result()
	if err2 != nil {
		return err
	}

	log.C(ctx).Infof("Time interval checker took %d seconds", time.Now().Unix()-start)

	timeStampInt, err := strconv.ParseInt(timeStampStr, 10, 64)
	if err != nil {
		return err
	}

	interval2 := time.Now().UnixMilli() - timeStampInt
	isValid := interval2 >= rq.LimitValue
	if !isValid {
		log.Errorf("%s request too frequently!", rq.Mobile)
		return errno.ErrTimestampCount

	}
	return nil
}
