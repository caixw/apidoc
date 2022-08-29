// SPDX-License-Identifier: MIT

package openapi

import (
	"testing"

	"github.com/issue9/assert/v3"
	"github.com/issue9/version"
)

func TestLatestVersion(t *testing.T) {
	a := assert.New(t, false)

	a.True(version.SemVerValid(LatestVersion))
}

func TestOpenAPI_sanitize(t *testing.T) {
	a := assert.New(t, false)

	oa := &OpenAPI{Info: &Info{
		Title:   "title",
		Version: "3.3.3",
	},
		Paths: map[string]*PathItem{
			"/api": {
				Get: &Operation{
					Responses: map[string]*Response{
						"json": {
							Description: "desc",
						},
					},
				},
			},
		},
	}
	a.NotError(oa.sanitize())
	a.Equal(oa.OpenAPI, LatestVersion)
}

func TestExternalDocumentation_sanitize(t *testing.T) {
	a := assert.New(t, false)

	ed := &ExternalDocumentation{}
	a.Error(ed.sanitize())

	ed.URL = "url"
	a.Error(ed.sanitize())

	ed.URL = "https://example.com"
	a.NotError(ed.sanitize())
}

func TestTag_sanitize(t *testing.T) {
	a := assert.New(t, false)

	tag := &Tag{
		ExternalDocs: &ExternalDocumentation{},
	}
	a.Error(tag.sanitize())

	tag.Name = "name"
	a.Error(tag.sanitize())

	tag.ExternalDocs.URL = "https://example.com"
	a.NotError(tag.sanitize())
}
