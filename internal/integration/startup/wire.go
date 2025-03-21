//go:build wireinject

package startup

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"webook/internal/repository"
	"webook/internal/repository/article"
	"webook/internal/repository/cache"
	"webook/internal/repository/dao"
	"webook/internal/service"
	"webook/internal/web"
	ijwt "webook/internal/web/jwt"
	"webook/ioc"
)

var thirdProvider = wire.NewSet(InitRedis, InitTestDB, InitLogger)

func InitWebServer() *gin.Engine {
	wire.Build(
		// 最基础的初始化
		ioc.InitDB, ioc.InitRedis, ioc.InitLogger,

		dao.NewUserDao,
		dao.article.NewGORMArticleDao,

		cache.NewUserCache,
		cache.NewRedisCodeCache,

		repository.NewUserRepository,
		repository.NewCodeRepository,
		article.NewArticleRepository,

		service.NewUserService,
		service.NewCodeService,
		service.NewArticleService,
		ioc.InitSMSService,
		ioc.InitOAuth2WechatService,

		web.NewUserHandler,
		web.NewOAuth2WechatHandler,
		ijwt.NewRedisJWTHandler,
		web.NewArticleHandler,

		ioc.InitWebServer,
		ioc.InitMiddlewares,
	)
	return new(gin.Engine)
}

func InitArticleHandler() *web.ArticleHandler {
	wire.Build(thirdProvider,
		service.NewArticleService,
		web.NewArticleHandler,
		article.NewArticleRepository,
		dao.article.NewGORMArticleDao)
	return &web.ArticleHandler{}
}
