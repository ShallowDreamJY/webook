package ioc

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"webook/internal/config"
	"webook/internal/repository/dao"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:SFuuC4vUQ3qzaD0@tcp(120.55.48.65:13306)/webook"))
	if err != nil {
		// 只会在初始化过程中 panic
		// panic相当于整个goroutine结束
		// 初始化出错，应用不要启动
		panic(err)
	}
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}

func InitRedis() redis.Cmdable {
	redisClient := redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr,
	})
	return redisClient
}
