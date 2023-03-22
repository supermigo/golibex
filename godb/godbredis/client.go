package godbredis

import (
	"github.com/go-redis/redis/v8"
	"github.com/supermigo/golibex/observability/metrics"
	tracer "github.com/supermigo/golibex/observability/tracing"
)

func newClientOptions(opt *Option) *redis.Options {
	return &redis.Options{
		Addr:               opt.Addr[0],
		Username:           opt.Username,
		Password:           opt.Password,
		DB:                 opt.DB,
		MaxRetries:         opt.MaxRetries,
		MinRetryBackoff:    opt.MinRetryBackoff,
		MaxRetryBackoff:    opt.MaxRetryBackoff,
		DialTimeout:        opt.DialTimeout,
		ReadTimeout:        opt.ReadTimeout,
		WriteTimeout:       opt.WriteTimeout,
		PoolSize:           opt.PoolSize,
		MinIdleConns:       opt.MinIdleConns,
		MaxConnAge:         opt.MaxConnAge,
		PoolTimeout:        opt.PoolTimeout,
		IdleTimeout:        opt.IdleTimeout,
		IdleCheckFrequency: opt.IdleCheckFrequency,
		TLSConfig:          opt.TLSConfig,
	}
}

func NewClient(opt *Option, metrics metrics.Provider, tracer tracer.Provider) *RedisContainer {
	baseClient := redis.NewClient(newClientOptions(opt))
	c := &RedisContainer{
		Redis:     baseClient,
		Opt:       *opt,
		redisType: RedisTypeClient,
		metrics:   metrics,
		tracer:    tracer,
	}

	if metrics != nil {
		metricsHook := newMetricHook(c, metrics)
		baseClient.AddHook(metricsHook)
	}
	// baseClient.AddHook(newMetricHook(c))
	if opt.EnableTracer {
		baseClient.AddHook(newTracingHook(c))
	}
	return c
}
