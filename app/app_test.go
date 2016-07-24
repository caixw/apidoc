// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package app

import (
	"strings"
	"testing"

	"golang.org/x/text/language"

	"github.com/issue9/assert"
	"github.com/issue9/is"
	"github.com/issue9/version"
)

// 对一些堂量的基本检测。
func TestConsts(t *testing.T) {
	a := assert.New(t)

	a.True(version.SemVerValid(Version))
	a.True(len(Name) > 0)
	a.True(is.URL(RepoURL))
	a.True(is.URL(OfficialURL))
	a.True(len(ConfigFilename) > 0).True(strings.IndexAny(ConfigFilename, "/\\") < 0)
	a.True(len(DefaultTitle) > 0)
	a.True(len(DefaultGroupName) > 0).True(strings.IndexAny(DefaultGroupName, "/\\") < 0)
	a.True(len(Profile) > 0).True(strings.IndexAny(Profile, "/\\") < 0)

	tag, err := language.Parse(DefaultLocale)
	a.NotError(err).NotEqual(tag, language.Und)
}
