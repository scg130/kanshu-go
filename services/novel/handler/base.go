package handler

import (
	"context"

	"github.com/opentracing/opentracing-go"
)

func StartSpan(c context.Context, operate string) (sp opentracing.Span, ctx context.Context) {
	span := opentracing.SpanFromContext(c)
	sp = span.Tracer().StartSpan(operate, opentracing.ChildOf(span.Context()))
	ctx = opentracing.ContextWithSpan(c, sp)
	return
}
