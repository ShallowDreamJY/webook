package repository

import (
	"context"
	"database/sql"
	"github.com/gin-gonic/gin"
	"time"
	"webook/internal/domain"
	"webook/internal/repository/cache"
	"webook/internal/repository/dao"
)

var (
	ErrUserDuplicated = dao.ErrUserDuplicated
	ErrUserNotFound   = dao.ErrUserNotFound
)

type UserRepository interface {
	Create(ctx context.Context, u domain.User) error
	FindById(ctx context.Context, id int64) (domain.User, error)
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindByPhone(ctx *gin.Context, phone string) (domain.User, error)
}

type CacheUserRepository struct {
	dao   dao.UserDao
	cache cache.UserCache
}

func NewUserRepository(dao dao.UserDao, c cache.UserCache) UserRepository {
	return &CacheUserRepository{
		dao:   dao,
		cache: c,
	}
}

func (r *CacheUserRepository) Create(ctx context.Context, u domain.User) error {

	return r.dao.Insert(ctx, r.domainToEntity(u))
}

func (r *CacheUserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	u, err := r.cache.Get(ctx, id)
	// 缓存命中
	if err == nil {
		return u, err
	}
	// 缓存没命中，查数据库
	// 缓存加载报错，继续从数据库加载
	// 选择继续从数据库加载，做好兜底，以免Redis真的崩了，保护住数据库，例如限流
	// 选择不加载，则用户体验差一点
	ue, err := r.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}

	u = r.entityToDomain(ue)

	go func() {
		err = r.cache.Set(ctx, u)
		if err != nil {
			// 打日志，做监控
		}
	}()

	return u, err
}

func (r *CacheUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return r.entityToDomain(u), nil
}

func (r *CacheUserRepository) FindByPhone(ctx *gin.Context, phone string) (domain.User, error) {
	u, err := r.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return r.entityToDomain(u), nil
}

func (r *CacheUserRepository) domainToEntity(u domain.User) dao.User {
	return dao.User{
		Id: u.Id,
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		Password: u.Password,
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
		Ctime: u.Ctime.UnixMilli(),
	}
}

func (r *CacheUserRepository) entityToDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email.String,
		Password: u.Password,
		Phone:    u.Phone.String,
		Ctime:    time.UnixMilli(u.Ctime),
	}
}
