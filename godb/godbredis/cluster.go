package godbredis

import (
	"github.com/go-redis/redis/v8"
	"github.com/supermigo/golibex/observability/metrics"
	tracer "github.com/supermigo/golibex/observability/tracing"
)

func newClusterOptions(opt *Option) *redis.ClusterOptions {
	return &redis.ClusterOptions{
		Addrs:              opt.Addr,
		Username:           opt.Username,
		Password:           opt.Password,
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

func NewClusterClient(opt *Option, metrics metrics.Provider, tracer tracer.Provider) *RedisContainer {
	baseClient := redis.NewClusterClient(newClusterOptions(opt))
	c := &RedisContainer{
		Redis:     baseClient,
		Opt:       *opt,
		redisType: RedisTypeCluster,
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
