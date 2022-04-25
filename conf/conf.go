package conf

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/yunsonbai/ysab/tools"
)

var usage = `Usage: ysab [Options]

Options:
  -r  Rounds of request to run, total requests equal r * n
  -n  Number of simultaneous requests, 0<n<=900, depends on machine performance
  -m  HTTP method, one of GET, POST, PUT, DELETE, Head, Default is GET
  -u  Url of request, use " please
      eg: 
      -u 'https://yunsonbai.top/?name=yunson'
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
  -F  Use express mode ([0]/1). Select 1, qps may be higher, but Time detail cannot be counted
  -v  Show verison
  -urlsfile  The urls file path. If you set this Option, -u,-d,-r will be invalid
      eg:
      -urlfile /tmp/urls.txt
`

type Config struct {
	Round       int    // 请求多少轮, 只对Url有效
	N           int    // 并发数
	Url         string // 需要请求的url, 与UrlFilePath只能一个有效
	UrlFilePath string //url文件路径
	UrlNum      int
	Headers     map[string]string
	TimeOut     int    //单次请求超时时间
	Method      string // 请求方法
	Fast        int
	Body        string
	StartTime   int64
	EndTime     int64
}

type headersSlice []string

var (
	Conf = Config{
		Headers: make(map[string]string),
	}
	fbr     *bufio.Reader
	headers headersSlice
)

func (h *headersSlice) String() string {
	return fmt.Sprintf("%s", *h)
}

func (h *headersSlice) Set(value string) error {
	*h = append(*h, value)
	return nil
}

func useConfile(filepath string) {
	configFile, err := os.Open(filepath)
	defer configFile.Close()
	if err != nil {
		panic(err)
	}
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&Conf)
	if err != nil {
		panic(err)
	}
}

func confError(err error) {
	fmt.Println(usage)
	fmt.Println(err)
	os.Exit(1)
}

func arrangeOptions() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage)
	}
	methods := [5]string{"GET", "POST", "PUT", "DELETE", "HEAD"}
	help := flag.Bool("h", false, "")
	m := flag.String("m", "GET", "")
	flag.Var(&headers, "H", "")
	body := flag.String("d", "", "")
	round := flag.Int("r", 0, "")
	fast := flag.Int("F", 0, "")
	version := flag.Bool("v", false, "")
	n := flag.Int("n", 0, "")
	url := flag.String("u", "", "")
	timeout := flag.Int("t", 10, "")
	urlsfile := flag.String("urlsfile", "", "")
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}
	if *version {
		fmt.Println("version is", VERSION)
		os.Exit(0)
	}
	methoderr := "(-m) method not allowed."
	for _, v := range methods {
		if v == *m {
			methoderr = ""
			break
		}
	}
	if methoderr != "" {
		confError(errors.New(methoderr))
	}
	Conf.Method = *m
	for _, v := range headers {
		headers := tools.KeyValueRexpGetKV(v)
		if len(headers) != 3 {
			confError(errors.New(
				fmt.Sprintf("(-H) header of %s error.", v)))
		}
		Conf.Headers[headers[1]] = headers[2]
	}
	Conf.Body = *body
	if *n <= 0 {
		confError(errors.New("(-n) Number must be greater than 0."))
	}
	Conf.N = *n
	if *round <= 0 {
		confError(errors.New("(-r) Round must be greater than 0."))
	}
	Conf.Round = *round
	Conf.Url = *url
	if *timeout <= 0 {
		confError(errors.New("(-t) timeout must be greater than 0."))
	}
	Conf.TimeOut = *timeout
	if *fast < 0 || *fast > 1 {
		confError(errors.New("(-F) must be 0 or 1."))
	}
	Conf.Fast = *fast
	if *url == "" && *urlsfile == "" {
		confError(errors.New("-u or -urlsfile must choice one."))
	}
	Conf.UrlFilePath = *urlsfile
	Conf.Url = tools.ReplaceQmarks(Conf.Url, "")
}

func init() {
	arrangeOptions()
	Conf.UrlNum = Conf.N * Conf.Round
	if Conf.UrlFilePath != "" {
		fi, err := os.Open(Conf.UrlFilePath)
		defer fi.Close()
		if err != nil {
			panic(err)
		}
		fbr = bufio.NewReader(fi)
		count := 0
		for {
			line, _, err := fbr.ReadLine()
			if err == io.EOF {
				break
			}
			if string(line) != "" {
				count++
			}
		}
		Conf.UrlNum = count
	}
	if Conf.N > 900 {
		Conf.N = 900
	}
	if Conf.TimeOut <= 0 || Conf.TimeOut > 60 {
		Conf.TimeOut = 60
	}
	Conf.StartTime = tools.GetNowUnixNano()
}
