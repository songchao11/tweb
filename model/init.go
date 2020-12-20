package model

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"tweb/global"
	"tweb/model/sys"
	"tweb/utils"
)

func Init() error {
	//todo 改成配置
	//dsn := "root:123456@tcp(127.0.0.1:3306)/tweb?charset=utf8&parseTime=True&loc=Local"
	dsn := "root:123456@tcp(mysql:3306)/tweb?charset=utf8&parseTime=True&loc=Local"
	var err error
	global.TDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	global.TDB = global.TDB.Debug()
	if err = migrateModels(global.TDB); err != nil {
		return err
	}
	return nil
}

func InitAdmin() error {
	salted, salt := utils.GenSaltedPasswd("123456")
	sysUser := sys.SysUser{
		Account:  "admin",
		RealName: "超管",
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
