package http

import (
	"strings"
	"time"

	"github.com/valyala/fasthttp"
	"github.com/yunsonbai/ysab/summary"
	"github.com/yunsonbai/ysab/tools"
)

var (
	fastHttpClient fasthttp.Client
)

func init() {
	fastHttpClient = creteFastHttpClient()
}

func creteFastHttpClient() fasthttp.Client {

	return fasthttp.Client{
		Name:                     "Ysab",
		NoDefaultUserAgentHeader: true,
		MaxConnsPerHost:          config.N + 8,
		MaxIdleConnDuration:      10 * time.Second,
		// ReadTimeout:              5 * time.Second,
		// WriteTimeout:             5 * time.Second,
		MaxConnWaitTimeout: time.Duration(config.TimeOut) * time.Second,
	}
}

func fastDo(url, method, data string, headers map[string]string) summary.Res {
	var code, size int

	startT := tools.GetNowUnixNano()
	req := fasthttp.AcquireRequest()

	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI(url)
	req.Header.SetMethod(method)
	req.Header.SetContentType("application/x-www-form-urlencoded")
	for k, v := range headers {
		lk := strings.ToLower(k)
		if lk == "host" {
			req.SetHost(v)
		} else if lk == "content-type" {
			req.Header.SetContentType(v)
		} else {
			req.Header.Set(k, v)
		}
	}
	requestBody := []byte(data)
	req.SetBody(requestBody)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err := fastHttpClient.Do(req, resp); err != nil {
		endT := tools.GetNowUnixNano()
		return summary.Res{
			Size:         2,
			TimeStamp:    endT,
			TotalUseTime: endT - startT,
			Code:         503,
			ConnTime:     0,
			DNSTime:      0,
			ReqTime:      0,
			DelayTime:    0,
			ResTime:      0,
		}
	}
	code = resp.StatusCode()
	if resp.Header.ContentLength() > -1 {
		size = resp.Header.ContentLength()
	} else {
		size = 0
	}
	resp.Body()
	endT := tools.GetNowUnixNano()
	return summary.Res{
		Size:         int64(size),
		TimeStamp:    endT,
		TotalUseTime: endT - startT,
		Code:         code,
		ConnTime:     0,
		DNSTime:      0,
		ReqTime:      0,
		DelayTime:    0,
		ResTime:      0,
	}
}
