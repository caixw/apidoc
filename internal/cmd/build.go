// SPDX-License-Identifier: MIT

package cmd

import (
	"flag"
	"io"
	"time"

	"github.com/caixw/apidoc/v7/build"
	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/locale"
)

var buildFlagSet *flag.FlagSet

func initBuild() {
	buildFlagSet = command.New("build", doBuild, buildUsage(locale.CmdBuildUsage))
}

func doBuild(io.Writer) error {
	h := core.NewMessageHandler(newHandlerFunc())
	defer h.Stop()

	build.LoadConfig(h, getPath(buildFlagSet)).Build(time.Now())
	return nil
}

// 人命令行尾部获取路径参数，或是在未指定的情况下，采用当前目录。
func getPath(fs *flag.FlagSet) core.URI {
	if fs != nil && 0 != fs.NArg() {
		return core.FileURI(fs.Arg(0))
	}
	return core.FileURI("./")
}
