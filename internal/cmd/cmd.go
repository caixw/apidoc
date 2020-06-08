// SPDX-License-Identifier: MIT

// Package cmd 提供子命令的相关功能
package cmd

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/issue9/cmdopt"
	"github.com/issue9/term/v2/colors"
	"golang.org/x/text/message"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/locale"
)

// 命令行输出的表格中，每一列为了对齐填补的空格数量。
const tail = 3

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

type uri core.URI

func (u uri) Get() interface{} {
	return string(u)
}

func (u *uri) Set(v string) error {
	*u = uri(core.FileURI(v))
	return nil
}

func (u *uri) String() string {
	return core.URI(*u).String()
}

// Init 初始化 cmdopt.CmdOpt 实例
func Init(out io.Writer) *cmdopt.CmdOpt {
	command := cmdopt.New(
		out,
		flag.ExitOnError,
		locale.Sprintf(locale.CmdUsage, core.Name),
		locale.Sprintf(locale.CmdUsageFooter, core.OfficialURL, core.RepoURL),
		locale.Sprintf(locale.CmdUsageOptions),
		locale.Sprintf(locale.CmdUsageCommands),
		func(name string) string {
			return locale.Sprintf(locale.CmdNotFound, name)
		})

	command.Help("help", locale.Sprintf(locale.CmdHelpUsage))
	initBuild(command)
	initDetect(command)
	initLang(command)
	initLocale(command)
	initSyntax(command)
	initVersion(command)
	initMock(command)
	initStatic(command)
	initLSP(command)

	return command
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
