
![ysab](https://github.com/yunsonbai/ysab/blob/master/ysab2.jpeg)

ysab is a tool, that can help you to get some performance parameters of your http server.
It can help you to send multiple urls with different parameters.

## Installation
* step 1:
    * Linux: wget https://github.com/yunsonbai/ysab/releases/download/v0.4.1/ysab_linux_0.4.1.tgz
    * MacOS: wget https://github.com/yunsonbai/ysab/releases/download/v0.4.1/ysab_mac_0.4.1.tgz
* step 2:
    * mv ysab_x_x ysab
    * chmod 777 ysab

* step 3:
    * mv ysab /usr/bin/


## Usage
```
Options:
  -r  Round of request to run.
  -n  Number of request to run concurrently, n>0 if n>900 n will be set to 900.
  -m  HTTP method, one of GET, POST, PUT, DELETE.
  -u  Url of request.
  -H  Add Arbitrary header line.
      For examples:
      -H "Accept: text/html". Set Accept to header.
      -H "Uid: yunson" -H "Content-Type: application/json". Set two fields to header.
  -t  Timeout for each request in seconds. Default is 10s.
  -d  HTTP request body. 
      For examples:
      '{"a": "a"}'
  -h  This help
  -v  Show verison
  -urlsfile  The urls file path. If you set this Option, -u,-d,-r will be invalid.
      For examples:
      -urlfile /tmp/urls.txt
```

## Some examples
* e1: ysab -n 900 -r 30 -u http://10.10.10.10:8080/test
* e2: ysab -n 900 -urlsfile ./examples/urls.txt
* e3: ysab -n 900 -r 30 -u http://10.10.10.10:8080/add -d '{"name": "yunson"}'
* e4: ysab -n 900 -urlsfile -m POST ./examples/urls2.txt

## Result show
```
(http://10.10.10.10:8080/test is API, it is writed by gin. The api will respone "hello world".)
[yunson ~]# ysab -n 900 -r 30 -u http://10.10.10.10:8080/test

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

## Note
* use -urlsfile
```
You can use -urlsfile to send multiple requests with different body.
cmd example:
	ysab -n 800 -urlsfile ./examples/urls.txt

urls.txt exaple:
	examples/urls.txt
	You can use create_urls.py to create a urls.txt file.
```
* use -u
```
example:
    ysab -n 900 -r 30 -u "http://10.121.130.218:8080/test"
    ysab -n 900 -r 30 -u 'http://10.121.130.218:8080/test'
    ysab -n 900 -r 30 -u http://10.121.130.218:8080/test
```

## Other
```
Thanks to Jason-Liu-Dream(https://github.com/Jason-Liu-Dream).
```
