
![ysab](https://github.com/yunsonbai/ysab/blob/master/ysab2.jpeg)

ysab 是一个可以帮助你获取http服务器压力测试性能指标的工具，有点像Apache的ab，不同的是，它可以帮你发送携带不同参数的请求，这样你就可以便捷地重放线上的真实请求。

[English](./README-ENGLISH.md)

## 安装
* mac

    * curl -L -o install_mac.sh https://github.com/yunsonbai/ysab/releases/download/install-tool/install_mac.sh && sh install_mac.sh && rm -rf install_mac.sh

    * 如果报权限问题请执行: curl -L -o install_mac.sh https://github.com/yunsonbai/ysab/releases/download/install-tool/install_mac.sh && sudo sh install_mac.sh && rm -rf install_mac.sh

        ```如果安装完后不能输入 ysab 命令，可以重启终端```

* linux

    * curl -L -o install_linux.sh https://github.com/yunsonbai/ysab/releases/download/install-tool/install_linux.sh && sh install_linux.sh && rm -rf install_linux.sh

    * 如果报权限问题请执行: curl -L -o install_linux.sh https://github.com/yunsonbai/ysab/releases/download/install-tool/install_linux.sh && sudo sh install_linux.sh && rm -rf install_linux.sh

* arm
    * 如果需要运行arm版本，可以clone一份代码，在arm机上build一下(go build -o ysab)，然后把可执行文件ysab放到/usr/local/bin/下即可

## 参数说明
* ysab -h

```
Options:
  -r  压测轮数，总的请求量是 r * n
  -n  并发数，最大900，最小1
  -T  运行时间(单位秒, 默认0), T大于0时, r无效
  -m  HTTP method, 可选值 GET，POST，PUT，DELETE，Head，默认GET
  -u  Url of request, 如果有特殊符号需要用引号
      例如: 
      -u "https://yunsonbai.top/?name=yunson"
      -u 'https://yunsonbai.top/?name=yunson'
      -u https://yunsonbai.top/?name=yunson
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
  -urlsfile  包含所有请求信息的文件，如果设置了该参数, -u,-d 将会失效
      例如:
      -urlsfile /tmp/urls.txt
```

* 注意: -urlsfile 是实现发送携带不同参数请求的关键参数，文件详细内容，可参照examples/post_urls.txt 和 examples/get_urls.txt

## 一些例子
1. GET请求，300个协程一起处理，每个协程做2轮

    ysab -n 300 -r 2 -u 'http://10.10.10.10:8080/test'

2. POST请求，400个协程一起处理，每个协程做2轮
  
    ysab -n 300 -r 2 -m POST -u 'http://10.10.10.10:8080/add' -d '{"name": "yunson"}'

3. GET请求，400个协程一起处理，持续100秒
  
    ysab -n 300 -T 100 -u 'http://10.10.10.10:8080/test'

4. POST请求，400个协程一起处理，持续100秒测试

    ysab -n 300 -T 100 -m POST -u 'http://10.10.10.10:8080/add' -d '{"name": "yunson"}'

5. GET请求，400个协程一起处理urls.txt中的连接，处理2次urls.txt中的连接

    ysab -r 2 -n 400 -urlsfile ./examples/urls.txt

6. GET请求，400个协程一起处理urls.txt中的连接，一直循环执行文件中的连接，持续100秒

    ysab -T 100 -n 400 -urlsfile ./examples/urls.txt

7. POST请求，400个协程一起处理urls.txt中的连接，一直循环执行文件中的连接，持续100秒

    ysab -T 100 -n 400 -m POST -urlsfile ./examples/urls.txt

## 结果展示
```
(http://10.10.10.10:8080/test 是一个借助gin完成的测试API, 限速. 这个 API 的 response 是 "hello world".)

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

## 关于 http code
* 2xx: Success
* != 2xx: Faild
    * 5xx:
        * 500: Server Error
        * 503: May be connection refused or connection reset by peer, you need to check your server.
    * other: [http code](https://en.wikipedia.org/wiki/List_of_HTTP_status_codes)

## 关于引号的使用
以下三个方式都可以, 若有特殊符号建议用''包裹
* ysab -n 900 -r 50 -u "http://10.121.130.218:8080/test"
* ysab -n 900 -r 50 -u 'http://10.121.130.218:8080/test'
* ysab -n 900 -r 50 -u http://10.121.130.218:8080/test

## 其他
推荐使用 -urlsfile, 你可以使用 -urlsfile 发送携带不同 body 或 url 的请求

* 样例:
	* ysab -r 1 -n 500 -urlsfile ./examples/urls.txt
    * ysab -r 2 -n 500 -m POST -urlsfile ./examples/post_urls.txt
    * ysab -T 30 -n 500 -m POST -urlsfile ./examples/post_urls.txt

* urls.txt example:
	* examples/xx_urls.txt
	* You can use create_urls.py to create a urls.txt file.


## 鸣谢
* [Jason-Liu-Dream](https://github.com/Jason-Liu-Dream)
* [zbing3](https://github.com/zbing3)
* [cugbtang](https://github.com/cugbtang)
