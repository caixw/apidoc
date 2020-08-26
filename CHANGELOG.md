# CHANGELOG

## [Unreleased]

### Added

- 添加对 Julia 和 Nim 的支持；
- 添加 core.Error.WithField、core.Error.WithLocation 和 core.Error.AddTypes 方法；
- 添加 ErrorType 用于表示该语法错误的类型；
- LSP 的 textDocument/publishDiagnostics 添加了 DiagnosticTag 的支持；
- 添加 core.Searcher 接口；
- 添加了对 textDocument/references 的支持；
- 添加了对 textDocument/definition 的支持；
- 添加自定义通知 apidoc/outline 用于向客户端展示文档的基本列表；
- 添加自定义服务 apidoc/refreshOutline 用于客户端主动请求刷新文档列表；
- 配置文件添加 ignores 字段，用于忽略不需要解析目录；

### Fixed

- 修正同一个服务可以有多个相同接口的 BUG；
- 当用户第一次打开项目时，不会发送错误通知给用户的错误；
- 修正 textDocument/publishDiagnostics 的行为错误；

### Changed

- internal/ast.Base 现在直接包含 core.Location 对象；

## [v7.1.0]

### Added

- 添加对 lua、dart、typescript、zig 和 lisp 及其方言的支持；
- 为 version 子命令添加 `kind` 选项；
- 为 core.Range 和 core.Location 添加了 Equal 方法用于判断值是否相等；
- 添加对 textDocument/foldingRange 的支持；
- 添加对 textDocument/semanticTokens 的支持；

### Fixed

- 无法识别 string.email 的错误；
- 修正文档中关于 apidoc.xml-namespaces 字段的描述错误；
- 修正 ast.APIDoc.Parse 未记录 URI 的错误；
- 修正 textDocument/hover 可能导致整个程序 panic 的错误；
- 修正 textDocument/didChange 解析错误；
- 修正 textDocument/publishDiagnostics 无法正确发送的错误；

### Changed

- 调整 build.Config 对路径项的处理方式，读取时转换为经对路径，保存时尽量转换为相对路径；
- 当文档中存在不能识别的元素时，不再 panic；
- 解析文档时在不碰到语法错误时不再退出解析进程；
- 优化解析 xml 的过程；

## [v7.0.0]

### Added

- 为 param 数据添加 array-style 字段，用以描述数组在查询参数的表现形式；
- 添加 core 包；
- 采用与 LSP 相同的方式定位错误信息；
- 采用新的编码与解码方式(internal/ast)；
- 添加 lsp 子命令；
- 代码可通过 build.Output.Version 修改文档中的版本号；
- 新的文档标签 xml-cdata 用于标记该内容在 xml 格式下是否需要以 CDATA 的形式展示；
- 添加地命名空间的支持；
- 添加全局的 XMLNamespaces 字段，用于指定可用的命名空间；
- build.Output 添加了两个配置项：Namespace 和 NamespacePrefix ；
- @xml-wrapped 添加两种语法表示；
- 重新添加 Pack 和 Unpack，用于打包和解包文档内容；
- 数据类型现在可以指定子类型，比如 `string.email` 表示邮箱地址，也必定是字符串；

### Changed

- 可以通过 SetLocale 随时设置本地化信息，同时也添加一系列相关的操作函数；
- Locales() 现在返回一个 map 副本；
- 去掉了 core.Block.Raw 字段；
- 分离程序版本和文档版本；

### Fixed

- 解决 View 返回的内容无法找到 xsl 文件的错误；

### Changed

- 新增 build 包，统一处理文档的提取与输出，去除了 input、output 和 config.go 的相关内容；

### Removed

- 删除 xml-ns 元素；

## [v6.0.1]

### Fixed

- 修正 Chrome 与 Safari 无法正确显示文档的错误；
- 修正命令行 `apidoc static` 导致 panic 的错误；

## [v6.0.0]

### Added

- 添加 Valid 方法，用于验证文档是否正确；
- 为 Request 和 Param 添加了五个用于描述 XML 数据的字段；
- 为 apidoc 元素添加 mimetype，用于指定全局可用的 mimetype 值；
- 为 Request 和 Param 添加 attr 字段，用于在 content-type 为 xml 的请求中表示属性；
- 添加 mock 子命令，用于生成 mock 数据的功能；
- 添加 static 子命令，用于打开查看文档的服务；
- 添加导出为 openapi 的功能；

### Changed

- 现在使用空值表示 doc.None，而不是字符串 `none`；
- 解析文档时，如果类型为枚举，则会验证枚举类型是否与其关联的类型兼容；

## [v5.2.1]

### Changed

- favicon 现在只支持 SVG 格式的图片；

### Fixed

- 修正 Pack 可能将二进制等文件进行打包的错误；

## [v5.2.0]

### Added

- 添加 Test 全局函数，用于对文档语言进行测试；
- 添加 Config 对象，用于处理配置文件；
- 添加 Pack 系列函数，用于将内容打包成一个文件；
- 添加 static 包，管理所有与 docs 下静态文件相关的功能；
- 添加 static.Type 类型，用于指定 static.Pack 和 static.FolderHandler 的文件类型；
- 添加导出为 openapi 的功能；

### Changed

- Make、MakeBuffer 标记为过时函数，不再推荐使用；

## [v5.1.0]

### Added

- 添加保存文档至内存的功能，Go 项目集成更加方便；

## [v5.0.0]

### Changed

- 改为 XML 作为文档的格式；
