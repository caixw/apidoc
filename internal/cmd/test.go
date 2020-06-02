// SPDX-License-Identifier: MIT

package cmd

import (
	"flag"
	"io"

	"github.com/caixw/apidoc/v7/build"
	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/locale"
)

var testFlagSet *flag.FlagSet
var testDir uri = "./"

func initTest() {
	testFlagSet = command.New("test", locale.Sprintf(locale.CmdTestUsage), test)
	testFlagSet.Var(&buildDir, "d", locale.Sprintf(locale.FlagTestDirUsage))
}

func test(w io.Writer) error {
	h := core.NewMessageHandler(messageHandle)
	defer h.Stop()

	if cfg := build.LoadConfig(h, core.URI(testDir)); cfg != nil {
		cfg.Test()
	}
	return nil
}
