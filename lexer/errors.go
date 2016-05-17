// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package lexer

import "fmt"

// SyntaxError 语法错误
type SyntaxError struct {
	Line    int
	File    string
	Message string
}

func (err *SyntaxError) Error() string {
	return fmt.Sprintf("在[%v:%v]出现语法错误[%v]", err.File, err.Line, err.Message)
}
