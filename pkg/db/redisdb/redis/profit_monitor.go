package redis

import (
	"context"
	"encoding/json"
	"errors"

	"go.uber.org/zap"

	constsR "queueJob/pkg/constant/redis"
	"queueJob/pkg/db/mysql"
	"queueJob/pkg/db/table"
	"queueJob/pkg/zlogger"
)

func GetProfitMonitor(ctx context.Context) *table.GameProfitMonitorConfig {
	// Try to get data from Redis
	cachedData, err := Get(constsR.GameProfitMonitor)
	if errors.Is(err, Nil) {
		// Load data from the database
		config := &table.GameProfitMonitorConfig{}
		if err = mysql.LiveDB.WithContext(ctx).First(config).Error; err != nil {
			zlogger.Errorw("GetProfitMonitor, get db config error", zap.Error(err))
			return nil
		}

		// Cache the data in Redis
		dataBytes, err := json.Marshal(config)
		if err != nil {
			zlogger.Errorw("GetProfitMonitor, parse config error", zap.Error(err))
			return nil
		}

		if err = Set(constsR.GameProfitMonitor, string(dataBytes), 0); err != nil {
			zlogger.Errorw("GetProfitMonitor, set config error", zap.Error(err))
			return nil
		}
		return config
	} else if err != nil {
		zlogger.Errorw("GetProfitMonitor, read redis error", zap.Error(err))
		return nil
	} else {
		config := &table.GameProfitMonitorConfig{}
		if err = json.Unmarshal([]byte(cachedData), config); err != nil {
			zlogger.Errorw("GetProfitMonitor, parse config error", zap.Error(err))
			return nil
		}
		return config
	}
}
