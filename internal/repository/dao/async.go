package dao

import (
	"context"
	"github.com/ecodeclub/ekit/sqlx"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

var ErrWaitingSMSNotFound = gorm.ErrRecordNotFound

type AsyncSmsDao interface {
	Insert(ctx context.Context, sms AsyncSms) error
	GetWaitingSMS(ctx context.Context) (AsyncSms, error)
	MarkSuccess(ctx context.Context, id int64) error
	MarkFailed(ctx context.Context, id int64) error
}

type AsyncSms struct {
	Id       int64
	Config   sqlx.JsonColumn[SmsConfig]
	RetryCnt int
	RetryMax int
	Status   int
	Ctime    int64
	Utime    int64
}

type SmsConfig struct {
	TplId   string
	Args    []string
	Numbers []string
}

const (
	asyncStatusWaiting = iota
	asyncStatusFailed
	asyncStatusSuccess
)

type GORMAsyncSmsDAO struct {
	db *gorm.DB
}

func (g *GORMAsyncSmsDAO) MarkSuccess(ctx context.Context, id int64) error {
	now := time.Now().UnixMilli()
	return g.db.WithContext(ctx).Model(&AsyncSms{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status": asyncStatusSuccess,
			"utime":  now,
		}).Error
}

func (g *GORMAsyncSmsDAO) MarkFailed(ctx context.Context, id int64) error {
	now := time.Now().UnixMilli()
	return g.db.WithContext(ctx).Model(&AsyncSms{}).
		Where("id = ? and `retry_cnt` >= `retry_max`", id).
		Updates(map[string]any{
			"status": asyncStatusFailed,
			"utime":  now,
		}).Error
}

func (g *GORMAsyncSmsDAO) GetWaitingSMS(ctx context.Context) (AsyncSms, error) {
	var s AsyncSms
	err := g.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now().UnixMilli()
		endTime := now - time.Second.Milliseconds()
		err := tx.Clauses(clause.Locking{Strength: "UPDATAE"}).
			Where("utime < ? and status = ?", endTime, asyncStatusWaiting).First(&s).Error
		if err != nil {
			return err
		}
		err = tx.Model(&AsyncSms{}).
			Where("id = ?", s.Id).
			Updates(map[string]any{
				"retry_cnt": gorm.Expr("retry_cnt+1"),
				"utime":     now,
			}).Error
		return err
	})
	return s, err
}

func NewGORMAsyncSmsDAO(db *gorm.DB) *GORMAsyncSmsDAO {
	return &GORMAsyncSmsDAO{
		db: db,
	}
}

func (g *GORMAsyncSmsDAO) Insert(ctx context.Context, sms AsyncSms) error {
	return g.db.Create(&sms).Error
}
