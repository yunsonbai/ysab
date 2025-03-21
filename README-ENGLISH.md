
![ysab](https://github.com/yunsonbai/ysab/blob/master/ysab2.jpeg)

ysab is a tool that can help you get some performance parameters of http server stress test.
It can help you send requests with different parameters, so you can easily replay the real request online.

## Installation
* mac
    * curl -L -o install_mac.sh https://github.com/yunsonbai/ysab/releases/download/install-tool/install_mac.sh && sh install_mac.sh && rm -rf install_mac.sh

    * If report a permission problem, please execute: curl -L -o install_mac.sh https://github.com/yunsonbai/ysab/releases/download/install-tool/install_mac.sh && sudo sh install_mac.sh && rm -rf install_mac.sh

        ```If you cannot enter the ysab command after installation, you can restart the terminal.```

* linux
    * curl -L -o install_linux.sh https://github.com/yunsonbai/ysab/releases/download/install-tool/install_linux.sh && sh install_linux.sh && rm -rf install_linux.sh

    * If report a permission problem, please execute: curl -L -o install_linux.sh https://github.com/yunsonbai/ysab/releases/download/install-tool/install_linux.sh && sudo sh install_linux.sh && rm -rf install_linux.sh

* arm
    * If you need to run the arm version, you can clone a copy of the code, build it on the arm machine (go build -o ysab), and then put the executable file ysab under /usr/local/bin/

## Usage
* ysab -h

```
Options:
  -r  Rounds of request to run, total requests equal r * n
  -n  Number of simultaneous requests, 0<n<=900, depends on machine performance.
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
  -d  HTTP request body. 
      eg:
      '{"a": "a"}'
  -b  reader buf size(KB), default 256
  -h  This help
  -v  Show verison
  -urlsfile  The urls file path. If you set this Option, -u,-d will be invalid
      eg:
      -urlsfile /tmp/urls.txt
```

* Note: -urlsfile is the key parameter for sending requests with different parameters. For details of the file, please refer to examples/post_urls.txt and examples/get_urls.txt

## Some examples
1. GET request, 300 goroutines are processed together, each coroutine does 2 rounds
  
    ysab -n 300 -r 2 -u 'http://10.10.10.10:8080/test'

2. POST request, 400 goroutines are processed together, each coroutine does 2 rounds

    ysab -n 300 -r 2 -m POST -u 'http://10.10.10.10:8080/add' -d '{"name": "yunson"}'

3. GET request, 400 goroutines are processed together, lasting 100 seconds
  
    ysab -n 300 -T 100 -u 'http://10.10.10.10:8080/test'

4. POST request, 400 goroutines are processed together, and the test lasts for 100 seconds
    ysab -n 300 -T 100 -m POST -u 'http://10.10.10.10:8080/add' -d '{"name": "yunson"}'

5. GET request, 400 goroutines process the connection in urls.txt together, and process the connection in urls.txt twice
  
    ysab -r 2 -n 400 -urlsfile ./examples/urls.txt

6. GET request, 400 goroutines process the connection in urls.txt together, and loop the connection in the file for 100 seconds.

    ysab -T 100 -n 400 -urlsfile ./examples/urls.txt

7. POST request, 400 goroutines process the connection in urls.txt together, and loop the connection in the file for 100 seconds.

    ysab -T 100 -n 400 -m POST -urlsfile ./examples/urls.txt

## Result show
```
(http://10.10.10.10:8080/test is API, it is writed by gin, limit rate. The api will respone "hello world".)

[yunson ~]# ysab -n 900 -r 3 -u 'http://10.10.10.10:8080/test'

Summary:
  Complete requests:		2700
  Failed requests:		2550
  Time taken (s):		2.203996471
  Total data size (Byte):	0
  Data size/request (Byte):	0
  Max use time (ms):		2076
  Min use time (ms):		3
  Average use time (ms):	139.997
  Requests/sec:			1225.047333571706
  SuccessRequests/sec:		68.05818519842812

Percentage of waiting time (ms):
    10.00%:		12
    25.00%:		26
    50.00%:		46
    75.00%:		69
    90.00%:		122
    95.00%:		1291
    99.00%:		1992
    99.90%:		2040
    99.99%:		2052


Time detail (ms)
  item		min		mean		max
  dns		0		0		0
  conn		0		11.714		78
  wait		3		69.509		1010
  resp		0		10.53		50

Response Time histogram (code: requests):
  200:		150
  429:		2550
```

## about http code
* 2xx: Success
* != 2xx: Faild
    * 5xx:
        * 500: Server Error
        * 503: May be connection refused or connection reset by peer, you need to check your server.
    * other: [http code](https://en.wikipedia.org/wiki/List_of_HTTP_status_codes)


## About quotation marks
The following three methods are available. If there are special symbols, it is recommended to wrap them with ''

1. ysab -n 900 -r 50 -u "http://10.121.130.218:8080/test"
2. ysab -n 900 -r 50 -u 'http://10.121.130.218:8080/test'
3. ysab -n 900 -r 50 -u http://10.121.130.218:8080/test

## Other
You can use -urlsfile to send multiple requests with different body.

* cmd example:

	* ysab -r 1 -n 500 -urlsfile ./examples/urls.txt
  
  * ysab -r 2 -n 500 -m POST -urlsfile ./examples/post_urls.txt
  
  * ysab -T 30 -n 500 -m POST -urlsfile ./examples/post_urls.txt

* urls.txt example
	* examples/xx_urls.txt
	* You can use create_urls.py to create a urls.txt file.

## Acknowledgements
* [Jason-Liu-Dream](https://github.com/Jason-Liu-Dream)
* [zbing3](https://github.com/zbing3)
* [cugbtang](https://github.com/cugbtang)