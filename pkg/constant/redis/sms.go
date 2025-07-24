package redis

const (
	VerifyCodeCacheBasisKey    = "verify_code_%v_%v_%v"        // 验证码发送场景key v1: 验证场景 v2:地区码 v3:手机号
	VerifyCodeCache24HLimitKey = "verify_code_24h_limit_%v_%v" // 24H发送限制key v1: 验证场景 v2:地区码
	VerifyCodeCoolDownCacheKey = "verify_code_cool_down_%v_%v" // 验证码发送60秒冷却key v1: 验证场景 v2:地区码
)
