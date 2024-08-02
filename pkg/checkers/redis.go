package checkers

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisCheckConfigParams struct {
	url string
}

type RedisCheck struct {
	config CheckConfig
	params redisCheckConfigParams
}

func (r *RedisCheck) Enable(e ...bool) bool {
	if len(e) > 0 {
		r.config.Enabled = &e[0]
	}
	return *r.config.Enabled
}

func (r *RedisCheck) GetConfig() CheckConfig {
	return r.config
}

func (r *RedisCheck) SetConfig(config CheckConfig) error {
	params := redisCheckConfigParams{
		url: config.Params["url"].(string),
	}
	r.config = config
	r.params = params
	return nil
}

func (r *RedisCheck) Check() Status {
	ctx, cancel := context.WithTimeout(context.Background(), r.config.Timeout)
	defer cancel()
	opts, err := redis.ParseURL(r.params.url)
	if err != nil {
		e := err.Error()
		return Status{
			IsAlive: false,
			Error:   &e,
			Latency: 0,
			Type:    REDIS,
		}
	}
	opts.DialTimeout = r.config.Timeout
	opts.PoolTimeout = r.config.Timeout
	opts.ReadTimeout = r.config.Timeout
	opts.WriteTimeout = r.config.Timeout
	t := time.Now()
	c := redis.NewClient(opts)
	defer c.Close()
	s := c.Ping(ctx)
	if s.Err() != nil {
		e := s.Err().Error()
		return Status{
			IsAlive: false,
			Error:   &e,
			Latency: time.Since(t),
			Type:    REDIS,
		}
	}
	s.Result()
	return Status{
		IsAlive: true,
		Error:   nil,
		Latency: time.Since(t),
		Type:    REDIS,
	}
}
