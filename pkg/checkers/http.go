package checkers

import (
	"fmt"
	"time"

	"github.com/valyala/fasthttp"
)

type httpCheckConfigParams struct {
	url     string
	headers map[string]any
}

type HTTPCheck struct {
	config CheckConfig
	params httpCheckConfigParams
}

func (h *HTTPCheck) Enable(e ...bool) bool {
	if len(e) > 0 {
		h.config.Enabled = &e[0]
	}
	return *h.config.Enabled
}

func (h *HTTPCheck) GetConfig() CheckConfig {
	return h.config
}

func (h *HTTPCheck) SetConfig(config CheckConfig) error {
	if config.Params["headers"] == nil {
		config.Params["headers"] = make(map[string]string)
	}
	params := httpCheckConfigParams{
		url:     config.Params["url"].(string),
		headers: config.Params["headers"].(map[string]any),
	}
	h.config = config
	h.params = params
	return nil
}

func (h *HTTPCheck) Check() Status {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(h.params.url)
	for k, v := range h.params.headers {
		req.Header.Set(k, fmt.Sprintf("%s", v))
	}
	req.SetTimeout(h.config.Timeout)
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)
	t := time.Now()
	err := fasthttp.Do(req, resp)
	if err != nil {
		return Status{
			IsAlive: false,
			Error:   err,
			Latency: time.Since(t),
			Type:    HTTP,
		}
	}
	if resp.StatusCode() >= fasthttp.StatusOK && resp.StatusCode() < fasthttp.StatusBadRequest {
		return Status{
			IsAlive: true,
			Latency: time.Since(t),
			Type:    HTTP,
			Error:   nil,
		}
	}
	return Status{
		IsAlive: false,
		Error:   fmt.Errorf("not connected"),
		Latency: time.Since(t),
		Type:    HTTP,
	}
}
