package article

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type ArticleDao interface {
	Insert(ctx context.Context, art Article) (int64, error)
	UpdateById(ctx context.Context, art Article) error
	Sync(ctx context.Context, art Article) (int64, error)
	Upsert(ctx context.Context, art PublishArticle) error
	SyncStatus(ctx context.Context, id int64, author int64, status uint8) error
	Transaction(ctx context.Context, bizFunc func(txDAO ArticleDao) error) error
}

// 制作库
type Article struct {
	Id      int64  `gorm:"primaryKey;autoIncrement"`
	Title   string `gorm:"type=varchar(1024)"`
	Content string `gorm:"type=BLOB"`
	//AuthorId int64  `gorm:"index=aid_ctime"`
	//Ctime int64 `gorm:"index=aid_ctime"`
	AuthorId int64 `gorm:"index"`
	// tatus domain.ArticleStatus `gorm:"type:uint8"`
	Status uint8
	Ctime  int64
	Utime  int64
}

type GORMArticleeDao struct {
	db *gorm.DB
}

func (dao *GORMArticleeDao) SyncStatus(ctx context.Context, id int64, author int64, status uint8) error {
	now := time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&Article{}).
			Where("id = ? AND author_id = ?", id, author).
			Updates(map[string]any{
				"status": status,
				"utime":  now,
			})
		if res.Error != nil {
			// 数据库有问题
			return res.Error
		}
		if res.RowsAffected == 0 {
			// 没查到数据，id错或者author对不上
			return fmt.Errorf("没查到数据 ，可能有人攻击或者误操作 id ： %d， author： %d", id, author)
		}
		return tx.Model(&Article{}).
			Where("id = ?", id).
			Updates(map[string]any{
				"status": status,
				"utime":  now,
			}).Error

	})
}

func (dao *GORMArticleeDao) Transaction(ctx context.Context, bizFunc func(txDAO ArticleDao) error) error {
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txDAO := NewGORMArticleDao(tx)
		return bizFunc(txDAO)
	})
}

func (dao *GORMArticleeDao) Upsert(ctx context.Context, art PublishArticle) error {
	now := time.Now().Unix()
	art.Ctime = now
	art.Utime = now
	err := dao.db.Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]any{
			"title":   art.Title,
			"status":  art.Status,
			"content": art.Content,
			"utime":   art.Utime,
		}),
	}).Create(&art).Error
	return err
}

func (dao *GORMArticleeDao) Sync(ctx context.Context, art Article) (int64, error) {
	var id = art.Id
	err := dao.db.Transaction(func(tx *gorm.DB) error {
		var err error
		txDAO := NewGORMArticleDao(tx)
		if id > 0 {
			err = txDAO.UpdateById(ctx, art)
		} else {
			id, err = txDAO.Insert(ctx, art)
		}
		if err != nil {
			return err
		}
		return txDAO.Upsert(ctx, PublishArticle{Article: art})
	})
	return id, err
}

func (dao *GORMArticleeDao) UpdateById(ctx context.Context, art Article) error {
	now := time.Now().UnixMilli()
	art.Utime = now
	res := dao.db.WithContext(ctx).Model(&Article{}).
		Where("id = ? AND author_id = ?", art.Id, art.AuthorId).Updates(map[string]any{
		"title":   art.Title,
		"content": art.Content,
		"status":  art.Status,
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
