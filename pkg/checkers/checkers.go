package checkers

import (
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	PORT       = "port"
	HTTP       = "http"
	POSTGRESQL = "postgresql"
	REDIS      = "redis"
)

type Status struct {
	Latency time.Duration
	Error   *string
	Type    string
	IsAlive bool
}

type CheckConfig struct {
	Type    string         `mapstructure:"type"`
	Timeout time.Duration  `mapstructure:"timeout"`
	Params  map[string]any `mapstructure:"params"`
	Enabled *bool          `mapstructure:"enabled"`
}

var defualtEnabled bool = true

type (
	CheckConfigs map[string]CheckConfig
	Checkers     map[string]Checker
	Statuses     map[string]Status
)

type Checker interface {
	// Sets configuration
	SetConfig(config CheckConfig) error
	// Gets configuration
	GetConfig() (config CheckConfig)
	// Checks if checker is enabled or not and changes enable value optionaly
	Enable(...bool) bool
	// Do the actual check
	Check() Status
}

func ConfigsToCheckers(configs CheckConfigs) (Checkers, error) {
	checkers := make(Checkers)
	caser := cases.Title(language.English)
	for name, c := range configs {
		if c.Enabled == nil {
			c.Enabled = &defualtEnabled
		}
		hName := caser.String(name)
		switch c.Type {
		case PORT:
			p := &PortCheck{}
			if err := p.SetConfig(c); err != nil {
				return checkers, err
			}
			checkers[hName] = p
		case HTTP:
			h := &HTTPCheck{}
			if err := h.SetConfig(c); err != nil {
				return checkers, err
			}
			checkers[hName] = h
		case POSTGRESQL:
			pg := &PGCheck{}
			if err := pg.SetConfig(c); err != nil {
				return checkers, err
			}
			checkers[hName] = pg
		case REDIS:
			r := &RedisCheck{}
			if err := r.SetConfig(c); err != nil {
				return checkers, err
			}
			checkers[hName] = r
		default:
			log.Warn().Msgf("Invalid check type %v, Ignoring...", c.Type)
			continue
		}
	}
	return checkers, nil
}

func RunChecks(chs Checkers) Statuses {
	s := make(Statuses)
	var wg sync.WaitGroup

	for name, ch := range chs {
		if ch.Enable() {
			wg.Add(1)
			go func(name string, ch Checker) {
				defer wg.Done()
				s[name] = ch.Check()
			}(name, ch)
		}
	}
	wg.Wait()
	return s
}
