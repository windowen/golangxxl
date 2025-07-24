package constant

import "fmt"

const (
	ConfigPath               = "/config/config.yaml"
	CtxApiToken              = "api-token"
	LocalHost                = "127.0.0.1"
	ServerAPICommonConfigKey = "ServerAPICommonConfig"
)

const (
	LangChinese = "zh" // 中文
	LangEnglish = "en" // 英语
	LangID      = "id" // 印尼语
	LangVi      = "vi" // 越南语
	LangPT      = "pt" // 巴西语
)

var (
	LangPack = []string{LangChinese, LangEnglish, LangID, LangVi, LangPT}
)

// 基础响应
const (
	SuccessResponseCode = 10000                      // 成功响应code
	ErrorResponseCode   = SuccessResponseCode + iota // 失败响应code
	LoginResponseCode                                // 需要登录时响应code
	KicOutResponseCode                               // 被踢出
)

const (
	OperationId          = "operationId"
	OpUserId             = "opUserId"
	Token                = "token"
	RpcCustomHeader      = "customHeader" // rpc中间件自定义ctx参数
	ConnId               = "connId"
	CountryCode          = "countryCode"
	OpUserPlatform       = "platform"
	Language             = "language"
	Authorization        = "Authorization"
	RefreshAuthorization = "Refresh-Authorization"
	OpTourist            = "isTourist"       // 是否游客
	AgoraSignature       = "Agora-Signature" // 声网签名
	Location             = "timeZone"        // 时区
	RefreshToken         = "refreshToken"    // 刷新令牌
	DeviceId             = "deviceId"        // 刷新令牌
)

const (
	RpcOperationID = OperationId
	RpcOpUserID    = OpUserId
	RpcOpUserType  = "opUserType"
)

const (
	NormalUser = 1
	AdminUser  = 2
)

const PasswordIteratorCount = 3

const (
	OtherLogins int32 = iota
	Phone
	Password
)

const (
	Zero          = 0
	No            = 0
	Yes           = 1
	StatusNormal  = 1 // 正常
	StatusDisable = 2 // 禁用
	StatusDel     = 3 // 删除
)

const (
	GameIntegratorHashKey        = "game_integrator_cache"
	GameIntegratorAceLotteryCode = "AceLottery"
	UrlAceGameRecord             = "/third/rest/third/u/openApi/v2/Query/GameRecord"

	AceStatusSuccess  = "200" // 标准错误码定义（HTTP语义兼容） 登录并进入游戏
	AceStatusSuccess0 = "0"   // 标准错误码定义（HTTP语义兼容） 登录并进入游戏
	DeviceType        = "WEB"
)

const (
	WagerStatusBetStr      = "BET"
	WagerStatusSettledStr  = "SETTLED"
	WagerStatusVoidStr     = "VOID"
	WagerStatusTipStr      = "TIP"
	WagerStatusROLLBACKStr = "ROLLBACK"
	WagerStatusCanceledStr = "CANCELED"
)

const (
	GameCountryCodeBR = "BR" // 巴西
	GameCountryCodeIN = "IN" // 印度
	GameCountryCodeUS = "US" // 美国
	GameCountryCodeID = "ID" // 印尼
	GameCountryCodeVN = "VN" // 越南
	GameCountryCodeTH = "TH" // 泰国
)

const (
	CurrencyCodeBRL  = "BRL"  // 巴西雷亚尔
	CurrencyCodeINR  = "INR"  // 印度卢比
	CurrencyCodeUSD  = "USD"  // 美元
	CurrencyCodeIDR  = "IDR"  // 印尼盾
	CurrencyCodeIDR2 = "IDR2" // 印尼盾 17少*1000的汇率
	CurrencyCodeVND  = "VND"  // 越南盾
	CurrencyCodeVND2 = "VND2" // 越南盾25.6的汇率换算
	CurrencyCodeTHB  = "THB"  // 泰铢
	CurrencyCodePTV  = "PTV"  // 越南盾 Original Vietnamese đồng PTV (1:1)	原版越南盾 (1:1)
)

var CountryToCurrency = map[string]string{
	GameCountryCodeBR: CurrencyCodeBRL,
	GameCountryCodeIN: CurrencyCodeINR,
	GameCountryCodeUS: CurrencyCodeUSD,
	GameCountryCodeID: CurrencyCodeIDR,
	GameCountryCodeVN: CurrencyCodeVND,
	GameCountryCodeTH: CurrencyCodeTHB,
}

func GetCurrencyByCountry(countryCode string) (string, error) {
	if currency, exists := CountryToCurrency[countryCode]; exists {
		return currency, nil
	}
	return "", fmt.Errorf("currency not found for country code: %s", countryCode)
}
