// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"io"
	"log"
	"os"

	"github.com/issue9/term/colors"
)

// 日志信息输出
var (
	info = log.New(&logWriter{out: os.Stdout, color: colors.Green, prefix: "[INFO] "}, "", 0)
	warn = log.New(&logWriter{out: os.Stderr, color: colors.Cyan, prefix: "[WARN] "}, "", 0)
	erro = log.New(&logWriter{out: os.Stderr, color: colors.Red, prefix: "[ERRO] "}, "", 0)
)

// 带色彩输出的控制台。
type logWriter struct {
	out    io.Writer
	color  colors.Color
	prefix string
}

func (w *logWriter) Write(bs []byte) (int, error) {
	colors.Fprint(w.out, w.color, colors.Default, w.prefix)
	return colors.Fprint(w.out, colors.Default, colors.Default, string(bs))
}
