// SPDX-License-Identifier: MIT

package locale

import "golang.org/x/text/language"

var zhHant = map[string]string{
	// 與 flag 包相關的處理
	FlagUsage: `%s 是壹個 RESTful API 文檔生成工具。

用法：
apidoc [options] [path]

參數：
%s

源代碼采用 MIT 開源許可證，發布於 %s
詳細信息可訪問官網 %s
`,
	FlagHUsage:            "顯示幫助信息",
	FlagVUsage:            "顯示版本信息",
	FlagLUsage:            "顯示所有支持的語言",
	FlagDUsage:            "根据目录下的内容生成配置文件",
	FlagVersionBuildWith:  "%s %s build with %s",
	FlagVersionCommitHash: "commit hash %s",

	VersionInCompatible: "當前程序與配置文件中指定的版本號不兼容",
	Complete:            "完成！文檔保存在：%s，總用時：%v",
	ConfigWriteSuccess:  "配置內容成功寫入 %s",

	// 錯誤信息，可能在地方用到
	ErrRequired:              "不能為空",
	ErrMustEmpty:             "只能為空",
	ErrInvalidFormat:         "格式不正確",
	ErrDirNotExists:          "目錄不存在",
	ErrUnsupportedInputLang:  "不支持的輸入語言：%s",
	ErrNotFoundEndFlag:       "找不到結束符號",
	ErrNotFoundSupportedLang: "該目錄下沒有支持的語言文件",
	ErrUnknownTag:            "不認識的標簽",
	ErrDuplicateTag:          "重復的標簽",
	ErrUnsupportedEncoding:   "不支持的編碼方式",
	ErrDirIsEmpty:            "目錄下沒有需要解析的文件",
	ErrInvalidValue:          "無效的值",
	ErrInvalidOpenapi:        "openapi 內容錯誤：字段：%s；錯誤內容：%s",
	ErrPathNotMatchParams:    "地址參數不匹配",
	ErrDuplicateValue:        "重復的值",
	ErrMessage:               "%s 位於 %s",

	// logs
	InfoPrefix:    "[信息] ",
	WarnPrefix:    "[警告] ",
	ErrorPrefix:   "[錯誤] ",
	SuccessPrefix: "[成功] ",
}

func init() {
	addLocale(language.MustParse("zh-Hant"), zhHant)
}
