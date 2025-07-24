package kafkaconsumer

import (
	"context"
	"encoding/json"

	"queueJob/pkg/queue"
	rpcClient "queueJob/pkg/rpcclient"
	"queueJob/pkg/zlogger"
)

var ujr = &userJoinRoom{}

type userJoinRoom struct{}

func (o *userJoinRoom) handleMessages(msg []byte) {
	ctx := context.Background()
	req := &queue.LiveRoomUserJoinReq{}
	if err := json.Unmarshal(msg, req); err != nil {
		zlogger.Errorf("userJoinRoom handleMessages | err: %v", err)
		return
	}

	zlogger.Debugf("userJoinRoom| handleMessages |req:%v|", req)

	// 发送通知
	if err := rpcClient.ServiceClientsInstance.LiveClient.UserJoinRoom(ctx, req); err != nil {
		zlogger.Errorf("userJoinRoom UserJoinRoom | err: %v", err)
		return
	}
}
