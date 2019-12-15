// SPDX-License-Identifier: MIT

package openapi

import (
	"testing"

	"github.com/issue9/assert"
	"github.com/issue9/version"
)

func TestLatestVersion(t *testing.T) {
	a := assert.New(t)

	a.True(version.SemVerValid(LatestVersion))
}

func TestOpenAPI_Sanitize(t *testing.T) {
	a := assert.New(t)

	oa := &OpenAPI{Info: &Info{
		Title:   "title",
		Version: "3.3.3",
	},
		Paths: map[string]*PathItem{
			"/api": {},
		},
	}
	a.NotError(oa.Sanitize())
	a.Equal(oa.OpenAPI, LatestVersion)
}
