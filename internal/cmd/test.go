// SPDX-License-Identifier: MIT

package cmd

import (
	"flag"
	"io"

	build2 "github.com/caixw/apidoc/v6/build"
	"github.com/caixw/apidoc/v6/core"
	"github.com/caixw/apidoc/v6/internal/locale"
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
