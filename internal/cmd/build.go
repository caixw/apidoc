// SPDX-License-Identifier: MIT

package cmd

import (
	"io"
	"time"

	"github.com/issue9/cmdopt"

	"github.com/caixw/apidoc/v7/build"
	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/locale"
)

var buildDir = uri("./")

func initBuild(command *cmdopt.CmdOpt) {
	fs := command.New("build", locale.Sprintf(locale.CmdBuildUsage), doBuild)
	fs.Var(&buildDir, "d", locale.Sprintf(locale.FlagBuildDirUsage))
}

func doBuild(io.Writer) error {
	start := time.Now()

	cfg, err := build.LoadConfig(core.URI(buildDir))
	if err != nil {
		return err
	}

	h := core.NewMessageHandler(messageHandle)
	defer h.Stop()

	cfg.Build(h)
	h.Locale(core.Info, locale.Complete, cfg.Output.Path, time.Since(start))
	return nil
}
