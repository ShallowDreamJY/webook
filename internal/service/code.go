package service

import (
	"context"
	"fmt"
	"math/rand"
	"webook/internal/repository"
	"webook/internal/service/sms"
)

const codeTplId = "1877556"

var (
	ErrCodeVerifyTooManyTimes = repository.ErrVerifyTooManyTimes
	ErrCodeSendTooMany        = repository.ErrCodeSendTooMany
)

type CodeService interface {
	Send(ctx context.Context, biz string, phone string) error
	Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error)
}

type codeService struct {
	repo   repository.CodeRepository
	smsSvc sms.Service
}

func NewCodeService(repo repository.CodeRepository, smsSvc sms.Service) CodeService {
	return &codeService{
		repo:   repo,
		smsSvc: smsSvc,
	}
}

func (svc *codeService) Send(ctx context.Context, biz string, phone string) error {
	code := svc.generateCode()
	err := svc.repo.Store(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	err = svc.smsSvc.Send(ctx, codeTplId, []string{code}, phone)
	return err
}

func (svc *codeService) Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error) {
	return svc.repo.Verify(ctx, biz, phone, inputCode)
}

func (svc *codeService) generateCode() string {
	// 生成 num 在 （0，999999）之间的六位数
	num := rand.Intn(1000000)
	// 不够6位数加上前导0
	return fmt.Sprintf("%06d", num)
}
