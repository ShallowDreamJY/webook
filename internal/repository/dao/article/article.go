package article

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type ArticleDao interface {
	Insert(ctx context.Context, art Article) (int64, error)
	UpdateById(ctx context.Context, art Article) error
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

func (dao *GORMArticleeDao) UpdateById(ctx context.Context, art Article) error {
	now := time.Now().UnixMilli()
	art.Utime = now
	res := dao.db.WithContext(ctx).Model(&Article{}).
		Where("id = ? AND author_id = ?", art.Id, art.AuthorId).Updates(map[string]any{
		"title":   art.Title,
		"content": art.Content,
		"utime":   art.Utime,
	})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("更新失败，可能是创作者非法 id %d， author_id %d", art.Id, art.AuthorId)
	}
	return res.Error
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
