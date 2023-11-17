package worker

import (
	"bufio"
	"fmt"
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
	rwg sync.WaitGroup
	// urlChanel = make(chan [2]string, 30000)
	urlChanel = make(chan [2]string, conf.Conf.UrlChNum)
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
	if conf.Conf.UrlFilePath == "" {
		req = yshttp.GetReq(conf.Conf.Url, method, conf.Conf.Body, conf.Conf.Headers)
	}
	for {
		data, ok := <-urlChanel
		if !ok {
			return
		}
		summary.ResChanel <- wf(req, data[0], data[1], conf.Conf.Headers, readBuf)
	}
}

func addTaskByFile(useDuration bool) {
	var url string
	var body string
	var curR uint32
	var count int64
	var over bool
	file, err := os.Open(conf.Conf.UrlFilePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	fbr := bufio.NewReader(file)
	endT := time.Now().Unix() + int64(conf.Conf.Duration)
	for {
		line, _, err := fbr.ReadLine()
		if err == io.EOF {
			if useDuration {
				if time.Now().Unix() >= endT {
					over = true
				}
			} else {
				curR++
				if curR >= conf.Conf.Round {
					over = true
				}
			}

			if over {
				conf.Conf.UrlNum = count
				summary.ResChanel <- summary.ResStruct{EndMk: true}
				fmt.Println("curR:", curR, "count:", count, "--body:", body,
					"conf.Conf.UrlNum:", conf.Conf.UrlNum)
				break
			}

			if _, err := file.Seek(0, io.SeekStart); err != nil {
				panic(err)
			}
			continue
		}

		reqData := ystools.GetReqData(string(line))
		url = reqData.Url
		body = reqData.Body
		if url != "" {
			data := [2]string{url, body}
			count++
			urlChanel <- data
		}
	}
}

func addTaskByCmd(useDuration bool) {
	totalUrlNum := int64(conf.Conf.Round) * int64(conf.Conf.N)
	var count int64
	var over bool
	endT := time.Now().Unix() + int64(conf.Conf.Duration)
	for {

		if useDuration {
			if time.Now().Unix() >= endT {
				over = true
			}
		} else {
			if count >= totalUrlNum {
				over = true
			}
		}

		if over {
			conf.Conf.UrlNum = count
			summary.ResChanel <- summary.ResStruct{EndMk: true}
			fmt.Println(
				"count:", count, "--body:", conf.Conf.Body,
				"conf.Conf.UrlNum:", conf.Conf.UrlNum)
			break
		}

		data := [2]string{conf.Conf.Url, conf.Conf.Body}
		count++
		urlChanel <- data
	}
}

func addTask() {
	var useDuration bool
	if conf.Conf.Duration > 0 {
		useDuration = true
	}
	if conf.Conf.UrlFilePath != "" {
		addTaskByFile(useDuration)
	} else {
		addTaskByCmd(useDuration)
	}
	for {
		time.Sleep(time.Duration(50) * time.Millisecond)
		if len(urlChanel) == 0 {
			close(urlChanel)
			return
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
	for index := 0; index < int(conf.Conf.N); index++ {
		go func() {
			worker(conf.Conf.Method)
		}()
	}
	rwg.Wait()
}
