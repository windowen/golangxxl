package errs

var (
	ErrArgs           = NewCodeError(ArgsError, "ArgsError")
	ErrNoPermission   = NewCodeError(NoPermissionError, "NoPermissionError")
	ErrDatabase       = NewCodeError(DatabaseError, "DatabaseError")
	ErrInternalServer = NewCodeError(ServerInternalError, "ServerInternalError")
	ErrNetwork        = NewCodeError(NetworkError, "NetworkError")

	ErrData   = NewCodeError(DataError, "DataError")
	ErrUser   = NewCodeError(UserError, "UserError")
	ErrToken  = NewCodeError(TokenError, "TokenError")
	ErrUpload = NewCodeError(UploadError, "UploadError")

	ErrDiamondNotEnough = NewCodeError(DiamondNotEnough, "diamond not enough")
)
