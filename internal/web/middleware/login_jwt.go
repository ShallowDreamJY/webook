package middleware

import (
	"encoding/gob"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"time"
	ijwt "webook/internal/web/jwt"
)

// LoginJWTMiddlewareBuilder JWT登录校验
type LoginJWTMiddlewareBuilder struct {
	paths []string
	ijwt.Handler
	cmd redis.Cmdable
}

func NewLoginJWTMiddlewareBuilder(jwtHdl ijwt.Handler) *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{
		Handler: jwtHdl,
	}
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
		tokenStr := l.ExtractToken(ctx)
		claims := &ijwt.UserClaims{}
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
		// 判断是否处于登出状态
		err = l.CheckSession(ctx, claims.Ssid)
		if err != nil {
			// redis有问题，或者已经登出
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Hour * 1))
		tokenStr, err = token.SignedString([]byte("R5iN7GRD73oWwBRLgJYJiIIei5bGahtX"))
		if err != nil {
			log.Println("jwt 续约失败")
		}
		ctx.Header("x-jwt-token", tokenStr)
		ctx.Set("claims", claims)
	}
}
