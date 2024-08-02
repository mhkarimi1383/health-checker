package checkers

import (
	"net"
	"strconv"
	"time"
)

type portCheckConfigParams struct {
	host string
	port int
}

type PortCheck struct {
	config CheckConfig
	params portCheckConfigParams
}

func (p *PortCheck) Enable(e ...bool) bool {
	if len(e) > 0 {
		p.config.Enabled = &e[0]
	}
	return *p.config.Enabled
}

func (p *PortCheck) GetConfig() CheckConfig {
	return p.config
}

func (p *PortCheck) SetConfig(config CheckConfig) error {
	params := portCheckConfigParams{
		host: config.Params["host"].(string),
		port: config.Params["port"].(int),
	}
	p.config = config
	p.params = params
	return nil
}

func (p *PortCheck) Check() Status {
	t := time.Now()
	c, err := net.DialTimeout("tcp", net.JoinHostPort(p.params.host, strconv.Itoa(p.params.port)), p.config.Timeout)
	if err != nil {
		e := err.Error()
		return Status{
			IsAlive: false,
			Error:   &e,
			Latency: time.Since(t),
			Type:    PORT,
		}
	}
	if c != nil {
		defer c.Close()
		return Status{
			IsAlive: true,
			Error:   nil,
			Latency: time.Since(t),
			Type:    PORT,
		}
	}
	e := "not connected"
	return Status{
		IsAlive: false,
		Error:   &e,
		Latency: time.Since(t),
		Type:    PORT,
	}
}
