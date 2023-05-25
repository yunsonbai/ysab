package summary

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"text/template"
)

var (
	htmlTemplate = `
Summary:
  Complete requests:		{{ .CompleteRequests }}
  Failed requests:		{{ .FailedRequests }}
  Time taken (s):		{{ .TimeToken }}
  Total data size (Byte):	{{ .TotalDataSize }}
  Data size/request (Byte):	{{ .AvgDataSize }}
  Max use time (ms):		{{.MaxUseTime}}
  Min use time (ms):		{{.MinUseTime}}
  Average use time (ms):	{{.AvgUseTime}}
  Requests/sec:			{{ .RequestsPerSec }}
  SuccessRequests/sec:		{{ .SuccessRequestsPerSec }}
  TransferRate/sec (Byte):	{{ .STransferRatePerSec }}

Percentage of waiting time (ms):
{{ formatMap .WaitingTimeDetail }}

Time detail (ms)
  item		min		mean		max
  dns		{{.MinDNS}}		{{.AvgDNS}}		{{.MaxDNS}}
  conn		{{.MinConn}}		{{.AvgConn}}		{{.MaxConn}}
  wait		{{.MinDelay}}		{{.AvgDelay}}		{{.MaxDelay}}
  resp		{{.MinRes}}		{{.AvgRes}}		{{.MaxRes}}		  

Status code detail:
{{ formatMap .CodeDetail }}
`
)

func formatMap(data map[string]int) string {
	var keys []string

	for k, _ := range data {
		keys = append(keys, k)

	}
	sort.Strings(keys)
	res := new(bytes.Buffer)
	for _, k := range keys {
		res.WriteString(fmt.Sprintf("  %s:\t\t%d\n", k, data[k]))
	}
	return res.String()
}

var tmplFuncMap = template.FuncMap{
	"formatMap": formatMap,
}

func Print(summaryData SummaryDataStruct) {

	tmpl, err := template.New("test").Funcs(tmplFuncMap).Parse(htmlTemplate)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(os.Stdout, summaryData)
	if err != nil {
		panic(err)
	}

}
