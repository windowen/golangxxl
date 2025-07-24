package queue

import "fmt"

// LiveRoomStop 主播下播
type LiveRoomStop struct {
	RoomId   int `json:"room_id"`   // 直播间id
	AnchorId int `json:"anchor_id"` // 主播id
	SceneId  int `json:"scene_id"`  // 直播场次id
}

// LiveRoomUserMinuteDelayPaid 付费直播分钟扣款延迟队列
type LiveRoomUserMinuteDelayPaid struct {
	RoomId   int `json:"roomId"`   // 直播间id
	AnchorId int `json:"anchorId"` // 主播id
	SceneId  int `json:"sceneId"`  // 直播场次id
	UserId   int `json:"userId"`   // 用户id
}

// LiveRoomPayDiamond 直播房间钻石消费
type LiveRoomPayDiamond struct {
	BillNo       string `json:"bill_no"`       // 订单号
	UserId       int    `json:"user_id"`       // 用户ID
	RoomId       int    `json:"room_id"`       // 房间ID
	FamilyId     int    `json:"family_id"`     // 家族ID
	AnchorId     int    `json:"anchor_id"`     // 主播ID
	SceneId      int    `json:"scene_id"`      // 场景历史ID
	Category     int    `json:"category"`      // 收入类型
	ProjectId    int    `json:"project_id"`    // 项目ID
	ProjectNum   int    `json:"project_num"`   // 项目数量
	UnitPrice    int    `json:"unit_price"`    // 单价
	ProjectTotal int    `json:"project_total"` // 总金额
	IsDivide     int    `json:"is_divide"`     // 是否分成
}

type PayInCache struct {
	UserId    int    `json:"userId"`
	OrderId   string `json:"orderId"`
	Code      int    `json:"code"`
	PayUrl    string `json:"payUrl"`
	PayQrCode string `json:"payQrCode"`
}

func (n PayInCache) String() string {
	return fmt.Sprintf("UserId: %d, OrderId: %s, Code: %d, PayUrl: %s, PayQrCode: %s",
		n.UserId, n.OrderId, n.Code, n.PayUrl, n.PayQrCode)
}

// LiveRoomTransferPayDelay 直播转付费
type LiveRoomTransferPayDelay struct {
	RoomId   int `json:"roomId"`   // 直播间id
	AnchorId int `json:"anchorId"` // 主播id
	SceneId  int `json:"sceneId"`  // 直播场次id
}

type StreamerReceiveDiamond struct {
	AnchorId     int `json:"anchorId"`     // 主播id
	RoomId       int `json:"room_id"`      // 房间ID
	ProjectTotal int `json:"projectTotal"` // 总金额
}

type LiveRoomRobotDelayReq struct {
	RoomId  int `json:"roomId"`  // 直播间id
	SceneId int `json:"sceneId"` // 直播场次id
}

func (n StreamerReceiveDiamond) String() string {
	return fmt.Sprintf("AnchorId: %d, ProjectTotal: %d, RoomId: %d", n.AnchorId, n.ProjectTotal, n.RoomId)
}
