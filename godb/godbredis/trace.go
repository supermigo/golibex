package godbredis

import (
	"context"
	"github.com/go-redis/redis/extra/rediscmd"
	"github.com/go-redis/redis/v8"
	tracer "github.com/supermigo/golibex/observability/tracing"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
	"strings"
)

type tracingHook struct {
	container *RedisContainer
}

func newTracingHook(container *RedisContainer) *tracingHook {
	return &tracingHook{container: container}
}

func (th *tracingHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	//TODO add enable
	newCtx, span := tracer.GetTracer("redis").Start(
		ctx, "redis:"+cmd.FullName(),
		tracer.WithSpanKind(tracer.SpanKindClient),
	)
	addrs := strings.Join(th.container.Opt.Addr, ",")
	span.SetAttributes(
		semconv.DBSystemRedis,
		semconv.NetHostPortKey.String(addrs),
		semconv.DBOperationKey.String(cmd.Name()),
	)
	return newCtx, nil
}
func (th *tracingHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	//TODO add enable
	span := tracer.SpanFromContext(ctx)
	if span == nil {
		return nil
	}
	if err := cmd.Err(); err != nil && err != redis.Nil {
		recordError(ctx, span, err)
		span.SetStatus(codes.Error, err.Error())
	} else {
		span.SetStatus(codes.Ok, "success")
	}
	span.End()
	return nil
}

func (th *tracingHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	//TODO add enable
	summary, cmdsString := rediscmd.CmdsString(cmds)

	newCtx, span := tracer.GetTracer("redis").Start(
		ctx, "redis-pipeline:"+summary,
		tracer.WithSpanKind(tracer.SpanKindClient),
	)
	addrs := strings.Join(th.container.Opt.Addr, ",")
	span.SetAttributes(
		semconv.DBSystemRedis,
		semconv.NetHostPortKey.String(addrs),
		semconv.DBOperationKey.String(cmdsString),
	)

	return newCtx, nil
}

func (th *tracingHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	//TODO add enable
	span := tracer.SpanFromContext(ctx)
	if span == nil {
		return nil
	}
	span.SetStatus(codes.Ok, "success")
	if err := cmds[0].Err(); err != nil && err != redis.Nil {
		recordError(ctx, span, err)
		span.SetStatus(codes.Error, err.Error())
	}
	span.End()
	return nil
}

func recordError(ctx context.Context, span tracer.Span, err error) {
	if err != redis.Nil {
		span.RecordError(err)
	}
}
