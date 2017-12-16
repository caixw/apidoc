// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package locale

import "github.com/caixw/apidoc/vars"

func init() {
	locales["zh-Hant"] = map[string]string{
		SyntaxError:  "在[%v:%v]出現語法錯誤[%v]",
		OptionsError: "配置文件[%v]中配置項[%v]錯誤[%v]",

		// 與 flag 包相關的處理
		FlagUsage: `%v 是壹個 RESTful API 文檔生成工具。

參數：
%v

源代碼采用 MIT 開源許可證，發布於 %v
詳細信息可訪問官網 %v
`,
		FlagHUsage:              "顯示幫助信息",
		FlagVUsage:              "顯示版本信息",
		FlagLUsage:              "顯示所有支持的語言",
		FlagGUsage:              "創建壹個默認的配置文件",
		FlagWDUsage:             "指定工作目錄，默認為當前目錄",
		FlagPprofUsage:          "指定壹種調試輸出類型，可以為 cpu 或是 mem",
		FlagVersionBuildWith:    "%v %v build with %v\n",
		FlagVersionCommitHash:   "commit hash %v\n",
		FlagSupportedLangs:      "目前支持以下語言 %v\n",
		FlagConfigWritedSuccess: "配置內容成功寫入 %v",
		FlagPprofWritedSuccess:  "pprof 的相關數據已經寫入到 %v",
		FlagInvalidPprrof:       "無效的 pprof 參數",

		VersionInCompatible: "當前程序與配置文件中指定的版本號不兼容",
		Complete:            "完成！文檔保存在：%v，總用時：%v",

		// 錯誤信息，可能在地方用到
		ErrRequired:              "不能為空",
		ErrInvalidFormat:         "格式不正確",
		ErrDirNotExists:          "目錄不存在",
		ErrMkdirError:            "創建目錄時發生以下錯誤：%v",
		ErrInvalidBlockType:      "無效的 block.Type 值：%v",
		ErrUnsupportedInputLang:  "無效的輸入語言：%v",
		ErrNotFoundEndFlag:       "找不到結束符號",
		ErrNotFoundSupportedLang: "該目錄下沒有支持的語言文件",
		ErrUnknownTag:            "不認識的標簽：%v",
		ErrDuplicateTag:          "重復的標簽：%v",
		ErrSuccessNotEmpty:       vars.APISuccess + " 不能为空",
		ErrTagArgTooMuch:         "標簽：%v 指定了太多的參數",
		ErrTagArgNotEnough:       "標簽：%v 參數不夠",
		ErrSecondArgMustURL:      vars.APILicense + " 第二個參數必須為 URL",

		// logs
		InfoPrefix:  "[信息] ",
		WarnPrefix:  "[警告] ",
		ErrorPrefix: "[錯誤] ",
	}
}
