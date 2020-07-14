// SPDX-License-Identifier: MIT

package cmd

import (
	"io"

	"github.com/issue9/cmdopt"

	"github.com/caixw/apidoc/v7/build"
	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/locale"
)

var syntaxDir uri = uri(core.FileURI("./"))

func initSyntax(command *cmdopt.CmdOpt) {
	fs := command.New("syntax", locale.Sprintf(locale.CmdSyntaxUsage), syntax)
	fs.Var(&syntaxDir, "d", locale.Sprintf(locale.FlagSyntaxDirUsage))
}

func syntax(w io.Writer) error {
	cfg, err := build.LoadConfig(syntaxDir.URI())
	if err != nil {
		return err
	}

	h := core.NewMessageHandler(messageHandle)
	defer h.Stop()

	cfg.CheckSyntax(h)
	h.Locale(core.Succ, locale.TestSuccess)
	return nil
}
