package srverr

import "errors"

var (
	// -------------- common ---------------
	ErrVerify          = errors.New("无效的验证码")
	ErrSystemException = errors.New("系统繁忙, 请重试")
	ErrInvalidPEM      = errors.New("invalid pem")
	ErrInvalidToken    = errors.New("无效的token,认证失败")
	ErrTokenNotFound   = errors.New("没有发现有效的token")
	ErrForbidden       = errors.New("没有访问权限")
	ErrDownload        = errors.New("下载失败")
	ErrUpload          = errors.New("文件上传失败")
	ErrFileMax         = errors.New("文件太大")
	ErrInvalidParams   = errors.New("请求参数无效或缺失")
	ErrLoginFailed    = errors.New("用户名或密码错误")
	ErrUpdateFailed   = errors.New("更新失败")
	ErrDeleteFailed   = errors.New("删除失败")

	// -------------- admin ---------------
	ErrInvalidCreator = errors.New("invalid admin, do you really an admin ？")
	ErrOwnPower   = errors.New("您没有权限执行此操作")
	ErrAdminFailed    = errors.New("login failed, admin not existed or password is false")
	ErrInvalidPower   = errors.New("invalid power")
	// -------------- user ----------------
	ErrRegisterChecker = errors.New("审核人不存在或没有权限")
	ErrRegisterExisted = errors.New("您注册的用户记录已存在")
	// -------------- class ---------------
	ErrInvalidClass = errors.New("系统没有该班级的记录")
)
