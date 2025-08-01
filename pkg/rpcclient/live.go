package rpcclient

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"queueJob/pkg/common/config"
	"queueJob/pkg/gozero/discov"
	"queueJob/pkg/gozero/zrpc"
	"queueJob/pkg/protobuf/live"
	"queueJob/pkg/queue"
	"queueJob/pkg/tools/cast"
	"queueJob/pkg/tools/errs"
	"queueJob/pkg/tools/mw"
	"queueJob/pkg/zlogger"
)

type LiveClient struct {
	live.LiveServerClient
}

func newLiveClient() *LiveClient {
	rpcKey := strings.ToLower(fmt.Sprintf("%s:///%s", config.Config.Etcd.Schema, config.Config.RpcName.LiveRPCName))
	cl := zrpc.MustNewClient(zrpc.RpcClientConf{
		Etcd: discov.EtcdConf{
			Hosts: config.Config.Etcd.Addr,
			Key:   rpcKey,
		},
	}, zrpc.WithDialOption(mw.AddUserType()), zrpc.WithDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())))

	return &LiveClient{live.NewLiveServerClient(cl.Conn())}
}

// CreateUserChatUuid 生成用户聊天室id
func (l *LiveClient) CreateUserChatUuid(ctx context.Context, UserId int) (string, error) {
	res, err := l.CreateUser(ctx, &live.CreateUserReq{
		Username: strconv.Itoa(UserId),
		Password: strconv.Itoa(UserId),
	})
	if err := errs.Unwrap(err); err != nil {
		return "", err
	}

	return res.Uuid, nil
}

// InsufficientBalance 余额不足通知
func (l *LiveClient) InsufficientBalance(ctx context.Context, userId, roomId int) error {
	_, err := l.InsufficientBalanceNotice(ctx, &live.InsufficientBalanceNoticeReq{
		UserId: cast.ToInt32(userId),
		RoomId: cast.ToInt32(roomId),
	})
	if err := errs.Unwrap(err); err != nil {
		zlogger.Errorf("InsufficientBalance InsufficientBalanceNotice |userId:%v| err: %v", userId, err)

		return err
	}

	return nil
}

// UpgradeNotifyWrap 升级通知
func (l *LiveClient) UpgradeNotifyWrap(ctx context.Context, userId, level, roomId int) error {
	_, err := l.UpgradeNotify(ctx, &live.UpgradeNotifyReq{
		RoomId: cast.ToInt32(roomId),
		UserId: cast.ToInt32(userId),
		Level:  cast.ToInt32(level),
	})

	if err = errs.Unwrap(err); err != nil {
		zlogger.Errorw("UpgradeNotifyWrap error", zap.Int("userId", userId), zap.Error(err))
		return err
	}

	return nil
}

// LiveMinutePaidIncomeNotify 直播间分钟扣费收益通知
func (l *LiveClient) LiveMinutePaidIncomeNotify(ctx context.Context, roomId, anchorId, userId, amount int, chatRoomId string) error {
	_, err := l.LiveMinutePaidIncome(ctx, &live.LiveMinutePaidIncomeReq{
		RoomId:     cast.ToInt32(roomId),
		ChatRoomId: chatRoomId,
		UserId:     cast.ToInt32(userId),
		AnchorId:   cast.ToInt32(anchorId),
		Amount:     cast.ToInt32(amount),
	})

	if err = errs.Unwrap(err); err != nil {
		zlogger.Errorw("LiveMinutePaidIncomeNotify error", zap.Int("userId", userId), zap.Error(err))
		return err
	}

	return nil
}

// RobotJoinChatRoom 加入房间
func (l *LiveClient) RobotJoinChatRoom(ctx context.Context, userId, roomId int) error {
	_, err := l.RobotJoinRoom(ctx, &live.RobotJoinRoomReq{
		RoomId: cast.ToInt32(roomId),
		UserId: cast.ToInt32(userId),
	})

	if err = errs.Unwrap(err); err != nil {
		zlogger.Errorf("RobotJoinChatRoom RobotJoinRoom |userId:%v,roomId:%v| err: %v", userId, roomId, err)
		return err
	}

	return nil
}

// RobotQuitChatRoom 退出房间
func (l *LiveClient) RobotQuitChatRoom(ctx context.Context, userId, roomId int, leaveNotify bool) error {
	_, err := l.InternalLeaveRoom(ctx, &live.LeaveRoomReq{
		RoomId:      cast.ToInt32(roomId),
		UserId:      cast.ToInt32(userId),
		LeaveNotify: cast.ToInt32(leaveNotify),
	})

	//	opts := options.Find().SetSort(bson.D{{"createtime", bson.D{{"$numberLong", -1}}}}).SetLimit(limit).SetSkip(skip)

	if err = errs.Unwrap(err); err != nil {
		zlogger.Errorf("RobotQuitChatRoom InternalLeaveRoom |userId:%v,chatRoomId:%v| err: %v", userId, roomId, err)
		return err
	}

	return nil
}

// UserJoinRoom 加入房间
func (l *LiveClient) UserJoinRoom(ctx context.Context, req *queue.LiveRoomUserJoinReq) error {
	_, err := l.QueueUserJoinRoom(ctx, &live.QueueUserJoinRoomReq{
		RoomId:     cast.ToInt32(req.RoomId),
		UserId:     cast.ToInt32(req.UserId),
		IsTraveler: cast.ToInt32(req.IsTraveler),
	})

	if err = errs.Unwrap(err); err != nil {
		zlogger.Errorf("UserJoinRoom QueueUserJoinRoom |userId:%v,roomId:%v| err: %v", req.UserId, req.RoomId, err)
		return err
	}

	return nil
}
