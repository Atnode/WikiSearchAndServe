<div>
  <h3 align="center"><img src="static/img/logo.png"/><br>WikiSearchAndServe</h3>
</div>

## Introduction
This project is a fork of [gozimhttpd](https://github.com/akhenakh/gozim) with only the search function and with some corrections.  
This tool also offers a simpler, modern and easy to use interface.  
This fork also has multilingual interface support.  

## How to build?
On Ubuntu/Debian you need those packages to compile WikiSearchAndServe:  
```
apt-get install git liblzma-dev mercurial build-essential golang
```

For the indexer bleve to work properly it's recommended that you use leveldb as storage:  
```
go get -u -v -tags all github.com/blevesearch/bleve/...
```

WikiSearchAndServe is using go.rice to embed html/css in the binary install the rice command:  
```
go get github.com/GeertJohan/go.rice
go get github.com/GeertJohan/go.rice/rice
go install github.com/GeertJohan/go.rice
go install github.com/GeertJohan/go.rice/rice
```

Get and build the WikiSearchAndServe executable:  
```bash
git clone https://github.com/EducaBox/WikiSearchAndServe
cd WikiSearchAndServe
go build
```

Then run this command to embed the files:  
```
rice append --exec WikiSearchAndServe
```

## Running
You need to build a index file to run WikiSearchAndServe.  
To start WikiSearchAndServe: `./WikiSearchAndServe -path=yourzimfile.zim -index=yourzimfile.idx`

***

Copyright (c) 2019 EducaBox

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.