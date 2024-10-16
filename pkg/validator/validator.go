// Package validator defines streaming custom binding validators used by gin.
package validator

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"regexp"
)

var mobileRegex = regexp.MustCompile(`^1[3-9]\d{9}$`)

func isMobile(fl validator.FieldLevel) bool {
	mobile := fl.Field().String()
	return mobileRegex.MatchString(mobile)
}

func isCode(fl validator.FieldLevel) bool {
	code := fl.Field().String()
	return len(code) == 6
}

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("isMobile", isMobile)
		_ = v.RegisterValidation("isCode", isCode)
	}
}
