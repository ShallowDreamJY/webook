package service

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"webook/internal/domain"
	"webook/internal/repository"
)

var (
	ErrUserDuplicated        = repository.ErrUserDuplicated
	ErrInvalidUserOrPassword = errors.New("invalid user or password")
	ErrUserNotFound          = repository.ErrUserNotFound
)

type UserService interface {
	SignUp(ctx context.Context, u domain.User) error
	Login(ctx context.Context, email string, password string) (domain.User, error)
	FindOrCreate(ctx *gin.Context, phone string) (domain.User, error)
	FindOrCreateByWechat(ctx context.Context, wechatInfo domain.WechatInfo) (domain.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (svc *userService) SignUp(ctx context.Context, u domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return svc.repo.Create(ctx, u)
}

func (svc *userService) Login(ctx context.Context, email string, password string) (domain.User, error) {
	// 查询用户是否存在
	u, err := svc.repo.FindByEmail(ctx, email)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		// 打日志
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return u, nil
}

func (svc *userService) FindOrCreate(ctx *gin.Context, phone string) (domain.User, error) {
	// 先查询该phone是否注册过用户
	u, err := svc.repo.FindByPhone(ctx, phone)
	// 判断是否存在该用户
	if err != repository.ErrUserNotFound {
		// nil 和 不是usernotfound（有用户）的走这里
		return u, err
	}
	// 系统资源不足时，降级处理，则慢路径不走
	if ctx.Value("降级") == true {
		return domain.User{}, errors.New("系统降级")
	}
	// 没有用户则Craete
	u = domain.User{
		Phone: phone,
	}
	err = svc.repo.Create(ctx, u)
	if err != nil && err == repository.ErrUserDuplicated {
		return u, err
	}
	return svc.repo.FindByPhone(ctx, phone)
}

func (svc *userService) FindOrCreateByWechat(ctx context.Context, wechatInfo domain.WechatInfo) (domain.User, error) {
	// 先查询该phone是否注册过用户
	u, err := svc.repo.FindByWechat(ctx, wechatInfo.OpenId)
	// 判断是否存在该用户
	if err != repository.ErrUserNotFound {
		// nil 和 不是usernotfound（有用户）的走这里
		return u, err
	}
	// 系统资源不足时，降级处理，则慢路径不走
	if ctx.Value("降级") == true {
		return domain.User{}, errors.New("系统降级")
	}
	// 没有用户则Craete
	u = domain.User{
		WechatInfo: wechatInfo,
	}
	err = svc.repo.Create(ctx, u)
	if err != nil && err == repository.ErrUserDuplicated {
		return u, err
	}
	return svc.repo.FindByWechat(ctx, wechatInfo.OpenId)
}
