// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package core

import (
	"bytes"
	"testing"

	"github.com/issue9/assert"
)

func cstyle(data []byte) ([]rune, int) {
	index := bytes.Index(data, []byte("/*"))
	end := bytes.Index(data, []byte("*/"))
	if index < 0 || end < index {
		return nil, -1
	}

	return []rune(string(data[index+2 : end])), end + 2
}

func TestScanFile(t *testing.T) {
	a := assert.New(t)
	ds := &docs{items: []*Doc{}}

	scanFile(ds, cstyle, "./testcode/php1.php")

	a.Equal(len(ds.items), 2)
	a.Equal(ds.items[0].Group, "php1").
		Equal(ds.items[0].Method, "post").
		Equal(ds.items[0].URL, "/api/php1/post")
	a.Equal(ds.items[1].Group, "php1").
		Equal(ds.items[1].Method, "get").
		Equal(ds.items[1].URL, "/api/php1/get")
}

func TestScanFiles(t *testing.T) {
	a := assert.New(t)

	paths := []string{
		"./testcode/php1.php",
		"./testcode/php2.php",
	}
	items, err := ScanFiles(paths, cstyle)
	a.NotError(err).NotNil(items)
	a.Equal(4, len(items))
	for _, v := range items {
		switch {
		case v.URL == "/api/php1/get":
			a.Equal(v.Method, "get")
		case v.URL == "/api/php2/post":
			a.Equal(v.Method, "post")
			a.Equal(v.Group, "php2")
		}
	}
}
