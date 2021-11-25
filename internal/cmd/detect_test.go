// SPDX-License-Identifier: MIT

package cmd

import (
	"bytes"
	"testing"

	"github.com/issue9/assert/v2"
	"gopkg.in/yaml.v2"

	"github.com/caixw/apidoc/v7/build"
	"github.com/caixw/apidoc/v7/internal/ast"
	"github.com/caixw/apidoc/v7/internal/docs"
)

func TestCmdDetect(t *testing.T) {
	a := assert.New(t, false)

	buf := new(bytes.Buffer)
	path := docs.Dir().Append("example")

	cmd := Init(buf)
	resetPrinters()
	err := cmd.Exec([]string{"detect", "-d", path.String()})
	a.NotError(err)
	cfg := &build.Config{}
	a.NotError(yaml.Unmarshal(buf.Bytes(), cfg))
	a.Equal(cfg.Version, ast.Version)

	cmd = Init(buf)
	resetPrinters()
	err = cmd.Exec([]string{"detect", "-d", path.String(), "-w"})
	a.NotError(err)
	cfg2 := &build.Config{}
	data, err := build.LoadConfig(path)
	a.NotError(err).NotNil(data)
	a.NotError(yaml.Unmarshal(buf.Bytes(), cfg2))
	a.Equal(cfg, cfg2)
}
