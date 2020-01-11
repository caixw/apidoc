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

	// 错误信息，可能在地方用到
	ErrRequired:              "不能为空",
	ErrInvalidFormat:         "格式不正确",
	ErrDirNotExists:          "目录不存在",
	ErrNotFoundEndFlag:       "找不到结束符号",
	ErrNotFoundSupportedLang: "该目录下没有支持的语言文件",
	ErrDirIsEmpty:            "目录下没有需要解析的文件",
	ErrInvalidValue:          "无效的值",
	ErrPathNotMatchParams:    "地址参数不匹配",
	ErrDuplicateValue:        "重复的值",
	ErrMessage:               "%s 位于 %s",
	ErrNotFound:              "未找到该值",
	ErrReadRemoteFile:        "读取远程文件 %s 时返回状态码 %d",

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
