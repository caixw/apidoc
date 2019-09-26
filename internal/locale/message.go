// SPDX-License-Identifier: MIT

package locale

// 各个语言需要翻译的所有字符串
const (
	// 与 flag 包相关的处理
	FlagUsage = `%s 是一个 RESTful API 文档生成工具。

用法：
apidoc [options] [path]

参数：
%s

源代码采用 MIT 开源许可证，发布于 %s
详细信息可访问官网 %s
`
	FlagHUsage            = "显示帮助信息"
	FlagVUsage            = "显示版本信息"
	FlagLanguagesUsage    = "显示所有支持的语言"
	FlagDUsage            = "根据目录下的内容生成配置文件"
	FlagVersionBuildWith  = "%s %s build with %s"
	FlagVersionCommitHash = "commit hash %s"

	VersionInCompatible = "当前程序与配置文件中指定的版本号不兼容"
	Complete            = "完成！文档保存在：%s，总用时：%v"
	ConfigWriteSuccess  = "配置内容成功写入 %s"

	// 错误信息，可能在地方用到
	ErrRequired              = "不能为空"
	ErrMustEmpty             = "只能为空"
	ErrInvalidFormat         = "格式不正确"
	ErrDirNotExists          = "目录不存在"
	ErrUnsupportedInputLang  = "不支持的输入语言：%s"
	ErrNotFoundEndFlag       = "找不到结束符号"
	ErrNotFoundSupportedLang = "该目录下没有支持的语言文件"
	ErrUnknownTag            = "不认识的标签"
	ErrDuplicateTag          = "重复的标签"
	ErrUnsupportedEncoding   = "不支持的编码方式"
	ErrDirIsEmpty            = "目录下没有需要解析的文件"
	ErrInvalidValue          = "无效的值"
	ErrInvalidOpenapi        = "openapi 内容错误：字段：%s；错误内容：%s"
	ErrDuplicateRoute        = "重复的路由项"
	ErrPathNotMatchParams    = "地址参数不匹配"
	ErrPathInvalid           = "地址格式错误"
	ErrDuplicateValue        = "重复的值"
	ErrMessage               = "%s 位于 %s:%s"
	ErrInvalidMethod         = "无效的请求方法: %s"
	ErrMethodExists          = "该请求方法已经存在"
	ErrInvalidTag            = "无效的标签 %s"
	ErrInvalidURL            = "无效的 URL %s"
	ErrInvalidEmail          = "无效的邮件地址 %s"
	ErrInvalidType           = "无效的类型名称 %s"
	ErrInvalidStatus         = "无效的 HTTP 状态值 %s"
	ErrInvalidVersionFormat  = "无效的版本号格式"
	ErrDuplicateEnum         = "重复的枚举值 %s"
	ErrNeedProperty          = "对象必须要有属性值"
	ErrPathSyntaxError       = "路由项语法错误"

	// logs
	InfoPrefix    = "[INFO] "
	WarnPrefix    = "[WARN] "
	ErrorPrefix   = "[ERRO] "
	SuccessPrefix = "[SUCC] "
)
