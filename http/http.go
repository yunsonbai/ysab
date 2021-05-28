package http

import (
	"bytes"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptrace"
	netulr "net/url"
	"time"
	"ysab/conf"
	"ysab/summary"
	"ysab/tools"
)

const (
	clientsN int = 2
)

var (
	HttpClients []*http.Client
	config      = conf.Conf
)

func init() {
	for i := 0; i < clientsN; i++ {
		HttpClients = append(HttpClients, creteHttpClient())
	}
}

func creteHttpClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			MaxConnsPerHost:     config.N/clientsN + 128,
			MaxIdleConnsPerHost: config.N/clientsN + 128,
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
	// req.Header.Set("Accept-Encoding", "gzip, deflate")

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

	client := HttpClients[rand.Intn(clientsN)]
	response, err := client.Do(req)

	tEnd := tools.Now()
	if response != nil {
		if response.ContentLength > -1 {
			size = response.ContentLength
		} else {
			size = 0
		}
		code = response.StatusCode
		bSize := 32 * 1024
		if int64(bSize) > size {
			if size < 1 {
				bSize = 1
			} else {
				bSize = int(size)
			}
		}
		buf := make([]byte, bSize)
		io.CopyBuffer(ioutil.Discard, response.Body, buf)
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

func Head(url, data string, headers map[string]string) summary.Res {
	return do(url, "HEAD", data, headers)
}

func Get(url, data string, headers map[string]string) summary.Res {
	return do(url, "GET", data, headers)
}
func Post(url, data string, headers map[string]string) summary.Res {
	return do(url, "POST", data, headers)
}
func Put(url, data string, headers map[string]string) summary.Res {
	return do(url, "PUT", data, headers)
}
func Delete(url, data string, headers map[string]string) summary.Res {
	return do(url, "DELETE", data, headers)
}
