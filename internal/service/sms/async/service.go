package async

import (
	"context"
	"time"
	"webook/internal/domain"
	"webook/internal/repository"
	"webook/pkg/logger"

	"webook/internal/service/sms"
)

type AsyncSMSService struct {
	svc  sms.Service
	repo repository.AsyncSmsRepository
	l    logger.LoggerV1
}

func NewAsyncSMSService(svc sms.Service,
	repo repository.AsyncSmsRepository,
	l logger.LoggerV1) *AsyncSMSService {
	res := &AsyncSMSService{
		svc:  svc,
		repo: repo,
		l:    l,
	}
	go func() {
		res.StartAsyncCycle()
	}()
	return res
}

func (s *AsyncSMSService) StartAsyncCycle() {
	for {
		s.AsyncSend()
	}
}

func (s *AsyncSMSService) AsyncSend() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	as, err := s.repo.PreemptWaitSMS(ctx)
	cancel()
	switch err {
	// 查询到了待发送的sms，发送
	case nil:
		ctx, cancel = context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		err = s.svc.Send(ctx, as.TplId, as.Args, as.Numbers...)
		res := err == nil
		// 发送失败
		if err != nil {
			s.l.Error("异步发送信息失败",
				logger.Bool("res", res),
				logger.Int64("id", as.Id))
			return
		}
		// 记录发送结果
		err = s.repo.ReportScheduleResult(ctx, as.Id, res)
		// 记录失败
		if err != nil {
			s.l.Error("执行异步发送短信成功，标记数据库失败",
				logger.Error(err),
				logger.Bool("res", res),
				logger.Int64("id", as.Id))
		}
	case repository.ErrWaitingSMSNotFound:
		time.Sleep(time.Second)
	default:
		s.l.Error("抢占异步发送短信任务失败",
			logger.Error(err))
		time.Sleep(time.Second)
	}
}

func (s *AsyncSMSService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	if s.NeedAsync() {
		err := s.repo.Add(ctx, domain.AsyncSms{
			TplId:    tplId,
			Args:     args,
			Numbers:  numbers,
			RetryMax: 3,
		})
		return err
	}
	return s.svc.Send(ctx, tplId, args, numbers...)
}

func (s *AsyncSMSService) NeedAsync() bool {
	// TODO:判断是否服务崩溃的具体实现
	return true
}
