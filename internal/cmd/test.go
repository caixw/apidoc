// SPDX-License-Identifier: MIT

package cmd

import (
	"flag"
	"fmt"
	"io"

	"github.com/caixw/apidoc/v5"
	"github.com/caixw/apidoc/v5/internal/locale"
	"github.com/caixw/apidoc/v5/message"
)

var testFlagSet *flag.FlagSet

func init() {
	testFlagSet = command.New("test", test, testUsage)
}

func test(w io.Writer) error {
	h := message.NewHandler(newHandlerFunc())
	apidoc.LoadConfig(h, getPath(testFlagSet)).Test()
	h.Stop()
	return nil
}

func testUsage(w io.Writer) error {
	_, err := fmt.Fprintln(w, locale.Sprintf(locale.CmdTestUsage))
	return err
}
