package repository

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
	"webook/internal/domain"
	"webook/internal/repository/cache"
	"webook/internal/repository/dao"
	daomocks "webook/internal/repository/dao/mocks"
	cachemocks "webook/internal/repository/user/mocks"
)

func TestCacheUserRepository_FindById(t *testing.T) {
	now := time.Now()
	// 去掉毫秒以外的部分
	now = time.UnixMilli(now.UnixMilli())
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache)
		// 输入
		id int64
		// 输出
		wantUser domain.User
		wantErr  error
	}{
		{
			name: "缓存未命中，DB查找成功",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UserCache) {
				c := cachemocks.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), int64(15)).Return(domain.User{}, cache.ErrKeyNotExist)
				c.EXPECT().Set(gomock.Any(), gomock.Any()).Return(nil)

				d := daomocks.NewMockUserDao(ctrl)
				d.EXPECT().FindById(gomock.Any(), int64(15)).
					Return(dao.User{
						Id: 15,
						Email: sql.NullString{
							String: "123@qq.com",
							Valid:  true,
						},
						Password: "123456",
						Phone: sql.NullString{
							String: "13412341234",
							Valid:  true,
						},
						Ctime: now.UnixMilli(),
						Utime: now.UnixMilli(),
					}, nil)
				return d, c
			},
			id: 15,
			wantUser: domain.User{
				Id:       15,
				Email:    "123@qq.com",
				Password: "123456",
				Phone:    "13412341234",
				Ctime:    now,
			},
			wantErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			r := NewUserRepository(tc.mock(ctrl))
			u, err := r.FindById(context.Background(), tc.id)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, u)
			time.Sleep(1 * time.Second)
		})
	}

}
