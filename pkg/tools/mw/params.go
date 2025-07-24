package mw

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"queueJob/pkg/constant"
)

type Param struct {
	c           *gin.Context
	platform    string
	language    string
	operationID string
	countryCode string
	userID      int
}

type Option func(*Param)

func WithPlatform(platform string) Option {
	return func(p *Param) {
		p.platform = platform
	}
}

func WithLanguage(language string) Option {
	return func(p *Param) {
		p.language = language
	}
}

func WithOperationID(operationID string) Option {
	return func(p *Param) {
		p.operationID = operationID
	}
}

func WithCountryCode(countryCode string) Option {
	return func(p *Param) {
		p.countryCode = countryCode
	}
}

func WithUserID(userID int) Option {
	return func(p *Param) {
		p.userID = userID
	}
}

// 应用 Option 的函数
func setRequireParamsWithOpts(c *gin.Context, opts ...Option) {
	// 检查是否已有 Param 对象
	param, exists := c.Get("param")
	if !exists {
		param = &Param{c: c}
	}

	p := param.(*Param)

	// 应用所有 Option
	for _, opt := range opts {
		opt(p)
	}

	// 设置 gin.Context 中的参数
	c.Set(constant.OperationId, p.operationID)
	c.Set(constant.RpcOpUserID, p.userID)
	c.Set(constant.RpcOpUserType, []string{strconv.Itoa(1)})
	c.Set(constant.RpcCustomHeader, []string{constant.RpcOpUserType})
	c.Set(constant.OpUserPlatform, p.platform)
	c.Set(constant.Language, p.language)
	c.Set(constant.CountryCode, p.countryCode)

	// 更新后的 Param
	c.Set("param", p)
}
