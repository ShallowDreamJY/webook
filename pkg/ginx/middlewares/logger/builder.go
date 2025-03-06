package logger

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"io"
	"time"
)

type MiddlewareBuilder struct {
	allowReqBody  bool
	allowRespBody bool
	loggerFunc    func(ctx context.Context, al *AccessLog)
}

type AccessLog struct {
	// http 请求的方法
	Method string
	// url 整个请求的url
	Url      string
	ReqBody  string
	RespBody string
	Status   int
	duration string
}

type responseWriter struct {
	al *AccessLog
	gin.ResponseWriter
}

func NewBuilder(fn func(ctx context.Context, al *AccessLog)) *MiddlewareBuilder {
	return &MiddlewareBuilder{
		loggerFunc: fn,
	}
}

func (b *MiddlewareBuilder) AllowReqBody() *MiddlewareBuilder {
	b.allowReqBody = true
	return b
}

func (b *MiddlewareBuilder) AllowRespBody() *MiddlewareBuilder {
	b.allowRespBody = true
	return b
}

func (b *MiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		url := ctx.Request.URL.String()
		if len(url) > 1024 {
			url = url[:1024]
		}
		al := &AccessLog{
			Method: ctx.Request.Method,
			Url:    url,
		}
		if b.allowReqBody && ctx.Request.Body != nil {
			// 读取
			body, _ := io.ReadAll(ctx.Request.Body)
			// 放回去，但是消耗性能
			ctx.Request.Body = io.NopCloser(bytes.NewReader(body))
			if len(body) > 1024 {
				body = body[:1024]
			}
			al.ReqBody = string(body)
		}

		//if b.allowRespBody {
		//	ctx.Writer = responseWriter{
		//		al:             al,
		//		ResponseWriter: ctx.Writer,
		//	}
		//}

		defer func() {
			al.duration = time.Since(start).String()
			b.loggerFunc(ctx, al)
		}()

		// 执行到业务逻辑
		ctx.Next()
	}
}

func (w responseWriter) WriteHeader(statusCode int) {
	w.al.Status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w responseWriter) Write(data []byte) (int, error) {
	w.al.RespBody = string(data)
	return w.ResponseWriter.Write(data)
}

func (w responseWriter) WriteString(data string) (int, error) {
	w.al.RespBody = data
	return w.ResponseWriter.WriteString(data)
}
