package user

import (
	"context"
	"github.com/rosas99/streaming/pkg/api/usercenter/v1"
	"github.com/rosas99/streaming/pkg/log"
	"github.com/rosas99/streaming/pkg/token"
)

// Authorize verifies the provided token and returns the user ID associated with it.
func (b *userBiz) Authorize(ctx context.Context, rq *v1.AuthzRequest) (*v1.AuthzResponse, error) {
	username, err := token.Parse(rq.Token, token.GetConfigKey())
	log.C(ctx).Infof("Authorize username: %s", username)
	if err != nil {
		return nil, err
	}
	return &v1.AuthzResponse{UserId: username}, nil
}
