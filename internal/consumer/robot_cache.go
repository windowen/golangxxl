package consumer

import (
	"errors"

	"liveJob/pkg/constant"
	constsR "liveJob/pkg/constant/redis"
	"liveJob/pkg/db/redisdb/redis"
	"liveJob/pkg/tools/cast"
	"liveJob/pkg/tools/strhelper"
	"liveJob/pkg/utils"
	"liveJob/pkg/zlogger"
)

// User 机器人信息
type User struct {
	UserId   int    `json:"userId"`
	Nickname string `json:"nickname"`
	Sex      int    `json:"sex"`
	LevelId  int    `json:"levelId"`
	Avatar   string `json:"avatar"`
}

// RobotConfig 机器人配置
type RobotConfig struct {
	RoomId                  int  `json:"roomId"`                  // 直播间id、全局配置默认0
	IsOpen                  bool `json:"isOpen"`                  // 是否开启配置
	RoomMaxRobots           int  `json:"roomMaxRobots"`           // 房间最大机器人数量
	MinStayTime             int  `json:"minStayTime"`             // 机器人最小停留时间 单位/秒
	MaxStayTime             int  `json:"maxStayTime"`             // 机器人最大停留时间 单位/秒
	JoinIncreaseViewerCount bool `json:"joinIncreaseViewerCount"` // 进入是否增加观众人数
	QuitLessenViewerCount   bool `json:"quitLessenViewerCount"`   // 退出是否减少观众人数
	MinJoinInterval         int  `json:"minJoinInterval"`         // 机器人最小加入时间间隔 单位/秒
	MaxJoinInterval         int  `json:"maxJoinInterval"`         // 机器人最大加入时间间隔 单位/秒
}

// GetRobotConfig 获取机器人配置
func GetRobotConfig(roomId int) *RobotConfig {
	roomRobotCfg := findRobotConfig(roomId)
	if roomRobotCfg == nil {
		roomRobotCfg = findRobotConfig(constant.Zero)
	}

	if roomRobotCfg == nil {
		return nil
	}

	// 是否开启
	if !roomRobotCfg.IsOpen {
		return nil
	}

	return roomRobotCfg
}

func findRobotConfig(roomId int) *RobotConfig {
	result, err := redis.HGet(constsR.RoomRobotConfigHash, cast.ToString(roomId))
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		}

		zlogger.Errorf("findRobotConfig |roomId:%v| err: %v", roomId, err)
		return nil
	}

	robotConfig := &RobotConfig{}
	if err = strhelper.Json2Struct(result, robotConfig); err != nil {
		zlogger.Errorf("GetRobotConfig Json2Struct |roomId:%v| err: %v", roomId, err)
		return nil
	}

	return robotConfig
}

// RandomAvailableRobot 随机获取机器人
func RandomAvailableRobot() (*User, error) {
	result, err := redis.SPop(constsR.RoomRobotSet)
	if err != nil {
		return nil, err
	}

	if result == "" {
		return nil, errors.New("failed to get robot")
	}

	// 获取机器人信息
	userCache, err := getUserCache(cast.ToInt(result))
	if err != nil {
		return nil, err
	}

	return &User{
		UserId:   userCache.Id,
		Nickname: userCache.Nickname,
		Sex:      userCache.Sex,
		LevelId:  utils.CompareMax(userCache.LevelId, userCache.SetLevelId),
		Avatar:   userCache.Avatar,
	}, nil
}
