apidoc
[![Build Status](https://img.shields.io/endpoint.svg?url=https%3A%2F%2Factions-badge.atrox.dev%2Fcaixw%2Fapidoc%2Fbadge%3Fref%3Dmaster&style=flat)](https://actions-badge.atrox.dev/caixw/apidoc/goto?ref=master)
[![Go version](https://img.shields.io/badge/Go-1.13-brightgreen.svg?style=flat)](https://golang.org)
[![Go Report Card](https://goreportcard.com/badge/github.com/caixw/apidoc)](https://goreportcard.com/report/github.com/caixw/apidoc)
[![license](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat)](https://opensource.org/licenses/MIT)
[![codecov](https://codecov.io/gh/caixw/apidoc/branch/master/graph/badge.svg)](https://codecov.io/gh/caixw/apidoc)
======

apidoc 是一个简单的 RESTful API 文档生成工具，它从代码注释中提取特定格式的内容，生成文档。
目前支持支持以下语言：C#、C/C++、D、Erlang、Go、Groovy、Java、JavaScript、Pascal/Delphi、
Perl、PHP、Python、Ruby、Rust、Scala 和 Swift。

具体文档可参考：<https://apidoc.tools>

```go
/**
 * <api method="GET" summary="获取所有的用户信息">
 *     <path path="/users">
 *         <query name="page" type="number" default="0" summary="显示第几页的内容" />
 *         <query name="size" type="number" default="20" summary="每页显示的数量" />
 *     </path>
 *     <tag>user</tag>
 *     <server>users</server>
 *     <response status="200" type="object" mimetype="application/json">
 *         <param name="count" type="int" optional="false" summary="符合条件的所有用户数量" />
 *         <param name="users" type="object" array="true" summary="用户列表">
 *             <param name="id" type="int" summary="唯一 ID" />
 *             <param name="name" type="string" summary="姓名" />
 *         </param>
 *         <example mimetype="application/json">
 *         <![CDATA[
 *         {
 *             "count": 500,
 *             "users": [
 *                 {"id":1, "name": "管理员2"},
 *                 {"id":2, "name": "管理员2"}
 *             ],
 *         }
 *         ]]>
 *         </example>
 *     </response>
 *     <response status="500" mimetype="application/json" type="object">
 *         <param name="code" type="int" summary="错误代码" />
 *         <param name="msg" type="string" summary="错误内容" />
 *     </response>
 * </api>
 */
func login(w http.ResponseWriter, r *http.Request) {
    // TODO
}
```

使用
----

在 <https://github.com/caixw/apidoc/releases> 提供了主流系统下可用软件，可直接下载使用。
如果你使用的系统不在此列，则需要手动下载编译。

支持多种本地化语言，默认情况下会根据当前系统所使用的语言进行调整。
也可以通过设置环境变更 `LANG` 指定一个本地化信息。*nix 系统也可以使用以下命令：

```shell
LANG=lang apidoc
```

将其中的 lang 设置为你需要的语言。

集成
----

若需要将 apidoc 当作包集成到其它 Go 程序中，可参考以下代码：

```go
// 初始本地化内容
apidoc.Init(language.MustParse("zh-Hans"))

// 可以自定义实现具体的错误处理方式
h := message.NewHandler(...)

output := &output.Options{...}
inputs := []*input.Options{
    &input.Options{},
}

apidoc.Build(h, output, inputs...)
```

参与开发
----

请阅读 [CONTRIBUTING.md](CONTRIBUTING.md) 文件的相关内容。

版权
----

本项目源码采用 [MIT](https://opensource.org/licenses/MIT) 开源授权许可证，完整的授权说明可在 [LICENSE](LICENSE) 文件中找到。

文档内容的版权由各个文档各自表述。
