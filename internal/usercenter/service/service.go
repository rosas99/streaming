package service

import (
	"github.com/rosas99/streaming/internal/usercenter/biz"
	"github.com/rosas99/streaming/pkg/api/usercenter/v1"
	"github.com/rosas99/streaming/pkg/auth"
)

type UserCenterService struct {
	Auth *auth.Authz
	v1.UnimplementedUserCenterServer
	biz biz.IBiz
}

func NewUserCenterService(biz biz.IBiz, a *auth.Authz) *UserCenterService {
	return &UserCenterService{Auth: a, biz: biz}
}
