package http

import (
	"bytes"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptrace"
	netulr "net/url"
	"time"
	"ysab/conf"
	"ysab/summary"
	"ysab/tools"
)

var (
	httpClient *http.Client
	config     = conf.Conf
)

const (
	MaxIdleConnections int = 1000
)

func init() {
	httpClient = creteHttpClient()
}

func creteHttpClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: MaxIdleConnections,
			DisableKeepAlives:   false,
			DisableCompression:  false,
		},
		Timeout: time.Duration(config.TimeOut) * time.Second,
	}
	return client
}

func do(url string, method string, headers map[string]string, bodydata string) summary.Res {
	var code int
	var size, tmpt int64
	var dnsStart, connStart, respStart, reqStart, delayStart int64
	var dnsDuration, connDuration, respDuration, reqDuration, delayDuration int64

	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(bodydata)))
	if err != nil {
		return summary.Res{}
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	trace := &httptrace.ClientTrace{
		DNSStart: func(info httptrace.DNSStartInfo) {
			dnsStart = tools.GetNowUnixNano()
		},
		DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
			dnsDuration = tools.GetNowUnixNano() - dnsStart
		},
		GetConn: func(h string) {
			connStart = tools.GetNowUnixNano()
		},
		GotConn: func(connInfo httptrace.GotConnInfo) {
			tmpt = tools.GetNowUnixNano()
			if !connInfo.Reused {
				connDuration = tmpt - connStart
			}
			reqStart = tmpt
		},
		WroteRequest: func(w httptrace.WroteRequestInfo) {
			tmpt = tools.GetNowUnixNano()
			reqDuration = tmpt - reqStart
			delayStart = tmpt
		},
		GotFirstResponseByte: func() {
			tmpt = tools.GetNowUnixNano()
			delayDuration = tmpt - delayStart
			respStart = tmpt
		},
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	tStart := tools.GetNowUnixNano()
	response, err := httpClient.Do(req)
	tEnd := tools.Now()
	if err == nil {
		if response.ContentLength > -1 {
			size = response.ContentLength
		} else {
			size = 0
		}
		code = response.StatusCode
		io.Copy(ioutil.Discard, response.Body)
		response.Body.Close()
	} else {
		if err, ok := err.(*netulr.Error); ok {
			if err2, ok := err.Err.(net.Error); ok {
				if err2.Timeout() {
					code = 504
				}
			}
		} else {
			code = 500
		}
	}
	respDuration = tEnd.UnixNano() - respStart
	return summary.Res{
		Size:         int(size),
		TimeStamp:    int(tEnd.UnixNano()),
		TotalUseTime: float64((tEnd.UnixNano() - tStart) / 10e5),
		Code:         code,
		ConnTime:     float64(connDuration / 10e5),
		DNSTime:      float64(dnsDuration / 10e5),
		ReqTime:      float64(reqDuration / 10e5),
		DelayTime:    float64(delayDuration / 10e5),
		ResTime:      float64(respDuration / 10e5),
	}

}

func Get(url string, headers map[string]string, data string) summary.Res {
	return do(url, "GET", headers, data)
}

func Post(url string, headers map[string]string, data string) summary.Res {
	return do(url, "POST", headers, data)
}
func Put(url string, headers map[string]string, data string) summary.Res {
	return do(url, "PUT", headers, data)
}
func Delete(url string, headers map[string]string, data string) summary.Res {
	return do(url, "DELETE", headers, data)
}
