## Ysab
```
ysab is a tool, that can help you to get some performance parameters of your http server.
```

## Installation
* step 1:
    * Linux: wget https://github.com/yunsonbai/ysab/releases/download/v0.1/ysab_Linux_0.1.tgz 
    * MacOS: wget https://github.com/yunsonbai/ysab/releases/download/v0.1/ysab_macOS_0.1.tgz
* step 2:
    * tar -zxvf ysab_*_*.tgz

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

## Note
* use urlsfile
```
You can use -urlsfile to send multiple requests with different body.
cmd example:
	ysab -n 800 -urlsfile ./examples/urls.txt
urls.txt exaple:
	examples/urls.txt
	You can use create_urls.py to create a urls.txt file.
```