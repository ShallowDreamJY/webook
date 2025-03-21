package ioc

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
	"webook/internal/web"
	ijwt "webook/internal/web/jwt"
	"webook/internal/web/middleware"
	"webook/pkg/ginx/middlewares/logger"
	"webook/pkg/ginx/middlewares/ratelimit"
	logger2 "webook/pkg/logger"
)

func InitWebServer(mdls []gin.HandlerFunc, hdl *web.UserHandler,
	oauth2WechatHdl *web.OAuth2WechatHandler,
	articleHdl *web.ArticleHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	hdl.RegisterUserRoutes(server)
	oauth2WechatHdl.RegisterRoutes(server)
	articleHdl.RegisterArticleRoutes(server)
	return server
}

func InitMiddlewares(redisClient redis.Cmdable,
	jwtHdl ijwt.Handler,
	l logger2.LoggerV1) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		cors.New(cors.Config{
			//AllowOrigins:     []string{"https://localhost:3000"},
			AllowMethods:     []string{"PUT", "PATCH", "POST", "GET"},
			AllowHeaders:     []string{"Content-Type", "Authorization"},
			AllowCredentials: true,
			ExposeHeaders:    []string{"x-jwt-token", "x-refresh-token"},
			AllowOriginFunc: func(origin string) bool {
				if strings.Contains(origin, "localhost") {
					return true
				}
				return strings.Contains(origin, "my.com")
			},
			MaxAge: 12 * time.Hour,
		}),
		middleware.NewLoginJWTMiddlewareBuilder(jwtHdl).
			IgnorePaths("/users/signup").
			IgnorePaths("/users/login").
			IgnorePaths("/users/login_sms/code/send").
			IgnorePaths("/users/login_sms").
			IgnorePaths("/users/refresh_token").
			IgnorePaths("/oauth2/wechat/authurl").
			Build(),
		ratelimit.NewBuilder(redisClient, time.Second, 100).Build(),
		logger.NewBuilder(func(ctx context.Context, al *logger.AccessLog) {
			l.Debug("HTTP请求", logger2.Field{Key: "al", Value: al})
		}).AllowReqBody().AllowRespBody().Build(),
	}
}
