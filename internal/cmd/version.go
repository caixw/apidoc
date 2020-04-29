// SPDX-License-Identifier: MIT

package cmd

import (
	"fmt"
	"io"
	"runtime"
	"strings"

	"github.com/caixw/apidoc/v7"
	"github.com/caixw/apidoc/v7/internal/ast"
	"github.com/caixw/apidoc/v7/internal/locale"
	"github.com/caixw/apidoc/v7/internal/vars"
)

func initVersion() {
	command.New("version", version, buildUsage(locale.CmdVersionUsage))
}

func version(w io.Writer) error {
	goVersion := strings.TrimLeft(runtime.Version(), "go")
	msg := locale.Sprintf(locale.Version, apidoc.Version(), ast.Version, vars.CommitHash(), goVersion)
	_, err := fmt.Fprintln(w, msg)
	return err
}
