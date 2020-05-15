// SPDX-License-Identifier: MIT

package locale

var cmnHant = map[string]string{
	// 與 flag 包相關的處理
	CmdUsage: `%s 是壹個 RESTful API 文檔生成工具

用法：
apidoc cmd [args]

cmd 為子命令，args 為傳遞給子命令的參數，目前支持以下子命令：
%s

源代碼采用 MIT 開源許可證，發布於 %s
詳細信息可訪問官網 %s`,
	CmdHelpUsage:    "顯示幫助信息",
	CmdVersionUsage: "顯示版本信息",
	CmdLangUsage:    "顯示所有支持的語言",
	CmdLocaleUsage:  "顯示所有支持的本地化內容",
	CmdDetectUsage: `根據目錄下的內容生成配置文件

用法：
apidoc detect [options] [path]

options 可以是以下參數：
%s

path 表示目錄的路徑，或不指定，表示使用當前工作目錄 ./ 代替。`,
	CmdTestUsage: "測試語法的正確性",
	CmdMockUsage: `啟用 mock 服務

用法：
apidoc mock [options] [path]

options 可以是以下參數：
%s

path 表示文檔路徑，或不指定，則使用當前工作目錄 ./ 代替。`,
	CmdBuildUsage: "生成文檔內容",
	CmdStaticUsage: `啟用靜態文件服務

用法：
apidoc static [options] [path]

options 可以是以下參數：
%s

path 表示需要展示的文檔路徑，為空表示沒有需要展示的文檔。`,
	CmdLSPUsage: `啟動 language server protocol 服務

用法：
apidoc lsp [options]

options 可以是以下參數
%s`,
	Version:                    "版本：%s\n文檔：%s\nLSP：%s\nGo：%s",
	CmdNotFound:                "子命令 %s 未找到\n",
	FlagMockPortUsage:          "指定 mock 服務的端口號",
	FlagMockServersUsage:       "指定 mock 服務時，文檔中 server 變量對應的路由前綴",
	FlagDetectRecursive:        "detect 子命令是否檢測子目錄的值",
	FlagStaticPortUsage:        "指定 static 服務的端口號",
	FlagStaticDocsUsage:        "指定 static 服務的文件夾",
	FlagStaticStylesheetUsage:  "指定 static 是否只啟用樣式文件內容",
	FlagStaticContentTypeUsage: "指定 static 的 content-type 值，不指定，則根據擴展名自動獲取",
	FlagStaticURLUsage:         "指定 static 服務中文檔的輸出地址",
	FlagLSPPortUsage:           "指定 language server protocol 服務的端口號",
	FlagLSPModeUsage:           "指定 language server protocol 的運行方式，可以是 websocket、tcp 和 udp",
	FlagLSPHeaderUsage:         "指定 language server protocol 傳遞內容是否帶報頭信息",

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

	// 文檔樹中各個字段的介紹
	UsageAPIDoc:            "用於描述整個文檔的相關內容，只能出現壹次。",
	UsageAPIDocAPIDoc:      "文檔的版本要號",
	UsageAPIDocLang:        "文檔內容的本地化 ID，比如 <samp>zh-Hans</samp>、<samp>en-US</samp> 等。",
	UsageAPIDocLogo:        "文檔的圖標，僅可使用 SVG 格式圖標。",
	UsageAPIDocCreated:     "文檔的創建時間",
	UsageAPIDocVersion:     "文檔的版本號，需要遵守 semver 的約定。",
	UsageAPIDocTitle:       "文檔的標題",
	UsageAPIDocDescription: "文檔的整體描述內容",
	UsageAPIDocContact:     "文檔作者的聯系方式",
	UsageAPIDocLicense:     "文檔的版權信息",
	UsageAPIDocTags:        "文檔中定義的所有標簽",
	UsageAPIDocServers:     "API 基地址列表，每個 API 最少應該有壹個 server。",
	UsageAPIDocAPIs:        "文檔中的 API 文檔",
	UsageAPIDocResponses:   "文檔中所有 API 文檔都需要支持的返回內容",
	UsageAPIDocMimetypes:   "文檔所支持的 mimetype",

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

	UsageLink:     "用於描述鏈接信息，壹般轉換為 HTML 的 a 標簽。",
	UsageLinkText: "鏈接的字面文字",
	UsageLinkURL:  "鏈接指向的文本",

	UsageContact:      "用於描述聯系方式",
	UsageContactName:  "聯系人的名稱",
	UsageContactURL:   "聯系人的 URL",
	UsageContactEmail: "聯系人的電子郵件",

	UsageCallback:            "定義回調信息",
	UsageCallbackMethod:      "回調的請求方法",
	UsageCallbackPath:        "定義回調的地址",
	UsageCallbackSummary:     "簡要介紹",
	UsageCallbackDeprecated:  "在此版本之後將會被棄用",
	UsageCallbackDescription: "該接口的詳細介紹",
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
	UsageParamXMLAttr:     "是否作為父元素的屬性，僅作用於 XML 元素。是否作為父元素的屬性，僅用於 XML 的請求。",
	UsageParamXMLExtract:  "將當前元素的內容作為父元素的內容，要求父元素必須為 <var>object</var>。",
	UsageParamXMLNS:       "XML 標簽的命名空間",
	UsageParamXMLNSPrefix: "XML 標簽的命名空間名稱前綴",
	UsageParamXMLWrapped:  "如果當前元素的 @array 為 <var>true</var>，是否將其包含在 wrapped 指定的標簽中。",
	UsageParamName:        "值的名稱",
	UsageParamType:        "值的類型，可以是 <var>string</var>、<var>number</var>、<var>bool</var> 和 <var>object</var>",
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
	UsageRequestType:        "值的類型，可以是 <var>none</var>、<var>string</var>、<var>number</var>、<var>bool</var>、<var>object</var> 和空值；空值表示不輸出任何內容。",
	UsageRequestDeprecated:  "表示在大於等於該版本號時不再啟作用",
	UsageRequestArray:       "是否為數組",
	UsageRequestItems:       "子類型，比如對象的子元素。",
	UsageRequestSummary:     "簡要介紹",
	UsageRequestStatus:      "狀態碼。在 request 中，該值不可用，否則為必填項。",
	UsageRequestEnums:       "當前參數可用的枚舉值",
	UsageRequestDescription: "詳細介紹，為 HTML 內容。",
	UsageRequestMimetype:    "媒體類型，比如 application/json 等。",
	UsageRequestExamples:    "示例代碼",
	UsageRequestHeaders:     "傳遞的報頭內容",

	UsageRichtext:     "富文本內容",
	UsageRichtextType: "指定內容的格式",
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
	UsageXMLNS:      "XML 標簽的命名空間",
	UsageXMLPrefix:  "XML 標簽的命名空間名稱前綴",
	UsageXMLWrapped: "如果當前元素的 <code>@array</code> 為 <var>true</var>，是否將其包含在 wrapped 指定的標簽中。",

	// 基本类型
	UsageString: "普通的字符串類型",

	// 錯誤信息，可能在地方用到
	ErrInvalidUTF8Character:      "無效的 UTF8 字符",
	ErrInvalidXML:                "无效的 XML 文檔",
	ErrIsNotAPIDoc:               "並非有效的 apidoc 的文檔格式",
	ErrInvalidContentTypeCharset: "報頭 ContentType 中指定的字符集無效 ",
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
