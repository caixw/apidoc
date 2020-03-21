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

	uri, err := getPath(testFlagSet)
	if err != nil {
		return err
	}
	build2.LoadConfig(h, uri).Test()
	return nil
}
