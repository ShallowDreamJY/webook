package tencent

import (
	"context"
	"fmt"
	"github.com/ecodeclub/ekit"
	"github.com/ecodeclub/ekit/slice"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20190711"
)

type Service interface {
	Send(ctx context.Context, biz string, args []string, numbers ...string) error
}

type ServiceImpl struct {
	appId     *string
	signature *string
	client    *sms.Client
}

func NewServiceImpl(client *sms.Client, appId string, signaure string) Service {
	return &ServiceImpl{
		appId:     ekit.ToPtr[string](appId),
		signature: ekit.ToPtr[string](signaure),
		client:    client,
	}
}

func (s *ServiceImpl) Send(ctx context.Context, biz string, args []string, numbers ...string) error {
	req := sms.NewSendSmsRequest()
	req.SmsSdkAppid = s.appId
	req.Sign = s.signature
	req.TemplateID = ekit.ToPtr[string](biz)
	req.PhoneNumberSet = s.toStringPtrSlice(numbers)
	req.TemplateParamSet = s.toStringPtrSlice(args)
	resp, err := s.client.SendSms(req)
	if err != nil {
		return err
	}
	for _, status := range resp.Response.SendStatusSet {
		if status.Code == nil || *(status.Code) != "Ok" {
			return fmt.Errorf("发送短信失败 %s， %s ", *status.Code, *status.Message)
		}
	}
	return nil
}

func (s *ServiceImpl) toStringPtrSlice(src []string) []*string {
	return slice.Map[string, *string](src, func(idx int, src string) *string {
		return &src
	})
}
