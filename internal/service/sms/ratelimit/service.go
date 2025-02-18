package ratelimit

import (
	"context"
	"fmt"
	"webook/internal/service/sms"
	"webook/pkg/ratelimit"
)

var errLimited = fmt.Errorf("触发限流")

type RateLimitSMSService struct {
	svc     sms.Service
	limiter ratelimit.Limiter
}

func NewRateLimitSMSService(svc sms.Service, l ratelimit.Limiter) sms.Service {
	return &RateLimitSMSService{
		svc:     svc,
		limiter: l,
	}
}

func (s *RateLimitSMSService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	//TODO implement me
	limited, err := s.limiter.Limit(ctx, "sms:tecent")
	if err != nil {
		return fmt.Errorf("短信服务判断是否限流出现问题，%w", err)
	}
	if limited {
		return errLimited
	}
	err = s.svc.Send(ctx, tplId, args, numbers...)
	return err
}
