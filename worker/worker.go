package worker

import (
	"bufio"
	"io"
	"os"
	"sync"
	"time"
	"ysab/conf"
	yshttp "ysab/http"
	"ysab/summary"
	ystools "ysab/tools"
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
		wf = yshttp.Post
	case "DELETE":
		wf = yshttp.Post
	case "HEAD":
		wf = yshttp.Head
	default:
		return
	}
	for {
		data, ok := <-urlChanel
		if !ok {
			return
		}
		summary.ResChanel <- wf(data[0], config.Headers, data[1])
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
			if i == config.UrlNum {
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
		go func() {
			for {
				time.Sleep(time.Duration(10) * time.Millisecond)
				if len(urlChanel) == 0 {
					close(urlChanel)
					return
				}
			}
		}()
	}
}

func StartWork() {
	go addTask()
	rwg.Add(config.N + 1)
	go func() {
		summary.HandleRes()
		rwg.Done()
	}()
	n := config.N
	for index := 0; index < n; index++ {
		go func() {
			worker(config.Method)
			rwg.Done()
		}()
	}
	rwg.Wait()
}
