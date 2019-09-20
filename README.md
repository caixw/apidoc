apidoc [![Build Status](https://travis-ci.org/caixw/apidoc.svg?branch=master)](https://travis-ci.org/caixw/apidoc)
[![Go version](https://img.shields.io/badge/Go-1.13-brightgreen.svg?style=flat)](https://golang.org)
[![Go Report Card](https://goreportcard.com/badge/github.com/caixw/apidoc)](https://goreportcard.com/report/github.com/caixw/apidoc)
[![license](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat)](https://opensource.org/licenses/MIT)
[![codecov](https://codecov.io/gh/caixw/apidoc/branch/master/graph/badge.svg)](https://codecov.io/gh/caixw/apidoc)
======

apidoc 是一个简单的 RESTful API 文档生成工具，它从代码注释中提取特定格式的内容，生成文档。
目前支持支持以下语言：C#、C/C++、D、Erlang、Go、Groovy、Java、JavaScript、Pascal/Delphi、
Perl、PHP、Python、Ruby、Rust、Scala 和 Swift。

具体文档可参考：https://apidoc.tools

```go
/**
 * @api get /users 获取所有的用户信息
 * @apiTags users
 * @apiQuery page int 显示第几页的内容
 * @apiQuery size int 每页显示的数量
 *
 * @apiResponse 200 object application/json ok
 * @apiParam count int required 符合条件的所有用户数量
 * @apiParam users array.object required 用户列表。
 * @apiExample application/json
 * {
 *     "count": 500,
 *     "users": [
 *         {"id":1, "username": "admin1", "name": "管理员2"},
 *         {"id":2, "username": "admin2", "name": "管理员2"}
 *     ],
 * }
 * @apiExample application/xml
 * <users count="500">
 *     <user id="1" username="admin1" name="管理员1" />
 *     <user id="2" username="admin2" name="管理员2" />
 * </users>
 */
func login(w http.ResponseWriter, r *http.Request) {
    // TODO
}
```



### 安装

```shell
go get github.com/caixw/apidoc
```

支持多种本地化语言，默认情况下会根据当前系统所使用的语言进行调整。若需要手动指定，
windows 可以设置一个 `LANG` 环境变量指定，*nix 系统可以使用以下命令：

```shell
LANG=lang apidoc
```

将其中的 lang 设置为你需要的语言。



### 集成

若需要将 apidoc 当作包集成到其它 Go 程序中，可参考以下代码：

```go
// 初始本地化内容
apidoc.InitLocale()

// 可以自定义实现具体的错误处理方式
h := errors.NewHandler()

erro := log.NewLogger()
output := &output.Options{...}
inputs := []*input.Options{
    &input.Options{},
}

apidoc.Do(h, output, inputs...)
```



### 参与开发

请阅读 [CONTRIBUTING.md](CONTRIBUTING.md) 文件的相关内容。



### 版权

本项目采用 [MIT](https://opensource.org/licenses/MIT) 开源授权许可证，完整的授权说明可在 [LICENSE](LICENSE) 文件中找到。
