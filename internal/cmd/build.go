// SPDX-License-Identifier: MIT

package cmd

import (
	"flag"
	"fmt"
	"io"
	"time"

	"github.com/caixw/apidoc/v5"
	"github.com/caixw/apidoc/v5/internal/locale"
	"github.com/caixw/apidoc/v5/message"
)

var buildFlagSet *flag.FlagSet

func init() {
	buildFlagSet = command.New("build", build, buildUsage)
}

func build(w io.Writer) error {
	var path = "./"
	if 0 != buildFlagSet.NArg() {
		path = buildFlagSet.Arg(0)
	}

	h := message.NewHandler(newHandlerFunc())
	apidoc.LoadConfig(h, path).Do(time.Now())

	h.Stop()
	return nil
}

func buildUsage(w io.Writer) error {
	_, err := fmt.Fprintln(w, locale.Sprintf(locale.CmdBuild))
	return err
}
