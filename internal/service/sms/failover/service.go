package failover

import (
	"context"
	"errors"
	"log"
	"sync/atomic"
	"webook/internal/service/sms"
)

type FailoverSMSService struct {
	svcs []sms.Service
	idx  uint64
}

func (f *FailoverSMSService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	for _, svc := range f.svcs {
		err := svc.Send(ctx, tplId, args, numbers...)
		// 发送成功
		if err == nil {
			return nil
		}
		log.Println(err)
	}
	return errors.New("所有短信服务商全部发送失败")
}

func NewFailoverSMSService(svcs []sms.Service) sms.Service {
	return &FailoverSMSService{
		svcs: svcs,
		idx:  0,
	}
}

//func (f *FailoverSMSService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
//	for _, svc := range f.svcs {
//		err := svc.Send(ctx, tplId, args, numbers...)
//		// 发送成功
//		if err == nil {
//			return nil
//		}
//		log.Println(err)
//	}
//	return errors.New("所有短信服务商全部发送失败")
//}

func (f *FailoverSMSService) SendV1(ctx context.Context, tplId string, args []string, numbers ...string) error {
	idx := atomic.AddUint64(&f.idx, 1)
	length := uint64(len(f.svcs))
	for i := idx; i < length+idx; i++ {
		err := f.svcs[i%length].Send(ctx, tplId, args, numbers...)
		// 发送成功
		if err == nil {
			return nil
		}
		log.Println(err)
	}
	return errors.New("所有短信服务商全部发送失败")
}
