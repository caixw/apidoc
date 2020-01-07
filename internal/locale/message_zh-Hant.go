// SPDX-License-Identifier: MIT

package locale

import "golang.org/x/text/language"

var zhHant = map[string]string{
	// 與 flag 包相關的處理
	CmdUsage: `%s 是壹個 RESTful API 文檔生成工具

用法：
apidoc cmd [args]

cmd 為子命令，args 為傳遞給子命令的參數，目前支持以下子命令。
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
	Version:                    "版本：%s\n文檔：%s\n提交：%s\nGo：%s",
	CmdNotFound:                "子命令 %s 未找到\n",
	FlagMockPortUsage:          "指定 mock 服務的端口號",
	FlagMockServersUsage:       "指定 mock 服務時，文檔中 server 變量對應的路由前綴",
	FlagDetectRecursive:        "detect 子命令是否檢測子目錄的值",
	FlagStaticPortUsage:        "指定 static 服務的端口號",
	FlagStaticDocsUsage:        "指定 static 服務的文件夾",
	FlagStaticStylesheetUsage:  "指定 static 是否只啟用樣式文件內容",
	FlagStaticContentTypeUsage: "指定 static 的 content-type 值，不指定，則根據擴展名自動獲取",
	FlagStaticURLUsage:         "指定 static 服務中文檔的輸出地址",

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

	// 錯誤信息，可能在地方用到
	ErrRequired:              "不能為空",
	ErrInvalidFormat:         "格式不正確",
	ErrDirNotExists:          "目錄不存在",
	ErrUnsupportedInputLang:  "不支持的輸入語言：%s",
	ErrNotFoundEndFlag:       "找不到結束符號",
	ErrNotFoundSupportedLang: "該目錄下沒有支持的語言文件",
	ErrDirIsEmpty:            "目錄下沒有需要解析的文件",
	ErrInvalidValue:          "無效的值",
	ErrPathNotMatchParams:    "地址參數不匹配",
	ErrDuplicateValue:        "重復的值",
	ErrMessage:               "%s 位於 %s",
	ErrNotFound:              "未找到該值",

	// logs
	InfoPrefix:    "[信息] ",
	WarnPrefix:    "[警告] ",
	ErrorPrefix:   "[錯誤] ",
	SuccessPrefix: "[成功] ",
}

func init() {
	addLocale(language.MustParse("zh-Hant"), zhHant)
}
