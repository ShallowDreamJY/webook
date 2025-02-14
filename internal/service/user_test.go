package service

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"
	"webook/internal/domain"
	"webook/internal/repository"
	repomocks "webook/internal/repository/mocks"
)

func Test_userService_Login(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) repository.UserRepository

		// 输入
		// ctx context.Context
		email    string
		password string

		// 输出
		wantUser domain.User
		wantErr  error
	}{
		{
			name: "登陆成功",

			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{
						Email:    "123@qq.com",
						Password: "$2a$10$Jyh.tUWc8nio4dJY1cnbcu40HD2RblTvsCA3vz/xl9CdGTOA3RR6.",
						Phone:    "12345678999",
						Ctime:    now,
					}, nil)
				return repo
			},
			email:    "123@qq.com",
			password: "123456",
			wantUser: domain.User{
				Email:    "123@qq.com",
				Password: "$2a$10$Jyh.tUWc8nio4dJY1cnbcu40HD2RblTvsCA3vz/xl9CdGTOA3RR6.",
				Phone:    "12345678999",
				Ctime:    now,
			},
			wantErr: nil,
		},
		{
			name: "用户不存在",

			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{}, repository.ErrUserNotFound)
				return repo
			},
			email:    "123@qq.com",
			password: "123456",
			wantUser: domain.User{},
			wantErr:  ErrInvalidUserOrPassword,
		},
		{
			name: "DB错误",

			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{}, errors.New("db error"))
				return repo
			},
			email:    "123@qq.com",
			password: "123456",
			wantUser: domain.User{},
			wantErr:  errors.New("db error"),
		},
		{
			name: "密码错误",

			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").
					Return(domain.User{
						Email:    "123@qq.com",
						Password: "$2a$10$Jyh.tUWc8nio4dJY1cnbcu40HD2RblTvsCA3vz/xl9CdGTOA3RR6.",
						Phone:    "12345678999",
						Ctime:    now,
					}, nil)
				return repo
			},
			email:    "123@qq.com",
			password: "123456wrong",
			wantUser: domain.User{},
			wantErr:  ErrInvalidUserOrPassword,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svc := NewUserService(tc.mock(ctrl))
			u, err := svc.Login(context.Background(), tc.email, tc.password)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, u)
		})
	}
}

func TestEncrypted(t *testing.T) {
	res, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err == nil {
		t.Log(string(res))
	}
}
