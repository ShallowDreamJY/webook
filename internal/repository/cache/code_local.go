package cache

import (
	"context"
	"errors"
	"fmt"
	lru "github.com/hashicorp/golang-lru"
	"sync"
	"time"
)

type LocalCodeCache struct {
	cache      *lru.Cache
	lock       sync.Mutex
	expireTime time.Duration
}

type codeItem struct {
	code   string
	cnt    int
	expire time.Time
}

func NewLocalCodeCache() *LocalCodeCache {
	return &LocalCodeCache{}
}

func (l *LocalCodeCache) Set(ctx context.Context, biz, phone, code string) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	key := l.key(biz, phone)
	now := time.Now()
	val, ok := l.cache.Get(key)
	if !ok {
		l.cache.Add(key, codeItem{
			code:   code,
			cnt:    3,
			expire: now.Add(l.expireTime),
		})
		return nil
	}
	itm, ok := val.(codeItem)
	if !ok {
		return errors.New("系统错误")
	}
	if itm.expire.Sub(now) > time.Minute*9 {
		return ErrSendCodeTooMany
	}
	l.cache.Add(key, codeItem{
		code:   code,
		cnt:    3,
		expire: now.Add(l.expireTime),
	})
	return nil
}

func (l *LocalCodeCache) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	l.lock.Lock()
	defer l.lock.Unlock()

	key := l.key(biz, phone)
	now := time.Now()
	val, ok := l.cache.Get(key)
	if !ok {
		return false, ErrUnknownForCode
	}
	itm, ok := val.(codeItem)
	if !ok {
		return false, errors.New("系统错误")
	}
	if itm.expire.Sub(now) < 0 || itm.cnt <= 0 {
		return false, ErrVerifyTooManyTimes
	}
	// 输入错误
	if itm.code != inputCode {
		l.cache.Add(key, codeItem{
			code:   itm.code,
			cnt:    itm.cnt - 1,
			expire: itm.expire,
		})
		return false, nil
	}
	return true, nil
}

func (l *LocalCodeCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}
