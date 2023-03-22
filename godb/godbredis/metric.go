package godbredis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/supermigo/golibex/observability/metrics"
	"time"
)

type redisMetricKey string

const (
	keyRequestStart     redisMetricKey = "requestStart"
	keyPipeRequestStart redisMetricKey = "pipeRequestStart"
)

const (
	metricsRequestTotalName         = "redis_request_total"
	metricsRequestDurationName      = "redis_request_duration"
	metricsRequestDurationRangeName = "redis_request_duration_range"
	metricsRequestErrorName         = "redis_request_error"
)

var (
	metricRequestTotal          metrics.Counter
	metricsRequestDuration      metrics.Gauge
	metricsRequestDurationRange metrics.Histogram
	metricsRequestError         metrics.Counter
	labelValues                 = []string{"host", "operation"}
	histogramBukets             = []float64{.003, .01, .02, .05, 1} // ms unit
)

var _ redis.Hook = &metricHook{}

type metricHook struct {
	container *RedisContainer
	enable    bool
}

func newMetricHook(container *RedisContainer, metrics metrics.Provider) *metricHook {
	if metrics != nil {
		metricRequestTotal = metrics.NewCounter(metricsRequestTotalName, labelValues...)
		metricsRequestDuration = metrics.NewGauge(metricsRequestDurationName, labelValues...)
		metricsRequestDurationRange = metrics.NewHistogram(metricsRequestDurationRangeName, histogramBukets, labelValues...)
		metricsRequestError = metrics.NewCounter(metricsRequestErrorName, labelValues...)
	}

	return &metricHook{
		container: container,
		enable:    metrics != nil,
	}
}

func (h *metricHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	if !h.enable {
		return ctx, nil
	}
	// host := h.container.Opt.Addr[0]
	// stats := collectors.RedisCollector().OnStart(host, cmd.Name())
	now := time.Now()
	ctx = context.WithValue(ctx, keyRequestStart, now)
	return ctx, nil
}

func (h *metricHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	if !h.enable {
		return nil
	}

	// 万一有人改写了数据，这里只能打日志再退出
	v := ctx.Value(keyRequestStart)
	if v == nil {
		return nil
	}

	// stats, ok := v.(*nssredis.StatsHolder)

	// if !ok {
	// 	h.logger.Errorf("convert to redis.Stats failed from %s", keyRequestStart)
	// 	return nil
	// }
	start, ok := v.(time.Time)
	if !ok {
		return nil
	}

	host := h.container.Opt.Addr[0]
	opertation := cmd.Name()

	// 忽略key不存在情况
	if cmd.Err() != nil && cmd.Err() != redis.Nil {
		// collectors.RedisCollector().OnError(stats, cmd.Err())
		metricsRequestError.With("host", host, "operation", opertation).Add(1)
	}

	// collectors.RedisCollector().OnComplete(stats, cmd.Err() == nil)
	du := float64(time.Now().Nanosecond()-start.Nanosecond()) / 1e6
	metricRequestTotal.With("host", host, "operation", opertation).Add(1)
	metricsRequestDuration.With("host", host, "operation", opertation).Set(du)
	metricsRequestDurationRange.With("host", host, "operation", opertation).Observe(du)
	return nil
}

func (h *metricHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	if !h.enable {
		return ctx, nil
	}
	// statses := make(map[redis.Cmder]*nssredis.StatsHolder, len(cmds))
	// for i := range cmds {
	// 	host := h.container.Opt.Addr[0]
	// 	statses[cmds[i]] = collectors.RedisCollector().OnStart(host, cmds[i].Name())
	// }
	// ctx = context.WithValue(ctx, keyPipeRequestStart, statses)
	now := time.Now()
	ctx = context.WithValue(ctx, keyRequestStart, now)
	return ctx, nil
}

func (h *metricHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	if !h.enable {
		return nil
	}

	// 万一有人改写了数据，这里只能打日志再退出
	v := ctx.Value(keyPipeRequestStart)
	if v == nil {
		return nil
	}

	// statses, ok := v.(map[redis.Cmder]*nssredis.StatsHolder)

	// if !ok {
	// 	h.logger.Errorf("convert to redis.Stats failed from %s", keyPipeRequestStart)
	// 	return nil
	// }
	start, ok := v.(time.Time)
	if !ok {
		return nil
	}
	du := float64(time.Now().Nanosecond()-start.Nanosecond()) / 1e6
	for i := range cmds {
		host := h.container.Opt.Addr[0]
		opertation := cmds[i].Name()
		// 忽略key不存在情况
		if cmds[i].Err() != nil && cmds[i].Err() != redis.Nil {
			// collectors.RedisCollector().OnError(statses[cmds[i]], cmds[i].Err())
			metricsRequestError.With("host", host, "operation", opertation).Add(1)
		}
		// collectors.RedisCollector().OnComplete(statses[cmds[i]], cmds[i].Err() == nil)
		metricRequestTotal.With("host", host, "operation", opertation).Add(1)
		metricsRequestDuration.With("host", host, "operation", opertation).Set(du)
		metricsRequestDurationRange.With("host", host, "operation", opertation).Observe(du)
	}
	return nil
}
