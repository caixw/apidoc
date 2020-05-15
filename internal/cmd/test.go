// SPDX-License-Identifier: MIT

package cmd

import (
	"flag"
	"io"

	build2 "github.com/caixw/apidoc/v7/build"
	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/locale"
)

var testFlagSet *flag.FlagSet

func initTest() {
	testFlagSet = command.New("test", test, buildUsage(locale.CmdTestUsage))
}

func test(w io.Writer) error {
	h := core.NewMessageHandler(newHandlerFunc())
	defer h.Stop()

	build2.LoadConfig(h, getPath(testFlagSet)).Test()
	return nil
}
