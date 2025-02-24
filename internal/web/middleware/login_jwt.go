package middleware

import (
	"encoding/gob"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"strings"
	"time"
	"webook/internal/web"
)

// LoginJWTMiddlewareBuilder JWT登录校验
type LoginJWTMiddlewareBuilder struct {
	paths []string
}

func NewLoginJWTMiddlewareBuilder() *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{}
}

func (l *LoginJWTMiddlewareBuilder) IgnorePaths(path string) *LoginJWTMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	gob.Register(time.Now())
	return func(ctx *gin.Context) {
		// 不需要登录校验的API
		for _, path := range l.paths {
			if ctx.Request.RequestURI == path {
				return
			}
		}
		// 获取JWT
		tokenHeader := ctx.GetHeader("Authorization")
		if tokenHeader == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		segs := strings.SplitN(tokenHeader, " ", 2)
		if len(segs) != 2 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := segs[1]
		claims := &web.UserClaims{}
		// ParseWithClaims里一定要传入指针
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("R5iN7GRD73oWwBRLgJYJiIIei5bGahtX"), nil
		})
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if token == nil || !token.Valid || claims.Uid == 0 {
			// 没登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// 监控userAgent是否一致
		//if claims.UserAgent != ctx.Request.UserAgent() {
		//	// 存在安全风险，userAgent不一致
		//	// 加入监控
		//	ctx.AbortWithStatus(http.StatusUnauthorized)
		//	return
		//}
		claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour * 1))
		tokenStr, err = token.SignedString([]byte("R5iN7GRD73oWwBRLgJYJiIIei5bGahtX"))
		if err != nil {
			log.Println("jwt 续约失败")
		}
		ctx.Header("x-jwt-token", tokenStr)
		ctx.Set("claims", claims)
	}
}
