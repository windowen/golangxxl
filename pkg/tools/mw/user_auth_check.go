package mw

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"queueJob/pkg/constant"
	"queueJob/pkg/db/cache"
	"queueJob/pkg/jwt"
	"queueJob/pkg/tools/apiresp"
	"queueJob/pkg/tools/errs"
	"queueJob/pkg/zlogger"
)

func UserAuthCheck(rdb redis.UniversalClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := strings.TrimSpace(c.GetHeader(constant.Authorization))

		if authorization == "" {
			apiresp.GinError(c, errs.ErrNoPermission.WithDetail("no_token"))
			c.Abort()
			return
		}

		jwtTools := jwt.NewJWT(rdb)

		claims, err := jwtTools.ParseToken(c, authorization)
		if err != nil {
			zlogger.Errorf("token parsing failed ,err: %v", err)
			apiresp.GinError(c, errs.ErrNoPermission.WithDetail("token_not_exists"))
			c.Abort()
			return
		}

		userCache := cache.NewUserCache(rdb)
		// 检查用户状态是否正常
		userCacheInfo, err := userCache.Get(c, claims.UserId)
		if err != nil {
			zlogger.Errorf("failed to get user cache when checking user status ,err: %v", err)
			apiresp.GinError(c, errs.ErrNoPermission.WithDetail("token_not_exists"))
			c.Abort()
			return
		}

		// 检查令牌过期情况并根据需要刷新
		if claims.RegisteredClaims.ExpiresAt.Time.Before(time.Now().Add(24 * time.Hour)) {
			// Token即将过期，生成新的Token
			newToken, err := jwtTools.GenToken(c, claims.UserId)
			if err != nil {
				apiresp.GinError(c, errs.ErrNoPermission.WithDetail("token_not_exists"))
				c.Abort()
				return
			}

			// 返回新的令牌
			c.Header(constant.RefreshAuthorization, newToken)
		}

		if userCacheInfo.Status != constant.UserStatusNormal {
			apiresp.GinError(c, errs.ErrNoPermission.WithDetail("user_disabled"))
			c.Abort()
			return
		}

		setRequireParamsWithOpts(c,
			WithUserID(claims.UserId),
		)

		c.Next()
	}
}
