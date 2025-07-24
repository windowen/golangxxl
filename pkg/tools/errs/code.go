package errs

// 通用错误码.
const (
	NoError       = 0     // 无错误
	DatabaseError = 90002 // redis/mysql等db错误
	NetworkError  = 90004 // 网络错误
	DataError     = 90007 // 数据错误

	// 通用错误码.
	ServerInternalError = 500  // 服务器内部错误
	ArgsError           = 1001 // 输入参数错误
	NoPermissionError   = 1002 // 权限不足

	// 账号错误码.
	UserError = 1101 // UserID不存在 或未注册

	UploadError = 1210 // 文件上传失败

	// token错误码.
	TokenError = 1201

	// DiamondNotEnough 钻石消费
	DiamondNotEnough = 2201 // 钻石不足
)
