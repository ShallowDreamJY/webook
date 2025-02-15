package dao

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestGORMUserDao_Insert(t *testing.T) {
	testCacses := []struct {
		name string

		mock func(t *testing.T) *sql.DB

		ctx  context.Context
		user User

		wantErr error
	}{
		{
			name: "插入成功",
			mock: func(t *testing.T) *sql.DB {
				mockDB, mock, err := sqlmock.New()
				res := sqlmock.NewResult(3, 1)
				mock.ExpectExec("INSERT INTO `users` .*").
					WillReturnResult(res)
				require.NoError(t, err)
				return mockDB
			},

			ctx: context.Background(),
			user: User{
				Email: sql.NullString{
					String: "test@test.com",
					Valid:  true,
				},
			},
			wantErr: nil,
		},
	}

	for _, tc := range testCacses {
		t.Run(tc.name, func(t *testing.T) {

			db, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      tc.mock(t),
				SkipInitializeWithVersion: true,
			}), &gorm.Config{
				DisableAutomaticPing:   true,
				SkipDefaultTransaction: true,
			})
			d := NewUserDao(db)
			err = d.Insert(tc.ctx, tc.user)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
