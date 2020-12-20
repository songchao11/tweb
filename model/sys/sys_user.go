package sys

import (
	"gorm.io/gorm"
	"tweb/global"
)

type SysUser struct {
	gorm.Model

	Account  string //用户账号
	RealName string //用户真实姓名
	Password string //密码
	Salt     string //盐
	Phone    string //手机号
}

func (sysUser SysUser) New() error {
	var err error
	tx := global.TDB.Begin()
	defer func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	exist, err := sysUser.checkExistByAccount(tx, sysUser.Account)
	if err != nil {
		return err
	}
	if exist {
		return global.ErrAccountCollision
	}

	if err = tx.Model(sysUser).Create(&sysUser).Error; err != nil {
		return err
	}

	return nil
}

func (sysUser SysUser) checkExistByAccount(db *gorm.DB, account string) (bool, error) {
	var cnt int64
	if err := db.Model(sysUser).Where("account = ?", account).Count(&cnt).Error; err != nil {
		return false, err
	}
	return cnt > 0, nil
}

func (sysUser SysUser) GetSysUserByAccount(account string) (*SysUser, error) {
	if err := global.TDB.Model(sysUser).Where("account = ?", account).Take(&sysUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, global.ErrSysUserNotExist
		}
		return nil, err
	}
	return &sysUser, nil
}
