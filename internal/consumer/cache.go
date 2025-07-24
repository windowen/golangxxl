package consumer

import (
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"

	constsR "queueJob/pkg/constant/redis"
	"queueJob/pkg/db/redisdb/redis"
	"queueJob/pkg/zlogger"
)

const (
	DelLockKeyScript     = "if redis.call('get', KEYS[1]) == ARGV[1] then return redis.call('del', KEYS[1]) else return 0 end"
	TryGetLockSleepTimes = time.Millisecond * 5
)

type UserCacheInfo struct {
	Id             int    `json:"id"`
	CountryCode    string `json:"countryCode"`
	AreaCode       string `json:"areaCode"`
	Mobile         string `json:"mobile"`
	Email          string `json:"email"`
	Nickname       string `json:"nickname"`
	Avatar         string `json:"avatar"`
	Sign           string `json:"sign"`
	Birthday       string `json:"birthday"`
	Sex            int    `json:"sex"`
	Feeling        int    `json:"feeling"`
	Country        string `json:"country"`
	Area           string `json:"area"`
	Profession     int    `json:"profession"`
	Category       int    `json:"category"`
	InviteCode     string `json:"inviteCode"`
	ParentId       int    `json:"parentId"`
	LevelId        int    `json:"levelId"`    // 自然等级
	SetLevelId     int    `json:"setLevelId"` // 后台设置等级
	Remark         string `json:"remark"`
	Status         int    `json:"status"`
	Password       string `json:"password"`
	PayPassword    string `json:"payPassword"`
	RoomId         int    `json:"roomId"`
	ChatUuid       string `json:"chatUuid"`
	GmStatus       int    `json:"gmStatus"`       // 超管状态 1- 开启 2-关闭
	IsFamilyMaster int    `json:"isFamilyMaster"` // 是否是家族长
	FamilyId       int    `json:"familyId"`       // 家族id
	FamilyMasterId int    `json:"familyMasterId"` // 家族长id
	MountsId       int    `json:"mountsId"`       // 坐骑id
}

type RoomCacheInfo struct {
	Id                int    `json:"id"`
	CountryCode       string `json:"countryCode"`
	UserId            int    `json:"userId"`
	Title             string `json:"title"`
	Tags              string `json:"tags"` // 直播间标签
	Cover             string `json:"cover"`
	VideoClarity      int    `json:"videoClarity"`
	PayRules          int    `json:"payRules"`
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
	RobotLastJoinTime int    `json:"robotLastJoinTime"` // 机器人最后加入时间
}

func getUserCache(userId int) (userCache *UserCacheInfo, err error) {
	userCache = &UserCacheInfo{}

	cacheUser, err := redis.GetKey(fmt.Sprintf(constsR.UserCacheInfoKey, userId))
	if err != nil {
		return
	}

	if err = json.Unmarshal([]byte(cacheUser), userCache); err != nil {
		zlogger.Errorf("getUserCache | userId:%v, cacheUser:%v | err: %v", userId, cacheUser, err)
	}
	return
}

// setUserCache 将用户信息保存到 Redis 中
func setUserCache(userId int, value *UserCacheInfo) error {
	cacheKey := fmt.Sprintf(constsR.UserCacheInfoKey, userId)

	marshal, err := json.Marshal(value)
	if err != nil {
		zlogger.Errorf("marshal user cache info error: %v", err)
		return fmt.Errorf("failed to save user data to redis: %v", err)
	}

	if err = redis.Set(cacheKey, string(marshal), 0); err != nil {
		zlogger.Errorf("failed to save user data to redis failed, err: %v", err)
		return fmt.Errorf("failed to save user data to redis: %v", err)
	}

	return nil
}

func getRoomCache(roomId int) (roomCache *RoomCacheInfo, err error) {
	roomCache = &RoomCacheInfo{}
	data, err := redis.GetKey(fmt.Sprintf(constsR.RoomCacheInfoKey, roomId))
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(data), roomCache)
	return
}

// setRoomCache 将房间信息保存到 Redis 中
func setRoomCache(roomId int, value *RoomCacheInfo) error {
	cacheKey := fmt.Sprintf(constsR.RoomCacheInfoKey, roomId)

	marshal, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to save room data to redis: %v", err)
	}

	if err = redis.Set(cacheKey, string(marshal), 0); err != nil {
		return fmt.Errorf("failed to save room data to redis: %v", err)
	}

	return nil
}

// TryGetDistributedLock 尝试获取分布式锁
// @param lockKey 锁
// @param requestId 请求标识
// @param expireTime 超时时间(ms)
// @param acquireTimeout 获取锁超时时间(ms)
// @return 是否获取成功 释放锁的函数
func tryGetDistributedLock(lockKey, requestId string, expireTime, acquireTimeout uint32) (bool, func()) {
	endTime := time.Now().UnixNano() + int64(acquireTimeout)*int64(time.Millisecond)
	for time.Now().UnixNano() < endTime {
		result, err := redis.SetNX(lockKey, requestId, time.Millisecond*time.Duration(expireTime))
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
				if err = redis.Eval(DelLockKeyScript, []string{lockKey}, requestId); err != nil {
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
