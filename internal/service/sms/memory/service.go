package memory

import (
	"context"
	"fmt"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s Service) Send(ctx context.Context, biz string, args []string, numbers ...string) error {
	//TODO implement me
	fmt.Println(args)
	return nil
}
