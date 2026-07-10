package service

import (
	"gorm.io/gorm"

	"go-admin/app/wiki/models"
)

// Migrate 自动建表并幂等注入种子数据(等级体系 + 初始分类/条目)
func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&models.WikiCategory{},
		&models.WikiEntry{},
		&models.LevelSystem{},
		&models.Level{},
	); err != nil {
		return err
	}
	if err := seedLevels(db); err != nil {
		return err
	}
	return seedContent(db)
}
