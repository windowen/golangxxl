package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	constsR "liveJob/pkg/constant/redis"
	"liveJob/pkg/zlogger"
)

type RoomCache struct {
	rdb     redis.UniversalClient // Redis客户端
	expires time.Duration         // 缓存过期时间(默认0不过期)
}

// NewRoomCache 创建一个 RoomCache 实例
func NewRoomCache(rdb redis.UniversalClient) *RoomCache {
	return &RoomCache{
		rdb:     rdb,
		expires: 0 * time.Second,
	}
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
}

func NewRoomCacheInfo(
	Id int,
	countryCode string,
	userId int,
	title string,
	cover string,
	videoClarity int,
	giftRatio int,
	platformRatio int,
	familyRatio int,
	status int,
	chatRoomId string,
	liveStatus int,
	summary string,
	sceneHistoryId int,
	gameId int,
	startLiveTime int,
) *RoomCacheInfo {
	return &RoomCacheInfo{
		Id:             Id,
		CountryCode:    countryCode,
		UserId:         userId,
		Title:          title,
		Cover:          cover,
		VideoClarity:   videoClarity,
		GiftRatio:      giftRatio,
		PlatformRatio:  platformRatio,
		FamilyRatio:    familyRatio,
		Status:         status,
		ChatRoomId:     chatRoomId,
		LiveStatus:     liveStatus,
		Summary:        summary,
		SceneHistoryId: sceneHistoryId,
		GameId:         gameId,
		StartLiveTime:  startLiveTime,
	}
}

// Save 将房间信息保存到 Redis 中
func (rc *RoomCache) Save(ctx context.Context, roomId int, value *RoomCacheInfo) error {
	cacheKey := fmt.Sprintf(constsR.RoomCacheInfoKey, roomId)

	marshal, err := json.Marshal(value)
	if err != nil {
		zlogger.Errorf("marshal user cache info error: %v", err)
		return fmt.Errorf("failed to save user data to redis: %v", err)
	}

	err = rc.rdb.Set(ctx, cacheKey, marshal, rc.expires).Err()
	if err != nil {
		zlogger.Errorf("failed to save room data to redis failed, err: %v", err)
		return fmt.Errorf("failed to save room data to redis: %v", err)
	}

	return nil
}

// Update 更新 Redis 中的用户信息
func (rc *RoomCache) Update(ctx context.Context, roomId int, value *RoomCacheInfo) error {
	return rc.Save(ctx, roomId, value)
}

// Delete 删除 Redis 中的用户信息
func (rc *RoomCache) Delete(ctx context.Context, roomId int) error {
	err := rc.rdb.Del(ctx, fmt.Sprintf(constsR.RoomCacheInfoKey, roomId)).Err()
	if err != nil {
		zlogger.Errorf("deletion of room data failed, err: %v", err)
		return fmt.Errorf("deletion of room data failed: %v", err)
	}
	return nil
}

// Get 从 Redis 中获取房间信息
func (rc *RoomCache) Get(ctx context.Context, roomId int) (*RoomCacheInfo, error) {
	var roomCacheInfo *RoomCacheInfo

	data, err := rc.rdb.Get(ctx, fmt.Sprintf(constsR.RoomCacheInfoKey, roomId)).Result()
	if err != nil {
		zlogger.Errorf("failed to get room data from redis, err: %v", err)
		return nil, fmt.Errorf("failed to get room data db redis: %v", err)
	}

	if err := json.Unmarshal([]byte(data), &roomCacheInfo); err != nil {
		zlogger.Errorf("unable to deserialize room data failed, err: %v", err)
		return nil, fmt.Errorf("unable to deserialize room data: %v", err)
	}

	return roomCacheInfo, nil
}

// Renew 更新缓存有效期
func (rc *RoomCache) Renew(ctx context.Context, roomId int) error {
	err := rc.rdb.Expire(ctx, fmt.Sprintf(constsR.RoomCacheInfoKey, roomId), rc.expires).Err()
	if err != nil {
		zlogger.Errorf("failed to set room %v expiration time, err: %v", roomId, err)
		return fmt.Errorf("failed to set room expiration time, err: %v", err)
	}

	return nil
}

// incrInt 整数自增
func (rc *RoomCache) incrInt(ctx context.Context, key string, value int64) (int64, error) {
	return rc.rdb.IncrBy(ctx, key, value).Result()
}

// decrInt 整数自减
func (rc *RoomCache) decrInt(ctx context.Context, key string, value int64) (int64, error) {
	return rc.rdb.DecrBy(ctx, key, value).Result()
}

type RoomUpdateOption func(*RoomCacheInfo)

// UpdateRoomOption 更新部分缓存
func (rc *RoomCache) UpdateRoomOption(ctx context.Context, roomId int, opts ...RoomUpdateOption) error {
	// 获取现有的 RoomCacheInfo
	roomInfo, err := rc.Get(ctx, roomId)
	if err != nil {
		return err
	}

	// 应用更新选项
	for _, opt := range opts {
		opt(roomInfo)
	}

	return rc.Save(ctx, roomId, roomInfo)
}

func UpdateRoomCountryCode(countryCode string) RoomUpdateOption {
	return func(info *RoomCacheInfo) {
		info.CountryCode = countryCode
	}
}

func UpdateRoomUserId(userId int) RoomUpdateOption {
	return func(info *RoomCacheInfo) {
		info.UserId = userId
	}
}

func UpdateRoomTitle(title string) RoomUpdateOption {
	return func(info *RoomCacheInfo) {
		info.Title = title
	}
}

func UpdateRoomCover(cover string) RoomUpdateOption {
	return func(info *RoomCacheInfo) {
		info.Cover = cover
	}
}

func UpdateRoomVideoClarity(videoClarity int) RoomUpdateOption {
	return func(info *RoomCacheInfo) {
		info.VideoClarity = videoClarity
	}
}

func UpdateRoomPayRules(payRules int) RoomUpdateOption {
	return func(info *RoomCacheInfo) {
		info.PayRules = payRules
	}
}

func UpdateRoomTrialDuration(trialDuration int) RoomUpdateOption {
	return func(info *RoomCacheInfo) {
		info.TrialDuration = trialDuration
	}
}

func UpdateRoomUnitPrice(unitPrice int) RoomUpdateOption {
	return func(info *RoomCacheInfo) {
		info.UnitPrice = unitPrice
	}
}

func UpdateRoomGiftRatio(giftRatio int) RoomUpdateOption {
	return func(info *RoomCacheInfo) {
		info.GiftRatio = giftRatio
	}
}

func UpdateRoomPlatformRatio(platformRatio int) RoomUpdateOption {
	return func(info *RoomCacheInfo) {
		info.PlatformRatio = platformRatio
	}
}

func UpdateRoomFamilyRatio(familyRatio int) RoomUpdateOption {
	return func(info *RoomCacheInfo) {
		info.FamilyRatio = familyRatio
	}
}

func UpdateRoomStatus(status int) RoomUpdateOption {
	return func(info *RoomCacheInfo) {
		info.Status = status
	}
}

func UpdateRoomLiveStatus(status int) RoomUpdateOption {
	return func(info *RoomCacheInfo) {
		info.LiveStatus = status
	}
}

func UpdateRoomSummary(summary string) RoomUpdateOption {
	return func(info *RoomCacheInfo) {
		info.Summary = summary
	}
}

func UpdateRoomSceneHistoryId(sceneHistoryId int) RoomUpdateOption {
	return func(info *RoomCacheInfo) {
		info.SceneHistoryId = sceneHistoryId
	}
}

func UpdateRoomGameId(gameId int) RoomUpdateOption {
	return func(info *RoomCacheInfo) {
		info.GameId = gameId
	}
}

func UpdateRoomStartLiveTime(startLiveTime int) RoomUpdateOption {
	return func(info *RoomCacheInfo) {
		info.StartLiveTime = startLiveTime
	}
}
