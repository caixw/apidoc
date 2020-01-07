// SPDX-License-Identifier: MIT

package cmd

import (
	"flag"
	"io"

	"github.com/caixw/apidoc/v6"
	"github.com/caixw/apidoc/v6/internal/locale"
	"github.com/caixw/apidoc/v6/message"
)

var testFlagSet *flag.FlagSet

func initTest() {
	testFlagSet = command.New("test", test, buildUsage(locale.CmdTestUsage))
}

func test(w io.Writer) error {
	h := message.NewHandler(newHandlerFunc())
	defer h.Stop()

	apidoc.LoadConfig(h, getPath(testFlagSet)).Test()
	return nil
}
