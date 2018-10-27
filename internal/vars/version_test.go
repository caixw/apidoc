// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package vars

import (
	"testing"

	"github.com/issue9/assert"
	v "github.com/issue9/version"
)

// 对一些堂量的基本检测。
func TestVersion(t *testing.T) {
	a := assert.New(t)

	a.True(v.SemVerValid(mainVersion))
	a.True(v.SemVerValid(version))
	a.True(v.SemVerValid(Version()))
}
