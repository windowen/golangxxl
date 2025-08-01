package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"gorm.io/gorm"

	constsR "queueJob/pkg/constant/redis"
	"queueJob/pkg/db/mysql"
	"queueJob/pkg/db/table"
	"queueJob/pkg/tools/errs"
	"queueJob/pkg/zlogger"
)

const (
	DelLockKeyScript     = "if redis.call('get', KEYS[1]) == ARGV[1] then return redis.call('del', KEYS[1]) else return 0 end"
	TryGetLockSleepTimes = time.Millisecond * 5
)

type UserCacheInfo struct {
	Id                 int       `json:"id"`
	CountryCode        string    `json:"countryCode"`
	AreaCode           string    `json:"areaCode"`
	Mobile             string    `json:"mobile"`
	Email              string    `json:"email"`
	Nickname           string    `json:"nickname"`
	Avatar             string    `json:"avatar"`
	Sign               string    `json:"sign"`
	Birthday           string    `json:"birthday"`
	Sex                int       `json:"sex"`
	Feeling            int       `json:"feeling"`
	Country            string    `json:"country"`
	Area               string    `json:"area"`
	Profession         int       `json:"profession"`
	Category           int       `json:"category"`
	InviteCode         string    `json:"inviteCode"`
	ParentId           int       `json:"parentId"`
	LevelId            int       `json:"levelId"`    // 自然等级
	SetLevelId         int       `json:"setLevelId"` // 后台设置等级
	Remark             string    `json:"remark"`
	Status             int       `json:"status"`
	Password           string    `json:"password"`
	PayPassword        string    `json:"payPassword"`
	RoomId             int       `json:"roomId"`
	ChatUuid           string    `json:"chatUuid"`
	GmStatus           int       `json:"gmStatus"`       // 超管状态 1- 开启 2-关闭
	IsFamilyMaster     int       `json:"isFamilyMaster"` // 是否是家族长
	FamilyId           int       `json:"familyId"`       // 家族id
	FamilyMasterId     int       `json:"familyMasterId"` // 家族长id
	MountsId           int       `json:"mountsId"`       // 坐骑id
	RegisterTime       time.Time `json:"registerTime"`   // 注册时间
	LoginIp            string    `json:"loginIp"`        // 登录ip
	NobleId            int       `json:"nobleId"`        // 贵族id
	AgentRebateStatus  int       `json:"agentRebateStatus"`
	InviteRebateStatus int       `json:"inviteRebateStatus"`
	BetStatus          int       `json:"betStatus"`
	DrawStatus         int       `json:"drawStatus"`
	ChannelCode        int       `json:"channelCode"`
	ChannelName        string    `json:"channelName"`   // 渠道商名字
	FirstRecharge      int       `json:"firstRecharge"` // 首充金额 0 没有首充过
}

type RoomCacheInfo struct {
	Id                int    `json:"id"`
	CountryCode       string `json:"countryCode"`
	UserId            int    `json:"userId"`
	Title             string `json:"title"`
	Tags              string `json:"tags"` // 直播间标签
	Cover             string `json:"cover"`
	VideoClarity      int    `json:"videoClarity"`
	PayRules          int    `json:"payRules"` // 付费规则 1-免费 2-分钟付费 3-场次付费
	TrialDuration     int    `json:"trialDuration"`
	UnitPrice         int    `json:"unitPrice"`
	GiftRatio         int    `json:"giftRatio"`
	PlatformRatio     int    `json:"platformRatio"`
	FamilyRatio       int    `json:"familyRatio"`
	Status            int    `json:"status"`
	ChatRoomId        string `json:"chatRoomId"` // 是播间临时聊天室id
	LiveStatus        int    `json:"liveStatus"`
	Summary           string `json:"summary"`        // 简介
	SceneHistoryId    int    `json:"sceneHistoryId"` // 场次id
	GameId            int    `json:"gameId"`
	StartLiveTime     int    `json:"startLiveTime"`
	PaidPurviewStatus int    `json:"paidPurviewStatus"` // 付费权限
	IsOpenVibrator    int    `json:"isOpenVibrator"`    // 是否开启跳蛋
	Sort              int    `json:"sort"`              // 置顶排序
	BottomSort        int    `json:"bottomSort"`        // 置底排序
	LastStartLiveTime int    `json:"lastStartLiveTime"` // 最后开播时间
	LastEndLiveTime   int    `json:"lastEndLiveTime"`   // 最后下播时间
	LevelId           int    `json:"levelId"`           // 自然等级
	SetLevelId        int    `json:"setLevelId"`        // 后台设置等级
	IsOpenWheel       int    `json:"isOpenWheel"`       // 是否开启转盘
	WheelConfig       string `json:"wheelConfig"`       // 转盘配置
	RobotLastJoinTime int    `json:"robotLastJoinTime"` // 机器人最后加入时间
}

type WagerCacheInfo struct {
	UserId             int             `json:"userId"`             // 用户id
	WagerCode          string          `json:"wagerCode"`          // 赌注id
	ChangeAmount       decimal.Decimal `json:"changeAmount"`       // 变动金额 法币
	WalletKey          string          `json:"walletKey"`          // 钱包Key
	BatchNum           string          `json:"batchNum"`           // 活动钱包相关的游戏类型
	ActChangeAmount    decimal.Decimal `json:"actChangeAmount"`    // 活动变动金额 法币
	CreatedAt          time.Time       `json:"createdAt"`          // 创建时间
	GameOrderNo        string          `json:"gameOrderNo"`        // 订单号
	Currency           string          `json:"currency"`           // 货币
	Remark             string          `json:"remark"`             // 备注
	GameProvider       string          `json:"gameProvider"`       // 请求里的游戏场馆
	ReqGameType        string          `json:"reqGameType"`        // 请求里的游戏类型
	ReqGameName        string          `json:"reqGameName"`        // 游戏名称
	TradeType          int             `json:"tradeType"`          // 交易类型 1 下注，2 结算
	Status             int             `json:"status"`             // 状态 (1 pending, 2 settled)
	WagerAmount        decimal.Decimal `json:"wagerAmount"`        // 下注金额 法币
	SettlementAmount   decimal.Decimal `json:"settlementAmount"`   // 结算金额 法币
	Exchange           decimal.Decimal `json:"exchange"`           // 汇率
	SettleWalletAmount decimal.Decimal `json:"settleWalletAmount"` // 结算钱包金额 法币
	SettleActAmount    decimal.Decimal `json:"settleActAmount"`    // 结算活动金额 法币
	SettledAt          time.Time       `json:"settledAt"`          // 结算时间
}

func (u *UserCacheInfo) IsEmpty() bool {
	return u.Id == 0 || u == nil
}

// String 方法实现
func (w WagerCacheInfo) String() string {
	data, err := json.MarshalIndent(w, "", "  ")
	if err != nil {
		return fmt.Sprintf("WagerCacheInfo: error marshaling to JSON: %v", err)
	}

	return string(data)
}

// NobleCacheInfo 贵族缓存信息
type NobleCacheInfo struct {
	NobleId          int      `json:"nobleId"`          // 贵族id
	NobleLevel       int      `json:"nobleLevel"`       // 贵族等级
	NobleName        string   `json:"nobleName"`        // 贵族等级名称
	Sort             int      `json:"sort"`             // 排序
	FirstChargePrice int      `json:"firstChargePrice"` // 首充价格
	FirstChargeBonus int      `json:"firstChargeBonus"` // 首充赚送
	RenewalPrice     int      `json:"renewalPrice"`     // 续费价格
	RenewalBonus     int      `json:"renewalBonus"`     // 续费赚送
	NobleIcon        string   `json:"nobleIcon"`        // 贵族图标
	AvatarFrameIcon  string   `json:"avatarFrameIcon"`  // 头像框图标
	NobleLabelIcon   string   `json:"nobleLabelIcon"`   // 贵族标签图标
	Privileges       []string `json:"privileges"`       // 特权key
	ExclusiveVehicle int      `json:"exclusiveVehicle"` // 座驾id
	LevelUpSpeed     int      `json:"levelUpSpeed"`     // 升级加速(1-100)
	NobleAnimation   string   `json:"nobleAnimation"`   // 动画
	Gifts            string   `json:"gifts"`            // 道具集合
}

func GetUserCache(userId int) (userCache *UserCacheInfo, err error) {
	userCache = &UserCacheInfo{}

	cacheUser, err := GetKey(fmt.Sprintf(constsR.UserCacheInfoKey, userId))
	if err != nil {
		return
	}

	if err = json.Unmarshal([]byte(cacheUser), userCache); err != nil {
		zlogger.Debugf("getUserCache | userId:%v, cacheUser:%v | err: %v", userId, cacheUser, err)
	}

	return
}

// SetUserCache 将用户信息保存到 Redis 中
func SetUserCache(userId int, value *UserCacheInfo) error {
	cacheKey := fmt.Sprintf(constsR.UserCacheInfoKey, userId)

	marshal, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to save user data to redis: %v", err)
	}

	if err = Set(cacheKey, string(marshal), 0); err != nil {
		return fmt.Errorf("failed to save user data to redis: %v", err)
	}

	return nil
}

func GetRoomCache(roomId int) (roomCache *RoomCacheInfo, err error) {
	roomCache = &RoomCacheInfo{}
	data, err := GetKey(fmt.Sprintf(constsR.RoomCacheInfoKey, roomId))
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(data), roomCache)
	return
}

// SetRoomCache 将房间信息保存到 Redis 中
func SetRoomCache(roomId int, value *RoomCacheInfo) error {
	cacheKey := fmt.Sprintf(constsR.RoomCacheInfoKey, roomId)

	zlogger.Debugw("房间缓存更新", zap.String("CacheKey", cacheKey), zap.Any("RoomCacheInfo", value))

	marshal, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to save room data to redis: %v", err)
	}

	if err = Set(cacheKey, string(marshal), 0); err != nil {
		return fmt.Errorf("failed to save room data to redis: %v", err)
	}

	return nil
}

// GetWagerCache 从 Redis 中获取用户赌注信息
func GetWagerCache(userId int, wagerCode string) (*WagerCacheInfo, error) {
	cacheKey := fmt.Sprintf(constsR.UserWagerInfo, userId, wagerCode)
	wagerCacheInfo := &WagerCacheInfo{}

	// 从 Redis 获取缓存数据
	data, err := GetKey(cacheKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get wager data from redis: %v", err)
	}

	// 反序列化缓存数据
	if err = json.Unmarshal([]byte(data), wagerCacheInfo); err != nil {
		zlogger.Errorw("parse json error",
			zap.Int("userId", userId),
			zap.String("wagerCode", wagerCode),
			zap.String("data", data),
			zap.Error(err))
		return nil, fmt.Errorf("unable to deserialize wager data: %v", err)
	}

	return wagerCacheInfo, nil
}

// DeleteWagerCache 删除 Redis 中的用户赌注信息
func DeleteWagerCache(userId int, wagerCode string) {
	err := DelKey(fmt.Sprintf(constsR.UserWagerInfo, userId, wagerCode))
	if err != nil {
		zlogger.Errorf("failed to delete wager data from redis, err: %v", err)
		return
	}
	return
}

// FindNobleCacheInfo 获取贵族配置信息
func FindNobleCacheInfo(nobleId int) (*NobleCacheInfo, error) {
	guardCache := &NobleCacheInfo{}

	result, err := HGet(constsR.NobleCfgCacheHash, strconv.Itoa(nobleId))
	if err != nil && !errors.Is(err, Nil) {
		zlogger.Errorw("FindNobleCacheInfo", zap.Error(err))
		return nil, err
	}

	if result == "" {
		return nil, errors.New("FindNobleCacheInfo failed to obtain gift information")
	}

	if err = json.Unmarshal([]byte(result), guardCache); err != nil {
		zlogger.Errorw("FindNobleCacheInfo Unmarshal byte", zap.Int("noble id", nobleId), zap.Error(err))
		return nil, err
	}

	return guardCache, nil
}

// TryGetDistributedLock 尝试获取分布式锁
// @param lockKey 锁
// @param requestId 请求标识
// @param expireTime 超时时间(ms)
// @param acquireTimeout 获取锁超时时间(ms)
// @return 是否获取成功 释放锁的函数
func TryGetDistributedLock(lockKey, requestId string, expireTime, acquireTimeout uint32) (bool, func()) {
	endTime := time.Now().UnixNano() + int64(acquireTimeout)*int64(time.Millisecond)
	for time.Now().UnixNano() < endTime {
		result, err := SetNX(lockKey, requestId, time.Millisecond*time.Duration(expireTime))
		if err != nil {
			zlogger.Errorw("tryGetDistributedLock error",
				zap.String("lockKey", lockKey),
				zap.String("requestId", requestId),
				zap.Uint32("expireTime", expireTime),
				zap.Error(err),
			)
			return false, nil
		}
		if result {
			return true, func() {
				if err = Eval(DelLockKeyScript, []string{lockKey}, requestId); err != nil {
					zlogger.Errorw("ReleaseDistributedLock error",
						zap.String("lockName", lockKey),
						zap.String("requestId", requestId),
						zap.Error(err),
					)
					return
				}

				return
			}
		}
		time.Sleep(TryGetLockSleepTimes)
	}
	return false, nil
}

// TryGetBlockChainLock 尝试获取关于区块链操作的分布式锁
func TryGetBlockChainLock(userId int, chainType int) bool {
	key := fmt.Sprintf("blockchain:collect:%d:%d", userId, chainType)
	requestId := time.Now().Format("2006-01-02 15:04:05.000")

	result, err := SetNX(key, requestId, 1*time.Hour)
	if err != nil {
		zlogger.Errorw("TryGetBlockChainLock error",
			zap.String("lockKey", key),
			zap.String("requestId", requestId),
			zap.Error(err),
		)
		return false
	}

	if result {
		return true
	}
	return false
}

func ReleaseBlockChainLock(userId int, chainType int) {
	key := fmt.Sprintf("blockchain:collect:%d:%d", userId, chainType)
	if err := DelKey(key); err != nil {
		zlogger.Errorw("ReleaseDistributedLock error",
			zap.String("lockName", key),
			zap.Error(err),
		)
		return
	}

	return
}

// MarkUserRechargedToday 记录用户今天充值
func MarkUserRechargedToday(key string) {
	expiration := 24 * time.Hour // 过期时间设为 24 小时
	err := Set(key, "1", expiration)
	if err != nil {
		zlogger.Errorw("MarkUserRechargedToday, mark user recharge error",
			zap.String("key", key), zap.Error(err))
	}
}

// HasUserRechargedToday 检查用户今天是否充值过
func HasUserRechargedToday(redisKey string) bool {
	_, err := Get(redisKey)
	if err != nil {
		return false
	}
	return true
}

// MarkKeyToday 设置key并设置超时时间为当天24点
func MarkKeyToday(key string) {
	now := time.Now()
	// 获取第二天零点时间
	expireAt := time.Date(
		now.Year(), now.Month(), now.Day()+1,
		0, 0, 0, 0, now.Location(),
	)
	// 计算从现在到第二天0点的时间间隔
	expiration := expireAt.Sub(now)

	err := Set(key, "1", expiration)
	if err != nil {
		zlogger.Errorw("MarkKeyToday, set key with expiration error",
			zap.String("key", key), zap.Error(err))
	}
}

// HasKeyToday 检查Key是否存在
func HasKeyToday(redisKey string) bool {
	_, err := Get(redisKey)
	if err != nil {
		return false
	}
	return true
}

func GetUserContinueRecharge(day, userId int) bool {
	prevBit, err := GetBit(fmt.Sprintf(constsR.BitmapUserRecharge, userId), int64(day-1))
	if err == nil {
		if prevBit == 0 {
			return false
		}
		return true
	}

	return false
}

func SetUserContinueRecharge(day, userId int) {
	// 记录当天充值
	_, err := SetBit(fmt.Sprintf(constsR.BitmapUserRecharge, userId), int64(day), 1)
	if err != nil {
		zlogger.Errorw("SetUserContinueRecharge, mark user continue recharge error",
			zap.Int("uid", userId), zap.Int("day", day), zap.Error(err))
	}
}

func GetFinanceCoinExchange(ctx context.Context, query interface{}, args ...interface{}) (*table.FinanceCoinExchange, error) {
	var coinExchange table.FinanceCoinExchange

	// 查询是否已有该渠道的今日统计数据
	err := mysql.LiveDB.WithContext(ctx).Where(query, args...).First(&coinExchange).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		zlogger.Errorw("activeData get record error", zap.Error(err))
		return nil, errs.Wrap(err)
	}

	return &coinExchange, nil
}

// FindCoinExchangeByCoinCode 获取汇率
func FindCoinExchangeByCoinCode(ctx context.Context, code string) (*table.FinanceCoinExchange, error) {
	var data table.FinanceCoinExchange
	if err := mysql.LiveDB.WithContext(ctx).Where("to_coin_code=?", code).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.Wrap(gorm.ErrRecordNotFound)
		}
		return nil, errs.Wrap(err)
	}
	return &data, nil
}
