package http

import (
	"bytes"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptrace"
	netulr "net/url"
	"strings"
	"time"

	"github.com/yunsonbai/ysab/summary"
	"github.com/yunsonbai/ysab/tools"
)

const (
	clientsN int = 2
)

var (
	HttpClients []*http.Client
)

func init() {
	for i := 0; i < clientsN; i++ {
		HttpClients = append(HttpClients, creteHttpClient())
	}
}

func creteHttpClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			MaxConnsPerHost:     config.N/clientsN + 8,
			MaxIdleConnsPerHost: config.N/clientsN + 8,
			DisableKeepAlives:   false,
			DisableCompression:  false,
		},
		Timeout: time.Duration(config.TimeOut) * time.Second,
	}
	return client
}

func do(url, method, bodydata string, headers map[string]string) summary.Res {
	var code int
	var size, tmpt int64
	var dnsStart, connStart, respStart, reqStart, delayStart int64
	var dnsDuration, connDuration, respDuration, reqDuration, delayDuration int64
	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(bodydata)))
	if err != nil {
		return summary.Res{}
	}
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	for k, v := range headers {
		if strings.ToLower(k) == "host" {
			req.Host = v
		} else {
			req.Header.Set(k, v)
		}
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

	// client := HttpClients[rand.Intn(clientsN)]
	response, err := HttpClients[rand.Intn(clientsN)].Do(req)

	tEnd := tools.Now()
	if response != nil {
		if response.ContentLength > -1 {
			size = response.ContentLength
		} else {
			size = 0
		}
		code = response.StatusCode
		bSize := 32768 // 32 * 1024
		if size < 32768 {
			if size < 1 {
				bSize = 1
			} else {
				bSize = int(size)
			}
		}
		io.CopyBuffer(ioutil.Discard, response.Body, make([]byte, bSize))
		response.Body.Close()
	} else {
		code = 503
		if err, ok := err.(*netulr.Error); ok {
			if err.Timeout() {
				code = 504
			}
		}
	}

	respDuration = tEnd.UnixNano() - respStart

	return summary.Res{
		Size:         size,
		TimeStamp:    tEnd.UnixNano(),
		TotalUseTime: (tEnd.UnixNano() - tStart),
		Code:         code,
		ConnTime:     connDuration,
		DNSTime:      dnsDuration,
		ReqTime:      reqDuration,
		DelayTime:    delayDuration,
		ResTime:      respDuration,
	}

}
