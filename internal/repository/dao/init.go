package dao

import "gorm.io/gorm"

func InitTables(db *gorm.DB) error {
	err := db.AutoMigrate(&User{}, &AsyncSms{}, &Article{})
	if err != nil {
		panic(err)
	}
	return err
}
