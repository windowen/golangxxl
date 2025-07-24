package elastic

// GameRecordES 游戏记录ES存储结构体
type GameRecordES struct {
	UserID            int64   `json:"user_id"`              // 用户ID
	UserName          string  `json:"user_name"`            // 用户名
	UserPlatform      string  `json:"user_platform"`        // 用户平台
	BetTxnID          string  `json:"bet_txn_id"`           // 投注交易ID
	BetAmount         float64 `json:"bet_amount"`           // 投注金额
	SettledTxnID      string  `json:"settled_txn_id"`       // 结算交易ID
	SettledAmount     float64 `json:"settled_amount"`       // 结算金额
	RefundTxnID       string  `json:"refund_txn_id"`        // 退款交易ID
	RefundAmount      float64 `json:"refund_amount"`        // 退款金额
	NetAmount         float64 `json:"net_amount"`           // 净输赢金额
	ValidBetAmount    float64 `json:"valid_bet_amount"`     // 有效投注金额
	Currency          string  `json:"currency"`             // 货币类型
	ExchangeRate      float64 `json:"exchange_rate"`        // 汇率
	BetAmountUSD      float64 `json:"bet_amount_usd"`       // 美元投注金额
	SettledAmountUSD  float64 `json:"settled_amount_usd"`   // 美元结算金额
	RefundAmountUSD   float64 `json:"refund_amount_usd"`    // 美元退款金额
	ValidBetAmountUSD float64 `json:"valid_bet_amount_usd"` // 美元有效投注金额
	NetAmountUSD      float64 `json:"net_amount_usd"`       // 美元净输赢金额
	GameCode          string  `json:"game_code"`            // 游戏代码
	GameName          string  `json:"game_name"`            // 游戏名称
	GameType          string  `json:"game_type"`            // 游戏类型
	GameExtName       string  `json:"game_ext_name"`        // 游戏扩展名称
	GameRoundID       string  `json:"game_round_id"`        // 游戏局号
	IntegratorCode    string  `json:"integrator_code"`      // 接入商代码
	VenueCode         string  `json:"venue_code"`           // 场馆代码
	RoundID           string  `json:"round_id"`             // 局号
	Nums              string  `json:"nums"`                 // nums	string	投注號碼
	NumsName          string  `json:"numsName"`             // numsName	string	投注說明 lotteryNum "numsName": "北部越南彩",
	LotteryNum        string  `json:"lotteryNum"`           // lotteryNum	string	開獎號碼
	State             int64   `json:"state"`                // 訂單狀態:(1未結算、2結算中、3已結算 4 用戶撤單 5系統撤單
	CurOdd            float64 `json:"curOdd"`               // curOdd	BigDecimal	目前使用的賠率  "curOdd": "1.97",
	IsWin             int64   `json:"isWin"`                // isWin	Integer	是否中獎(0-未中獎，1-已中獎 ，2-和局)gameStatus
	GameStatus        int64   `json:"gameStatus"`           //  按这四个数字枚举 -是否中獎(1-全部，2-未中獎，3-已中獎 ，4-和局 gameStatus
	// SubRoundList      string  `json:"sub_round_list"`       // 回合详情  "isWin": 1,
	Status       string `json:"status"`         // 状态
	BetAt        *int64 `json:"bet_at"`         // 投注时间
	SettledAt    *int64 `json:"settled_at"`     // 结算时间
	RefundAt     *int64 `json:"refund_at"`      // 退款时间
	CancelAt     *int64 `json:"cancel_at"`      // 取消时间
	CreatedAt    *int64 `json:"created_at"`     // 创建时间
	UpdatedAt    *int64 `json:"updated_at"`     // 更新时间
	CreatedAtISO string `json:"created_at_iso"` // ISO格式创建时间
	UpdatedAtISO string `json:"updated_at_iso"` // ISO格式更新时间
}

type SubRound struct {
	PlatformTxId string  `json:"platformTxId"` // 平台交易ID
	RoundId      string  `json:"roundId"`      // 局号
	Turnover     float64 `json:"turnover"`     // 输赢
	BetAmount    float64 `json:"betAmount"`    // 投注金额
	WinAmount    float64 `json:"winAmount"`    // 赢金额
}
