package template

import (
	"github.com/gin-gonic/gin"
	"github.com/rosas99/streaming/internal/pkg/core"
	v1 "github.com/rosas99/streaming/pkg/api/sms/v1"
	"strconv"
	// custom gin validators.
	_ "github.com/rosas99/streaming/pkg/validator"
)

func (b *Controller) Update(c *gin.Context) {
	var r v1.UpdateTemplateRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	i, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	ret, err := b.svc.UpdateTemplate(c, i, &r)
	if err != nil {
		return
	}
	core.WriteResponse(c, nil, ret)

}
