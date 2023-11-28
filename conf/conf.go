package conf

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/yunsonbai/ysab/tools"
)

var usage = `Usage: ysab [Options]

Options:
  -r  Rounds of request to run, total requests equal r * n
  -n  Number of simultaneous requests, 0<n<=900, depends on machine performance
  -T  Running time (seconds, default 0), r is invalid when T is greater than 0.
  -m  HTTP method, one of GET, POST, PUT, DELETE, Head, Default is GET
  -u  Url of request, quotes are recommended if there are envoy symbols
      eg: 
      -u 'https://yunsonbai.top/?name=yunson'
      -u "https://yunsonbai.top/?name=yunson"
      -u https://yunsonbai.top/?name=yunson
  -H  Add Arbitrary header line
      eg:
      -H "Accept: text/html", Set Accept to header
      -H "Host: yunsonbai.top", Set Host to header
      -H "Uid: yunson" -H "Content-Type: application/json", Set two fields to header
  -t  Timeout for each request in seconds, Default is 10
  -d  HTTP request body
      eg:
      '{"a": "a"}'
  -h  This help
  -v  Show verison
  -urlsfile  The urls file path. If you set this Option, -u,-d will be invalid
      eg:
      -urlsfile /tmp/urls.txt
`

type Config struct {
	N           uint16
	Round       uint32 // 轮次
	Duration    int    // 持续时间
	UrlNum      int64
	TimeOut     int64
	Url         string
	UrlFilePath string
	Headers     map[string]string
	Method      string
	Body        string
	StartTime   int64
	EndTime     int64
	TimeBase    int64
	UrlChNum    int
	SumResChNum int
}

type headersSlice []string

var (
	Conf = Config{
		Headers: make(map[string]string),
	}
	headers headersSlice
)

func (h *headersSlice) String() string {
	return fmt.Sprintf("%s", *h)
}

func (h *headersSlice) Set(value string) error {
	*h = append(*h, value)
	return nil
}

// func useConfile(filepath string) {
// 	configFile, err := os.Open(filepath)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer configFile.Close()
// 	jsonParser := json.NewDecoder(configFile)
// 	err = jsonParser.Decode(&Conf)
// 	if err != nil {
// 		panic(err)
// 	}
// }

func confError(err error) {
	fmt.Println(usage)
	fmt.Println(err)
	os.Exit(1)
}

func arrangeOptions() {
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
	}
	methods := [5]string{"GET", "POST", "PUT", "DELETE", "HEAD"}
	help := flag.Bool("h", false, "")
	flag.Var(&headers, "H", "")

	round := flag.Int("r", 0, "")
	version := flag.Bool("v", false, "")
	n := flag.Int("n", 0, "")
	body := flag.String("d", "", "")
	duration := flag.Int("T", 0, "")
	method := flag.String("m", "GET", "")
	url := flag.String("u", "", "")
	timeout := flag.Int64("t", 10, "")
	urlFilePath := flag.String("urlsfile", "", "")
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}
	if *version {
		fmt.Println("ysab version", VERSION)
		os.Exit(0)
	}
	Conf.Body = *body
	Conf.Duration = *duration
	Conf.Method = *method
	Conf.Url = *url
	Conf.TimeOut = *timeout
	Conf.UrlFilePath = *urlFilePath
	methoderr := "(-m) method not allowed."
	for _, v := range methods {
		if v == Conf.Method {
			methoderr = ""
			break
		}
	}
	if methoderr != "" {
		confError(errors.New(methoderr))
	}
	for _, v := range headers {
		headers := tools.KeyValueRexpGetKV(v)
		if len(headers) != 3 {
			confError(fmt.Errorf("(-H) header of %s error", v))
		}
		Conf.Headers[headers[1]] = headers[2]
	}

	if *n <= 0 || *n > 65535 {
		confError(errors.New("(-n) Number must be 0<n<65535"))
	}
	Conf.N = uint16(*n)

	if *round < 0 {
		confError(errors.New("(-r) Round must be greater than 0"))
	}
	Conf.Round = uint32(*round)
	if Conf.Duration < 0 {
		confError(errors.New("(-T) Time duration must be greater than 0"))
	}
	if Conf.Round == 0 && Conf.Duration == 0 {
		confError(errors.New("-T or -r you must choice one"))
	}

	if Conf.TimeOut <= 0 || Conf.TimeOut > 60 {
		confError(errors.New("(-t) timeout must be 0<t<=60"))
	}
	if Conf.Url == "" && Conf.UrlFilePath == "" {
		confError(errors.New("-u or -urlsfile must choice one"))
	}
	Conf.Url = tools.ReplaceQmarks(Conf.Url, "")
}

func init() {
	arrangeOptions()
	if Conf.N > 900 {
		Conf.N = 900
	}
	if Conf.TimeOut <= 0 || Conf.TimeOut > 60 {
		Conf.TimeOut = 60
	}
	Conf.TimeBase = 1000000
	Conf.TimeOut = Conf.TimeOut * Conf.TimeBase
	// Conf.StartTime = time.Now().UnixMicro()
	numCPU := runtime.NumCPU()
	Conf.UrlChNum = int(Conf.N) * numCPU
	Conf.SumResChNum = int(Conf.N) * numCPU / 2
	if Conf.UrlChNum < 4000 {
		Conf.UrlChNum = 4000
	}
	if Conf.SumResChNum < 2000 {
		Conf.SumResChNum = 2000
	}
}
