package summary

import (
	"math"
	"sync"
	"ysab/conf"
	"ysab/tools"
)

var (
	AnalysisData  sync.Map
	ResChanel     = make(chan Res, 50000)
	RunOverSignal = make(chan int, 1)
	summaryData   = SummaryData{
		CodeDetail: make(map[int]int),
		QpsDetail:  make(map[int]int),
		MinConn:    float64(config.TimeOut),
		MinDNS:     float64(config.TimeOut),
		MinDelay:   float64(config.TimeOut),
		MinReq:     float64(config.TimeOut),
		MinUseTime: float64(config.TimeOut),
		MinRes:     float64(config.TimeOut),
	}
	tmpSummaryData = TmpSummaryData{}
	config         = conf.Conf
)

type Res struct {
	Size         int
	TimeStamp    int
	TotalUseTime float64
	Code         int
	ConnTime     float64
	DNSTime      float64
	ReqTime      float64
	DelayTime    float64
	ResTime      float64
}

type TmpSummaryData struct {
	Tlt50mNum  int
	Tlt100mNum int
	Tlt300mNum int
	Tlt500mNum int
	Tgt500mNum int
}

type SummaryData struct {
	CompleteRequests int
	FailedRequests   int
	TotalDataSize    int
	AvgDataSize      int
	QpsDetail        map[int]int
	RequestsPerSec   float64

	Tlt50mPercent  string
	Tlt100mPercent string
	Tlt300mPercent string
	Tlt500mPercent string
	Tgt500mPercent string

	MinUseTime float64
	MaxUseTime float64
	AvgUseTime float64
	CodeDetail map[int]int

	AvgConn  float64
	MaxConn  float64
	MinConn  float64
	AvgDNS   float64
	MaxDNS   float64
	MinDNS   float64
	AvgReq   float64
	MaxReq   float64
	MinReq   float64
	AvgDelay float64
	MaxDelay float64
	MinDelay float64
	AvgRes   float64
	MaxRes   float64
	MinRes   float64
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
		if _, ok := summaryData.CodeDetail[code]; ok {
			summaryData.CodeDetail[code]++
		} else {
			summaryData.CodeDetail[code] = 1
		}
		if _, ok := summaryData.QpsDetail[res.TimeStamp]; ok {
			summaryData.QpsDetail[res.TimeStamp]++
		} else {
			summaryData.QpsDetail[res.TimeStamp] = 1
		}
		if code != 200 {
			summaryData.FailedRequests++
		}
		if res.TotalUseTime <= 50 {
			tmpSummaryData.Tlt50mNum++
		}
		if res.TotalUseTime <= 100 {
			tmpSummaryData.Tlt100mNum++
		}
		if res.TotalUseTime <= 300 {
			tmpSummaryData.Tlt300mNum++
		}
		if res.TotalUseTime <= 500 {
			tmpSummaryData.Tlt500mNum++
		} else {
			tmpSummaryData.Tgt500mNum++
		}

		summaryData.AvgUseTime += res.TotalUseTime
		summaryData.AvgConn += res.ConnTime
		summaryData.AvgDNS += res.DNSTime
		summaryData.AvgDelay += res.DelayTime
		summaryData.AvgReq += res.ReqTime
		summaryData.AvgRes += res.ResTime

		summaryData.MinUseTime = math.Min(res.TotalUseTime, summaryData.MinUseTime)
		summaryData.MinConn = math.Min(res.ConnTime, summaryData.MinConn)
		summaryData.MinDNS = math.Min(res.DNSTime, summaryData.MinDNS)
		summaryData.MinDelay = math.Min(res.DelayTime, summaryData.MinDelay)
		summaryData.MinReq = math.Min(res.ReqTime, summaryData.MinReq)
		summaryData.MinRes = math.Min(res.ResTime, summaryData.MinRes)

		summaryData.MaxUseTime = math.Max(res.TotalUseTime, summaryData.MaxUseTime)
		summaryData.MaxConn = math.Max(res.ConnTime, summaryData.MaxConn)
		summaryData.MaxDNS = math.Max(res.DNSTime, summaryData.MaxDNS)
		summaryData.MaxDelay = math.Max(res.DelayTime, summaryData.MaxDelay)
		summaryData.MaxReq = math.Max(res.ReqTime, summaryData.MaxReq)
		summaryData.MaxRes = math.Max(res.ResTime, summaryData.MaxRes)

	}
	summaryData.Tlt50mPercent = tools.FloatToPercent(
		float64(tmpSummaryData.Tlt50mNum) / float64(config.UrlNum))
	summaryData.Tlt100mPercent = tools.FloatToPercent(
		float64(tmpSummaryData.Tlt100mNum) / float64(config.UrlNum))
	summaryData.Tlt300mPercent = tools.FloatToPercent(
		float64(tmpSummaryData.Tlt300mNum) / float64(config.UrlNum))

	summaryData.Tlt500mPercent = tools.FloatToPercent(
		float64(tmpSummaryData.Tlt500mNum) / float64(config.UrlNum))

	summaryData.Tgt500mPercent = tools.FloatToPercent(
		float64(tmpSummaryData.Tgt500mNum) / float64(config.UrlNum))

	summaryData.AvgUseTime = tools.Decimal2(summaryData.AvgUseTime / float64(config.UrlNum))
	summaryData.AvgConn = tools.Decimal2(summaryData.AvgConn / float64(config.UrlNum))
	summaryData.AvgDNS = tools.Decimal2(summaryData.AvgDNS / float64(config.UrlNum))
	summaryData.AvgDelay = tools.Decimal2(summaryData.AvgDelay / float64(config.UrlNum))
	summaryData.AvgReq = tools.Decimal2(summaryData.AvgReq / float64(config.UrlNum))
	summaryData.AvgRes = tools.Decimal2(summaryData.AvgRes / float64(config.UrlNum))
	summaryData.AvgDataSize = summaryData.TotalDataSize / config.UrlNum

	qpsL := len(summaryData.QpsDetail)
	min_v := 1000000000.0
	max_v := 0.0
	for _, v := range summaryData.QpsDetail {
		min_v = math.Min(float64(v), min_v)
		max_v = math.Max(float64(v), max_v)
		summaryData.RequestsPerSec += float64(v)
	}
	if qpsL >= 3 {
		summaryData.RequestsPerSec = summaryData.RequestsPerSec - (min_v + max_v)
		qpsL = qpsL - 2
	}
	summaryData.RequestsPerSec = summaryData.RequestsPerSec / float64(qpsL)

	Print(summaryData)
}
