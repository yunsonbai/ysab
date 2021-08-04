
![ysab](https://github.com/yunsonbai/ysab/blob/master/ysab2.jpeg)

ysab 是一个可以帮助你获取http服务器压力测试性能指标的工具，有点像Apache的ab。不同的是，它可以帮你发送携带不同参数的请求，这样你就可以便捷地重放线上的真实请求。

[English](./README-ENGLISH.md)

## 安装
* mac

wget https://github.com/yunsonbai/ysab/releases/download/install-tool/install_mac -O install_mac && sh install_mac && rm -rf install_mac

如果报权限问题请执行:
wget https://github.com/yunsonbai/ysab/releases/download/install-tool/install_mac -O install_mac && sudo sh install_mac && rm -rf install_mac

如果安装完后不能输入 ysab 命令，可以重启终端或者执行 source /etc/profile

* linux

wget https://github.com/yunsonbai/ysab/releases/download/install-tool/install_linux -O install_linux && sh install_linux && rm -rf install_linux

如果报权限问题请执行:
wget https://github.com/yunsonbai/ysab/releases/download/install-tool/install_linux -O install_linux && sudo sh install_linux && rm -rf install_linux

## 参数说明
* ysab -h

```
Options:
  -r  压测轮数，总的请求量是 r * n
  -n  并发数，最大900，最小1
  -m  HTTP method, 可选值 GET，POST，PUT，DELETE，Head，默认GET
  -u  Url of request, 使用 " 括起来
      例如: 
      -u "https://yunsonbai.top/?name=yunson"
  -H  添加请求头
      例如:
      -H "Accept: text/html"  设置 Accept
      -H "Host: yunsonbai.top"  设置 Host
      -H "Uid: yunson" -H "Content-Type: application/json" 设置Uid和Content-Type
  -t  每个请求的超时时间，单位为秒，默认10
  -d  请求体 
      例如:
      '{"a": "a"}'
  -h  帮助
  -v  显示版本号
  -urlsfile  包含所有请求信息的文件，如果设置了该参数, -u,-d,-r 将会失效
      例如:
      -urlfile /tmp/urls.txt
```

* 注意: -urlsfile 是实现发送携带不同参数请求的关键参数，文件详细内容，可参照examples/post_urls.txt 和 examples/get_urls.txt

## 一些例子
* 1: ysab -n 900 -r 2 -u http://10.10.10.10:8080/test
* 2: ysab -n 900 -urlsfile ./examples/get_urls.txt
* 3: ysab -n 900 -r 2 -m POST -u http://10.10.10.10:8080/add -d '{"name": "yunson"}'
* 4: ysab -n 900 -m POST -urlsfile ./examples/post_urls.txt

## 结果展示
```
(http://10.10.10.10:8080/test 是一个借助gin完成的测试 API. 这个 API 的 response 是 "hello world".)

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

## 关于 http code
* 2xx: Success
* != 2xx: Faild
    * 5xx:
        * 500: Server Error
        * 503: May be connection refused or connection reset by peer, you need to check your server.
    * other: [http code](https://en.wikipedia.org/wiki/List_of_HTTP_status_codes)

## 注意
* 推荐使用 -urlsfile
```
你可以使用 -urlsfile 发送携带不同 body 或 url 的请求

样例:
	ysab -n 500 -urlsfile ./examples/get_urls.txt
    ysab -n 500 -m POST -urlsfile ./examples/post_urls.txt

urls.txt example:
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

## 鸣谢
* [Jason-Liu-Dream](https://github.com/Jason-Liu-Dream)
* [zbing3](https://github.com/zbing3)
