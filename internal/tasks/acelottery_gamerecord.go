package tasks

import (
	ctxos "context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"queueJob/pkg/constant"
	redis2 "queueJob/pkg/constant/redis"
	"queueJob/pkg/context"
	"queueJob/pkg/db/redisdb/redis"
	"queueJob/pkg/db/table"
	"queueJob/pkg/tools/passwordhelper"
	"queueJob/pkg/tools/strhelper"
	"queueJob/pkg/xxl"
	"queueJob/pkg/zlogger"
)

func AceLotteryGameRecord(ctx *context.Context, _ *xxl.RunReq) (msg string) {
	ctx.Trace = fmt.Sprintf("Ace_Lottery_Game_Record_%s", ctx.Trace)
	nowTime := time.Now()
	currentDate := nowTime.Format(time.DateOnly)
	ctos := *ctx.Ctx

	zlogger.Infof("AceLotteryGameRecord begin %v ", nowTime)

	// 需要先登陆
	gameIntegrator, err := getGameIntegratorByCode(ctos, constant.GameIntegratorAceLotteryCode)
	if err != nil {
		zlogger.Errorf("AceLotteryBetting getGameIntegratorByCode | err:%v | params:%v", err, gameIntegrator)
		return "failed"
	}

	zlogger.Infof("AceLotteryGameRecord integratorCode:%v begin: %v ", currentDate, gameIntegrator)

	endDate := BJNowTime()
	lastDate := endDate - 1000 // 默认从当前时间前推1秒开始

	//lastDateStr, err := redis.Get(redis2.StatsAceLotteryTime)
	//if errors.Is(err, redis.Nil) {
	//	if setStatus := redis.Set(redis2.StatsAceLotteryTime, fmt.Sprintf("%d", endDate), 0); setStatus != nil {
	//		return "failed"
	//	}
	//} else if err != nil {
	//	return "failed"
	//} else {
	//	// redis 中已有 lastDateStr，可以根据需要解析成 int64 使用
	//	lastDate, _ = strconv.ParseInt(lastDateStr, 10, 64)
	//}

	// lastDate = 1750470639000   // 默认从当前时间前推1秒开始
	// lastDate = 1750750621000   // 默认从当前时间前推1秒开始

	fixedStr := fmt.Sprintf("operatorId=%s&secretKey=%s&",
		gameIntegrator.AgentCode,
		gameIntegrator.AgentKey,
	)

	sign := passwordhelper.Md5Encrypt(fixedStr)

	postCurrent := 1

	params := map[string]interface{}{
		"operatorId": gameIntegrator.AgentCode,    // 字符串类型
		"sign":       sign,                        // 数字类型（例如int、int64等）
		"size":       100,                         // 数字类型（例如int、int64等）
		"lang":       "en",                        // 数字类型（例如int、int64等）
		"startTime":  fmt.Sprintf("%d", lastDate), // 数字类型（例如int、int64等）
		"endTime":    fmt.Sprintf("%d", endDate),  // 数字类型（例如int、int64等）
	}

	zlogger.Infof("sys stats data sync mysql %v begin", postCurrent)
	zlogger.Infof("sys stats data sync mysql %v begin", params)

	//
	//for {
	//	// 1. 组装请求体
	//	params["current"] = postCurrent // 把当前页放进去
	//	dataByte, _ := json.Marshal(params)
	//
	//	// 2. 发送请求
	//	result, err := httpclient.ProxyPostJson(
	//		gameIntegrator.LoginUrl+constant.UrlAceGameRecord,
	//		dataByte, nil)
	//	if err != nil {
	//		zlogger.Errorf("request error: %v", err)
	//		return
	//	}
	//
	//	// 3. 抽取分页信息
	//	root := gjson.ParseBytes(result)
	//
	//	rest := gjson.Get(string(result), "code").String()
	//	if rest != constant.AceStatusSuccess && rest != constant.AceStatusSuccess0 {
	//		zlogger.Errorf("AceLotteryBetting | result:%v | ", string(result))
	//		return "failed"
	//	}
	//
	//	data := gjson.Get(string(result), "data")
	//	if data.Exists() {
	//		recordsJson := gjson.Get(string(result), "data").Get("records")
	//		if recordsJson.Exists() {
	//			arrays := gjson.Get(string(result), "data").Get("records").Array()
	//			if len(arrays) == 0 {
	//				zlogger.Errorf("AceLotteryBetting len(arrays) == 0")
	//				return "success"
	//			}
	//			records := parseGameRecords(ctos, arrays)
	//			for _, rOne := range records {
	//				// fmt.Printf("Parsed Record: %+v\n", rOne)
	//				// 进一步处理
	//				kafka.PublicKey(rOne.GameRoundID, kafka.GameRecordTopic, rOne)
	//			}
	//		}
	//
	//	}
	//
	//	current := root.Get("data.current").Int() // 当前页
	//	pages := root.Get("data.pages").Int()     // 总页数
	//
	//	zlogger.Infof("AceLotteryBetting resp:%v begin: %v ", lastDateStr, string(result))
	//
	//	// 5. 判断是否还有下一页
	//	if current >= pages {
	//		break // 已经处理完最后一页
	//	}
	//	postCurrent++ // 否则继续下一页
	//}

	// "total\":105,\"size\":100,\"current\":1,\"pages\":2}} "}

	if setStatus := redis.Set(redis2.StatsAceLotteryTime, fmt.Sprintf("%d", endDate), 0); setStatus != nil {
		return "failed"
	}

	return "success"

}

// getGameIntegratorByCode 获取游戏场馆信息
func getGameIntegratorByCode(ctx ctxos.Context, integratorCode string) (*table.GameIntegrator, error) {
	integratorData := &table.GameIntegrator{}
	key := constant.GameIntegratorHashKey

	data, err := redis.HGet(key, constant.GameIntegratorAceLotteryCode)

	if err == nil {
		// Redis 命中缓存，直接返回
		if err = strhelper.Json2Struct(data, integratorData); err != nil {
			zlogger.Errorf("getGameIntegratorByCode Json2Struct |roomId:%v| err: %v", key, err)
			return integratorData, nil
		}
		return integratorData, nil
	}

	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("redis error: %w", err)
	}

	zlogger.Infof("sys stats data sync mysql %v begin", data)

	//db := mysql.LiveDB.WithContext(ctx)
	//
	//err = db.Where("code = ?", integratorCode).First(&integratorData).Error
	//
	//if err != nil {
	//	zlogger.Warnf("getGameIntegratorByCode query data failed | integratorCode:%v | err:%v", integratorCode, err)
	//}

	if integratorData == nil {
		return nil, errors.New("query data is empty")
	}

	// 假设 integratorData 是 *table.GameIntegrator
	jsonData, err := json.Marshal(integratorData)
	if err != nil {
		zlogger.Errorf("json.Marshal error: %v", err)
		return nil, fmt.Errorf("failed to set cache: %w", err)
	}

	// err = redis.HSet(ctx, key, integratorCode, jsonData).Err()
	// if err != nil {
	// 	zlogger.Errorf("HSet redis error: %v", err)
	// }

	if err = redis.HSet(key, integratorCode, jsonData); err != nil {
		zlogger.Errorf("json.Marshal error: %v", err)
		return nil, fmt.Errorf("failed to set cache: %w", err)
	}

	return integratorData, nil
}

// BJNowTime 北京当前时间
func BJNowTime() int64 {
	// 获取北京时间, 在 windows系统上 time.LoadLocation 会加载失败, 最好的办法是用 time.FixedZone, es 中的时间为: "2019-03-01T21:33:18+08:00"
	var beiJinLocation *time.Location
	var err error

	beiJinLocation, err = time.LoadLocation("Asia/Shanghai")
	if err != nil {
		beiJinLocation = time.FixedZone("CST", 8*3600)
	}

	nowTime := time.Now().In(beiJinLocation)

	// nowTime := time.Now()
	currentMillis := nowTime.UnixNano() / int64(time.Millisecond)
	// millisStr := fmt.Sprintf("%d", currentMillis)

	fmt.Println("Current time in milliseconds:", currentMillis)
	fmt.Println("Current time in nowTime:", nowTime)

	return currentMillis
}

//
//func parseGameRecords(ctx ctxos.Context, arrays []gjson.Result) []elastic.GameRecordES {
//	var records []elastic.GameRecordES
//
//	for _, dataOne := range arrays {
//		var record elastic.GameRecordES
//
//		// userId
//		var winMoney64, exchange64 float64
//		if val := dataOne.Get("operatorAccount"); val.Exists() {
//			if val := dataOne.Get("operatorAccount").String(); val != "" {
//				// record.UserID, err = strconv.ParseInt(val, 10, 64)
//				num, err := strconv.Atoi(val)
//				if err != nil {
//					zlogger.Warnf("operatorAccount 字段无法转为 int: %s", val)
//				} else {
//					record.UserID = int64(num)
//					// 获取用户信息
//					userCacheInfo, err := redis.GetUserCache(num)
//					if err != nil {
//						zlogger.Errorf("AceLotteryBetting userCache.Get err:%+v", err)
//						return nil
//					}
//					record.UserName = userCacheInfo.Nickname
//					currency, err := constant.GetCurrencyByCountry(userCacheInfo.CountryCode)
//					if err != nil {
//						zlogger.Errorf("AceLotteryBetting getUserCurrency err:%+v", err)
//						return nil
//					}
//
//					// todo 财务模块需要 下注扣款 放在回调中 存投注记录
//					exchange, err := redis.FindCoinExchangeByCoinCode(ctx, currency)
//					if err != nil {
//						zlogger.Errorw("get exchange error", zap.String("currency", currency), zap.Error(err))
//						return nil
//					}
//
//					// 游戏记录中使用 用户注册地的法币符号。但是投注中用的美元
//					record.Currency = currency
//
//					exchange64 = exchange.ExchangeRate.Truncate(4).InexactFloat64()
//					record.ExchangeRate = exchange64
//
//				}
//				// 用户投注状态
//
//			} else {
//				zlogger.Warnf("Missing userId for serialNumber: %s", record.UserID)
//				continue
//			}
//		}
//		winInt := 0
//		record.GameStatus = 2
//
//		// 是否中獎(0-未中獎，1-已中獎 ，2-和局) SettledAmount     float64 `json:"settled_amount"`       // 结算金额
//		// Integer	是否中獎(1-未中獎，2-已中獎 ，2-和局)gameStatus 游戏状态， 按这四个数字枚举 -是否中獎(1-全部，2-未中獎，3-已中獎 ，4-和局
//		if val := dataOne.Get("isWin"); val.Exists() {
//			winInt = int(val.Int())
//			record.IsWin = int64(winInt)
//			// continue
//			if winInt == 0 {
//				record.GameStatus = 2
//			} else if winInt == 1 {
//				record.GameStatus = 3
//			} else if winInt == 2 {
//				record.GameStatus = 4
//			}
//		}
//
//		if val := dataOne.Get("winMoney"); val.Exists() {
//			winMoney64 = val.Float()
//			record.SettledAmountUSD = winMoney64
//			record.SettledAmount = winMoney64 * exchange64
//			// continue
//		}
//
//		// serialNumber（必须字段）
//		if val := dataOne.Get("serialNumber"); val.Exists() {
//			if val := dataOne.Get("serialNumber").String(); val != "" {
//				record.GameRoundID = val
//			} else {
//				zlogger.Warnf("Missing serialNumber, skipping record")
//				continue
//			}
//		}
//
//		// issueNo	string	期號.  // 局号
//		if val := dataOne.Get("issueNo"); val.Exists() {
//			if val := dataOne.Get("issueNo").String(); val != "" {
//				record.RoundID = val
//			} else {
//				zlogger.Warnf("Missing issueNo, skipping record")
//				continue
//			}
//		}
//
//		// gameType // 游戏代码
//		if val := dataOne.Get("gamePlayCode"); val.Exists() {
//			if val := dataOne.Get("gamePlayCode").String(); val != "" {
//				record.GameCode = val
//			} else {
//				zlogger.Warnf("Missing gamePlayCode for GameCode: %s", record.GameCode)
//				continue
//			}
//		}
//		record.IntegratorCode = constant.GameIntegratorAceLotteryCode
//		// record.Currency = constant.CurrencyCodeUSD
//		record.UserPlatform = constant.DeviceType
//		// 游戏类型
//		if val := dataOne.Get("gameCode"); val.Exists() {
//			if val := dataOne.Get("gameCode").String(); val != "" {
//				record.GameType = val
//				// "venue_code": "JDB", "gameCode": "1FNVN"
//				record.VenueCode = val
//			} else {
//				zlogger.Warnf("Missing gameCode for VenueCode: %s", record.GameCode)
//				continue
//			}
//		}
//
//		// 游戏名称
//		if val := dataOne.Get("gamePlayName"); val.Exists() {
//			if val := dataOne.Get("gamePlayName").String(); val != "" {
//				record.GameName = val
//			} else {
//				zlogger.Warnf("Missing gamePlayName for GameName: %s", record.GameCode)
//				continue
//			}
//		}
//
//		// 游戏扩展名称 "numsName": "北部越南彩"
//		if val := dataOne.Get("numsName"); val.Exists() {
//			if val := dataOne.Get("numsName").String(); val != "" {
//				record.GameExtName = val
//			} else {
//				zlogger.Warnf("Missing numsName for GameExtName: %s", record.GameCode)
//				continue
//			}
//		}
//
//		// amount
//		if val := dataOne.Get("betMoneyTotal"); val.Exists() {
//			if val := dataOne.Get("betMoneyTotal").Float(); val != 0 {
//				record.BetAmountUSD = val
//				record.BetAmount = val * exchange64
//			} else {
//				zlogger.Warnf("missing betMoneyTotal for BetAmountUSD: %s", record.BetAmount)
//				continue
//			}
//		}
//
//		// amount
//		if val := dataOne.Get("curOdd"); val.Exists() {
//			if val := dataOne.Get("curOdd").Float(); val != 0 {
//				record.CurOdd = val
//			} else {
//				zlogger.Warnf("curOdd is zero or missing for curOdd: %s", record.CurOdd)
//				continue
//			}
//		}
//
//		// amount NetAmount         float64 `json:"net_amount"`           // 净输赢金额
//		if val := dataOne.Get("winOrLossAmount"); val.Exists() {
//			if val := dataOne.Get("winOrLossAmount").Float(); val != 0 {
//				// SettledAmountUSD  float64 `json:"settled_amount_usd"`   // 美元结算金额
//				// if winInt == 1 {
//				// 	// record.SettledAmount = float64(valInt)
//				// 	zlogger.Debugf("isWin == 1 for serialNumber: %s", record.SettledAmount)
//				// 	record.SettledAmountUSD = winMoney64
//				// }
//				record.NetAmountUSD = val
//				record.NetAmount = val * exchange64
//			} else {
//				zlogger.Warnf("missing winOrLossAmount for NetAmountUSD: %s", record.BetAmount)
//				continue
//			}
//		}
//
//		if val := dataOne.Get("betTime"); val.Exists() {
//			if valInt := val.Int(); valInt != 0 {
//				record.BetAt = &valInt
//				record.CreatedAt = &valInt
//				// timeFormat := utils.GetESTimeFormat(time.Now().Format("2006.01.02 15:04:05"))
//				timestamp := time.UnixMilli(valInt) // 将毫秒时间戳转为 time.Time
//				record.CreatedAtISO = utils.GetESTimeFormat(timestamp.Format("2006.01.02 15:04:05"))
//
//			} else {
//				zlogger.Warnf("missing betTime for BetAt: %s", record.BetAmount)
//				continue
//			}
//		}
//
//		if val := dataOne.Get("settleTime"); val.Exists() {
//			if valInt := val.Int(); valInt != 0 {
//				record.SettledAt = &valInt
//			} else {
//				zlogger.Warnf("missing settleTime for SettledAt: %s", record.BetAmount)
//				continue
//			}
//		}
//
//		if val := dataOne.Get("updateTime"); val.Exists() {
//			if valInt := val.Int(); valInt != 0 {
//				record.UpdatedAt = &valInt
//				timestamp := time.UnixMilli(valInt) // 将毫秒时间戳转为 time.Time
//				record.UpdatedAtISO = utils.GetESTimeFormat(timestamp.Format("2006.01.02 15:04:05"))
//			} else {
//				zlogger.Warnf("missing  updateTime for UpdatedAtISO: %s", record.BetAmount)
//				continue
//			}
//		}
//
//		// 訂單狀態:(1未結算、2結算中、3已結算 4 用戶撤單 5系統撤單
//		if val := dataOne.Get("state"); val.Exists() {
//			valInt := val.Int()
//			if valInt == 1 {
//				record.Status = constant.WagerStatusBetStr
//				zlogger.Debugf("state == 1 for Status: %s", record.SettledAmount)
//			} else if valInt == 2 {
//				record.Status = constant.WagerStatusBetStr
//				zlogger.Debugf("state == 2 for Status: %s", record.SettledAmount)
//			} else if valInt == 3 {
//				record.Status = constant.WagerStatusSettledStr
//				zlogger.Debugf("state == 3 for Status: %s", record.SettledAmount)
//			} else if valInt == 4 {
//				record.Status = constant.WagerStatusCanceledStr
//				zlogger.Debugf("Unexpected state4 value: %d for Status: %s", val, record.SettledAmount)
//			} else if valInt == 5 {
//				record.Status = constant.WagerStatusROLLBACKStr
//				zlogger.Debugf("state == 5 for Status: %s", record.SettledAmount)
//			} else {
//				zlogger.Debugf("Unexpected state value: %d for Status: %s", val, record.SettledAmount)
//				continue
//			}
//
//		}
//
//		// numsName	string	投注說明 lotteryNum
//		if val := dataOne.Get("numsName"); val.Exists() {
//			if valStr := val.String(); valStr != "" {
//				record.NumsName = valStr
//			} else {
//				zlogger.Warnf("missing numsName for numsName: %s", record.BetAmount)
//				continue
//			}
//		}
//		// nums	string	投注號碼
//		if val := dataOne.Get("nums"); val.Exists() {
//			if valStr := val.String(); valStr != "" {
//				record.Nums = valStr
//			} else {
//				zlogger.Warnf("missing nums for Nums: %s", record.BetAmount)
//				continue
//			}
//		}
//
//		// nums	string	投注號碼
//		if val := dataOne.Get("lotteryNum"); val.Exists() {
//			if valStr := val.String(); valStr != "" {
//				record.LotteryNum = valStr
//			} else {
//				zlogger.Warnf("missing lotteryNum for lotteryNum: %s", record.BetAmount)
//				continue
//			}
//		}
//
//		if winInt == 0 {
//			zlogger.Debugf("isWin == 0 for serialNumber: %s", record.SettledAmount)
//		} else if winInt == 1 {
//			// record.SettledAmount = float64(valInt)
//			record.Status = constant.WagerStatusSettledStr
//			zlogger.Debugf("isWin == 1 for serialNumber: %s", record.SettledAmount)
//		} else if winInt == 2 {
//			record.Status = constant.WagerStatusSettledStr
//			zlogger.Debugf("isWin == 2 for serialNumber: %s", record.SettledAmount)
//		}
//
//		records = append(records, record)
//	}
//
//	return records
//}
