package checkers

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type pgCheckConfigParams struct {
	url string
}

type PGCheck struct {
	config CheckConfig
	params pgCheckConfigParams
}

func (pg *PGCheck) Enable(e ...bool) bool {
	if len(e) > 0 {
		pg.config.Enabled = &e[0]
	}
	return *pg.config.Enabled
}

func (pg *PGCheck) GetConfig() CheckConfig {
	return pg.config
}

func (pg *PGCheck) SetConfig(config CheckConfig) error {
	params := pgCheckConfigParams{
		url: config.Params["url"].(string),
	}
	pg.config = config
	pg.params = params
	return nil
}

func (pg *PGCheck) Check() Status {
	ctx, cancel := context.WithTimeout(context.Background(), pg.config.Timeout)
	defer cancel()
	t := time.Now()
	c, err := pgx.Connect(ctx, pg.params.url)
	if err != nil {
		e := err.Error()
		return Status{
			IsAlive: false,
			Error:   &e,
			Latency: time.Since(t),
			Type:    POSTGRESQL,
		}
	}
	defer c.Close(ctx)
	err = c.Ping(ctx)
	if err != nil {
		e := err.Error()
		return Status{
			IsAlive: false,
			Error:   &e,
			Latency: time.Since(t),
			Type:    POSTGRESQL,
		}
	}
	return Status{
		IsAlive: true,
		Error:   nil,
		Latency: time.Since(t),
		Type:    POSTGRESQL,
	}
}
