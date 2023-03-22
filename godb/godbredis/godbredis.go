package godbredis

import (
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/supermigo/golibex/observability/metrics"
	tracer "github.com/supermigo/golibex/observability/tracing"
	"io"
)

const (
	RedisTypeClient   = "client"
	RedisTypeCluster  = "cluster"
	RedisTypeSentinel = "sentinel"
)

type Redis interface {
	redis.Cmdable
	io.Closer
}

func New(opt *Option, metrics metrics.Provider, tracer tracer.Provider) (*RedisContainer, error) {
	return newWithOption(opt, metrics, tracer)
}

func newWithOption(opt *Option, metrics metrics.Provider, tracer tracer.Provider) (*RedisContainer, error) {
	if err := checkOptions(opt); err != nil {
		return nil, err
	}

	var c *RedisContainer
	// 判断连接类型
	switch opt.ConnType {
	case RedisTypeClient:
		c = NewClient(opt, metrics, tracer)
	case RedisTypeCluster:
		c = NewClusterClient(opt, metrics, tracer)
	case RedisTypeSentinel:
		if len(opt.MasterNames) == 0 {
			err := errors.New("empty master name")
			return nil, err
		}
		c = NewSentinelClient(opt, metrics, tracer)
	default:
		err := errors.New("redis connection type need ")
		return nil, err
	}

	return c, nil
}

// RedisContainer 用来存储redis客户端及其额外信息
type RedisContainer struct {
	Redis
	Opt       Option
	redisType string
	metrics   metrics.Provider
	tracer    tracer.Provider
}
