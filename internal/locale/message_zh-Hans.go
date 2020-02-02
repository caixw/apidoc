// SPDX-License-Identifier: MIT

package locale

import "golang.org/x/text/language"

var zhHans = map[string]string{
	// 与 flag 包相关的处理
	CmdUsage: `%s 是一个 RESTful API 文档生成工具

用法：
apidoc cmd [args]

cmd 为子命令，args 为传递给子命令的参数，目前支持以下子命令：
%s

源代码采用 MIT 开源许可证，发布于 %s
详细信息可访问官网 %s`,
	CmdHelpUsage:    "显示帮助信息",
	CmdVersionUsage: "显示版本信息",
	CmdLangUsage:    "显示所有支持的语言",
	CmdLocaleUsage:  "显示所有支持的本地化内容",
	CmdDetectUsage: `根据目录下的内容生成配置文件

用法：
apidoc detect [options] [path]

options 可以是以下参数：
%s

path 表示目录的路径，或不指定，表示使用当前工作目录 ./ 代替。`,
	CmdTestUsage: "测试语法的正确性",
	CmdMockUsage: `启用 mock 服务

用法：
apidoc mock [options] [path]

options 可以是以下参数：
%s

path 表示文档路径，或不指定，则使用当前工作目录 ./ 代替。`,
	CmdBuildUsage: "生成文档内容",
	CmdStaticUsage: `启用静态文件服务

用法：
apidoc static [options] [path]

options 可以是以下参数：
%s

path 表示需要展示的文档路径，为空表示没有需要展示的文档。`,
	CmdLSPUsage: `启动 language server protocol 服务

用法：
apidoc lsp [options]

options 可以是以下参数
%s`,
	Version:                    "版本：%s\n文档：%s\n提交：%s\nGo：%s",
	CmdNotFound:                "子命令 %s 未找到\n",
	FlagMockPortUsage:          "指定 mock 服务的端口号",
	FlagMockServersUsage:       "指定 mock 服务时，文档中 server 变量对应的路由前缀",
	FlagDetectRecursive:        "detect 子命令是否检测子目录的值",
	FlagStaticPortUsage:        "指定 static 服务的端口号",
	FlagStaticDocsUsage:        "指定 static 服务的文件夹",
	FlagStaticStylesheetUsage:  "指定 static 是否只启用样式文件内容",
	FlagStaticContentTypeUsage: "指定 static 的 content-type 值，不指定，则根据扩展名自动获取",
	FlagStaticURLUsage:         "指定 static 服务中文档的输出地址",
	FlagLSPPortUsage:           "指定 language server protocol 服务的端口号",
	FlagLSPModeUsage:           "指定 language server protocol 的运行方式，可以是 http、websocket、tcp 和 udp",

	VersionInCompatible: "当前程序与配置文件中指定的版本号不兼容",
	Complete:            "完成！文档保存在：%s，总用时：%v",
	ConfigWriteSuccess:  "配置内容成功写入 %s",
	TestSuccess:         "语法没有问题！",
	LangID:              "ID",
	LangName:            "名称",
	LangExts:            "扩展名",
	LoadAPI:             "加载 API：%s %s",
	RequestAPI:          "访问 API：%s %s",
	DeprecatedWarn:      "%s %s 将于 %s 被废弃",
	GeneratorBy:         "当前文档由 %s 生成",
	ServerStart:         "服务启动，可通过 %s 访问",

	// 文档树中各个字段的介绍
	UsageAPIDoc:            "用于描述整个文档的相关内容，只能出现一次。",
	UsageAPIDocAPIDoc:      "文档的版本要号",
	UsageAPIDocLang:        "文档内容的本地化 ID，比如 <samp>zh-Hans</samp>、<samp>en-US</samp> 等。",
	UsageAPIDocLogo:        "文档的图标，仅可使用 SVG 格式图标。",
	UsageAPIDocCreated:     "文档的创建时间",
	UsageAPIDocVersion:     "文档的版本号，需要遵守 semver 的约定。",
	UsageAPIDocTitle:       "文档的标题",
	UsageAPIDocDescription: "文档的整体描述内容",
	UsageAPIDocContact:     "文档作者的联系方式",
	UsageAPIDocLicense:     "文档的版权信息",
	UsageAPIDocTags:        "文档中定义的所有标签",
	UsageAPIDocServers:     "API 基地址列表，每个 API 最少应该有一个 server。",
	UsageAPIDocAPIs:        "文档中的 API 文档",
	UsageAPIDocResponses:   "文档中所有 API 文档都需要支持的返回内容",
	UsageAPIDocMimetypes:   "文档所支持的 mimetype",

	UsageAPI:            "用于定义单个 API 接口的具体内容",
	UsageAPIVersion:     "表示此接口在该版本中添加",
	UsageAPIMethod:      "当前接口所支持的请求方法",
	UsageAPIID:          "接口的唯一 ID",
	UsageAPIPath:        "定义路径信息",
	UsageAPISummary:     "简要介绍",
	UsageAPIDescription: "该接口的详细介绍，为 HTML 内容。",
	UsageAPIRequests:    "定义可用的请求信息",
	UsageAPIResponses:   "定义可能的返回信息",
	UsageAPICallback:    "定义回调接口内容",
	UsageAPIDeprecated:  "在此版本之后将会被弃用",
	UsageAPIHeaders:     "传递的报头内容，如果是某个 mimetype 专用的，可以放在 request 元素中。",
	UsageAPITags:        "关联的标签",
	UsageAPIServers:     "关联的服务",

	UsageLink:     "用于描述链接信息，一般转换为 HTML 的 a 标签。",
	UsageLinkText: "链接的字面文字",
	UsageLinkURL:  "链接指向的文本",

	UsageContact:      "用于描述联系方式",
	UsageContactName:  "联系人的名称",
	UsageContactURL:   "联系人的 URL",
	UsageContactEmail: "联系人的电子邮件",

	UsageCallback:            "定义回调信息",
	UsageCallbackMethod:      "回调的请求方法",
	UsageCallbackPath:        "定义回调的地址",
	UsageCallbackSummary:     "简要介绍",
	UsageCallbackDescription: "该接口的详细介绍",
	UsageCallbackResponses:   "定义可能的返回信息",
	UsageCallbackRequests:    "定义可用的请求信息",
	UsageCallbackHeaders:     "传递的报头内容",

	UsageEnum:           "定义枚举类型的数所的枚举值",
	UsageEnumDeprecated: "该属性弃用的版本号",
	UsageEnumValue:      "枚举值",
	UsageEnumSummary:    "枚举值的说明",

	UsageExample:         "示例代码",
	UsageExampleMimetype: "特定于类型的示例代码",
	UsageExampleSummary:  "示例代码的概要信息",
	UsageExampleContent:  "示例代码的内容，需要使用 CDATA 包含代码。",

	UsageParam:            "参数类型，基本上可以作为 request 的子集使用。",
	UsageParamXMLAttr:     "是否作为父元素的属性，仅作用于 XML 元素。是否作为父元素的属性，仅用于 XML 的请求。",
	UsageParamXMLExtract:  "将当前元素的内容作为父元素的内容，要求父元素必须为 <var>object</var>。",
	UsageParamXMLNS:       "XML 标签的命名空间",
	UsageParamXMLNSPrefix: "XML 标签的命名空间名称前缀",
	UsageParamXMLWrapped:  "如果当前元素的 @array 为 <var>true</var>，是否将其包含在 wrapped 指定的标签中。",
	UsageParamName:        "值的名称",
	UsageParamType:        "值的类型，可以是 <var>string</var>、<var>number</var>、<var>bool</var> 和 <var>object</var>",
	UsageParamDeprecated:  "表示在大于等于该版本号时不再启作用",
	UsageParamDefault:     "默认值",
	UsageParamOptional:    "是否为可选的参数",
	UsageParamArray:       "是否为数组",
	UsageParamItems:       "子类型，比如对象的子元素。",
	UsageParamSummary:     "简要介绍",
	UsageParamEnums:       "当前参数可用的枚举值",
	UsageParamDescription: "详细介绍，为 HTML 内容。",
	UsageParamArrayStyle:  "以数组的方式展示数据",

	UsagePath:        "用于定义请求时与路径相关的内容",
	UsagePathPath:    "接口地址",
	UsagePathParams:  "地址中的参数",
	UsagePathQueries: "地址中的查询参数",

	UsageRequest:            "定义了请求和返回的相关内容",
	UsageRequestXMLAttr:     "是否作为父元素的属性，仅作用于 XML 元素。是否作为父元素的属性，仅用于 XML 的请求。",
	UsageRequestXMLExtract:  "将当前元素的内容作为父元素的内容，要求父元素必须为 <var>object</var>。",
	UsageRequestXMLNS:       "XML 标签的命名空间",
	UsageRequestXMLNSPrefix: "XML 标签的命名空间名称前缀",
	UsageRequestXMLWrapped:  "如果当前元素的 @array 为 <var>true</var>，是否将其包含在 wrapped 指定的标签中。",
	UsageRequestName:        "当 mimetype 为 <var>application/xml</var> 时，此值表示 XML 的顶层元素名称，否则无用。",
	UsageRequestType:        "值的类型，可以是 <var>none</var>、<var>string</var>、<var>number</var>、<var>bool</var>、<var>object</var> 和空值；空值表示不输出任何内容。",
	UsageRequestDeprecated:  "表示在大于等于该版本号时不再启作用",
	UsageRequestArray:       "是否为数组",
	UsageRequestItems:       "子类型，比如对象的子元素。",
	UsageRequestSummary:     "简要介绍",
	UsageRequestEnums:       "当前参数可用的枚举值",
	UsageRequestDescription: "详细介绍，为 HTML 内容。",
	UsageRequestMimetype:    "媒体类型，比如 application/json 等。",
	UsageRequestExamples:    "示例代码",
	UsageRequestHeaders:     "传递的报头内容",

	UsageRichtext:     "富文本内容",
	UsageRichtextType: "指定内容的格式",
	UsageRichtextText: "富文本的实际内容",

	UsageTag:           "用于对各个 API 进行分类",
	UsageTagName:       "标签的唯一 ID",
	UsageTagTitle:      "标签的字面名称",
	UsageTagDeprecated: "该标签在大于该版本时被弃用",

	UsageServer:            "用于指定各个 API 的服务器地址",
	UsageServerName:        "服务唯一 ID",
	UsageServerTitle:       "服务的字面名称",
	UsageServerDeprecated:  "服务在大于该版本时被弃用",
	UsageServerSummary:     "服务的摘要信息",
	UsageServerDescription: "服务的详细描述",

	// 错误信息，可能在地方用到
	ErrInvalidUTF8Character:      "无效的 UTF8 字符",
	ErrInvalidURIScheme:          "无效的 uri 协议",
	ErrInvalidXML:                "无效的 XML 文档",
	ErrIsNotAPIDoc:               "并非有效的 apidoc 的文档格式",
	ErrInvalidContentTypeCharset: "报头 ContentType 中指定的字符集无效 ",
	ErrInvalidContentLength:      "报头 ContentLength 无效",
	ErrBodyIsEmpty:               "请求的报文为空",
	ErrInvalidHeaderFormat:       "无效的报头格式",
	ErrRequired:                  "不能为空",
	ErrInvalidFormat:             "格式不正确",
	ErrDirNotExists:              "目录不存在",
	ErrNotFoundEndFlag:           "找不到结束符号",
	ErrNotFoundSupportedLang:     "该目录下没有支持的语言文件",
	ErrDirIsEmpty:                "目录下没有需要解析的文件",
	ErrInvalidValue:              "无效的值",
	ErrPathNotMatchParams:        "地址参数不匹配",
	ErrDuplicateValue:            "重复的值",
	ErrMessage:                   "%s 位于 %s",
	ErrNotFound:                  "未找到该值",
	ErrReadRemoteFile:            "读取远程文件 %s 时返回状态码 %d",
	ErrServerNotInitialized:      "服务未初始化",

	// logs
	InfoPrefix:    "[信息] ",
	WarnPrefix:    "[警告] ",
	ErrorPrefix:   "[错误] ",
	SuccessPrefix: "[成功] ",
}

func init() {
	addLocale(language.MustParse("zh-Hans"), zhHans)

	// 大部分的系统都采用 zh-cn 作为语言标记，
	// 但是 golang.org/x/text 现在不能将 zh-cn 自动转换成 zh-hans
	addLocale(language.MustParse("zh-cn"), zhHans)
}
