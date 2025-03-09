package repository

import (
	"context"
	"github.com/ecodeclub/ekit/sqlx"
	"webook/internal/domain"
	"webook/internal/repository/dao"
)

var ErrWaitingSMSNotFound = dao.ErrWaitingSMSNotFound

type AsyncSmsRepository interface {
	Add(ctx context.Context, sms domain.AsyncSms) error
	PreemptWaitSMS(ctx context.Context) (domain.AsyncSms, error)
	ReportScheduleResult(ctx context.Context, id int64, sucess bool) error
}

type asyncSmsRepository struct {
	dao dao.AsyncSmsDao
}

func NewAsyncSmsRepository(dao dao.AsyncSmsDao) AsyncSmsRepository {
	return &asyncSmsRepository{dao: dao}
}

func (a *asyncSmsRepository) Add(ctx context.Context, sms domain.AsyncSms) error {
	//将sms发送消息记录到db中，等待异步发送
	return a.dao.Insert(ctx, dao.AsyncSms{
		Id: sms.Id,
		Config: sqlx.JsonColumn[dao.SmsConfig]{
			Val: dao.SmsConfig{
				TplId:   sms.TplId,
				Args:    sms.Args,
				Numbers: sms.Numbers,
			},
			Valid: true,
		},
		RetryMax: sms.RetryMax,
	})

}

func (a *asyncSmsRepository) PreemptWaitSMS(ctx context.Context) (domain.AsyncSms, error) {
	//找到一个可以发送的sms信息
	as, err := a.dao.GetWaitingSMS(ctx)
	if err != nil {
		return domain.AsyncSms{}, err
	}
	return domain.AsyncSms{
		Id:       as.Id,
		TplId:    as.Config.Val.TplId,
		Args:     as.Config.Val.Args,
		Numbers:  as.Config.Val.Numbers,
		RetryMax: as.RetryMax,
	}, err
}

func (a *asyncSmsRepository) ReportScheduleResult(ctx context.Context, id int64, sucess bool) error {
	// 记录异步发送的结果
	if sucess {
		return a.dao.MarkSuccess(ctx, id)
	}
	return a.dao.MarkFailed(ctx, id)
	return nil
}
