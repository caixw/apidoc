// SPDX-License-Identifier: MIT

package locale

var cmnHans = map[string]string{
	// 与 flag 包相关的处理
	CmdUsage: "%s 是一个 RESTful API 文档生成工具\n",
	CmdUsageFooter: `详细文档可访问官网 %s
源码以 MIT 许可发布于 %s
`,
	CmdUsageOptions:  "选项：",
	CmdUsageCommands: "子命令：",
	CmdHelpUsage:     "显示帮助信息\n",
	CmdVersionUsage:  "显示版本信息\n",
	CmdLangUsage:     "显示所有支持的语言\n",
	CmdLocaleUsage:   "显示所有支持的本地化内容\n",
	CmdDetectUsage:   "根据目录下的内容生成配置文件\n",
	CmdSyntaxUsage:   "测试语法的正确性\n",
	CmdMockUsage: `启用 mock 服务

mock 服务会根据接口定义检测用户提交的数据是否合法，并生成随机的数据返回给用户。
对于数据只作检测是否合规，但是无法理解其内容，比如提交地址中添加了 size=20，
只会检测 20 的类型是否符合 size 的要求，但是不会只返回给用户 20 条数据。
`,
	CmdBuildUsage:  "生成文档内容\n",
	CmdStaticUsage: "启用静态文件服务\n",
	CmdLSPUsage:    "启动 language server protocol 服务\n",
	Version:        "版本：%s\n文档：%s\nLSP：%s\nopenapi：%s\nGo：%s",
	CmdNotFound:    "子命令 %s 未找到\n",

	FlagSyntaxDirUsage:         "以 `URI` 形式表示测试项目地址",
	FlagBuildDirUsage:          "以 `URI` 形式表示的项目地址",
	FlagMockPortUsage:          "指定 mock 服务的端口号",
	FlagMockServersUsage:       "指定 mock 服务时，文档中 server 名对应的路由前缀。",
	FlagMockIndentUsage:        "指定缩进内容",
	FlagMockSliceSizeUsage:     "生成数组大小的范围",
	FlagMockNumSliceUsage:      "生成数值类型的数据时的数值范围",
	FlagMockNumFloatUsage:      "生成的数值是否允许有浮点数存在",
	FlagMockPathUsage:          "指定文档的 `URI` 格式路径，根据此文档的内容生成 mock 数据。",
	FlagMockStringSizeUsage:    "生成字符串类型数据时字符串的长度范围",
	FlagMockStringAlphaUsage:   "生成的字符串中允许出现的字符",
	FlagMockUsernameSizeUsage:  "生成邮箱地址时，用户名的长度范围。",
	FlagMockEmailDomainsUsage:  "生成邮箱地址时所可用的域名列表，多个用半角逗号分隔。",
	FlagMockURLDomainsUsage:    "生成 URL 地址时所可用的域名列表，多个用半角逗号分隔。",
	FlagDetectRecursiveUsage:   "detect 子命令是否检测子目录的值",
	FlagDetectDirUsage:         "以 `URI` 形式表示检测项目地址",
	FlagDetectWrite:            "是否将配置内容写入文件，如果为 true，会将配置内容写入检测目录下的 .apidoc.yaml 文件。",
	FlagStaticPortUsage:        "指定 static 服务的端口号",
	FlagStaticDocsUsage:        "指定 static 服务静态文件所在的 `URI`",
	FlagStaticStylesheetUsage:  "指定 static 是否只启用样式文件内容",
	FlagStaticContentTypeUsage: "指定 static 的 content-type 值，不指定，则根据扩展名自动获取",
	FlagStaticURLUsage:         "指定 static 服务中文档的输出地址",
	FlagStaticPathUsage:        "指定 static 服务 `URI` 格式的文档路径，如果未指定，则不生成相关的文档内容。",
	FlagLSPPortUsage:           "指定 LSP 服务的端口号",
	FlagLSPModeUsage:           "指定 LSP 的运行方式，可以是 websocket、tcp 和 udp",
	FlagLSPHeaderUsage:         "指定 LSP 传递内容是否带报头信息",

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
	RequestRPC:          "访问 RPC：%s",
	UnimplementedRPC:    "未实现该 RPC 服务 %s",
	PackFileHeader:      "文档由 %s 自动生成，请勿手动修改！",

	// 文档树中各个字段的介绍
	UsageAPIDoc:              "用于描述整个文档的相关内容，只能出现一次。",
	UsageAPIDocAPIDoc:        "文档的版本要号",
	UsageAPIDocLang:          "文档内容的本地化 ID，比如 <var>zh-Hans</var>、<var>en-US</var> 等。",
	UsageAPIDocLogo:          "文档的图标，仅可使用 SVG 格式图标。",
	UsageAPIDocCreated:       "文档的创建时间",
	UsageAPIDocVersion:       "文档的版本号",
	UsageAPIDocTitle:         "文档的标题",
	UsageAPIDocDescription:   "文档的整体描述内容",
	UsageAPIDocContact:       "文档作者的联系方式",
	UsageAPIDocLicense:       "文档的版权信息",
	UsageAPIDocTags:          "文档中定义的所有标签",
	UsageAPIDocServers:       "API 基地址列表，每个 API 最少应该有一个 server。",
	UsageAPIDocAPIs:          "文档中的 API 文档",
	UsageAPIDocHeaders:       "文档中所有 API 都包含的公共报头",
	UsageAPIDocResponses:     "文档中所有 API 文档都需要支持的返回内容",
	UsageAPIDocMimetypes:     "文档所支持的 mimetype",
	UsageAPIDocXMLNamespaces: "针对 <var>application/xml</var> 类型的内容的命名空间设置",

	UsageXMLNamespace:       "为 <var>application/xml</var> 定义命名空间的相关属性",
	UsageXMLNamespacePrefix: "命名空间的前缀，如果为空，则表示作为默认命名空间，命局只能有一个默认命名空间。",
	UsageXMLNamespaceURN:    "命名空间的唯一标识，需要全局唯一，且区分大小写。",

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

	UsageLink:     "用于描述链接信息，一般转换为 HTML 的 <code>a</code> 标签。",
	UsageLinkText: "链接的字面文字",
	UsageLinkURL:  "链接指向的文本",

	UsageContact:      "用于描述联系方式",
	UsageContactName:  "联系人的名称",
	UsageContactURL:   "联系人的 URL",
	UsageContactEmail: "联系人的电子邮件",

	UsageCallback:            "定义接口的回调内容",
	UsageCallbackMethod:      "回调的请求方法",
	UsageCallbackPath:        "回调的请求地址",
	UsageCallbackSummary:     "简要介绍",
	UsageCallbackDeprecated:  "在此版本之后将会被弃用",
	UsageCallbackDescription: "对于回调的详细介绍",
	UsageCallbackResponses:   "定义可能的返回信息",
	UsageCallbackRequests:    "定义可用的请求信息",
	UsageCallbackHeaders:     "传递的报头内容",

	UsageEnum:            "定义枚举类型的数所的枚举值",
	UsageEnumDeprecated:  "该属性弃用的版本号",
	UsageEnumValue:       "枚举值",
	UsageEnumSummary:     "枚举值的说明",
	UsageEnumDescription: "枚举值的详细说明",

	UsageExample:         "示例代码",
	UsageExampleMimetype: "特定于类型的示例代码",
	UsageExampleSummary:  "示例代码的概要信息",
	UsageExampleContent:  "示例代码的内容，需要使用 CDATA 包含代码。",

	UsageParam:            "参数类型，基本上可以作为 request 的子集使用。",
	UsageParamName:        "值的名称",
	UsageParamType:        "值的类型",
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
	UsageRequestName:        "当 mimetype 为 <var>application/xml</var> 时，此值表示 XML 的顶层元素名称，否则无用。",
	UsageRequestType:        "值的类型",
	UsageRequestDeprecated:  "表示在大于等于该版本号时不再启作用",
	UsageRequestArray:       "是否为数组",
	UsageRequestItems:       "子类型，比如对象的子元素。",
	UsageRequestSummary:     "简要介绍",
	UsageRequestStatus:      "状态码。在 request 中，该值不可用，否则为必填项。",
	UsageRequestEnums:       "当前参数可用的枚举值",
	UsageRequestDescription: "详细介绍，为 HTML 内容。",
	UsageRequestMimetype:    "媒体类型，比如 <var>application/json</var> 等。",
	UsageRequestExamples:    "示例代码",
	UsageRequestHeaders:     "传递的报头内容",

	UsageRichtext:     "富文本内容",
	UsageRichtextType: "指定富文本内容的格式，目前支持 <var>html</var> 和 <var>markdown</var>。",
	UsageRichtextText: "富文本的实际内容",

	UsageTag:           "用于对各个 API 进行分类",
	UsageTagName:       "标签的唯一 ID",
	UsageTagTitle:      "标签的字面名称",
	UsageTagDeprecated: "该标签在大于该版本时被弃用",

	UsageServer:            "用于指定各个 API 的服务器地址",
	UsageServerName:        "服务唯一 ID",
	UsageServerTitle:       "服务的字面名称",
	UsageServerURL:         "服务的基地址，与该服务关联的 API，访问地址都是相对于此地址的。",
	UsageServerDeprecated:  "服务在大于该版本时被弃用",
	UsageServerSummary:     "服务的摘要信息",
	UsageServerDescription: "服务的详细描述",

	UsageXMLAttr:    "是否作为父元素的属性，仅作用于 XML 元素。是否作为父元素的属性，仅用于 XML 的请求。",
	UsageXMLExtract: "将当前元素的内容作为父元素的内容，要求父元素必须为 <var>object</var>。",
	UsageXMLCData:   "当前内容为 CDATA，与 <code>@xml-attr</code> 互斥。",
	UsageXMLPrefix:  "XML 标签的命名空间名称前缀",
	UsageXMLWrapped: `如果当前元素的 <code>@array</code> 为 <var>true</var>，则可以通过此值指定在 XML 格式中的名称。
	可以有三种格式：<ul>
	<li><samp>name</samp>：表示为数组添加一个父元素名称为 <var>name</var>；</li>
	<li><samp>name1&gt;name2</samp>：表示数组项的名称改为 <var>name2</var>，且添加一个父元素名为 <var>name1</var>；</li>
	<li><samp>&gt;name</samp>：表示将当前数组元素的名称改为 <var>name</var>；</li>
	</ul>`,

	// 基本类型
	UsageString:  "普通的字符串类型，特殊字符需要使用 XML 实体，比如 <samp>&lt;</samp> 需要使用 <samp>&amp;lt;</samp> 代替。",
	UsageNumber:  "普通的数值类型，比如：<samp>1</samp>、<samp>-11.1</samp> 等。",
	UsageBool:    "布尔值类型，取值为 <var>true</var> 或是 <var>false</var>。",
	UsageVersion: `版本号，格式遵守 <a href="https://semver.org/lang/zh-CN/">semver</a> 规则。比如：<samp>1.0.1</samp>、<samp>1.0.1+20200618</samp>。`,
	UsageDate:    `采用 <a href="https://tools.ietf.org/html/rfc3339">RFC3339</a> 格式表示的时间，比如：<samp>2019-12-16T00:35:48+08:00</samp>。`,
	UsageType: `用于表示数据的类型值，格式为 <code>primitive[.subtype]</code>，其中 <code>primitive</code> 为基本类型，而 <code>subtype</code> 为子类型，用于对 <code>primitive</code> 进行进一步的约束，当客户端无法处理整个类型时，可以按照 <code>primitive</code> 的类型处理。<br />
	目前支持以下几种类型：<ul>
	<li>空值；</li>
	<li><var>bool</var> 布尔值；</li>
	<li><var>object</var> 对象；</li>
	<li><var>number</var> 数值类型；</li>
	<li><var>number.int</var> 整数类型的数值；</li>
	<li><var>number.float</var> 浮点类型的数值；</li>
	<li><var>string</var> 字符串；</li>
	<li><var>string.url</var> URL 类型的字符串；</li>
	<li><var>string.email</var> email 类型的字符串；</li>
	</ul>`,

	// 以下是有关 build.Config 的字段说明
	UsageConfigVersion:               "此配置文件的所使用的文档版本",
	UsageConfigInputs:                "指定输入的数据，同一项目只能解析一种语言。",
	UsageConfigInputsLang:            "源文件类型。具体支持的类型可通过 -l 参数进行查找。",
	UsageConfigInputsDir:             "需要解析的源文件所在目录",
	UsageConfigInputsExts:            "只从这些扩展名的文件中查找文档",
	UsageConfigInputsRecursive:       "是否解析子目录下的源文件",
	UsageConfigInputsEncoding:        `编码，默认为 <var>utf-8</var>，值可以是 <a href="https://www.iana.org/assignments/character-sets/character-sets.xhtml">character-sets</a> 中的内容。`,
	UsageConfigOutput:                "控制输出行为",
	UsageConfigOutputType:            "输出的类型，目前可以 <var>apidoc+xml</var>、<var>openapi+json</var> 和 <var>openapi+yaml</var>。",
	UsageConfigOutputPath:            "指定输出的文件名，包含路径信息。",
	UsageConfigOutputTags:            "只输出与这些标签相关联的文档，默认为全部。",
	UsageConfigOutputStyle:           "为 XML 文件指定的 XSL 文件",
	UsageConfigOutputNamespace:       "是否输出命名空间",
	UsageConfigOutputNamespacePrefix: "如果输出了命名空间，还可以指定命名空间前缀。",

	// 错误信息，可能在地方用到
	ErrInvalidUTF8Character:      "无效的 UTF8 字符",
	ErrInvalidXML:                "无效的 XML 文档",
	ErrIsNotAPIDoc:               "并非有效的 apidoc 的文档格式",
	ErrInvalidContentTypeCharset: "报头 ContentType 中指定的字符集无效",
	ErrInvalidContentLength:      "报头 ContentLength 无效",
	ErrBodyIsEmpty:               "请求的报文为空",
	ErrInvalidHeaderFormat:       "无效的报头格式",
	ErrRequired:                  "不能为空",
	ErrInvalidFormat:             "格式不正确",
	ErrDirNotExists:              "目录不存在",
	ErrNotFoundEndFlag:           "找不到结束符号",
	ErrNotFoundEndTag:            "找不到结束标签",
	ErrNotFoundSupportedLang:     "该目录下没有支持的语言文件",
	ErrDirIsEmpty:                "目录下没有需要解析的文件",
	ErrInvalidValue:              "无效的值",
	ErrPathNotMatchParams:        "地址参数不匹配",
	ErrDuplicateValue:            "重复的值",
	ErrMessage:                   "%s 位于 %s",
	ErrNotFound:                  "未找到该值",
	ErrReadRemoteFile:            "读取远程文件 %s 时返回状态码 %d",
	ErrServerNotInitialized:      "服务未初始化",
	ErrInvalidLSPState:           "无效的 LSP 状态",
	ErrInvalidURIScheme:          "无效的 URI 协议",
	ErrFileNotFound:              "未找到文件 %s",

	// logs
	InfoPrefix:    "[信息] ",
	WarnPrefix:    "[警告] ",
	ErrorPrefix:   "[错误] ",
	SuccessPrefix: "[成功] ",
}
