apidoc [![Build Status](https://travis-ci.org/caixw/apidoc.svg?branch=master)](https://travis-ci.org/caixw/apidoc)
======

初步的文档设计如下：

```
 @api sumary
  description
 @apiURL api/test.json
 @apiMethods get/put
 @apiParam id int {optional} 用户ID
 @apiQuery order int
 @apiVersion 1.0
 @apiGroup groupName
 @apiRequest json/xml
  @apiHeader aa bb
  @apiHeader aa bb
  @apiParam abc int desc
  @apiExample json
 @apiStatus 200 json/xml summary
  @apiHeader aa bb
  @apiParam abc int desc
  @apiExample json
  {
      a:b
      c:d
  }
  @apiExample xml
  <root>
      <item key="a">b</item>
      <item key="c">d</item>
  </root>
@apiEnd
```

```
@api summary
description
```
 - summary 该api的简单描述
 - description 该api的详细介绍。


#### 命令行语法:
```shell
apidoc [options] src doc

options:
 -h    显示当前帮助信息；
 -v    显示apidoc和go程序的版本信息；
 -t    源文件类型；
 -r    是否搜索子目录，默认为true；
 -ext  需要监视的扩展名，若不指定，则只搜索与t参数指定的类型。

src:
 源文件所在的目录。
doc:
 产生的文档所在的目录。
```


### 安装

```shell
go get github.com/caixw/apidoc
```


### 文档

[![Go Walker](http://gowalker.org/api/v1/badge)](http://gowalker.org/github.com/caixw/apidoc)
[![GoDoc](https://godoc.org/github.com/caixw/apidoc?status.svg)](https://godoc.org/github.com/caixw/apidoc)


### 版权

本项目采用[MIT](http://opensource.org/licenses/MIT)开源授权许可证，完整的授权说明可在[LICENSE](LICENSE)文件中找到。
