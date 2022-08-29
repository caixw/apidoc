// SPDX-License-Identifier: MIT

package openapi

import (
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	"github.com/issue9/assert/v3"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/ast/asttest"
)

func TestJSON(t *testing.T) {
	a := assert.New(t, false)
	data, err := JSON(asttest.Get())
	a.NotError(err).NotNil(data)

	openapi := &OpenAPI{}
	a.NotError(json.Unmarshal(data, openapi)).
		Equal(3, len(openapi.Tags)).
		Equal(1, len(openapi.Paths)).
		Equal(openapi.ExternalDocs.URL, core.OfficialURL).
		NotEmpty(openapi.ExternalDocs.Description)

	path := openapi.Paths["/users"]
	a.NotNil(path)
	a.NotNil(path.Post).NotNil(path.Get).Nil(path.Patch)
	a.True(path.Post.Deprecated)
	a.Equal(path.Post.Summary, "summary")
	a.NotNil(path.Get).NotNil(path.Post)

	get := path.Get
	a.Equal(1, len(get.Responses))
	a.Equal(len(get.Responses[strconv.Itoa(http.StatusOK)].Headers), 1)
}

func TestYAML(t *testing.T) {
	a := assert.New(t, false)
	data, err := YAML(asttest.Get())
	a.NotError(err).NotNil(data)
}
