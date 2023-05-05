package worker

import (
	"bufio"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/yunsonbai/ysab/conf"
	yshttp "github.com/yunsonbai/ysab/http"
	"github.com/yunsonbai/ysab/summary"
	ystools "github.com/yunsonbai/ysab/tools"
)

var (
	rwg       sync.WaitGroup
	config    = conf.Conf
	urlChanel = make(chan [2]string, 30000)
)

func worker(method string) {
	wf := yshttp.Get
	switch method {
	case "GET":
		wf = yshttp.Get
	case "POST":
		wf = yshttp.Post
	case "PUT":
		wf = yshttp.Put
	case "DELETE":
		wf = yshttp.Delete
	case "HEAD":
		wf = yshttp.Head
	default:
		return
	}
	readBuf := make([]byte, 32*1024)
	var req *http.Request
	if config.UrlFilePath == "" {
		req = yshttp.GetReq(config.Url, method, config.Body, config.Headers)
	}
	for {
		data, ok := <-urlChanel
		if !ok {
			return
		}
		summary.ResChanel <- wf(req, data[0], data[1], config.Headers, readBuf)
	}
}

func addTask() {
	i := 0
	url := config.Url
	body := config.Body
	var fbr *bufio.Reader
	if config.UrlFilePath != "" {
		fi, _ := os.Open(config.UrlFilePath)
		defer fi.Close()
		fbr = bufio.NewReader(fi)
	}
	done := 0
	for {
		i++
		if config.UrlFilePath != "" {
			line, _, err := fbr.ReadLine()
			if err == io.EOF {
				done = 1
			}
			reqData := ystools.GetReqData(string(line))
			url = reqData.Url
			body = reqData.Body
		} else {
			if i == int(config.UrlNum) {
				done = 1
			}
		}
		if url != "" {
			data := [2]string{url, body}
			urlChanel <- data
		}
		if done == 1 {
			break
		}
	}
	if done == 1 {
		for {
			time.Sleep(time.Duration(50) * time.Millisecond)
			if len(urlChanel) == 0 {
				close(urlChanel)
				return
			}
		}
	}
}

func StartWork() {
	rwg.Add(1)
	go addTask()
	go func() {
		summary.HandleRes()
		rwg.Done()
	}()
	for index := 0; index < int(config.N); index++ {
		go func() {
			worker(config.Method)
		}()
	}
	rwg.Wait()
}
