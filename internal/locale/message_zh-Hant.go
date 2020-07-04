// SPDX-License-Identifier: MIT

package locale

var cmnHant = map[string]string{
	// 與 flag 包相關的處理
	CmdUsage: "%s 是壹個 RESTful API 文檔生成工具\n",
	CmdUsageFooter: `詳細文檔可訪問官網 %s
源碼以 MIT 許可發布於 %s
`,
	CmdUsageOptions:  "選項：",
	CmdUsageCommands: "子命令：",
	CmdHelpUsage:     "顯示幫助信息\n",
	CmdVersionUsage:  "顯示版本信息\n",
	CmdLangUsage:     "顯示所有支持的語言\n",
	CmdLocaleUsage:   "顯示所有支持的本地化內容\n",
	CmdDetectUsage:   "根據目錄下的內容生成配置文件\n",
	CmdSyntaxUsage:   "測試語法的正確性\n",
	CmdMockUsage: `啟用 mock 服務

mock 服務會根據接口定義檢測用戶提交的數據是否合法，並生成隨機的數據返回給用戶。
對於數據只作檢測是否合規，但是無法理解其內容，比如提交地址中添加了 size=20，
只會檢測 20 的類型是否符合 size 的要求，但是不會只返回給用戶 20 條數據。
`,
	CmdBuildUsage:  "生成文檔內容\n",
	CmdStaticUsage: "啟用靜態文件服務\n",
	CmdLSPUsage:    "啟動 language server protocol 服務\n",
	Version:        "版本：%s\n文檔：%s\nLSP：%s\nopenapi：%s\nGo：%s",
	CmdNotFound:    "子命令 %s 未找到\n",

	FlagSyntaxDirUsage:         "以 `URI` 形式表示的測試項目地址",
	FlagBuildDirUsage:          "以 `URI` 形式表示的項目地址",
	FlagMockPortUsage:          "指定 mock 服務的端口號",
	FlagMockServersUsage:       "指定 mock 服務時，文檔中 server 名對應的路由前綴。",
	FlagMockIndentUsage:        "指定縮進內容",
	FlagMockSliceSizeUsage:     "生成數組大小的範圍，格式為 [min,max]。",
	FlagMockNumSliceUsage:      "生成數值類型的數據時的數值範圍，格式為 [min,max]。",
	FlagMockNumFloatUsage:      "生成的數值是否允許有浮點數存在",
	FlagMockPathUsage:          "指定文檔的 `URI` 格式路徑，根據此文檔的內容生成 mock 數據。",
	FlagMockStringSizeUsage:    "生成字符串類型數據時字符串的長度範圍，格式為 [min,max]。",
	FlagMockStringAlphaUsage:   "生成的字符串中允許出現的字符",
	FlagMockUsernameSizeUsage:  "生成郵箱地址時，用戶名的長度範圍，格式為 [min,max]。",
	FlagMockEmailDomainsUsage:  "生成郵箱地址時所可用的域名列表，多個用半角逗號分隔。",
	FlagMockURLDomainsUsage:    "生成 URL 地址時所可用的域名列表，多個用半角逗號分隔。",
	FlagMockImagePrefixUsage:   "生成圖片類型數據的基地址",
	FlagMockDateRangeUsage:     "生成可用的日期範圍，格式為 [start,end]，start 和 end 均為 RFC3339 格式。",
	FlagDetectRecursiveUsage:   "detect 子命令是否檢測子目錄的值",
	FlagDetectDirUsage:         "以 `URI` 形式表示的檢測項目地址",
	FlagDetectWrite:            "是否將配置內容寫入文件，如果為 true，會將配置內容寫入檢測目錄下的 .apidoc.yaml 文件。",
	FlagStaticPortUsage:        "指定 static 服務的端口號",
	FlagStaticDocsUsage:        "指定 static 服務靜態文件所在的 `URI`",
	FlagStaticStylesheetUsage:  "指定 static 是否只啟用樣式文件內容",
	FlagStaticContentTypeUsage: "指定 static 的 content-type 值，不指定，則根據擴展名自動獲取",
	FlagStaticURLUsage:         "指定 static 服務中文檔的輸出地址",
	FlagStaticPathUsage:        "指定 static 服務 `URI` 格式的文檔路徑，如果未指定，則不生成相關的文檔內容。",
	FlagLSPPortUsage:           "指定 LSP 服務的端口號",
	FlagLSPModeUsage:           "指定 LSP 的運行方式，可以是 websocket、tcp 和 udp",
	FlagLSPHeaderUsage:         "指定 LSP 傳遞內容是否帶報頭信息",
	FlagVersionKindUsage:       "只顯示該類型的版本號，可以是 apidoc、doc、lsp、openapi 和 all",

	VersionInCompatible: "當前程序與配置文件中指定的版本號不兼容",
	Complete:            "完成！文檔保存在：%s，總用時：%v",
	ConfigWriteSuccess:  "配置內容成功寫入 %s",
	TestSuccess:         "語法沒有問題！",
	LangID:              "ID",
	LangName:            "名稱",
	LangExts:            "擴展名",
	LoadAPI:             "加載 API：%s %s",
	RequestAPI:          "訪問 API：%s %s",
	DeprecatedWarn:      "%s %s 將於 %s 被廢棄",
	GeneratorBy:         "當前文檔由 %s 生成",
	ServerStart:         "服務啟動，可通過 %s 訪問",
	RequestRPC:          "訪問 RPC：%s",
	UnimplementedRPC:    "未實現該 RPC 服務 %s",
	PackFileHeader:      "文檔由 %s 自動生成，請勿手動修改！",

	// 文檔樹中各個字段的介紹
	UsageAPIDoc:              "用於描述整個文檔的相關內容，只能出現壹次。",
	UsageAPIDocAPIDoc:        "文檔的版本要號",
	UsageAPIDocLang:          "文檔內容的本地化 ID，比如 <var>zh-Hans</var>、<var>en-US</var> 等。",
	UsageAPIDocLogo:          "文檔的圖標，僅可使用 SVG 格式圖標。",
	UsageAPIDocCreated:       "文檔的創建時間",
	UsageAPIDocVersion:       "文檔的版本號",
	UsageAPIDocTitle:         "文檔的標題",
	UsageAPIDocDescription:   "文檔的整體描述內容",
	UsageAPIDocContact:       "文檔作者的聯系方式",
	UsageAPIDocLicense:       "文檔的版權信息",
	UsageAPIDocTags:          "文檔中定義的所有標簽",
	UsageAPIDocServers:       "API 基地址列表，每個 API 最少應該有壹個 server。",
	UsageAPIDocAPIs:          "文檔中的 API 文檔",
	UsageAPIDocHeaders:       "文檔中所有 API 都包含的公共報頭",
	UsageAPIDocResponses:     "文檔中所有 API 文檔都需要支持的返回內容",
	UsageAPIDocMimetypes:     "文檔所支持的 mimetype",
	UsageAPIDocXMLNamespaces: "針對 <var>application/xml</var> 類型的內容的命名空間設置",

	UsageXMLNamespace:       "為 <var>application/xml</var> 定義命名空間的相關屬性",
	UsageXMLNamespacePrefix: "命名空間的前綴，如果為空，則表示作為默認命名空間，命局只能有壹個默認命名空間。",
	UsageXMLNamespaceURN:    "命名空間的唯壹標識，需要全局唯壹，且區分大小寫。",

	UsageAPI:            "用於定義單個 API 接口的具體內容",
	UsageAPIVersion:     "表示此接口在該版本中添加",
	UsageAPIMethod:      "當前接口所支持的請求方法",
	UsageAPIID:          "接口的唯壹 ID",
	UsageAPIPath:        "定義路徑信息",
	UsageAPISummary:     "簡要介紹",
	UsageAPIDescription: "該接口的詳細介紹，為 HTML 內容。",
	UsageAPIRequests:    "定義可用的請求信息",
	UsageAPIResponses:   "定義可能的返回信息",
	UsageAPICallback:    "定義回調接口內容",
	UsageAPIDeprecated:  "在此版本之後將會被棄用",
	UsageAPIHeaders:     "傳遞的報頭內容，如果是某個 mimetype 專用的，可以放在 request 元素中。",
	UsageAPITags:        "關聯的標簽",
	UsageAPIServers:     "關聯的服務",

	UsageLink:     "用於描述鏈接信息，壹般轉換為 HTML 的 <code>a</code> 標簽。",
	UsageLinkText: "鏈接的字面文字",
	UsageLinkURL:  "鏈接指向的文本",

	UsageContact:      "用於描述聯系方式",
	UsageContactName:  "聯系人的名稱",
	UsageContactURL:   "聯系人的 URL",
	UsageContactEmail: "聯系人的電子郵件",

	UsageCallback:            "定義接口的回調內容",
	UsageCallbackMethod:      "回調的請求方法",
	UsageCallbackPath:        "回調的請求地址",
	UsageCallbackSummary:     "簡要介紹",
	UsageCallbackDeprecated:  "在此版本之後將會被棄用",
	UsageCallbackDescription: "對於回調的詳細介紹",
	UsageCallbackResponses:   "定義可能的返回信息",
	UsageCallbackRequests:    "定義可用的請求信息",
	UsageCallbackHeaders:     "傳遞的報頭內容",

	UsageEnum:            "定義枚舉類型的數所的枚舉值",
	UsageEnumDeprecated:  "該屬性棄用的版本號",
	UsageEnumValue:       "枚舉值",
	UsageEnumSummary:     "枚舉值的說明",
	UsageEnumDescription: "枚舉值的詳細說明",

	UsageExample:         "示例代碼",
	UsageExampleMimetype: "特定於類型的示例代碼",
	UsageExampleSummary:  "示例代碼的概要信息",
	UsageExampleContent:  "示例代碼的內容，需要使用 CDATA 包含代碼。",

	UsageParam:            "參數類型，基本上可以作為 request 的子集使用。",
	UsageParamName:        "值的名稱",
	UsageParamType:        "值的類型",
	UsageParamDeprecated:  "表示在大於等於該版本號時不再啟作用",
	UsageParamDefault:     "默認值",
	UsageParamOptional:    "是否為可選的參數",
	UsageParamArray:       "是否為數組",
	UsageParamItems:       "子類型，比如對象的子元素。",
	UsageParamSummary:     "簡要介紹",
	UsageParamEnums:       "當前參數可用的枚舉值",
	UsageParamDescription: "詳細介紹，為 HTML 內容。",
	UsageParamArrayStyle:  "以數組的方式展示數據",

	UsagePath:        "用於定義請求時與路徑相關的內容",
	UsagePathPath:    "接口地址",
	UsagePathParams:  "地址中的參數",
	UsagePathQueries: "地址中的查詢參數",

	UsageRequest:            "定義了請求和返回的相關內容",
	UsageRequestName:        "當 mimetype 為 <var>application/xml</var> 時，此值表示 XML 的頂層元素名稱，否則無用。",
	UsageRequestType:        "值的類型",
	UsageRequestDeprecated:  "表示在大於等於該版本號時不再啟作用",
	UsageRequestArray:       "是否為數組",
	UsageRequestItems:       "子類型，比如對象的子元素。",
	UsageRequestSummary:     "簡要介紹",
	UsageRequestStatus:      "狀態碼。在 request 中，該值不可用，否則為必填項。",
	UsageRequestEnums:       "當前參數可用的枚舉值",
	UsageRequestDescription: "詳細介紹，為 HTML 內容。",
	UsageRequestMimetype:    "媒體類型，比如 <var>application/json</var> 等。",
	UsageRequestExamples:    "示例代碼",
	UsageRequestHeaders:     "傳遞的報頭內容",

	UsageRichtext:     "富文本內容",
	UsageRichtextType: "指定富文本內容的格式，目前支持 <var>html</var> 和 <var>markdown</var>。",
	UsageRichtextText: "富文本的實際內容",

	UsageTag:           "用於對各個 API 進行分類",
	UsageTagName:       "標簽的唯壹 ID",
	UsageTagTitle:      "標簽的字面名稱",
	UsageTagDeprecated: "該標簽在大於該版本時被棄用",

	UsageServer:            "用於指定各個 API 的服務器地址",
	UsageServerName:        "服務唯壹 ID",
	UsageServerTitle:       "服務的字面名稱",
	UsageServerURL:         "服務的基地址，與該服務關聯的 API，訪問地址都是相對於此地址的。",
	UsageServerDeprecated:  "服務在大於該版本時被棄用",
	UsageServerSummary:     "服務的摘要信息",
	UsageServerDescription: "服務的詳細描述",

	UsageXMLAttr:    "是否作為父元素的屬性，僅作用於 XML 元素。是否作為父元素的屬性，僅用於 XML 的請求。",
	UsageXMLExtract: "將當前元素的內容作為父元素的內容，要求父元素必須為 <var>object</var>。",
	UsageXMLCData:   "當前內容為 CDATA，與 <code>@xml-attr</code> 互斥。",
	UsageXMLPrefix:  "XML 標簽的命名空間名稱前綴",
	UsageXMLWrapped: `如果當前元素的 <code>@array</code> 為 <var>true</var>，則可以通過此值指定在 XML 格式中的名稱。
	可以有三種格式：<ul>
	<li><samp>name</samp>：表示為數組添加壹個父元素名稱為 <var>name</var>；</li>
	<li><samp>name1&gt;name2</samp>：表示數組項的名稱改為 <var>name2</var>，且添加壹個父元素名為 <var>name1</var>；</li>
	<li><samp>&gt;name</samp>：表示將當前數組元素的名稱改為 <var>name</var>；</li>
	</ul>`,

	// 基本类型
	UsageString:  "普通的字符串類型，特殊字符需要使用 XML 實體，比如 <samp>&lt;</samp> 需要使用 <samp>&amp;lt;</samp> 代替。",
	UsageNumber:  "普通的數值類型，比如：<samp>1</samp>、<samp>-11.1</samp> 等。",
	UsageBool:    "布爾值類型，取值為 <var>true</var> 或是 <var>false</var>。",
	UsageVersion: `版本號，格式遵守 <a href="https://semver.org/lang/zh-CN/">semver</a> 規則。比如：<samp>1.0.1</samp>、<samp>1.0.1+20200618</samp>。`,
	UsageDate:    `采用 <a href="https://tools.ietf.org/html/rfc3339">RFC3339</a> 格式表示的時間，比如：<samp>2019-12-16T00:35:48+08:00</samp>。`,
	UsageType: `用於表示數據的類型值，格式為 <code>primitive[.subtype]</code>，其中 <code>primitive</code> 為基本類型，而 <code>subtype</code> 為子類型，用於對 <code>primitive</code> 進行進壹步的約束，當客戶端無法處理整個類型時，可以按照 <code>primitive</code> 的類型處理。<br />
	目前支持以下幾種類型：<ul>
	<li>空值；</li>
	<li><var>bool</var> 布爾值；</li>
	<li><var>object</var> 對象；</li>
	<li><var>number</var> 數值類型；</li>
	<li><var>number.int</var> 整數類型的數值；</li>
	<li><var>number.float</var> 浮點類型的數值；</li>
	<li><var>string</var> 字符串；</li>
	<li><var>string.url</var> URL 類型的字符串；</li>
	<li><var>string.email</var> email 類型的字符串；</li>
	<li><var>string.image</var> 表示圖片地址的 URL，在 mock 中該類型會生成壹張指定大小的圖片，圖片大小由查詢參數 <code>width</code> 和 <code>height</code> 指定，圖片類型由報頭 <code>Accept</code> 指定，目前允許 <var>image/png</var>、<var>image/jpeg</var> 和 <var>image/gif</var> 三種類型；</li>
	<li><var>string.date</var> 表示 <a href="https://tools.ietf.org/html/rfc3339#section-5.6">RFC3339</a> 中的 <code>full-date</code> 日期格式，比如 <samp>2020-01-02</samp>；</li>
	<li><var>string.time</var> 表示 <a href="https://tools.ietf.org/html/rfc3339#section-5.6">RFC3339</a> 中的 <code>full-time</code> 時間格式，比如 <samp>15:16:17Z</samp>、<samp>15:16:17+08:00</samp>；</li>
	<li><var>string.date-time</var> 表示 <a href="https://tools.ietf.org/html/rfc3339#section-5.6">RFC3339</a> 中的 <code>date-time</code> 格式，比如 <samp>2020-01-02T15:16:17-08:00</samp>；</li>
	</ul>`,

	// 以下是有关 build.Config 的字段说明
	UsageConfigVersion:               "此配置文件的所使用的文档版本",
	UsageConfigInputs:                "指定輸入的數據，同壹項目只能解析壹種語言。",
	UsageConfigInputsLang:            "源文件類型。具體支持的類型可通過 -l 參數進行查找。",
	UsageConfigInputsDir:             "需要解析的源文件所在目錄",
	UsageConfigInputsExts:            "只從這些擴展名的文件中查找文檔",
	UsageConfigInputsRecursive:       "是否解析子目錄下的源文件",
	UsageConfigInputsEncoding:        `編碼，默認為 <var>utf-8</var>，值可以是 <a href="https://www.iana.org/assignments/character-sets/character-sets.xhtml">character-sets</a> 中的內容。`,
	UsageConfigOutput:                "控制輸出行為",
	UsageConfigOutputType:            "輸出的類型，目前可以 <var>apidoc+xml</var>、<var>openapi+json</var> 和 <var>openapi+yaml</var>。",
	UsageConfigOutputPath:            "指定輸出的文件名，包含路徑信息。",
	UsageConfigOutputTags:            "只輸出與這些標簽相關聯的文檔，默認為全部。",
	UsageConfigOutputStyle:           "為 XML 文件指定的 XSL 文件",
	UsageConfigOutputNamespace:       "是否輸出命名空間",
	UsageConfigOutputNamespacePrefix: "如果輸出了命名空間，還可以指定命名空間前綴。",

	// 錯誤信息，可能在地方用到
	ErrInvalidUTF8Character:      "無效的 UTF8 字符",
	ErrInvalidXML:                "无效的 XML 文檔",
	ErrIsNotAPIDoc:               "並非有效的 apidoc 的文檔格式",
	ErrInvalidContentTypeCharset: "報頭 ContentType 中指定的字符集無效",
	ErrInvalidContentLength:      "報頭 ContentLength 無效",
	ErrBodyIsEmpty:               "請求的報文為空",
	ErrInvalidHeaderFormat:       "無效的報頭格式",
	ErrRequired:                  "不能為空",
	ErrInvalidFormat:             "格式不正確",
	ErrDirNotExists:              "目錄不存在",
	ErrNotFoundEndFlag:           "找不到結束符號",
	ErrNotFoundEndTag:            "找不到結束標簽",
	ErrNotFoundSupportedLang:     "該目錄下沒有支持的語言文件",
	ErrDirIsEmpty:                "目錄下沒有需要解析的文件",
	ErrInvalidValue:              "無效的值",
	ErrPathNotMatchParams:        "地址參數不匹配",
	ErrDuplicateValue:            "重復的值",
	ErrMessage:                   "%s 位於 %s",
	ErrNotFound:                  "未找到該值",
	ErrReadRemoteFile:            "讀取遠程文件 %s 時返回狀態碼 %d",
	ErrServerNotInitialized:      "服務未初始化",
	ErrInvalidLSPState:           "無效的 LSP 狀態",
	ErrInvalidURIScheme:          "無效的 URI 協議",
	ErrFileNotFound:              "未找到文件 %s",

	// logs
	InfoPrefix:    "[信息] ",
	WarnPrefix:    "[警告] ",
	ErrorPrefix:   "[錯誤] ",
	SuccessPrefix: "[成功] ",
}
