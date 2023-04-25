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

	"github.com/yunsonbai/ysab/conf"
	"github.com/yunsonbai/ysab/summary"
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
			MaxConnsPerHost:     int(config.N)/clientsN + 128,
			MaxIdleConnsPerHost: int(config.N)/clientsN + 32,
			DisableKeepAlives:   false,
			DisableCompression:  false,
		},
		Timeout: time.Duration(config.TimeOut) * time.Second,
	}
	return client
}

func do(url, method, bodydata string, headers map[string]string, buf []byte) summary.ResStruct {
	var code int
	var size, tmpt int64
	var dnsStart, connStart, respStart, reqStart, delayStart int64
	var dnsDuration, connDuration, respDuration, reqDuration, delayDuration int64
	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(bodydata)))
	if err != nil {
		return summary.ResStruct{}
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
			dnsStart = time.Now().UnixMicro()
		},
		DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
			dnsDuration = time.Now().UnixMicro() - dnsStart
		},
		GetConn: func(h string) {
			connStart = time.Now().UnixMicro()
		},
		GotConn: func(connInfo httptrace.GotConnInfo) {
			tmpt = time.Now().UnixMicro()
			if !connInfo.Reused {
				if connStart <= 0 {
					connDuration = 0
				} else {
					connDuration = tmpt - connStart
				}
			}
			reqStart = tmpt
		},
		WroteRequest: func(w httptrace.WroteRequestInfo) {
			tmpt = time.Now().UnixMicro()
			reqDuration = tmpt - reqStart
			delayStart = tmpt
		},
		GotFirstResponseByte: func() {
			tmpt = time.Now().UnixMicro()
			delayDuration = tmpt - delayStart
			respStart = tmpt
		},
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	tStart := time.Now().UnixMicro()

	client := HttpClients[rand.Intn(clientsN)]
	response, err := client.Do(req)
	tEnd := time.Now()

	if response != nil {
		if response.ContentLength > -1 {
			size = response.ContentLength
		} else {
			size = 0
		}
		code = response.StatusCode
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

	respDuration = tEnd.UnixMicro() - respStart

	return summary.ResStruct{
		Size:         size,
		TimeStamp:    tEnd.UnixMicro(),
		TotalUseTime: tEnd.UnixMicro() - tStart,
		Code:         code,
		ConnTime:     connDuration,
		DNSTime:      dnsDuration,
		ReqTime:      reqDuration,
		DelayTime:    delayDuration,
		ResTime:      respDuration}

}

func Head(url, data string, headers map[string]string, readBuf []byte) summary.ResStruct {
	return do(url, "HEAD", data, headers, readBuf)
}

func Get(url, data string, headers map[string]string, readBuf []byte) summary.ResStruct {
	return do(url, "GET", data, headers, readBuf)
}
func Post(url, data string, headers map[string]string, readBuf []byte) summary.ResStruct {
	return do(url, "POST", data, headers, readBuf)
}
func Put(url, data string, headers map[string]string, readBuf []byte) summary.ResStruct {
	return do(url, "PUT", data, headers, readBuf)
}
func Delete(url, data string, headers map[string]string, readBuf []byte) summary.ResStruct {
	return do(url, "DELETE", data, headers, readBuf)
}
