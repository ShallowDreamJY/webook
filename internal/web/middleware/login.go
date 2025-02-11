package middleware

import (
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LoginMiddlewareBuilder struct {
	paths []string
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

func (l *LoginMiddlewareBuilder) IgnorePaths(path string) *LoginMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	gob.Register(time.Now())
	return func(ctx *gin.Context) {
		// 不需要登录校验的API
		for _, path := range l.paths {
			if ctx.Request.RequestURI == path {
				return
			}
		}
		//if ctx.Request.URL.Path == "/user/login" ||
		//	ctx.Request.URL.Path == "/user/signup" {
		//	return
		//}
		sess := sessions.Default(ctx)
		id := sess.Get("userId")
		if id != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		updateTime := sess.Get("update_time")
		sess.Set("userId", id)
		now := time.Now().UnixMilli()
		if updateTime == nil {
			updateTime = time.Now().UnixMilli()
			sess.Set("update_time", updateTime)
			sess.Options(sessions.Options{
				MaxAge: 60,
			})
			sess.Save()
			return
		}
		updateTimeVal, _ := updateTime.(int64)
		if now-updateTimeVal > 60*1000 {
			updateTime = time.Now().UnixMilli()
			sess.Set("update_time", updateTime)
			sess.Options(sessions.Options{
				MaxAge: 60,
			})
			sess.Save()
			return
		}

	}
}
