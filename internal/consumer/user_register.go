package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"strconv"

	"go.uber.org/zap"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"queueJob/pkg/db/redisdb/redis"

	"queueJob/pkg/db/mysql"
	"queueJob/pkg/db/table"
	rpcClient "queueJob/pkg/rpcclient"
	"queueJob/pkg/zlogger"
)

var userReg = &userRegister{}

type userRegister struct{}

func (ur *userRegister) handleMessages(ctx context.Context, msgList ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	defer func() {
		if err := recover(); err != nil {
			var buf [4096]byte
			n := runtime.Stack(buf[:], false)
			tmpStr := fmt.Sprintf("err=%v panic ==> %s\n", err, string(buf[:n]))
			zlogger.Error(tmpStr) // 记录到日志
		}
	}()

	for _, msg := range msgList {
		details := &table.User{}
		if err := json.Unmarshal(msg.Body, details); err != nil {
			zlogger.Errorw("userRegisterConsumer::handleMessages, unmarshal msg fail", zap.String("msgID", msg.MsgId), zap.Error(err))
			continue
		}

		user := &table.User{}
		err := mysql.LiveDB.WithContext(ctx).First(user, details.Id).Error
		if err = mysql.CheckErr(err); err != nil {
			zlogger.Errorw("userRegisterConsumer::handleMessages, get user error", zap.String("msgID", msg.MsgId), zap.String("userId", strconv.Itoa(details.Id)), zap.Error(err))
			continue
		}

		if user.Id == 0 || user.ChatUuid != "" {
			continue
		}

		// 生成 uuid
		uuid, err := rpcClient.ServiceClientsInstance.LiveClient.CreateUserChatUuid(ctx, details.Id)
		if err != nil {
			zlogger.Errorw("userRegisterConsumer::CreateUserChatUuid", zap.String("msgID", msg.MsgId), zap.Int("UserId", details.Id), zap.Error(err))
			continue
		}

		// 保存
		if err = mysql.LiveDB.WithContext(ctx).Model(details).Update("chat_uuid", uuid).Error; err != nil {
			zlogger.Errorw("userRegisterConsumer::Update", zap.String("msgID", msg.MsgId), zap.Int("UserId", details.Id), zap.Any("detail", details), zap.Error(err))
			continue
		}

		// 获取用户缓存
		userCache, err := redis.GetUserCache(user.Id)
		if err := redis.CheckErr(err); err != nil {
			zlogger.Errorw("userRegisterConsumer, get user cache error", zap.String("msgID", msg.MsgId), zap.Int("UserId", user.Id), zap.Error(err))
			continue
		}

		if !userCache.IsEmpty() {
			// 更新声网id
			userCache.ChatUuid = uuid
			if err = redis.SetUserCache(user.Id, userCache); err != nil {
				zlogger.Errorw("userRegisterConsumer, update user cache level error",
					zap.String("msgID", msg.MsgId), zap.Int("userId", user.Id), zap.Error(err))
				continue
			}
		}
	}

	return consumer.ConsumeSuccess, nil
}
