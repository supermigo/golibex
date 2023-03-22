package godbredis

import (
	"github.com/go-redis/redis/v8"
	"github.com/supermigo/golibex/observability/metrics"
	tracer "github.com/supermigo/golibex/observability/tracing"
)

func newSentinelOptions(opt *Option) *redis.FailoverOptions {
	return &redis.FailoverOptions{
		MasterName:         opt.MasterNames[0],
		SentinelAddrs:      opt.Addr,
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

func NewSentinelClient(opt *Option, metrics metrics.Provider, tracer tracer.Provider) *RedisContainer {
	baseClient := redis.NewFailoverClient(newSentinelOptions(opt))
	c := &RedisContainer{
		Redis:     baseClient,
		Opt:       *opt,
		redisType: RedisTypeSentinel,
		metrics:   metrics,
		tracer:    tracer,
	}
	// baseClient.AddHook(newMetricHook(c))
	if metrics != nil {
		metricsHook := newMetricHook(c, metrics)
		baseClient.AddHook(metricsHook)
	}
	// baseClient.AddHook(newTracingHook(c))
	return c
}
