// SPDX-License-Identifier: MIT

package main

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/issue9/assert"
	"github.com/issue9/term/colors"
	"github.com/issue9/utils"

	"github.com/caixw/apidoc/v5/internal/lang"
)

func TestDetect(t *testing.T) {
	a := assert.New(t)

	if utils.FileExists("./.apidoc.yaml") {
		a.NotError(os.Remove("./.apidoc.yaml"))
	}

	detect()
	a.FileExists("./.apidoc.yaml")
}

func TestParse(t *testing.T) {
	a := assert.New(t)

	if utils.FileExists("./apidoc.xml") {
		a.NotError(os.Remove("./apidoc.xml"))
	}

	// 测试语法，不生成文件
	parse(true)
	a.FileNotExists("./apidoc.xml")

	parse(false)
	a.FileExists("./apidoc.xml")
}

func TestLangs(t *testing.T) {
	a := assert.New(t)
	w := new(bytes.Buffer)

	lines := func(w *bytes.Buffer) []string {
		b := bufio.NewReader(w)
		lines := make([]string, 0, 100)
		for line, err := b.ReadString('\n'); err != io.EOF; line, err = b.ReadString('\n') {
			lines = append(lines, line)
		}

		return lines
	}

	langs(w, colors.Default, 3)
	ls := lines(w)
	a.Equal(len(ls), len(lang.Langs()))
	for _, l := range ls {
		cnt := strings.Count(l, strings.Repeat(" ", 3))
		a.True(cnt >= 2)
	}

	w.Reset()
	langs(w, colors.Default, 10)
	ls = lines(w)
	a.Equal(len(ls), len(lang.Langs()))
	for _, l := range ls {
		cnt := strings.Count(l, strings.Repeat(" ", 10))
		a.True(cnt >= 2)
	}
}

func TestGetPaths(t *testing.T) {
	a := assert.New(t)

	paths, err := getPaths()
	a.NotError(err).NotNil(paths)
	a.Equal(1, len(paths))

	abs, err := filepath.Abs("./")
	a.NotError(err).NotEmpty(abs)
	a.Equal(paths[0], abs)
}
