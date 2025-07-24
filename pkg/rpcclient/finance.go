package rpcclient

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"liveJob/pkg/common/config"
	"liveJob/pkg/gozero/discov"
	"liveJob/pkg/gozero/zrpc"
	"liveJob/pkg/protobuf/finance"
	"liveJob/pkg/tools/cast"
	"liveJob/pkg/tools/mw"
)

type FinanceClient struct {
	finance.FinanceServerClient
}

// 获取 FinanceServer 客户端
func newFinanceClient() *FinanceClient {
	rpcKey := strings.ToLower(fmt.Sprintf("%s:///%s", config.Config.Etcd.Schema, config.Config.RpcName.FinanceRPCName))
	cl := zrpc.MustNewClient(zrpc.RpcClientConf{
		Etcd: discov.EtcdConf{
			Hosts: config.Config.Etcd.Addr,
			Key:   rpcKey,
		},
	}, zrpc.WithDialOption(mw.AddUserType()), zrpc.WithDialOption(mw.GrpcClient()), zrpc.WithDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())))

	return &FinanceClient{finance.NewFinanceServerClient(cl.Conn())}
}

// RoomLiveMinuteDelayPaid 直播间延迟扣费
func (f *FinanceClient) RoomLiveMinuteDelayPaid(ctx context.Context, userId, unitPrice, roomId, anchorId int) error {
	_, err := f.PaymentDiamond(ctx, &finance.PayDiamondReq{
		UserId:     cast.ToInt32(userId),
		Diamond:    cast.ToInt32(unitPrice),
		ChangeType: finance.PaymentType_PT_LiveCharge,
		RoomId:     cast.ToInt32(roomId),
		StreamerId: cast.ToInt32(anchorId),
	})
	return err
}
