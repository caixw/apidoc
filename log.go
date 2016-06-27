// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package main

import (
	"log"

	"github.com/caixw/apidoc/app"
)

// 带色彩输出的控制台。
type syntaxWriter struct {
}

// io.Writer
func (w *syntaxWriter) Write(bs []byte) (size int, err error) {
	app.Error(string(bs))
	return len(bs), nil
}

func newSyntaxLog() *log.Logger {
	return log.New(&syntaxWriter{}, "", 0)
}
