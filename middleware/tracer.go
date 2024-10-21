package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

var t opentracing.Tracer

func init() {
	cfg := jaegercfg.Configuration{
		ServiceName: "runapp",
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
		},
	}

	host := os.Getenv("TRACER_HOST")
	port := os.Getenv("TRACER_PORT")
	if host == "" {
		panic("tracerAddr is invalid")
	}
	tracerAddr := fmt.Sprintf("%s:%s", host, port)

	sender, err := jaeger.NewUDPTransport(tracerAddr, 0)
	if err != nil {
		panic(err)
	}

	reporter := jaeger.NewRemoteReporter(sender)
	// Initialize tracer with a logger and a metrics factory
	tracer, _, err := cfg.NewTracer(
		jaegercfg.Reporter(reporter),
	)
	opentracing.SetGlobalTracer(tracer)
	t = tracer
}

func Tracer() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		span := t.StartSpan(ctx.Request.URL.Path)
		ctx.Set("span", span)
		defer span.Finish()

		// span1 := opentracing.StartSpan("111", opentracing.ChildOf(span.Context()))
		// span1.SetTag("key1", map[string]interface{}{"a": "b"})
		// time.Sleep(time.Microsecond * 20)
		// span1.Finish()

		// span2 := opentracing.StartSpan("222", opentracing.ChildOf(span.Context()))
		// time.Sleep(time.Microsecond * 200)
		// span3 := opentracing.StartSpan("333", opentracing.ChildOf(span2.Context()))
		// time.Sleep(time.Microsecond * 100)
		// span3.Finish()
		// span2.Finish()

		ctx.Next()
	}
}
