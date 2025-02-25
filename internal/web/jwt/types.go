package jwt

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Handler interface {
	ExtractToken(ctx *gin.Context) string
	SetJWTToken(ctx *gin.Context, uid int64, ssid string) error
	// SetRefreshToken(ctx *gin.Context, uid int64, ssid string) error
	SetLoginToken(ctx *gin.Context, uid int64) error
	ClearToken(ctx *gin.Context) error
	CheckSession(ctx *gin.Context, ssid string) error
}

type UserClaims struct {
	// 组合
	jwt.RegisteredClaims
	// 声明要放进去token里的数据
	Uid       int64
	Ssid      string
	UserAgent string
}

type RefreshClaims struct {
	Uid  int64
	Ssid string
	jwt.RegisteredClaims
}
