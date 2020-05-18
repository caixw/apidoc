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
	h := core.NewMessageHandler(messageHandle)
	defer h.Stop()

	build.LoadConfig(h, getPath(buildFlagSet)).Build(time.Now())
	return nil
}
