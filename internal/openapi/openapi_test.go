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

func TestOpenAPI_sanitize(t *testing.T) {
	a := assert.New(t)

	oa := &OpenAPI{Info: &Info{
		Title:   "title",
		Version: "3.3.3",
	},
		Paths: map[string]*PathItem{
			"/api": {
				Get: &Operation{
					Responses: map[string]*Response{
						"json": {},
					},
				},
			},
		},
	}
	a.NotError(oa.sanitize())
	a.Equal(oa.OpenAPI, LatestVersion)
}
