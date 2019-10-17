# 如何为 apidoc 贡献代码

apidoc 是一个基于 [MIT](https://opensource.org/licenses/MIT) 的开源软件。
欢迎大家共同参与开发。**若需要新功能，请先开 issue 讨论。**



## 本地化

本地化包含以下几个部分：
- `internal/locale` 主要包含了程序内各种语法错误以及命令行的提示信息；
- `docs/vx/locales.xsl` 包含展示界面中的本化元素；`vx` 表示版本信息，比如 `v5`、`v6` 等；
- `docs/index.*.xml` https://apidoc.tools 网站的内容，* 表示语言 ID，同时需要修改 `docs\locales.xml` 文件；


## 测试

xslt 可以通过 /cmd/xsltest 构建本地服务进行测试；


## 添加新编程语言支持

`internal/lang/lang.go` 文件中有所有语言模型的定义，若需要添加对新语言的支持，
在该文件中有详细的文档说明如何定义语言模弄。提交后请更新
[#11](https://github.com/caixw/apidoc/issues/11)
