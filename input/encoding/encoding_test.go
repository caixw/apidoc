// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package encoding

import (
	"io/ioutil"
	"testing"

	"github.com/issue9/assert"
)

func TestTransform(t *testing.T) {
	a := assert.New(t)

	u8, err := ioutil.ReadFile("./testdata/utf8")
	a.NotError(err).NotNil(u8)

	utf8, err := Transform("./testdata/utf8", "utf-8")
	a.NotError(err).NotNil(utf8)
	a.Equal(u8, utf8)

	gb18030, err := Transform("./testdata/gb18030", "gb18030")
	a.NotError(err).NotNil(gb18030)
	a.Equal(u8, gb18030)

	gbk, err := Transform("./testdata/gbk", "gbk")
	a.NotError(err).NotNil(gbk)
	a.Equal(u8, gbk)

	// 以错误的编码方式加载
	big5, err := Transform("./testdata/gbk", "big5")
	a.NotError(err).NotNil(big5)
	a.NotEqual(u8, big5)

	// 以不存在的编码加载
	notExists, err := Transform("./testdata/gbk", "not-exists")
	a.Error(err).Nil(notExists)
}
