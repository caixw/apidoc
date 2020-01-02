// SPDX-License-Identifier: MIT

package openapi

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v5/doc"
	"github.com/caixw/apidoc/v5/internal/vars"
)

func getTestDoc() *doc.Doc {
	return &doc.Doc{
		Title: "test",
		Description: doc.Richtext{
			Text: "<p>desc</p>",
			Type: doc.RichtextTypeHTML,
		},
		Tags: []*doc.Tag{{Name: "t1"}, {Name: "t2"}},
		Apis: []*doc.API{
			{
				Method: http.MethodGet,
				Tags:   []string{"t1", "t2"},
				Path:   &doc.Path{Path: "/users"},
			},
			{
				Method:     http.MethodPost,
				Tags:       []string{"t1", "t2"},
				Path:       &doc.Path{Path: "/users"},
				Deprecated: "1.0.1",
				Summary:    "summary",
			},
		},
	}
}

func TestJSON(t *testing.T) {
	a := assert.New(t)
	data, err := JSON(getTestDoc())
	a.NotError(err).NotNil(data)

	openapi := &OpenAPI{}
	a.NotError(json.Unmarshal(data, openapi)).
		Equal(2, len(openapi.Tags)).
		Equal(1, len(openapi.Paths)).
		Equal(openapi.ExternalDocs.URL, vars.OfficialURL).
		NotEmpty(openapi.ExternalDocs.Description)

	path := openapi.Paths["/users"]
	a.NotError(path)
	a.NotNil(path.Post).NotNil(path.Get).Nil(path.Patch)
	a.True(path.Post.Deprecated)
	a.Equal(path.Post.Summary, "summary")
}

func TestYAML(t *testing.T) {
	a := assert.New(t)
	data, err := YAML(getTestDoc())
	a.NotError(err).NotNil(data)
}
