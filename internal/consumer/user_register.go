package consumer

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"go.uber.org/zap"

	"liveJob/pkg/db/mysql"
	"liveJob/pkg/db/table"
	rpcClient "liveJob/pkg/rpcclient"
	"liveJob/pkg/zlogger"
)

var userReg = &userRegister{}

type userRegister struct{}

func (ur *userRegister) handleMessages(ctx context.Context, msgList ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	for _, msg := range msgList {
		zlogger.Debugw("UserRegister Received message", zap.String("msgID", msg.MsgId), zap.String("topic", msg.Topic), zap.String("body", string(msg.Body)))

		details := &table.User{}
		if err := json.Unmarshal(msg.Body, details); err != nil {
			zlogger.Errorw("userRegisterConsumer::handleMessages, unmarshal msg fail", zap.String("msgID", msg.MsgId), zap.Error(err))
			continue
		}

		var user table.User
		err := mysql.LiveDB.WithContext(ctx).First(&user, details.Id).Error
		if err = mysql.CheckErr(err); err != nil {
			zlogger.Errorw("userRegisterConsumer::handleMessages, get user error", zap.String("msgID", msg.MsgId), zap.String("userId", strconv.Itoa(details.Id)), zap.Error(err))
			continue
		}

		if user.IsEmpty() || user.ChatUuid != "" {
			continue
		}

		// 生成 uuid
		uuid, err := rpcClient.ServiceClientsInstance.LiveClient.CreateUserChatUuid(ctx, details.Id)
		if err != nil {
			zlogger.Errorw("userRegisterConsumer::CreateUserChatUuid", zap.String("msgID", msg.MsgId), zap.Error(err))
			continue
		}

		// 保存
		if err = mysql.LiveDB.WithContext(ctx).Model(details).Update("chat_uuid", uuid).Error; err != nil {
			zlogger.Errorw("userRegisterConsumer::Update", zap.String("msgID", msg.MsgId), zap.Any("detail", details), zap.Error(err))
			continue
		}
	}

	return consumer.ConsumeSuccess, nil
}
