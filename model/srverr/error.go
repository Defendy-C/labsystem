package srverr

import "errors"

var (
	// -------------- common ---------------
	ErrVerify          = errors.New("无效的验证码")
	ErrSystemException = errors.New("系统繁忙, 请重试")
	ErrInvalidPEM      = errors.New("invalid pem")
	ErrInvalidToken    = errors.New("invalid token")
	ErrForbidden       = errors.New("认证失败")
	ErrDownload        = errors.New("下载失败")
	ErrInvalidParams   = errors.New("请求参数无效或缺失")
	ErrLoginFailed    = errors.New("用户名或密码错误")

	// -------------- admin ---------------
	ErrInvalidCreator = errors.New("invalid admin, do you really an admin ？")
	ErrOwnPower   = errors.New("you haven't this power, please link to administrator")
	ErrAdminFailed    = errors.New("login failed, admin not existed or password is false")
	ErrInvalidPower   = errors.New("invalid power")
	// -------------- user ----------------
	ErrRegisterChecker = errors.New("not found checker")
	// -------------- class ---------------
	ErrInvalidClass = errors.New("valid class info")
)
