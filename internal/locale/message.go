// SPDX-License-Identifier: MIT

package locale

// 各个语言需要翻译的所有字符串
const (
	// 与 flag 包相关的处理
	CmdUsage = `%s 是一个 RESTful API 文档生成工具

用法：
apidoc cmd [args]

cmd 为子命令，args 为传递给子命令的参数，目前支持以下子命令：
%s

源代码采用 MIT 开源许可证，发布于 %s
详细信息可访问官网 %s`
	CmdHelpUsage    = "显示帮助信息"
	CmdVersionUsage = "显示版本信息"
	CmdLangUsage    = "显示所有支持的语言"
	CmdLocaleUsage  = "显示所有支持的本地化内容"
	CmdDetectUsage  = `根据目录下的内容生成配置文件

用法：
apidoc detect [options] [path]

options 可以是以下参数：
%s

path 表示目录的路径，或不指定，表示使用当前工作目录 ./ 代替。`
	CmdTestUsage = "测试语法的正确性"
	CmdMockUsage = `启用 mock 服务

用法：
apidoc mock [options] [path]

options 可以是以下参数：
%s

path 表示文档路径，或不指定，则使用当前工作目录 ./ 代替。`
	CmdBuildUsage  = "生成文档内容"
	CmdStaticUsage = `启用静态文件服务

用法：
apidoc static [options] [path]

options 可以是以下参数：
%s

path 表示需要展示的文档路径，为空表示没有需要展示的文档。`
	CmdLSPUsage = `启动 language server protocol 服务

用法：
apidoc lsp [options]

options 可以是以下参数
%s`
	Version                    = "版本：%s\n文档：%s\n提交：%s\nLSP：%s\nGo：%s"
	CmdNotFound                = "子命令 %s 未找到\n"
	FlagMockPortUsage          = "指定 mock 服务的端口号"
	FlagMockServersUsage       = "指定 mock 服务时，文档中 server 变量对应的路由前缀"
	FlagDetectRecursive        = "detect 子命令是否检测子目录的值"
	FlagStaticPortUsage        = "指定 static 服务的端口号"
	FlagStaticDocsUsage        = "指定 static 服务的文件夹"
	FlagStaticStylesheetUsage  = "指定 static 是否只启用样式文件内容"
	FlagStaticContentTypeUsage = "指定 static 的 content-type 值，不指定，则根据扩展名自动获取"
	FlagStaticURLUsage         = "指定 static 服务中文档的输出地址"
	FlagLSPPortUsage           = "指定 language server protocol 服务的端口号"
	FlagLSPModeUsage           = "指定 language server protocol 的运行方式，可以是 websocket、tcp 和 udp"
	FlagLSPHeaderUsage         = "指定 language server protocol 传递内容是否带报头信息"

	VersionInCompatible = "当前程序与配置文件中指定的版本号不兼容"
	Complete            = "完成！文档保存在：%s，总用时：%v"
	ConfigWriteSuccess  = "配置内容成功写入 %s"
	TestSuccess         = "语法没有问题！"
	LangID              = "ID"
	LangName            = "名称"
	LangExts            = "扩展名"
	LoadAPI             = "加载 API：%s %s"
	RequestAPI          = "访问 API：%s %s"
	DeprecatedWarn      = "%s %s 将于 %s 被废弃"
	GeneratorBy         = "当前文档由 %s 生成"
	ServerStart         = "服务启动，可通过 %s 访问"
	RequestRPC          = "访问 RPC：%s"

	// 文档树中各个字段的介绍
	UsageAPIDoc            = "usage-apidoc"
	UsageAPIDocAPIDoc      = "usage-apidoc-apidoc"
	UsageAPIDocLang        = "usage-apidoc-lang"
	UsageAPIDocLogo        = "usage-apidoc-logo"
	UsageAPIDocCreated     = "usage-apidoc-created"
	UsageAPIDocVersion     = "usage-apidoc-version"
	UsageAPIDocTitle       = "usage-apidoc-title"
	UsageAPIDocDescription = "usage-apidoc-description"
	UsageAPIDocContact     = "usage-apidoc-contact"
	UsageAPIDocLicense     = "usage-apidoc-license"
	UsageAPIDocTags        = "usage-apidoc-tags"
	UsageAPIDocServers     = "usage-apidoc-servers"
	UsageAPIDocAPIs        = "usage-apidoc-apis"
	UsageAPIDocResponses   = "usage-apidoc-responses"
	UsageAPIDocMimetypes   = "usage-apidoc-mimetypes"

	UsageAPI            = "usage-api"
	UsageAPIVersion     = "usage-api-version"
	UsageAPIMethod      = "usage-api-method"
	UsageAPIID          = "usage-api-id"
	UsageAPIPath        = "usage-api-path"
	UsageAPISummary     = "usage-api-summary"
	UsageAPIDescription = "usage-api-description"
	UsageAPIRequests    = "usage-api-requests"
	UsageAPIResponses   = "usage-api-responses"
	UsageAPICallback    = "usage-api-callback"
	UsageAPIDeprecated  = "usage-api-deprecated"
	UsageAPIHeaders     = "usage-api-headers"
	UsageAPITags        = "usage-api-tags"
	UsageAPIServers     = "usage-api-servers"

	UsageLink     = "usage-link"
	UsageLinkText = "usage-link-text"
	UsageLinkURL  = "usage-link-url"

	UsageContact      = "usage-contact"
	UsageContactName  = "usage-contact-name"
	UsageContactURL   = "usage-contact-url"
	UsageContactEmail = "usage-contact-email"

	UsageCallback            = "usage-callback"
	UsageCallbackMethod      = "usage-callback-method"
	UsageCallbackPath        = "usage-callback-path"
	UsageCallbackSummary     = "usage-callback-summary"
	UsageCallbackDescription = "usage-callback-description"
	UsageCallbackResponses   = "usage-callback-responses"
	UsageCallbackRequests    = "usage-callback-requests"
	UsageCallbackHeaders     = "usage-callback-headers"

	UsageEnum           = "usage-enum"
	UsageEnumDeprecated = "usage-enum-deprecated"
	UsageEnumValue      = "usage-enum-value"
	UsageEnumSummary    = "usage-enum-summary"

	UsageExample         = "usage-example"
	UsageExampleMimetype = "usage-example-mimetype"
	UsageExampleSummary  = "usage-example-summary"
	UsageExampleContent  = "usage-example-content"

	UsageParam            = "usage-param"
	UsageParamXMLAttr     = "usage-param-xml-attr"
	UsageParamXMLExtract  = "usage-param-xml-extract"
	UsageParamXMLNS       = "usage-param-xml-ns"
	UsageParamXMLNSPrefix = "usage-param-xml-ns-prefix"
	UsageParamXMLWrapped  = "usage-param-xml-wrapped"
	UsageParamName        = "usage-param-name"
	UsageParamType        = "usage-param-type"
	UsageParamDeprecated  = "usage-param-deprecated"
	UsageParamDefault     = "usage-param-default"
	UsageParamOptional    = "usage-param-optional"
	UsageParamArray       = "usage-param-array"
	UsageParamItems       = "usage-param-items"
	UsageParamSummary     = "usage-param-summary"
	UsageParamEnums       = "usage-param-enums"
	UsageParamDescription = "usage-param-description"
	UsageParamArrayStyle  = "usage-param-array-style"

	UsagePath        = "usage-path"
	UsagePathPath    = "usage-path-path"
	UsagePathParams  = "usage-path-params"
	UsagePathQueries = "usage-path-queries"

	UsageRequest            = "usage-request"
	UsageRequestXMLAttr     = "usage-request-xml-attr"
	UsageRequestXMLExtract  = "usage-request-xml-extract"
	UsageRequestXMLNS       = "usage-request-xml-ns"
	UsageRequestXMLNSPrefix = "usage-request-xml-ns-prefix"
	UsageRequestXMLWrapped  = "usage-request-xml-wrapped"
	UsageRequestName        = "usage-request-name"
	UsageRequestType        = "usage-request-type"
	UsageRequestDeprecated  = "usage-request-deprecated"
	UsageRequestArray       = "usage-request-array"
	UsageRequestItems       = "usage-request-items"
	UsageRequestSummary     = "usage-request-summary"
	UsageRequestEnums       = "usage-request-enums"
	UsageRequestDescription = "usage-request-description"
	UsageRequestMimetype    = "usage-request-mimetype"
	UsageRequestExamples    = "usage-request-examples"
	UsageRequestHeaders     = "usage-request-headers"

	UsageRichtext     = "usage-richtext"
	UsageRichtextType = "usage-richtext-type"
	UsageRichtextText = "usage-richtext-text"

	UsageTag           = "usage-tag"
	UsageTagName       = "usage-tag-id"
	UsageTagTitle      = "usage-tag-title"
	UsageTagDeprecated = "usage-tag-deprecated"

	UsageServer            = "usage-server"
	UsageServerName        = "usage-server-name"
	UsageServerTitle       = "usage-server-title"
	UsageServerDeprecated  = "usage-server-deprecated"
	UsageServerSummary     = "usage-server-summary"
	UsageServerDescription = "usage-server-description"

	// 错误信息，可能在地方用到
	ErrInvalidUTF8Character      = "无效的 UTF8 字符"
	ErrInvalidURIScheme          = "无效的 URI 协议"
	ErrInvalidXML                = "无效的 XML 文档"
	ErrIsNotAPIDoc               = "并非有效的 apidoc 的文档格式"
	ErrInvalidContentTypeCharset = "报头 ContentType 中指定的字符集无效 "
	ErrInvalidContentLength      = "报头 ContentLength 无效"
	ErrBodyIsEmpty               = "请求的报文为空"
	ErrInvalidHeaderFormat       = "无效的报头格式"
	ErrRequired                  = "不能为空"
	ErrInvalidFormat             = "格式不正确"
	ErrDirNotExists              = "目录不存在"
	ErrNotFoundEndFlag           = "找不到结束符号"
	ErrNotFoundSupportedLang     = "该目录下没有支持的语言文件"
	ErrDirIsEmpty                = "目录下没有需要解析的文件"
	ErrInvalidValue              = "无效的值"
	ErrPathNotMatchParams        = "地址参数不匹配"
	ErrDuplicateValue            = "重复的值"
	ErrMessage                   = "%s 位于 %s"
	ErrNotFound                  = "未找到该值"
	ErrReadRemoteFile            = "读取远程文件 %s 时返回状态码 %d"
	ErrServerNotInitialized      = "服务未初始化"
	ErrInvalidLSPState           = "无效的 LSP 状态"

	// logs
	InfoPrefix    = "[INFO] "
	WarnPrefix    = "[WARN] "
	ErrorPrefix   = "[ERRO] "
	SuccessPrefix = "[SUCC] "
)
