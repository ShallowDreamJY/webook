package ratelimit

import "context"

type Limiter interface {
	// Limited有没有触发限流，key就是限流对象
	// bool 代表是否限流，true则限流
	// error 代表是否出错
	Limit(ctx context.Context, key string) (bool, error)
}
