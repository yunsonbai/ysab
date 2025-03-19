package http

import (
	"bufio"
	"bytes"
	"io"
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
)

func init() {
	for i := 0; i < clientsN; i++ {
		HttpClients = append(HttpClients, creteHttpClient())
	}
}

func creteHttpClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			MaxConnsPerHost:     int(conf.Conf.N)/clientsN + 128,
			MaxIdleConnsPerHost: int(conf.Conf.N)/clientsN + 32,
			DisableKeepAlives:   false,
			DisableCompression:  false,
		},
		Timeout: time.Duration(conf.Conf.TimeOut) * time.Second,
	}
	return client
}

func do(req *http.Request, bufioReader *bufio.Reader) (sumRes summary.ResStruct) {
	var tmpt int64
	var dnsStart, connStart, respStart, reqStart, delayStart int64
	trace := &httptrace.ClientTrace{
		DNSStart: func(info httptrace.DNSStartInfo) {
			dnsStart = time.Now().UnixMicro()
		},
		DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
			sumRes.DNSTime = time.Now().UnixMicro() - dnsStart
		},
		GetConn: func(h string) {
			connStart = time.Now().UnixMicro()
		},
		GotConn: func(connInfo httptrace.GotConnInfo) {
			tmpt = time.Now().UnixMicro()
			if !connInfo.Reused {
				if connStart <= 0 {
					sumRes.ConnTime = 0
				} else {
					sumRes.ConnTime = tmpt - connStart
				}
			}
			reqStart = tmpt
		},

		WroteRequest: func(w httptrace.WroteRequestInfo) {
			tmpt = time.Now().UnixMicro()
			sumRes.ReqTime = tmpt - reqStart
			delayStart = tmpt
		},
		GotFirstResponseByte: func() {
			tmpt = time.Now().UnixMicro()
			sumRes.DelayTime = tmpt - delayStart
			respStart = tmpt
		},
	}
	newReq := req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	tStart := time.Now().UnixMicro()

	client := HttpClients[rand.Intn(clientsN)]
	response, err := client.Do(newReq)
	tEnd := time.Now()
	if response != nil {

		sumRes.Code = response.StatusCode
		defer response.Body.Close()
		bufioReader.Reset(response.Body)
		_, err = io.Copy(io.Discard, bufioReader)
		if err != nil {
			sumRes.Code = 500
		}
		tEnd = time.Now()
		if response.ContentLength > -1 {
			sumRes.Size = response.ContentLength
		} else {
			sumRes.Size = 0
		}

	} else {
		sumRes.Code = 503
		if err, ok := err.(*netulr.Error); ok {
			if err.Timeout() {
				sumRes.Code = 504
			}
		}
	}
	sumRes.TimeStamp = tEnd.UnixMicro()
	sumRes.TotalUseTime = tEnd.UnixMicro() - tStart
	sumRes.ResTime = tEnd.UnixMicro() - respStart

	return
}

func GetReq(url, method, bodydata string, headers map[string]string) (req *http.Request) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(bodydata)))
	if err != nil {
		return
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
	return
}

func Head(req *http.Request, url, data string, headers map[string]string, bufioReader *bufio.Reader) summary.ResStruct {
	if req == nil {
		req = GetReq(url, "HEAD", data, headers)
	}
	return do(req, bufioReader)
}

func Get(req *http.Request, url, data string, headers map[string]string, bufioReader *bufio.Reader) summary.ResStruct {
	if req == nil {
		req = GetReq(url, "GET", data, headers)
	}
	return do(req, bufioReader)
}
func Post(req *http.Request, url, data string, headers map[string]string, bufioReader *bufio.Reader) summary.ResStruct {
	if req == nil {
		req = GetReq(url, "POST", data, headers)
	}
	return do(req, bufioReader)
}
func Put(req *http.Request, url, data string, headers map[string]string, bufioReader *bufio.Reader) summary.ResStruct {
	if req == nil {
		req = GetReq(url, "PUT", data, headers)
	}
	return do(req, bufioReader)
}
func Delete(req *http.Request, url, data string, headers map[string]string, bufioReader *bufio.Reader) summary.ResStruct {
	if req == nil {
		req = GetReq(url, "DELETE", data, headers)
	}
	return do(req, bufioReader)
}
