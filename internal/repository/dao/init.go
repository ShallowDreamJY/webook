package dao

import (
	"gorm.io/gorm"
	"webook/internal/repository/dao/article"
)

func InitTables(db *gorm.DB) error {
	err := db.AutoMigrate(&User{}, &AsyncSms{}, &article.Article{})
	if err != nil {
		panic(err)
	}
	return err
}
