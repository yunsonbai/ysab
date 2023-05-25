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
	ResChanel     = make(chan ResStruct, 4000)
	RunOverSignal = make(chan int, 1)

	codeDetail     = make(map[int]int)
	summaryDataTmp = summaryDataTmpStruct{
		MinConn:    config.TimeOut,
		MinDNS:     config.TimeOut,
		MinDelay:   config.TimeOut,
		MinReq:     config.TimeOut,
		MinUseTime: config.TimeOut,
		MinRes:     config.TimeOut,
	}
	config    = conf.Conf
	waitTimes = make(map[int64]int)
)

type ResStruct struct {
	Size         int64
	TimeStamp    int64
	Code         int
	TotalUseTime int64
	ConnTime     int64
	DNSTime      int64
	ReqTime      int64
	DelayTime    int64
	ResTime      int64
}

type summaryDataTmpStruct struct {
	CompleteRequests uint32
	FailedRequests   uint32
	SuccessRequests  uint32
	TotalDataSize    int64

	MinUseTime int64 // 微妙级
	MaxUseTime int64 // 微妙级
	AvgUseTime int64 // 微妙级

	AvgConn  int64 // 微妙级
	MaxConn  int64 // 微妙级
	MinConn  int64 // 微妙级
	AvgDNS   int64 // 微妙级
	MaxDNS   int64 // 微妙级
	MinDNS   int64 // 微妙级
	AvgReq   int64 // 微妙级
	MaxReq   int64 // 微妙级
	MinReq   int64 // 微妙级
	AvgDelay int64 // 微妙级
	MaxDelay int64 // 微妙级
	MinDelay int64 // 微妙级
	AvgRes   int64 // 微妙级
	MaxRes   int64 // 微妙级
	MinRes   int64 // 微妙级
}

type SummaryDataStruct struct {
	CompleteRequests      uint32
	FailedRequests        uint32
	SuccessRequests       uint32
	TimeToken             string
	TotalDataSize         int64
	AvgDataSize           string
	RequestsPerSec        string
	SuccessRequestsPerSec string
	STransferRatePerSec   string

	MinUseTime        string
	MaxUseTime        string
	AvgUseTime        string
	CodeDetail        map[string]int
	WaitingTimeDetail map[string]int

	AvgConn  string
	MaxConn  string
	MinConn  string
	AvgDNS   string
	MaxDNS   string
	MinDNS   string
	AvgReq   string
	MaxReq   string
	MinReq   string
	AvgDelay string
	MaxDelay string
	MinDelay string
	AvgRes   string
	MaxRes   string
	MinRes   string
}

func microToSecond(t int64) float64 {
	return float64(t) / float64(config.TimeBase)
}

func microToMilli(t int64) int64 {
	return 1000 * t / config.TimeBase
}

func HandleRes() {
	var tkey int64
	var code int
	for {
		res, ok := <-ResChanel
		if !ok {
			break
		}
		summaryDataTmp.CompleteRequests++
		summaryDataTmp.TotalDataSize += res.Size
		if summaryDataTmp.CompleteRequests == config.UrlNum {
			close(ResChanel)
		}
		code = res.Code
		codeDetail[code]++
		if config.EndTime < res.TimeStamp {
			config.EndTime = res.TimeStamp
		}
		if code > 299 || code < 200 {
			summaryDataTmp.FailedRequests++
		} else {
			summaryDataTmp.SuccessRequests++
		}
		summaryDataTmp.AvgUseTime += res.TotalUseTime
		summaryDataTmp.AvgConn += res.ConnTime
		summaryDataTmp.AvgDNS += res.DNSTime
		summaryDataTmp.AvgDelay += res.DelayTime
		summaryDataTmp.AvgReq += res.ReqTime
		summaryDataTmp.AvgRes += res.ResTime

		summaryDataTmp.MinUseTime = tools.MinInt64(res.TotalUseTime, summaryDataTmp.MinUseTime)
		summaryDataTmp.MinConn = tools.MinInt64(res.ConnTime, summaryDataTmp.MinConn)
		summaryDataTmp.MinDNS = tools.MinInt64(res.DNSTime, summaryDataTmp.MinDNS)
		summaryDataTmp.MinDelay = tools.MinInt64(res.DelayTime, summaryDataTmp.MinDelay)
		summaryDataTmp.MinReq = tools.MinInt64(res.ReqTime, summaryDataTmp.MinReq)
		summaryDataTmp.MinRes = tools.MinInt64(res.ResTime, summaryDataTmp.MinRes)

		summaryDataTmp.MaxUseTime = tools.MaxInt64(res.TotalUseTime, summaryDataTmp.MaxUseTime)
		summaryDataTmp.MaxConn = tools.MaxInt64(res.ConnTime, summaryDataTmp.MaxConn)
		summaryDataTmp.MaxDNS = tools.MaxInt64(res.DNSTime, summaryDataTmp.MaxDNS)
		summaryDataTmp.MaxDelay = tools.MaxInt64(res.DelayTime, summaryDataTmp.MaxDelay)
		summaryDataTmp.MaxReq = tools.MaxInt64(res.ReqTime, summaryDataTmp.MaxReq)
		summaryDataTmp.MaxRes = tools.MaxInt64(res.ResTime, summaryDataTmp.MaxRes)

		tkey = microToMilli(res.TotalUseTime)
		// if tkey > 3000 {
		// 	tkey = tkey / 1000 * 1000
		// }
		waitTimes[tkey]++
	}

	summaryData := SummaryDataStruct{
		CompleteRequests:  summaryDataTmp.CompleteRequests,
		FailedRequests:    summaryDataTmp.FailedRequests,
		SuccessRequests:   summaryDataTmp.SuccessRequests,
		TotalDataSize:     summaryDataTmp.TotalDataSize,
		MinUseTime:        tools.FloatToStr3f(microToSecond(summaryDataTmp.MinUseTime)),
		MaxUseTime:        tools.FloatToStr3f(microToSecond(summaryDataTmp.MaxUseTime)),
		CodeDetail:        make(map[string]int),
		WaitingTimeDetail: make(map[string]int),

		MaxConn:  tools.FloatToStr3f(microToSecond(summaryDataTmp.MaxConn)),
		MinConn:  tools.FloatToStr3f(microToSecond(summaryDataTmp.MinConn)),
		MaxDNS:   tools.FloatToStr3f(microToSecond(summaryDataTmp.MaxDNS)),
		MinDNS:   tools.FloatToStr3f(microToSecond(summaryDataTmp.MinDNS)),
		MaxReq:   tools.FloatToStr3f(microToSecond(summaryDataTmp.MaxReq)),
		MinReq:   tools.FloatToStr3f(microToSecond(summaryDataTmp.MinReq)),
		MaxDelay: tools.FloatToStr3f(microToSecond(summaryDataTmp.MaxDelay)),
		MinDelay: tools.FloatToStr3f(microToSecond(summaryDataTmp.MinDelay)),
		MaxRes:   tools.FloatToStr3f(microToSecond(summaryDataTmp.MaxRes)),
		MinRes:   tools.FloatToStr3f(microToSecond(summaryDataTmp.MinRes))}

	summaryData.AvgUseTime = tools.FloatToStr3f(microToSecond(summaryDataTmp.AvgUseTime) / float64(config.UrlNum))
	summaryData.AvgConn = tools.FloatToStr3f(microToSecond(summaryDataTmp.AvgConn) / float64(config.UrlNum))
	summaryData.AvgDNS = tools.FloatToStr3f(microToSecond(summaryDataTmp.AvgDNS) / float64(config.UrlNum))
	summaryData.AvgDelay = tools.FloatToStr3f(microToSecond(summaryDataTmp.AvgDelay) / float64(config.UrlNum))
	summaryData.AvgReq = tools.FloatToStr3f(microToSecond(summaryDataTmp.AvgReq) / float64(config.UrlNum))
	summaryData.AvgRes = tools.FloatToStr3f(microToSecond(summaryDataTmp.AvgRes) / float64(config.UrlNum))
	summaryData.AvgDataSize = tools.FloatToStr3f(float64(summaryDataTmp.TotalDataSize) / float64(config.UrlNum))

	for k, v := range codeDetail {
		summaryData.CodeDetail[strconv.Itoa(k)] = v
	}
	t := microToSecond(config.EndTime - config.StartTime)
	summaryData.TimeToken = tools.FloatToStr3f(t)
	summaryData.RequestsPerSec = tools.FloatToStr3f(float64(config.UrlNum) / t)
	summaryData.SuccessRequestsPerSec = tools.FloatToStr3f(float64(summaryData.SuccessRequests) / t)
	summaryData.STransferRatePerSec = tools.FloatToStr3f(float64(summaryData.TotalDataSize) / t)

	tps := []float64{0.1, 0.25, 0.5, 0.75, 0.9, 0.95, 0.99, 0.999, 0.9999}
	tpsL := len(tps)
	tpsCount := make([]int, tpsL)
	tkeys := []int{}
	for i, v := range tps {
		tpsCount[i] = int(float64(config.UrlNum) * v)
	}
	for k, _ := range waitTimes {
		tkeys = append(tkeys, int(k))
	}
	tkeysL := len(tkeys)
	sort.Ints(tkeys)
	tmpN := 0
	j := 0
	for i := 0; i < tpsL; i++ {
		for {
			if j >= tkeysL {
				if _, ok := summaryData.WaitingTimeDetail[tools.FloatToPercent(tps[i])]; !ok {
					summaryData.WaitingTimeDetail[tools.FloatToPercent(tps[i])] = tkeys[tkeysL-1]
				}
				break
			}
			if tmpN >= tpsCount[i] {
				summaryData.WaitingTimeDetail[tools.FloatToPercent(tps[i])] = tkeys[j]
				break
			}
			tmpN = tmpN + waitTimes[int64(tkeys[j])]
			j++
		}
	}
	Print(summaryData)
}
