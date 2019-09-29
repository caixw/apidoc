// SPDX-License-Identifier: MIT

// Package term 终端处理方法
package term

import (
	"io"

	"github.com/issue9/term/colors"
	xmessage "golang.org/x/text/message"

	"github.com/caixw/apidoc/v5/internal/locale"
	"github.com/caixw/apidoc/v5/message"
)

// NewHandlerFunc 声明用于 message.Handler 的处理函数
func NewHandlerFunc(
	erroOut, warnOut, infoOut, succOut io.Writer,
	erroColor, warnColor, infoColor, succColor colors.Color) message.HandlerFunc {
	erroPrefix := locale.Sprintf(locale.ErrorPrefix)
	warnPrefix := locale.Sprintf(locale.WarnPrefix)
	infoPrefix := locale.Sprintf(locale.InfoPrefix)
	succPrefix := locale.Sprintf(locale.SuccessPrefix)

	return func(msg *message.Message) {
		switch msg.Type {
		case message.Erro:
			printMessage(erroOut, erroColor, erroPrefix, msg.Message)
		case message.Warn:
			printMessage(warnOut, warnColor, warnPrefix, msg.Message)
		case message.Succ:
			printMessage(succOut, succColor, succPrefix, msg.Message)
		default: // message.Info 采用相同的值
			printMessage(infoOut, infoColor, infoPrefix, msg.Message)
		}
	}
}

func printMessage(out io.Writer, color colors.Color, prefix, msg string) {
	if _, err := colors.Fprint(out, color, colors.Default, prefix); err != nil {
		panic(err)
	}
	Line(out, colors.Default, msg)
}

// Locale 向控制台输出一行本地化的内容
func Locale(out io.Writer, color colors.Color, key xmessage.Reference, v ...interface{}) {
	l := locale.Sprintf(key, v...)
	Line(out, color, l)
}

// Line 向控制台输出一行内容
func Line(out io.Writer, color colors.Color, v ...interface{}) {
	if _, err := colors.Fprintln(out, color, colors.Default, v...); err != nil {
		panic(err)
	}
}
