apidoc [![Build Status](https://travis-ci.org/caixw/apidoc.svg?branch=master)](https://travis-ci.org/caixw/apidoc)
======

提取特定格式的注释，生成api文档。目前支持以下标签：

- @api
- @apiParam
- @apiQuery
- @apiVersion
- @apiGroup
- @apiRequest
- @apiHeader
- @apiStatus
- @apiExample


#### 命令行语法:
```shell
apidoc [options] src doc

src:
 源文件所在的目录。
doc:
 产生的文档所在的目录。
```
详细内容可参参考程序-h参数。


### 安装

```shell
go get github.com/caixw/apidoc
```


### 版权

本项目采用[MIT](http://opensource.org/licenses/MIT)开源授权许可证，完整的授权说明可在[LICENSE](LICENSE)文件中找到。
