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
		ctx.Next()
	}
}
