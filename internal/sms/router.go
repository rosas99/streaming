package sms

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/rosas99/streaming/internal/pkg/client/usercenter"
	"github.com/rosas99/streaming/internal/pkg/core"
	"github.com/rosas99/streaming/internal/pkg/errno"
	"github.com/rosas99/streaming/internal/sms/controller/v1/message"
	"github.com/rosas99/streaming/internal/sms/controller/v1/template"
	mw "github.com/rosas99/streaming/internal/sms/middleware/auth"
	"github.com/rosas99/streaming/internal/sms/service"
)

func installRouters(g *gin.Engine, svc *service.SmsServerService) {
	// register 404 Handler.
	g.NoRoute(func(c *gin.Context) {
		core.WriteResponse(c, errno.ErrPageNotFound, nil)
	})

	// register pprof handler
	pprof.Register(g)

	// creates a v1 router group
	v1 := g.Group("/v1")
	{
		// get a grpc usercenter client
		client := usercenter.GetClient()
		// create template router group and adds an auth middleware.
		templatev1 := v1.Group("/template", mw.BasicAuth(client))
		{
			tl := template.New(svc)
			templatev1.POST("", tl.Create)
			templatev1.PUT("", tl.Update)
			templatev1.GET("/:id", tl.Get)
			templatev1.GET("", tl.List)
			templatev1.DELETE("/:id", tl.Delete)
		}

		// creates message router group
		msgv1 := v1.Group("/message")
		{
			ms := message.New(svc)
			msgv1.POST("/send", ms.Send)
			msgv1.POST("/verify", ms.CodeVerify)

			msgv1.POST("/report/ailiyun", ms.AiliReport)
			msgv1.POST("/interaction/ailiyun", ms.AILIYUNCallback)
		}

	}

}
