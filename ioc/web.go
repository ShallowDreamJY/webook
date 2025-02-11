package ioc

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
	"webook/internal/web"
	"webook/internal/web/middleware"
	"webook/pkg/ginx/middlewares/ratelimit"
)

func InitGin(mdls []gin.HandlerFunc, hdl *web.UserHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	hdl.RegisterUserRoutes(server)
	return server
}

func InitMiddlewares(redisClient redis.Cmdable) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		cors.New(cors.Config{
			//AllowOrigins:     []string{"https://localhost:3000"},
			AllowMethods:     []string{"PUT", "PATCH", "POST", "GET"},
			AllowHeaders:     []string{"Content-Type", "Authorization"},
			AllowCredentials: true,
			ExposeHeaders:    []string{"x-jwt-token"},
			AllowOriginFunc: func(origin string) bool {
				if strings.Contains(origin, "localhost") {
					return true
				}
				return strings.Contains(origin, "my.com")
			},
			MaxAge: 12 * time.Hour,
		}),
		middleware.NewLoginJWTMiddlewareBuilder().
			IgnorePaths("/users/signup").
			IgnorePaths("/users/login").
			IgnorePaths("/users/login_sms/code/send").
			IgnorePaths("/users/login_sms").
			Build(),
		ratelimit.NewBuilder(redisClient, time.Second, 100).Build(),
	}
}
