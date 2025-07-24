package constant

// 用户组redis key
const (
	RegisterValidation    = "register_Validation_%v_%v"  // 注册验证key
	ModifyPaymentPassword = "modify_payment_password_%v" // 修改支付密码key
	ForgetPassword        = "forget_password_%v_%v"      // 忘记密码key
	ModifyLoginPassword   = "modify_login_password_%v"   // 修改登录密码key
)

const (
	MsgLoginValidationPrefix = "Msg_Login_Validation_Prefix_%v_%v"
)

const (
	UniversalTemplateID = "Universal"
)
