//nolint:unused
package streamingx

import (
	"context"

	"github.com/golang-jwt/jwt/v4"

	"github.com/rosas99/streaming/internal/usercenter/model"
)

// 定义全局上下文中的键.
type (
	transCtx     struct{}
	noTransCtx   struct{}
	transLockCtx struct{}
	userIDCtx    struct{}
	traceIDCtx   struct{}
)

type (
	authKey        struct{}
	userKey        struct{}
	userMKey       struct{}
	accessTokenKey struct{}
)

// NewContext put auth info into context.
func NewContext(ctx context.Context, info *jwt.RegisteredClaims) context.Context {
	return context.WithValue(ctx, authKey{}, info)
}

// FromContext extract auth info from context.
func FromContext(ctx context.Context) (token *jwt.RegisteredClaims, ok bool) {
	token, ok = ctx.Value(authKey{}).(*jwt.RegisteredClaims)
	return
}

// NewUserID put userID into context.
func NewUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userKey{}, userID)
}

// FromUserID extract userID from context.
func FromUserID(ctx context.Context) string {
	userID, _ := ctx.Value(userKey{}).(string)
	return userID
}

// NewAccessToken put accessToken into context.
func NewAccessToken(ctx context.Context, accessToken string) context.Context {
	return context.WithValue(ctx, accessTokenKey{}, accessToken)
}

// FromAccessToken extract accessToken from context.
func FromAccessToken(ctx context.Context) string {
	accessToken, _ := ctx.Value(accessTokenKey{}).(string)
	return accessToken
}

// NewUserM put *UserM into context.
func NewUserM(ctx context.Context, user *model.UserM) context.Context {
	return context.WithValue(ctx, userMKey{}, user)
}

// FromUserM extract *UserM from extract.
func FromUserM(ctx context.Context) *model.UserM {
	user, _ := ctx.Value(userMKey{}).(*model.UserM)
	return user
}
