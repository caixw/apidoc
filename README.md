apidoc [![Build Status](https://travis-ci.org/caixw/apidoc.svg?branch=master)](https://travis-ci.org/caixw/apidoc)
======

apidoc 是一个简单的 RESTful API 文档生成工具，它从代码注释中提取特定格式的内容，生成文档。
目前支持支持以下语言：C#、C/C++、D、Erlang、Go、Groovy、Java、Javascript、Pascal/Delphi、
Perl、PHP、Python、Ruby、Rust、Scala 和 Swift。

具体文档可参考：http://apidoc.tools

```go
/**
 * @api get /users 获取所有的用户信息
 * @apiGroup users
 * @apiQuery page int 显示第几页的内容
 * @apiQuery size int 每页显示的数量
 *
 * @apiSuccess 200 ok
 * @apiParam count int 符合条件的所有用户数量
 * @apiParam users array 用户列表。
 * @apiExample json
 * {
 *     "count": 500,
 *     "users": [
 *         {"id":1, "username": "admin1", "name": "管理员2"},
 *         {"id":2, "username": "admin2", "name": "管理员2"}
 *     ],
 * }
 * @apiExample xml
 * <users count="500">
 *     <user id="1" username="admin1" name="管理员1" />
 *     <user id="2" username="admin2" name="管理员2" />
 * </users>
 */
func login(w http.ResponseWriter, r *http.Request){
    // TODO
}
```



### 安装

```shell
go get github.com/caixw/apidoc
```



### 集成

若需要将 apidoc 当作包集成到其它 Go 程序中，可分别引用 `input` 和 `output` 的相关函数：

```go
start := time.Now()
docs := doc.New()

tag, err := locale.GetLocale()
if err != nil {
    panic(err)
}
locale.SetLocale(tag)

// 分析文档内容
inputOptions := &input.Options{...}
docs, err := input.Parse(docs, inputOptions)
if err != nil {
    // TODO
}

// 输出内容
outputOptions := &output.Options{...}
outputOptions.Elapsed = time.Now().Sub(start)
if err = output.Render(docs, outputOptions); err != nil {
    // TODO
}
```



### 参与开发

请阅读 [CONTRIBUTING.md](CONTRIBUTING.md) 文件的相关内容。



### 版权

本项目采用 [MIT](https://opensource.org/licenses/MIT) 开源授权许可证，完整的授权说明可在 [LICENSE](LICENSE) 文件中找到。
