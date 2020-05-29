// SPDX-License-Identifier: MIT

// Package cmd 提供子命令的相关功能
package cmd

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/issue9/cmdopt"
	"github.com/issue9/term/v2/colors"
	"golang.org/x/text/message"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/locale"
)

// 命令行输出的表格中，每一列为了对齐填补的空格数量。
const tail = 3

var command *cmdopt.CmdOpt

var printers = map[core.MessageType]*printer{
	core.Erro: {
		out:    os.Stderr,
		color:  colors.Red,
		prefix: locale.ErrorPrefix,
	},
	core.Warn: {
		out:    os.Stderr,
		color:  colors.Cyan,
		prefix: locale.WarnPrefix,
	},
	core.Info: {
		out:    os.Stdout,
		color:  colors.Default,
		prefix: locale.InfoPrefix,
	},
	core.Succ: {
		out:    os.Stdout,
		color:  colors.Green,
		prefix: locale.SuccessPrefix,
	},
}

type printer struct {
	out    io.Writer
	color  colors.Color
	prefix message.Reference
}

// Init 初始化 cmdopt.CmdOpt 实例
func Init(out io.Writer) *cmdopt.CmdOpt {
	command = cmdopt.New(out, flag.ContinueOnError, usage, func(name string) string {
		return locale.Sprintf(locale.CmdNotFound, name)
	})

	command.Help("help", buildUsage(locale.CmdHelpUsage))
	initBuild()
	initDetect()
	initLang()
	initLocale()
	initTest()
	initVersion()
	initMock()
	initStatic()
	initLSP()

	return command
}

func usage(w io.Writer) error {
	cmds := strings.Join(command.Commands(), ",")
	msg := locale.Sprintf(locale.CmdUsage, core.Name, cmds, core.RepoURL, core.OfficialURL)
	_, err := fmt.Fprintln(w, msg)
	return err
}

func buildUsage(key message.Reference, v ...interface{}) cmdopt.DoFunc {
	return func(w io.Writer) error {
		_, err := fmt.Fprintln(w, locale.Sprintf(key, v...))
		return err
	}
}

func getFlagSetUsage(fs *flag.FlagSet) string {
	buf := new(bytes.Buffer)
	origin := fs.Output()
	fs.SetOutput(buf)
	fs.PrintDefaults()
	fs.SetOutput(origin)
	return buf.String()
}

// 从命令行尾部获取路径参数，或是在未指定的情况下，采用当前目录。
func getPath(fs *flag.FlagSet) core.URI {
	if fs != nil && 0 != fs.NArg() {
		return core.FileURI(fs.Arg(0))
	}
	return core.FileURI("./")
}

func messageHandle(msg *core.Message) {
	printers[msg.Type].print(msg.Message)
}

func (p *printer) print(msg interface{}) {
	if _, err := colors.Fprint(p.out, colors.Normal, p.color, colors.Default, locale.New(p.prefix)); err != nil {
		panic(err)
	}

	if _, err := fmt.Fprintln(p.out, msg); err != nil {
		panic(err)
	}
}
