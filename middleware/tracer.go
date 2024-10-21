package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
)

func Tracer() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		t := opentracing.GlobalTracer()
		span := t.StartSpan(ctx.Request.URL.Path)
		ctx.Set("span", span)
		defer span.Finish()

		span1 := opentracing.StartSpan("111", opentracing.ChildOf(span.Context()))
		span1.SetTag("key1", map[string]interface{}{"a": "b"})
		time.Sleep(time.Microsecond * 20)
		span1.Finish()

		span2 := opentracing.StartSpan("222", opentracing.ChildOf(span.Context()))
		time.Sleep(time.Microsecond * 200)
		span3 := opentracing.StartSpan("333", opentracing.ChildOf(span2.Context()))
		time.Sleep(time.Microsecond * 100)
		span3.Finish()
		span2.Finish()
		ctx.Next()
	}
}
