package consumer

import (
	"bytes"
	"context"
	"encoding/json"
	"text/template"
	"time"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"go.uber.org/zap"

	"queueJob/pkg/db/mysql"
	"queueJob/pkg/db/table"
	"queueJob/pkg/message"
	"queueJob/pkg/zlogger"
)

var notify = &siteNotify{}

type siteNotify struct{}

func (o *siteNotify) handleMessages(ctx context.Context, msgList ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	for _, msg := range msgList {
		zlogger.Infow("siteNotify::handleMessages, notification message",
			zap.String("msgID", msg.MsgId), zap.String("topic", msg.Topic), zap.String("body", string(msg.Body)))

		pMessage := &message.PackageMessage{}
		if err := json.Unmarshal(msg.Body, pMessage); err != nil {
			zlogger.Errorw("siteNotify::handleMessages, unmarshal msg fail", zap.String("msgID", msg.MsgId), zap.Error(err))
			continue
		}

		userCache, err := getUserCache(pMessage.UserId)
		if err != nil {
			zlogger.Errorw("siteNotify::handleMessages, get user cache error", zap.String("msgID", msg.MsgId), zap.Error(err))
			continue
		}

		templateData := &table.SysSiteTemplate{}
		if err = mysql.LiveDB.WithContext(ctx).Where("language_code = ? AND category = ?", pMessage.Language, pMessage.MessageType).
			First(templateData).Error; err != nil {
			zlogger.Errorw("siteNotify::handleMessages, get site template error", zap.String("msgID", msg.MsgId), zap.Any("detail", pMessage), zap.Error(err))
			continue
		}

		tpl, err := template.New("notify").Parse(templateData.TemplateContent)
		if err != nil {
			zlogger.Errorw("siteNotify::handleMessages, new template error", zap.String("msgID", msg.MsgId), zap.Error(err))
			continue
		}

		var content bytes.Buffer
		switch pMessage.MessageType {
		case message.RechargeSuccess, message.RechargeFailure, message.WithdrawalSuccess, message.WithdrawalFailure, message.DiamondExchange:
			if err = tpl.Execute(&content, pMessage.RechargeMsg); err != nil {
				zlogger.Errorw("siteNotify::handleMessages, execute template error", zap.String("msgID", msg.MsgId), zap.Error(err))
				continue
			}

			if err = mysql.LiveDB.WithContext(ctx).Create(&table.SiteTransactionMsg{
				Category:    int(pMessage.MessageType),
				UserId:      pMessage.UserId,
				CountryCode: userCache.CountryCode,
				Content:     content.String(),
				CreatedAt:   time.Now(),
			}).Error; err != nil {
				zlogger.Errorw("siteNotify::handleMessages, insert transaction message table fail", zap.String("msgID", msg.MsgId), zap.Any("detail", pMessage), zap.Error(err))
				continue
			}
		}
	}
	return consumer.ConsumeSuccess, nil
}
