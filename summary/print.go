package summary

import (
	"bytes"
	"fmt"
	"os"
	"text/template"
)

var (
	htmlTemplate = `
Summary:
  Complete requests:	{{ .CompleteRequests }}
  Failed requests:	{{ .FailedRequests }}
  Total data size:	{{ .TotalDataSize }}
  Data size/request:	{{ .AvgDataSize }}
  Max use time:		{{.MaxUseTime}} ms
  Min use time:		{{.MinUseTime}} ms
  Average use time:	{{.AvgUseTime}} ms
  Requests/sec:		{{ .RequestsPerSec }}

QPS time histogram (timestamp: requests):
{{ formatMap .QpsDetail }}

Use Time Percent:
  <=50ms:		{{ .Tlt50mPercent }}
  <=100ms:		{{ .Tlt100mPercent }}
  <=300ms:		{{ .Tlt300mPercent }}
  <=500ms:		{{ .Tlt500mPercent }}
  >500ms:		{{ .Tgt500mPercent }}

Code Time histogram (code: requests):
{{ formatMap .CodeDetail }}

Time detail (ms)
  item		min		mean		max
  dns		{{.MinDNS}}		{{.AvgDNS}}		{{.MaxDNS}}
  conn		{{.MinConn}}		{{.AvgConn}}		{{.MaxConn}}
  wait		{{.MinDelay}}		{{.AvgDelay}}		{{.MaxDelay}}
  resp		{{.MinRes}}		{{.AvgRes}}		{{.MaxRes}}		  
`
)

func formatMap(data map[int]int) string {
	res := new(bytes.Buffer)
	for k, v := range data {
		res.WriteString(fmt.Sprintf("  %d:\t\t%d\n", k, v))
	}
	return res.String()
}

var tmplFuncMap = template.FuncMap{
	"formatMap": formatMap,
}

func Print(summaryData SummaryData) {
	// fmt.Println("summaryData:", summaryData)

	tmpl, err := template.New("test").Funcs(tmplFuncMap).Parse(htmlTemplate)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(os.Stdout, summaryData)
	if err != nil {
		panic(err)
	}

}
