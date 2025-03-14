package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type ArticleDao interface {
	Insert(ctx context.Context, art Article) (int64, error)
}

// 制作库
type Article struct {
	Id      int64  `gorm:"primaryKey;autoIncrement"`
	Title   string `gorm:"type=varchar(1024)"`
	Content string `gorm:"type=BLOB"`
	//AuthorId int64  `gorm:"index=aid_ctime"`
	//Ctime int64 `gorm:"index=aid_ctime"`
	AuthorId int64 `gorm:"index"`
	Ctime    int64
	Utime    int64
}

type GORMArticleeDao struct {
	db *gorm.DB
}

func (dao *GORMArticleeDao) Insert(ctx context.Context, art Article) (int64, error) {
	now := time.Now().UnixMilli()
	art.Ctime = now
	art.Utime = now
	err := dao.db.WithContext(ctx).Create(&art).Error
	return art.Id, err
}

func NewGORMArticleDao(db *gorm.DB) ArticleDao {
	return &GORMArticleeDao{
		db: db,
	}
}
