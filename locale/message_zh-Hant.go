// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package locale

func init() {
	locales["zh-Hant"] = map[string]string{
		// 與 flag 包相關的處理
		FlagUsage: `%v 是壹個 RESTful API 文檔生成工具。

參數：
%v

源代碼采用 MIT 開源許可證，發布於 %v
詳細信息可訪問官網 %v
`,
		FlagHUsage:              "顯示幫助信息",
		FlagVUsage:              "顯示版本信息",
		FlagLanguagesUsage:      "顯示所有支持的語言",
		FlagGUsage:              "創建壹個默認的配置文件",
		FlagEncodingsUsage:      "顯示支持的編碼方式",
		FlagWDUsage:             "指定工作目錄，默認為當前目錄",
		FlagPprofUsage:          "指定壹種調試輸出類型，可以為 %s 或是 %s",
		FlagVersionBuildWith:    "%v %v build with %v\n",
		FlagVersionCommitHash:   "commit hash %v\n",
		FlagSupportedLanguages:  "目前支持以下語言 %v\n",
		FlagSupportedEncodings:  "目前支持以下編碼 %v\n",
		FlagConfigWritedSuccess: "配置內容成功寫入 %v",
		FlagPprofWritedSuccess:  "pprof 的相關數據已經寫入到 %v",
		FlagInvalidPprrof:       "無效的 pprof 參數",

		VersionInCompatible: "當前程序與配置文件中指定的版本號不兼容",
		Complete:            "完成！文檔保存在：%v，總用時：%v",

		// 錯誤信息，可能在地方用到
		ErrRequired:              "不能為空",
		ErrMustEmpty:             "隻能為空",
		ErrInvalidFormat:         "格式不正確",
		ErrDirNotExists:          "目錄不存在",
		ErrInvalidBlockType:      "無效的 block.Type 值：%v",
		ErrUnsupportedInputLang:  "無效的輸入語言：%v",
		ErrNotFoundEndFlag:       "找不到結束符號",
		ErrNotFoundSupportedLang: "該目錄下沒有支持的語言文件",
		ErrUnknownTag:            "不認識的標簽：%v",
		ErrDuplicateTag:          "重復的標簽：%v",
		ErrTagArgTooMuch:         "標簽：%v 指定了太多的參數",
		ErrTagArgNotEnough:       "標簽：%v 參數不夠",
		ErrUnsupportedEncoding:   "不支持的編碼方式：%v",
		ErrDirIsEmpty:            "目錄下沒有需要解析的文件",
		ErrInvalidValue:          "無效的值",
		ErrInvalidOpenapi:        "openapi 內容錯誤：字段：%s；錯誤內容：%s",
		ErrSyntax:                "在[%s:%d]出現語法錯誤[%s]",
		ErrConfig:                "配置文件[%s]中配置項[%s]錯誤[%s]",
		ErrApidocExists:          "相同組名的 @apidoc 標簽已經存在",
		ErrInvalidMethod:         "無效的請求方法",
		ErrMethodExists:          "該請求方法已經存在",
		ErrAPIMissingParam:       "@api 缺少參數",

		// logs
		InfoPrefix:  "[信息] ",
		WarnPrefix:  "[警告] ",
		ErrorPrefix: "[錯誤] ",
	}
}
