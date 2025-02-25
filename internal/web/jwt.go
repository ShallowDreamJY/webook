package web

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"net/http"
	"strings"
	"time"
)

type jwtHandler struct {
	// access_token key
	atKey []byte
	// refresh_token key
	rtKey []byte
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

func NewJwtHandler() jwtHandler {
	return jwtHandler{
		atKey: []byte("R5iN7GRD73oWwBRLgJYJiIIei5bGahtX"),
		rtKey: []byte("R5iN7GRD73oWwBRLgJYIkIIei5bGahtX"),
	}

}

func (h jwtHandler) setLoginToken(ctx *gin.Context, uid int64) error {
	ssid := uuid.New().String()
	err := h.setJWTToken(ctx, uid, ssid)
	if err != nil {
		return err
	}
	err = h.setRefreshToken(ctx, uid, ssid)
	if err != nil {
		return err
	}
	return nil
}

func (h jwtHandler) setJWTToken(ctx *gin.Context, uid int64, ssid string) error {
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(1))),
		},
		Uid:       uid,
		Ssid:      ssid,
		UserAgent: ctx.Request.UserAgent(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(h.atKey)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return err
	}
	ctx.Header("x-jwt-token", tokenStr)
	return nil
}

func (h jwtHandler) setRefreshToken(ctx *gin.Context, uid int64, ssid string) error {
	claims := RefreshClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
		Uid:  uid,
		Ssid: ssid,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(h.rtKey)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return err
	}
	ctx.Header("x-refresh-token", tokenStr)
	return nil
}

func (h jwtHandler) ExtractToken(ctx *gin.Context) string {
	tokenHeader := ctx.GetHeader("Authorization")
	if tokenHeader == "" {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return ""
	}
	segs := strings.SplitN(tokenHeader, " ", 2)
	if len(segs) != 2 {
		return ""
	}
	return segs[1]
}
