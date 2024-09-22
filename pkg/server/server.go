package server

import (
	"encoding/json"
	"html/template"
	"runtime"
	"time"

	"github.com/fasthttp/router"
	"github.com/go-co-op/gocron/v2"
	"github.com/li-jin-gou/http2curl"
	"github.com/mhkarimi1383/health-checker/pkg/checkers"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"
	mjson "github.com/tdewolff/minify/v2/json"
)

var (
	success      bool = false
	statusData   checkers.Statuses
	htmlTemplate = template.Must(template.ParseFiles("templates/index.html"))
	textTemplate = template.Must(template.ParseFiles("templates/index.txt"))
	minifier     = minify.New()
)

func init() {
	minifier.AddFunc("text/html", html.Minify)
	minifier.AddFunc("application/json", mjson.Minify)
}

type pageData struct {
	StatusData    checkers.Statuses
	Title         string
	OverallStatus bool
}

func minifyResponse(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		h(ctx)
		if ctx.QueryArgs().Peek("pretty") != nil {
			return
		}
		minified, err := minifier.Bytes(string(ctx.Response.Header.Peek("Content-Type")), ctx.Response.Body())
		if err != nil {
			log.Warn().Err(err).Msg("Minifing HTTP client response")
			return
		}
		ctx.Response.SetBody(minified)
	}
}

func Start(listenAddress string, interval time.Duration, chs checkers.Checkers) error {
	log.Info().Msg("Creating Scheduler")
	c, err := gocron.NewScheduler()
	if err != nil {
		return err
	}
	log.Info().Msg("Creating Scheduler Job")
	_, err = c.NewJob(
		gocron.DurationJob(
			interval,
		),
		gocron.NewTask(
			func(chs checkers.Checkers) {
				status := checkers.RunChecks(chs)
				success = true
				for name, s := range status {
					log.Info().Dur("latency", s.Latency).Bool("isAlive", s.IsAlive).Any("error", s.Error).Str("name", name).Str("type", s.Type).Msg("Health Check status")
					if !s.IsAlive {
						success = false
					}
				}
				statusData = status
			},
			chs,
		),
	)
	if err != nil {
		return err
	}
	log.Info().Msg("Starting Scheduler")
	c.Start()
	r := router.New()
	s := &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			begin := time.Now()
			r.Handler(ctx)
			end := time.Now()
			log.Info().Dur("duration", end.Sub(begin)).Bytes("method", ctx.Method()).Bytes("uri", ctx.URI().FullURI()).Str("remote_addr", ctx.RemoteAddr().String()).Bytes("user_agent", ctx.UserAgent()).Msg("Server Access Log")
			c, _ := http2curl.GetCurlCommandFastHttp(&ctx.Request)
			log.Trace().Str("curl", c.String()).Msg("Incomming Request cURL")
		},
		Name:        "Health Checker",
		Logger:      &log.Logger,
		Concurrency: runtime.GOMAXPROCS(0),
	}
	r.GET("/status", minifyResponse(getStatus))
	r.HEAD("/status", headStatus)
	return s.ListenAndServe(listenAddress)
}

func headStatus(ctx *fasthttp.RequestCtx) {
	if !success {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
	ctx.SetStatusCode(fasthttp.StatusNoContent)
}

func returnHTTPError(ctx *fasthttp.RequestCtx, err error) {
	log.Error().Err(err).Msg("Generating HTTP client response")
	ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	ctx.WriteString(fasthttp.StatusMessage(fasthttp.StatusInternalServerError))
}

func getStatus(ctx *fasthttp.RequestCtx) {
	accept := ctx.Request.Header.Peek(fasthttp.HeaderAccept)
	cType := getType(string(accept))
	data := pageData{
		StatusData:    statusData,
		Title:         "Health Checker",
		OverallStatus: success,
	}
	switch cType {
	case applicationJSONContentType:
		d, err := json.Marshal(data.StatusData)
		if err != nil {
			returnHTTPError(ctx, err)
			return
		}
		ctx.Write(d)
	case textHtmlContentType:
		err := htmlTemplate.Execute(ctx, data)
		if err != nil {
			returnHTTPError(ctx, err)
			return
		}
	case textPlainContentType:
		err := textTemplate.Execute(ctx, data)
		if err != nil {
			returnHTTPError(ctx, err)
			return
		}
	default:
		ctx.SetStatusCode(fasthttp.StatusNotAcceptable)
		ctx.WriteString(fasthttp.StatusMessage(fasthttp.StatusNotAcceptable))
		return
	}
	ctx.SetContentType(cType)
	if !success {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	}
}
