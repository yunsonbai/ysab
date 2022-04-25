package worker

import (
	"bufio"
	"io"
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

func worker(method string, fast int) {
	for {
		data, ok := <-urlChanel
		if !ok {
			return
		}
		summary.ResChanel <- yshttp.Request(
			data[0], data[1], method, config.Headers, fast)
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
			url, body = ystools.GetReqDataNew(string(line))
			// reqData := ystools.GetReqData(string(line))
			// url = reqData.Url
			// body = reqData.Body
		} else {
			if i == config.UrlNum {
				done = 1
			}
		}
		if url != "" {
			urlChanel <- [2]string{url, body}
		}
		if done == 1 {
			break
		}
	}
	if done == 1 {
		for {
			time.Sleep(time.Duration(200) * time.Millisecond)
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
	for index := 0; index < config.N; index++ {
		go func() {
			worker(config.Method, config.Fast)
		}()
	}
	rwg.Wait()
}
