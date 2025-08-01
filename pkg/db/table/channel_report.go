package table

import (
	"time"

	"github.com/shopspring/decimal"
)

type ChannelReport struct {
	Id                         int             `gorm:"column:id" json:"id"`                                                       // 渠道报表Id
	ChannelPartnerId           int             `gorm:"column:channel_partner_id" json:"channel_partner_id"`                       // 渠道商ID（对应channel_partner表id）
	ChannelCode                int             `gorm:"column:channel_code" json:"channel_code"`                                   // 渠道code，对应channel_manage表的code
	ChannelPartnerName         string          `gorm:"column:channel_partner_name" json:"channel_partner_name"`                   // 渠道商名称
	ChannelUrl                 string          `gorm:"column:channel_url" json:"channel_url"`                                     // 渠道链接
	CRegisterCount             int             `gorm:"column:c_register_count" json:"c_register_count"`                           // 扣过量的渠道平台(openInstall)统计的新增会员数量
	RegisterCount              int             `gorm:"column:register_count" json:"register_count"`                               // 渠道平台(openInstall)统计的新增会员数量
	RechargeMemberTimes        int             `gorm:"column:recharge_member_times" json:"recharge_member_times"`                 // 充值次数
	CRechargeMemberTimes       int             `gorm:"column:c_recharge_member_times" json:"c_recharge_member_times"`             // 扣过量的充值次数
	ValidRechargeTimes         int             `gorm:"column:valid_recharge_times" json:"valid_recharge_times"`                   // 有效充值次数
	FirstRechargeMemberCount   int             `gorm:"column:first_recharge_member_count" json:"first_recharge_member_count"`     // 首充会员人数
	TodayRegisterRechargeCount int             `gorm:"column:today_register_recharge_count" json:"today_register_recharge_count"` // 当日注册并且充值人数
	RechargeMemberCount        int             `gorm:"column:recharge_member_count" json:"recharge_member_count"`                 // 充值人数
	FirstRechargeMoney         decimal.Decimal `gorm:"column:first_recharge_money" json:"first_recharge_money"`                   // 首充金额
	H5RegisterSum              int             `gorm:"column:h5_register_sum" json:"h5_register_sum"`                             // 真实的每日累计会员数量(h5)
	AndroidRegisterSum         int             `gorm:"column:android_register_sum" json:"android_register_sum"`                   // 真实的每日累计会员数量(android)
	RegisterSum                int             `gorm:"column:register_sum" json:"register_sum"`                                   // 真实的每日累计会员数量(IOS)
	RechargeSumMoney           decimal.Decimal `gorm:"column:recharge_sum_money" json:"recharge_sum_money"`                       // 总充值金额
	CRechargeSumMoney          decimal.Decimal `gorm:"column:c_recharge_sum_money" json:"c_recharge_sum_money"`                   // 扣量后的总充值金额
	WithdrawSumMoney           decimal.Decimal `gorm:"column:withdraw_sum_money" json:"withdraw_sum_money"`                       // 总提现金额
	WithdrawMemberCount        int             `gorm:"column:withdraw_member_count" json:"withdraw_member_count"`                 // 提现人数
	ValidSumBet                decimal.Decimal `gorm:"column:valid_sum_bet" json:"valid_sum_bet"`                                 // 总有效投注
	AndroidDownloadCount       int             `gorm:"column:android_download_count" json:"android_download_count"`               // 安卓总下载数量
	CAndroidDownloadCount      int             `gorm:"column:c_android_download_count" json:"c_android_download_count"`           // 扣过量的显示的安卓总下载数量
	IosDownloadCount           int             `gorm:"column:ios_download_count" json:"ios_download_count"`                       // ios总下载数量
	CIosDownloadCount          int             `gorm:"column:c_ios_download_count" json:"c_ios_download_count"`                   // 扣过量的显示的ios总下载数量
	CAndroidInstallCount       int             `gorm:"column:c_android_install_count" json:"c_android_install_count"`             // 扣过量的显示的安卓安装数
	AndroidInstallCount        int             `gorm:"column:android_install_count" json:"android_install_count"`                 // 安卓安装数
	IosInstallCount            int             `gorm:"column:ios_install_count" json:"ios_install_count"`                         // ios安装数
	CIosInstallCount           int             `gorm:"column:c_ios_install_count" json:"c_ios_install_count"`                     // 扣过量的显示的ios安装数
	Visits                     int             `gorm:"column:visits" json:"visits"`                                               // 渠道访问量
	CVisits                    int             `gorm:"column:c_visits" json:"c_visits"`                                           // 扣过量的显示的渠道访问量
	ActiveUsers                int             `gorm:"column:active_users" json:"active_users"`                                   // 活跃用户数量
	H5ActiveUsers              int             `gorm:"column:h5_active_users" json:"h5_active_users"`                             // h5活跃用户数量
	H5NewDevices               int             `gorm:"column:h5_new_devices" json:"h5_new_devices"`                               // h5新增设备数量
	UniqueDevice               int             `gorm:"column:unique_device" json:"unique_device"`                                 // 活跃设备数量
	RetentionRate              int             `gorm:"column:retention_rate" json:"retention_rate"`                               // 1日留存数量
	RetentionRate7             int             `gorm:"column:retention_rate_7" json:"retention_rate_7"`                           // 7日留存数
	RetentionRate30            int             `gorm:"column:retention_rate_30" json:"retention_rate_30"`                         // 30日留存数
	H5RetentionNum             int             `gorm:"column:h5_retention_num" json:"h5_retention_num"`                           // h5的1日留存数量
	H5RetentionNum7            int             `gorm:"column:h5_retention_num_7" json:"h5_retention_num_7"`                       // h5的7日留存数
	H5RetentionNum30           int             `gorm:"column:h5_retention_num_30" json:"h5_retention_num_30"`                     // h5的30日留存数
	ExchangeDiamondCount       int             `gorm:"column:exchange_diamond_count" json:"exchange_diamond_count"`               // 钻石兑换总额
	ExchangeDiamondTimes       int             `gorm:"column:exchange_diamond_times" json:"exchange_diamond_times"`               // 钻石兑换次数(每个用户只累计一次）
	PlayGameTimes              int             `gorm:"column:play_game_times" json:"play_game_times"`                             // 玩游戏次数(每个用户只累计一次）
	PageStayTime               int             `gorm:"column:page_stay_time" json:"page_stay_time"`                               // 游戏页面停留时间，单位为秒
	GameLaunchCount            int             `gorm:"column:game_launch_count" json:"game_launch_count"`                         // 游戏启动次数
	GameAwardCount             int             `gorm:"column:game_award_count" json:"game_award_count"`                           // 游戏中奖次数
	HomepageBannerClicks       int             `gorm:"column:homepage_banner_clicks" json:"homepage_banner_clicks"`               // 首页 banner 点击次数
	RecommendedBannerClicks    int             `gorm:"column:recommended_banner_clicks" json:"recommended_banner_clicks"`         // 推荐 banner 点击次数
	PopularBannerClicks        int             `gorm:"column:popular_banner_clicks" json:"popular_banner_clicks"`                 // 热门 banner 点击次数
	GameBannerClicks           int             `gorm:"column:game_banner_clicks" json:"game_banner_clicks"`                       // 游戏 banner 点击次数
	CreatedDate                time.Time       `gorm:"column:created_date" json:"created_date"`                                   // 统计日期（如：2025-02-11）
	UpdatedAt                  time.Time       `gorm:"column:updated_at" json:"updated_at"`                                       // 更新时间
}

func (ChannelReport) TableName() string {
	return "channel_report"
}
