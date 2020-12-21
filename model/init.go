package model

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"tweb/global"
	"tweb/model/sys"
	"tweb/utils"
)

const (
	MYSQL = "mysql" //数据库类型

	//初始超管账号信息
	ADMIN_ACCOUNT  = "admin"
	ADMIN_PASSWORD = "123456"
	ADMIN_REALNAME = "超管"
)

func Init(dbType, addr string, debug bool) error {
	//todo 改成配置
	//dsn := "root:123456@tcp(127.0.0.1:3306)/tweb?charset=utf8&parseTime=True&loc=Local"
	//dsn := "root:123456@tcp(mysql:3306)/tweb?charset=utf8&parseTime=True&loc=Local"
	var err error
	if dbType == MYSQL {
		global.TDB, err = gorm.Open(mysql.Open(addr), &gorm.Config{})
	}
	if err != nil {
		return err
	}
	if debug {
		global.TDB = global.TDB.Debug()
	}
	if err = migrateModels(global.TDB); err != nil {
		return err
	}
	return nil
}

func InitAdmin() error {
	salted, salt := utils.GenSaltedPasswd(ADMIN_PASSWORD)
	sysUser := sys.SysUser{
		Account:  ADMIN_ACCOUNT,
		RealName: ADMIN_REALNAME,
		Password: salted,
		Salt:     salt,
	}
	err := sysUser.New()
	if err == global.ErrAccountCollision {
		return nil
	}
	if err != nil {
		return err
	}
	return nil
}
