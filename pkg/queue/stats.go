package queue

// EventType 定义统计事件的类型
type EventType int

const (
	// EventUserRegistrations 用户注册事件
	EventUserRegistrations EventType = iota + 1
	// EventPageStay 页面停留事件
	EventPageStay
	// EventGameLaunch 游戏启动次数事件
	EventGameLaunch
	// EventGameAward 游戏中奖次数事件
	EventGameAward
	// EventHomepageBannerClicks 首页 Banner 点击事件
	EventHomepageBannerClicks
	// EventRecommendedBannerClicks 推荐 Banner 点击事件
	EventRecommendedBannerClicks
	// EventPopularBannerClicks 热门 Banner 点击事件
	EventPopularBannerClicks
)

// EventTypeName 获取事件类型的名称
func EventTypeName(eventType EventType) string {
	switch eventType {
	case EventUserRegistrations:
		return "用户注册"
	case EventPageStay:
		return "页面停留"
	case EventGameLaunch:
		return "游戏启动"
	case EventGameAward:
		return "游戏中奖"
	case EventHomepageBannerClicks:
		return "首页 Banner 点击"
	case EventRecommendedBannerClicks:
		return "推荐 Banner 点击"
	case EventPopularBannerClicks:
		return "热门 Banner 点击"
	default:
		return "未知事件"
	}
}

// EventStats 统计事件
type EventStats struct {
	EventType EventType `json:"eventType"`           // 事件类型
	UserId    int       `json:"userId,omitempty"`    // 用户id
	PageName  string    `json:"pageName,omitempty"`  // 页面名称
	Timestamp int64     `json:"timestamp,omitempty"` // 操作时间
}
