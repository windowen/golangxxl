package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	constsR "liveJob/pkg/constant/redis"
	"liveJob/pkg/tools/cast"
	"liveJob/pkg/zlogger"
)

type UserCache struct {
	rdb     redis.UniversalClient // Redis客户端
	expires time.Duration         // 缓存过期时间(默认0不过期)
}

// NewUserCache 创建一个 UserCache 实例
func NewUserCache(rdb redis.UniversalClient) *UserCache {
	return &UserCache{
		rdb:     rdb,
		expires: 0 * time.Second,
	}
}

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

func NewUserCacheInfo(id int, countryCode string, areaCode string, mobile string, email string, nickname string, avatar string, sign string, sex int, birthday string, feeling int, country string, area string, profession int, category int, inviteCode string, parentId int, levelId int, setLevelId int, remark string, status int, password string, payPassword string, roomId int, chatUuid string, gmStatus int, isFamilyMaster int, familyId int, familyMasterId int, mountsId int) *UserCacheInfo {
	return &UserCacheInfo{Id: id, CountryCode: countryCode, AreaCode: areaCode, Mobile: mobile, Email: email, Nickname: nickname, Avatar: avatar, Sign: sign, Sex: sex, Birthday: birthday, Feeling: feeling, Country: country, Area: area, Profession: profession, Category: category, InviteCode: inviteCode, ParentId: parentId, LevelId: levelId, SetLevelId: setLevelId, Remark: remark, Status: status, Password: password, PayPassword: payPassword, RoomId: roomId, ChatUuid: chatUuid, GmStatus: gmStatus, IsFamilyMaster: isFamilyMaster, FamilyId: familyId, FamilyMasterId: familyMasterId, MountsId: mountsId}
}

// Save 将用户信息保存到 Redis 中
func (uc *UserCache) Save(ctx context.Context, userId int, value *UserCacheInfo) error {
	cacheKey := fmt.Sprintf(constsR.UserCacheInfoKey, cast.ToString(userId))

	marshal, err := json.Marshal(value)
	if err != nil {
		zlogger.Errorf("marshal user cache info error: %v", err)
		return fmt.Errorf("failed to save user data to redis: %v", err)
	}

	err = uc.rdb.Set(ctx, cacheKey, marshal, uc.expires).Err()
	if err != nil {
		zlogger.Errorf("failed to save user data to redis failed, err: %v", err)
		return fmt.Errorf("failed to save user data to redis: %v", err)
	}

	return nil
}

// Update 更新 Redis 中的用户信息
func (uc *UserCache) Update(ctx context.Context, userId int, value *UserCacheInfo) error {
	return uc.Save(ctx, userId, value)
}

// Delete 删除 Redis 中的用户信息
func (uc *UserCache) Delete(ctx context.Context, userId int) {
	err := uc.rdb.Del(ctx, fmt.Sprintf(constsR.UserCacheInfoKey, cast.ToString(userId))).Err()
	if err != nil {
		zlogger.Errorf("deletion of user data failed, err: %v", err)
	}
}

// Get 从 Redis 中获取用户信息
func (uc *UserCache) Get(ctx context.Context, userId int) (*UserCacheInfo, error) {
	var userCacheInfo *UserCacheInfo

	data, err := uc.rdb.Get(ctx, fmt.Sprintf(constsR.UserCacheInfoKey, cast.ToString(userId))).Result()
	if err := CheckErr(err); err != nil {
		zlogger.Errorf("failed to get user data from redis, err: %v", err)
		return nil, fmt.Errorf("failed to get user data from redis: %v", err)
	}

	if data == "" {
		zlogger.Errorf("roomCacheInfo | err: failed to get user cache")
		return nil, errors.New("failed to get room cache")
	}

	if err := json.Unmarshal([]byte(data), &userCacheInfo); err != nil {
		zlogger.Errorf("userCacheInfo json.Unmarshal |data:%v| err: %v", data, err)
		return nil, fmt.Errorf("unable to deserialize user data: %v", err)
	}

	return userCacheInfo, nil
}

func (uc *UserCache) Renew(ctx context.Context, userId interface{}) error {
	err := uc.rdb.Expire(ctx, fmt.Sprintf(constsR.UserCacheInfoKey, cast.ToString(userId)), uc.expires).Err()
	if err != nil {
		zlogger.Errorf("failed to set user information %v expiration time, err: %v", userId, err)
		return fmt.Errorf("failed to set user information expiration time, err: %v", err)
	}

	return nil
}

// incrInt 整数自增
func (uc *UserCache) incrInt(ctx context.Context, key string, value int64) (int64, error) {
	return uc.rdb.IncrBy(ctx, key, value).Result()
}

// decrInt 整数自减
func (uc *UserCache) decrInt(ctx context.Context, key string, value int64) (int64, error) {
	return uc.rdb.DecrBy(ctx, key, value).Result()
}

type UserUpdateOption func(info *UserCacheInfo)

// UpdateUserOption 更新部分缓存
func (uc *UserCache) UpdateUserOption(ctx context.Context, roomId int, opts ...UserUpdateOption) error {
	// 获取现有的 RoomCacheInfo
	roomInfo, err := uc.Get(ctx, roomId)
	if err != nil {
		return err
	}

	// 应用更新选项
	for _, opt := range opts {
		opt(roomInfo)
	}

	return uc.Save(ctx, roomId, roomInfo)
}

func UpdateUserAvatar(avatar string) UserUpdateOption {
	return func(info *UserCacheInfo) {
		info.Avatar = avatar
	}
}

func UpdateUserNickname(nickname string) UserUpdateOption {
	return func(info *UserCacheInfo) {
		info.Nickname = nickname
	}
}

func UpdateUserSign(sign string) UserUpdateOption {
	return func(info *UserCacheInfo) {
		info.Sign = sign
	}
}

func UpdateUserSex(sex int) UserUpdateOption {
	return func(info *UserCacheInfo) {
		info.Sex = sex
	}
}

func UpdateUserBirthday(birthday string) UserUpdateOption {
	return func(info *UserCacheInfo) {
		info.Birthday = birthday
	}
}

func UpdateUserFeeling(feeling int) UserUpdateOption {
	return func(info *UserCacheInfo) {
		info.Feeling = feeling
	}
}

func UpdateUserCountry(country string) UserUpdateOption {
	return func(info *UserCacheInfo) {
		info.Country = country
	}
}

func UpdateUserArea(area string) UserUpdateOption {
	return func(info *UserCacheInfo) {
		info.Area = area
	}
}

func UpdateUserProfession(profession int) UserUpdateOption {
	return func(info *UserCacheInfo) {
		info.Profession = profession
	}
}

func UpdateUserFamilyId(familyId int) UserUpdateOption {
	return func(info *UserCacheInfo) {
		info.FamilyId = familyId
	}
}

func UpdateUserFamilyMasterId(familyMasterId int) UserUpdateOption {
	return func(info *UserCacheInfo) {
		info.FamilyMasterId = familyMasterId
	}
}

func UpdateUserPayPassword(payPassword string) UserUpdateOption {
	return func(info *UserCacheInfo) {
		info.PayPassword = payPassword
	}
}

func UpdateUserPassword(password string) UserUpdateOption {
	return func(info *UserCacheInfo) {
		info.Password = password
	}
}

func UpdateUserAreaCode(areaCode string, isSetEmpty bool) UserUpdateOption {
	return func(info *UserCacheInfo) {
		if areaCode != "" || isSetEmpty {
			info.AreaCode = areaCode
		}
	}
}

func UpdateUserMobile(mobile string, isSetEmpty bool) UserUpdateOption {
	return func(info *UserCacheInfo) {
		if mobile != "" || isSetEmpty {
			info.Mobile = mobile
		}
	}
}

func UpdateUserEmail(email string, isSetEmpty bool) UserUpdateOption {
	return func(info *UserCacheInfo) {
		if email != "" || isSetEmpty {
			info.Email = email
		}
	}
}

func UpdateUserMountsId(mountsId int) UserUpdateOption {
	return func(info *UserCacheInfo) {
		info.MountsId = mountsId
	}
}
