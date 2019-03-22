
![ysab](https://yunsonbai.github.io/images/golang/ysab2.jpeg)

ysab is a tool, that can help you to get some performance parameters of your http server.
It can help you to send multiple urls with different parameters.

## Installation
* step 1:
    * Linux: wget https://github.com/yunsonbai/ysab/releases/download/v0.1/ysab_Linux_0.1.tgz 
    * MacOS: wget https://github.com/yunsonbai/ysab/releases/download/v0.1/ysab_macOS_0.1.tgz
* step 2:
    * tar -zxvf ysab_x_x.tgz

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
  Complete requests:	27000
  Failed requests:	0
  Total data size(ContentLength):	324000
  Data size/request:	12
  Max use time:		257 ms
  Min use time:		10 ms
  Average use time:	36.074 ms
  Requests/sec:		13500

QPS time histogram (timestamp: requests):
  1551254255:		14198
  1551254256:		12802


Use Time Percent:
  <=50ms:		  87.74%
  <=100ms:		  99.31%
  <=300ms:		  100.00%
  <=500ms:		  100.00%
  >500ms:		  0.00%

Code Time histogram (code: requests):
  200:		27000


Time detail (ms)
  item		min		mean		max
  dns		0		0		0
  conn		0		1.088		51
  wait		10		33.82		257
  resp		0		0.596		28
```

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
