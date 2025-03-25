//go:build wireinject

package startup

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"webook/internal/repository"
	"webook/internal/repository/article"
	"webook/internal/repository/cache"
	"webook/internal/repository/dao"
	articledao "webook/internal/repository/dao/article"
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
		articledao.NewReaderDAO,
		articledao.NewAuthorDAO,
		articledao.NewGORMArticleDao,

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
		articledao.NewGORMArticleDao,
		articledao.NewReaderDAO,
		articledao.NewAuthorDAO)
	return &web.ArticleHandler{}
}
