# CHANGELOG

## [Unreleased]

### Added

- 添加 Valid 方法，用于验证文档是否正确；
- 为 Request 和 Param 添加了五个用于描述 XML 数据的字段；
- 为 apidoc 元素添加 mimetype，用于指定全局可用的 mimetype 值；

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

### Changed

- Make、MakeBuffer 标记为过时函数，不再推荐使用；

## [v5.1.0]

### Added

- 添加保存文档至内存的功能，Go 项目集成更加方便；

## [v5.0.0]

### Changed

- 改为 XML 作为文档的格式；
