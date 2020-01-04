// SPDX-License-Identifier: MIT

package cmd

import (
	"flag"
	"io"
	"time"

	"github.com/caixw/apidoc/v5"
	"github.com/caixw/apidoc/v5/internal/locale"
	"github.com/caixw/apidoc/v5/message"
)

var buildFlagSet *flag.FlagSet

func initBuild() {
	buildFlagSet = command.New("build", build, buildUsage(locale.CmdBuildUsage))
}

func build(w io.Writer) error {
	h := message.NewHandler(newHandlerFunc())
	apidoc.LoadConfig(h, getPath(buildFlagSet)).Do(time.Now())
	h.Stop()
	return nil
}

// 人命令行尾部获取路径参数，或是在未指定的情况下，采用当前目录。
func getPath(fs *flag.FlagSet) string {
	if fs != nil && 0 != fs.NArg() {
		return fs.Arg(0)
	}
	return "./"
}
