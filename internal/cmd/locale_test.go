// SPDX-License-Identifier: MIT

package cmd

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v7"
)

func TestLocale(t *testing.T) {
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

	a.NotError(doLocale(w))
	ls := lines(w)
	a.Equal(len(ls), len(apidoc.Locales()))
	for _, l := range ls {
		cnt := strings.Count(l, strings.Repeat(" ", tail))
		a.True(cnt >= 1) // 至少两列
	}
}
