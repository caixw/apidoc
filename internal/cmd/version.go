// SPDX-License-Identifier: MIT

package cmd

import (
	"fmt"
	"io"
	"runtime"
	"strings"

	"github.com/caixw/apidoc/v5"
	"github.com/caixw/apidoc/v5/internal/locale"
	"github.com/caixw/apidoc/v5/internal/vars"
)

func initVersion() {
	command.New("version", version, buildUsage(locale.CmdVersionUsage))
}

func version(w io.Writer) error {
	goVersion := strings.TrimLeft(runtime.Version(), "go")
	msg := locale.Sprintf(locale.Version, apidoc.Version(), vars.DocVersion(), vars.CommitHash(), goVersion)
	_, err := fmt.Fprintln(w, msg)
	return err
}
