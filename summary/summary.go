package summary

import (
	"sort"
	"strconv"
	"sync"

	"github.com/yunsonbai/ysab/conf"
	"github.com/yunsonbai/ysab/tools"
)

var (
	AnalysisData  sync.Map
	ResChanel     = make(chan Res, 20000)
	RunOverSignal = make(chan int, 1)

	codeDetail  = make(map[int]int)
	summaryData = SummaryData{
		CodeDetail:        make(map[string]int),
		WaitingTimeDetail: make(map[string]int),
		MinConn:           int64(config.TimeOut),
		MinDNS:            int64(config.TimeOut),
		MinDelay:          int64(config.TimeOut),
		MinReq:            int64(config.TimeOut),
		MinUseTime:        int64(config.TimeOut),
		MinRes:            int64(config.TimeOut),
	}
	config    = conf.Conf
	waitTimes = make([]int, 0, config.UrlNum)
	B         = int64(1000000)
)

type Res struct {
	Size         int64
	TimeStamp    int64
	TotalUseTime int64
	Code         int
	ConnTime     int64
	DNSTime      int64
	ReqTime      int64
	DelayTime    int64
	ResTime      int64
}

type SummaryData struct {
	CompleteRequests      int
	FailedRequests        int
	SuccessRequests       int
	TimeTaken             float64
	TotalDataSize         int64
	AvgDataSize           int64
	RequestsPerSec        float64
	SuccessRequestsPerSec float64

	MinUseTime        int64
	MaxUseTime        int64
	AvgUseTime        float64
	totolTime         int64
	CodeDetail        map[string]int
	WaitingTimeDetail map[string]int

	AvgConn       float64
	totalAvgConn  int64
	MaxConn       int64
	MinConn       int64
	AvgDNS        float64
	totalAvgDNS   int64
	MaxDNS        int64
	MinDNS        int64
	AvgReq        float64
	totalAvgReq   int64
	MaxReq        int64
	MinReq        int64
	AvgDelay      float64
	totalAvgDelay int64
	MaxDelay      int64
	MinDelay      int64
	AvgRes        float64
	totalAvgRes   int64
	MaxRes        int64
	MinRes        int64
}

func HandleRes() {
	for {
		res, ok := <-ResChanel
		if !ok {
			break
		}
		summaryData.CompleteRequests++
		summaryData.TotalDataSize += res.Size
		if summaryData.CompleteRequests == config.UrlNum {
			close(ResChanel)
		}

		code := res.Code
		if _, ok := codeDetail[code]; ok {
			codeDetail[code]++
		} else {
			codeDetail[code] = 1
		}
		if config.EndTime < res.TimeStamp {
			config.EndTime = res.TimeStamp
		}
		if code > 299 || code < 200 {
			summaryData.FailedRequests++
		} else {
			summaryData.SuccessRequests++
		}
		summaryData.totolTime += res.TotalUseTime
		summaryData.totalAvgConn += res.ConnTime
		summaryData.totalAvgDNS += res.DNSTime
		summaryData.totalAvgDelay += res.DelayTime
		summaryData.totalAvgReq += res.ReqTime
		summaryData.totalAvgRes += res.ResTime

		summaryData.MinUseTime = tools.MinInt64(res.TotalUseTime, summaryData.MinUseTime)
		summaryData.MinConn = tools.MinInt64(res.ConnTime, summaryData.MinConn)
		summaryData.MinDNS = tools.MinInt64(res.DNSTime, summaryData.MinDNS)
		summaryData.MinDelay = tools.MinInt64(res.DelayTime, summaryData.MinDelay)
		summaryData.MinReq = tools.MinInt64(res.ReqTime, summaryData.MinReq)
		summaryData.MinRes = tools.MinInt64(res.ResTime, summaryData.MinRes)

		summaryData.MaxUseTime = tools.MaxInt64(res.TotalUseTime, summaryData.MaxUseTime)
		summaryData.MaxConn = tools.MaxInt64(res.ConnTime, summaryData.MaxConn)
		summaryData.MaxDNS = tools.MaxInt64(res.DNSTime, summaryData.MaxDNS)
		summaryData.MaxDelay = tools.MaxInt64(res.DelayTime, summaryData.MaxDelay)
		summaryData.MaxReq = tools.MaxInt64(res.ReqTime, summaryData.MaxReq)
		summaryData.MaxRes = tools.MaxInt64(res.ResTime, summaryData.MaxRes)
		waitTimes = append(waitTimes, int(res.TotalUseTime/B))
	}
	urlNum := int64(config.UrlNum)
	summaryData.AvgDataSize = summaryData.TotalDataSize / urlNum

	summaryData.AvgUseTime = tools.Decimal2(float64(summaryData.totolTime / urlNum / B))
	summaryData.AvgConn = tools.Decimal2(float64(summaryData.totalAvgConn / urlNum / B))
	summaryData.AvgDNS = tools.Decimal2(float64(summaryData.totalAvgDNS / urlNum / B))
	summaryData.AvgDelay = tools.Decimal2(float64(summaryData.totalAvgDelay / urlNum / B))
	summaryData.AvgReq = tools.Decimal2(float64(summaryData.totalAvgReq / urlNum / B))
	summaryData.AvgRes = tools.Decimal2(float64(summaryData.totalAvgRes / urlNum / B))

	summaryData.MaxConn = summaryData.MaxConn / B
	summaryData.MaxDNS = summaryData.MaxDNS / B
	summaryData.MaxDelay = summaryData.MaxDelay / B
	summaryData.MaxReq = summaryData.MaxReq / B
	summaryData.MaxRes = summaryData.MaxRes / B
	summaryData.MaxUseTime = summaryData.MaxUseTime / B
	summaryData.MinConn = summaryData.MinConn / B
	summaryData.MinDNS = summaryData.MinDNS / B
	summaryData.MinDelay = summaryData.MinDelay / B
	summaryData.MinReq = summaryData.MinReq / B
	summaryData.MinRes = summaryData.MinRes / B
	summaryData.MinUseTime = summaryData.MinUseTime / B

	for k, v := range codeDetail {
		summaryData.CodeDetail[strconv.Itoa(k)] = v
	}

	t := (float64(config.EndTime-config.StartTime) / 10e8)
	summaryData.TimeTaken = t
	summaryData.RequestsPerSec = float64(config.UrlNum) / t
	summaryData.SuccessRequestsPerSec = float64(summaryData.SuccessRequests) / t
	sort.Ints(waitTimes)
	waitTimesL := len(waitTimes)
	tps := []float64{0.1, 0.25, 0.5, 0.75, 0.9, 0.95, 0.99, 0.999, 0.9999}
	tpsL := len(tps)
	for i := 0; i < tpsL; i++ {
		summaryData.WaitingTimeDetail[tools.FloatToPercent(
			tps[i])] = int(waitTimes[int(float64(waitTimesL)*tps[i]-1)])
	}
	Print(summaryData)
}
