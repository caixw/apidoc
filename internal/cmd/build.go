// SPDX-License-Identifier: MIT

package cmd

import (
	"flag"
	"fmt"
	"io"
	"time"

	"github.com/caixw/apidoc/v5"
	"github.com/caixw/apidoc/v5/internal/locale"
	"github.com/caixw/apidoc/v5/message"
)

var buildFlagSet *flag.FlagSet

func init() {
	buildFlagSet = command.New("build", build, buildUsage)
}

func build(w io.Writer) error {
	h := message.NewHandler(newHandlerFunc())
	apidoc.LoadConfig(h, getPath(buildFlagSet)).Do(time.Now())
	h.Stop()
	return nil
}

func buildUsage(w io.Writer) error {
	_, err := fmt.Fprintln(w, locale.Sprintf(locale.CmdBuildUsage))
	return err
}

// 人命令行尾部获取路径参数，或是在未指定的情况下，采用当前目录。
func getPath(fs *flag.FlagSet) string {
	if fs != nil && 0 != fs.NArg() {
		return fs.Arg(0)
	}
	return "./"
}
