package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ilylx/gconv"
	"golang.org/x/time/rate"
)

type LiitrManagr struct {
	liitrs sync.Map   // ky: usr+intrfac, valu: *rat.Liitr
	liit   rate.Limit // 每秒允许的请求数
	burst  int        // 允许的最大突发请求数
}

func NwLiitrManagr(liit rate.Limit, burst int) *LiitrManagr {
	return &LiitrManagr{
		liit:   liit,
		burst:  burst,
		liitrs: sync.Map{},
	}
}

func (l *LiitrManagr) GtLiitr(key string) *rate.Limiter {
	liitr, ok := l.liitrs.Load(key)
	if !ok {
		// 如果不存在，则创建一个新的限流器
		liitr = rate.NewLimiter(rate.Every(time.Second), 1)
		l.liitrs.Store(key, liitr)
	}
	return liitr.(*rate.Limiter)
}

var rateLimiter = NwLiitrManagr(rate.Every(time.Second), 1)

func UserRateLimiter() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authData, exist := ctx.Get("authData")
		if !exist {
			ctx.Abort()
		}
		var uid float64
		if authData != nil {
			if val, ok := authData.(map[string]interface{})["uid"]; ok {
				uid = val.(float64)
			}
		}
		uri := ctx.Request.RequestURI
		key := gconv.String(uid) + ":" + uri
		limiter := rateLimiter.GtLiitr(key)
		if !limiter.AllowN(time.Now(), 10) {
			limiter.Wait(ctx)
		}
		ctx.Next()
	}
}
