package global

import "fmt"

//error code
const (
	ErrCodeSuccess      = 10000
	ErrCodeParamInvalid = 40002
	ErrCodeInternal     = 40004
	ErrCodePriviledge   = 40006
	ErrCodeSessionGone  = 40010
)

var (
	ErrAccountCollision = fmt.Errorf("账号已被使用")
	ErrSysUserNotExist  = fmt.Errorf("用户不存在")
)
