// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package syntax

import (
	"testing"

	"github.com/issue9/assert"
)

func TestSplit(t *testing.T) {
	a := assert.New(t)

	tag := []byte("@tag s1\ts2  \n  s3")

	bs := Split(tag, 1)
	a.Equal(bs, [][]byte{[]byte("@tag s1\ts2  \n  s3")})

	bs = Split(tag, 2)
	a.Equal(bs, [][]byte{[]byte("@tag"), []byte("s1\ts2  \n  s3")})

	bs = Split(tag, 3)
	a.Equal(bs, [][]byte{[]byte("@tag"), []byte("s1"), []byte("s2  \n  s3")})

	bs = Split(tag, 4)
	a.Equal(bs, [][]byte{[]byte("@tag"), []byte("s1"), []byte("s2"), []byte("s3")})

	// 不够
	bs = Split(tag, 5)
	a.Equal(bs, [][]byte{[]byte("@tag"), []byte("s1"), []byte("s2"), []byte("s3")})
}
