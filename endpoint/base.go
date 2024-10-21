package endpoint

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
)

func GetSpanTraceIDAndCtx(ctx *gin.Context) (span opentracing.Span, traceID jaeger.TraceID, c context.Context) {
	sp := ctx.MustGet("span")
	span = sp.(opentracing.Span)
	traceID = span.Context().(jaeger.SpanContext).TraceID()
	c = opentracing.ContextWithSpan(ctx, span)
	return
}
