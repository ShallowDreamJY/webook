package main

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	//db := initDB()
	initViper()
	server := InitWebServer()
	//rdb := initRedis()
	//u := initUser(db, rdb)
	// u.RegisterUserRoutes(server)
	//server := gin.Default()
	//server.GET("/hello", func(ctx *gin.Context) {
	//	ctx.String(http.StatusOK, "hello k8s")
	//})
	server.Run("localhost:8080")
}

func initViper() {
	cfile := pflag.String("config", "config/config.yaml", "config file path")
	pflag.Parse()
	viper.SetConfigFile(*cfile)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

//func initWebServer() *gin.Engine {
//	server := gin.Default()
//	// 跨域问题
//	server.Use(cors.New(cors.Config{
//		//AllowOrigins:     []string{"https://localhost:3000"},
//		AllowMethods:     []string{"PUT", "PATCH", "POST", "GET"},
//		AllowHeaders:     []string{"Content-Type", "Authorization"},
//		AllowCredentials: true,
//		ExposeHeaders:    []string{"x-jwt-token"},
//		AllowOriginFunc: func(origin string) bool {
//			if strings.Contains(origin, "localhost") {
//				return true
//			}
//			return strings.Contains(origin, "my.com")
//		},
//		MaxAge: 12 * time.Hour,
//	}))
//	// 限流
//	redisClient := redis.NewClient(&redis.Options{
//		Addr: config.Config.Redis.Addr,
//	})
//	server.Use(ratelimit.NewBuilder(redisClient, time.Second, 100).Build())
//	// 跨域问题
//	//server.Use(cors.New(cors.Config{
//	//	//AllowMethods:     []string{"PUT", "PATCH", "POST", "GET"},
//	//	AllowHeaders:     []string{"Content-Type", "Authorization"},
//	//	AllowCredentials: true,
//	//	AllowOriginFunc: func(origin string) bool {
//	//		return strings.HasPrefix(origin, "http://localhost")
//	//	},
//	//	MaxAge: 12 * time.Hour,
//	//}))
//	// session机制
//	//store, err := redis.NewStore(16, "tcp", "120.55.48.65:16379",
//	//	"", []byte("R5iN7GRD73oWwBRLgJYJiIIei5bGahtX"), []byte("qACtJYKaTOcAu3EGZH88JmJPvgw9bEzz"))
//	//if err != nil {
//	//	panic(err)
//	//}
//	store := memstore.NewStore([]byte("R5iN7GRD73oWwBRLgJYJiIIei5bGahtX"), []byte("qACtJYKaTOcAu3EGZH88JmJPvgw9bEzz"))
//	server.Use(sessions.Sessions("mysession", store))
//	// jwt机制
//	server.Use(middleware.NewLoginJWTMiddlewareBuilder().
//		IgnorePaths("/users/signup").
//		IgnorePaths("/users/login").
//		IgnorePaths("/users/login_sms/code/send").
//		IgnorePaths("/users/login_sms").
//		Build())
//	return server
//}
//
//func initRedis() redis.Cmdable {
//	redisClient := redis.NewClient(&redis.Options{
//		Addr: config.Config.Redis.Addr,
//	})
//	return redisClient
//}
//
//func initUser(db *gorm.DB, rdb redis.Cmdable) *web.UserHandler {
//	ud := dao.NewUserDao(db)
//	uc := cache.NewUserCache(rdb)
//	repo := repository.NewUserRepository(ud, uc)
//	userSvc := service.NewUserService(repo)
//	codeCache := cache.NewCodeCache(rdb)
//	codeRepo := repository.NewCodeRepository(codeCache)
//	// 基于内存的实现
//	smsSvc := memory.NewService()
//	codeSvc := service.NewCodeService(codeRepo, smsSvc)
//	u := web.NewUserHandler(userSvc, codeSvc)
//	return u
//}
//
//func initDB() *gorm.DB {
//	db, err := gorm.Open(mysql.Open("root:SFuuC4vUQ3qzaD0@tcp(120.55.48.65:13306)/webook"))
//	if err != nil {
//		// 只会在初始化过程中 panic
//		// panic相当于整个goroutine结束
//		// 初始化出错，应用不要启动
//		panic(err)
//	}
//	err = dao.InitTable(db)
//	if err != nil {
//		panic(err)
//	}
//
//	return db
//}
