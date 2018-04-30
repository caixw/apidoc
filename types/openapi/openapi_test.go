// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package openapi

import (
	"testing"

	"github.com/issue9/assert"
	"github.com/issue9/version"
)

func TestLatestVersion(t *testing.T) {
	a := assert.New(t)

	a.True(version.SemVerValid(latestVersion))
}

func TestOpenAPI_Sanitize(t *testing.T) {
	a := assert.New(t)

	oa := &OpenAPI{Info: &Info{
		Title:   "title",
		Version: "3.3.3",
	}}
	a.NotError(oa.Sanitize())
	a.Equal(oa.OpenAPI, latestVersion)
}
