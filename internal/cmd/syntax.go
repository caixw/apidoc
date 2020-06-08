// SPDX-License-Identifier: MIT

package cmd

import (
	"io"

	"github.com/issue9/cmdopt"

	"github.com/caixw/apidoc/v7/build"
	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/locale"
)

var syntaxDir uri = "./"

func initSyntax(command *cmdopt.CmdOpt) {
	fs := command.New("syntax", locale.Sprintf(locale.CmdSyntaxUsage), syntax)
	fs.Var(&syntaxDir, "d", locale.Sprintf(locale.FlagSyntaxDirUsage))
}

func syntax(w io.Writer) error {
	h := core.NewMessageHandler(messageHandle)
	defer h.Stop()

	if cfg := build.LoadConfig(h, core.URI(syntaxDir)); cfg != nil {
		cfg.CheckSyntax()
	}
	return nil
}
