package constant

type FlowType int

const (
	FlowAmountType1             FlowType = 1
	FlowAmountCreateRedEnvelope FlowType = 2
	FlowAmountGrabRedEnvelope   FlowType = 3
	FlowAmountRefundToOwner     FlowType = 4
	FlowAmountSignIn            FlowType = 5
)

func (s FlowType) String() string {
	var suit string
	switch s {
	case FlowAmountType1:
		suit = "积分竞猜"
	case FlowAmountCreateRedEnvelope:
		suit = "创建红包"
	case FlowAmountGrabRedEnvelope:
		suit = "拆红包"
	case FlowAmountRefundToOwner:
		suit = "红包余额退回"
	case FlowAmountSignIn:
		suit = "签到积分"
	}
	return suit
}
