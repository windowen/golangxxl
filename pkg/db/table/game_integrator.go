package table

import "time"

type GameIntegrator struct {
	Id               int       `gorm:"column:id" json:"id"`                                   // 游戏供应商配置表ID
	AgentCode        string    `gorm:"column:agent_code" json:"agent_code"`                   // 代理编码
	AgentKey         string    `gorm:"column:agent_key" json:"agent_key"`                     // 代理密钥
	AgentApiKey      string    `gorm:"column:agent_api_key" json:"agent_api_key"`             // 供应商api key
	Code             string    `gorm:"column:code" json:"code"`                               // 场馆编码
	Name             string    `gorm:"column:name" json:"name"`                               // 场馆名称
	LoginUrl         string    `gorm:"column:login_url" json:"login_url"`                     // 登录API地址
	LobbyUrl         string    `gorm:"column:lobby_url" json:"lobby_url"`                     // 游戏大厅地址
	DefaultUrl       string    `gorm:"column:default_url" json:"default_url"`                 // 默认地址
	GameListUrl      string    `gorm:"column:game_list_url" json:"game_list_url"`             // 游戏列表地址
	ProductListUrl   string    `gorm:"column:product_list_url" json:"product_list_url"`       // 产品列表 (3.6Product List)
	GameBrandListUrl string    `gorm:"column:game_brand_list_url" json:"game_brand_list_url"` // 游戏品牌列表地址
	CallbackUrl      string    `gorm:"column:callback_url" json:"callback_url"`               // 回调地址
	Icon             string    `gorm:"column:icon" json:"icon"`                               // 图标
	Status           int       `gorm:"column:status" json:"status"`                           // 状态0启用1禁用3-删除
	IsDel            int       `gorm:"column:is_del" json:"is_del"`                           // 是否删除状态0未删除1已删除
	CreatedAt        time.Time `gorm:"column:created_at" json:"created_at"`                   // 创建时间
	UpdatedAt        time.Time `gorm:"column:updated_at" json:"updated_at"`                   // 更新时间
}

func (GameIntegrator) TableName() string {
	return "game_integrator"
}
