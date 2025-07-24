package jwt

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"

	"queueJob/pkg/common/config"
	constsR "queueJob/pkg/constant/redis"
	"queueJob/pkg/db/cache"
	"queueJob/pkg/zlogger"
)

type JWT struct {
	key []byte
	rdb redis.UniversalClient
}

type MyCustomClaims struct {
	UserId  int `json:"userId"`
	Version int `json:"version"`
	jwt.RegisteredClaims
}

func NewJWT(redisClient redis.UniversalClient) *JWT {
	return &JWT{
		key: []byte(config.Config.Jwt.Key),
		rdb: redisClient,
	}
}

// generateToken 生成Token
func (j *JWT) generateToken(userId, version int) (string, error) {
	var (
		timeNow   = time.Now()
		expiresAt = timeNow.Add(constsR.TokenTTL)
	)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, MyCustomClaims{
		UserId:  userId,
		Version: version,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "",
			ExpiresAt: jwt.NewNumericDate(expiresAt), // 7天有效期
			IssuedAt:  jwt.NewNumericDate(timeNow),   // 签发时间
			NotBefore: jwt.NewNumericDate(timeNow),   // 生效时间
		},
	})

	tokenString, err := token.SignedString(j.key)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("Bearer %s", tokenString), nil
}

// ParseToken 解析token
func (j *JWT) ParseToken(ctx context.Context, tokenString string) (*MyCustomClaims, error) {
	// 移除 Bearer 前缀
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// 检查 token 是否为空
	if strings.TrimSpace(tokenString) == "" {
		return nil, errors.New("token is empty")
	}

	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return j.key, nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	// 检查 token 是否有效
	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		// 版本检查
		storedVersion, err := j.rdb.Get(ctx, fmt.Sprintf(constsR.UserLoginTokenVersion, claims.UserId)).Result()
		if err := cache.CheckErr(err); err != nil {
			return nil, fmt.Errorf("failed to get token version from Redis: %w", err)
		}

		intVersion, err := strconv.Atoi(storedVersion)
		if err != nil {
			return nil, fmt.Errorf("failed to parse token version from Redis: %w", err)
		}
		if intVersion != claims.Version {
			return nil, fmt.Errorf("token version mismatch")
		}

		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// GenToken 获取携带版本的Token
func (j *JWT) GenToken(ctx context.Context, userId int) (string, error) {
	// 获取用户token当前版本
	storedVersion, err := j.rdb.Get(ctx, fmt.Sprintf(constsR.UserLoginTokenVersion, userId)).Result()
	if err := cache.CheckErr(err); err != nil {
		return "", fmt.Errorf("failed to get token version: %v", err)
	}

	newVersion := 1
	if storedVersion != "" {
		if version, err := strconv.Atoi(storedVersion); err == nil {
			newVersion = version + 1
		}
	}

	token, err := j.generateToken(userId, newVersion)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}

	err = j.rdb.Set(ctx, fmt.Sprintf(constsR.UserLoginTokenVersion, userId), newVersion, constsR.TokenTTL).Err()
	if err := cache.CheckErr(err); err != nil {
		return "", fmt.Errorf("failed to save token: %v", err)
	}

	err = j.rdb.Set(ctx, fmt.Sprintf(constsR.UserLoginToken, userId), token, constsR.TokenTTL).Err()
	if err := cache.CheckErr(err); err != nil {
		return "", fmt.Errorf("failed to save token version: %v", err)
	}

	return token, nil
}

func (j *JWT) DelToken(ctx context.Context, userId int) {
	err := j.rdb.Del(ctx, fmt.Sprintf(constsR.UserLoginToken, userId)).Err()
	if err := cache.CheckErr(err); err != nil {
		zlogger.Errorf("failed to del user token cache, err: %v", err)
		return
	}

	err = j.rdb.Del(ctx, fmt.Sprintf(constsR.UserLoginTokenVersion, userId)).Err()
	if err := cache.CheckErr(err); err != nil {
		zlogger.Errorf("failed to del user token version cache, err: %v", err)
		return
	}
}
