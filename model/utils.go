package model

import (
	"gorm.io/gorm"
	"tweb/model/sys"
)

func migrateModels(db *gorm.DB) error {
	return db.Migrator().CreateTable(&sys.SysUser{})
}
