// SPDX-License-Identifier: MIT

package output

import (
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v6/doc/doctest"
)

func TestRender(t *testing.T) {
	a := assert.New(t)
	doc := doctest.Get()
	o := &Options{
		Path: "./apidoc.xml",
	}

	a.NotError(Render(doc, o))
}

func TestRender_openapiJSON(t *testing.T) {
	a := assert.New(t)
	doc := doctest.Get()

	o := &Options{
		Type: OpenapiJSON,
		Path: "./openapi.json",
	}

	a.NotError(Render(doc, o))
}

func TestBuffer(t *testing.T) {
	a := assert.New(t)
	doc := doctest.Get()

	o := &Options{}
	buf, err := Buffer(doc, o)
	a.NotError(err).NotNil(buf)

	// 返回错误
	o = &Options{Type: "not-exists"}
	buf, err = Buffer(doc, o)
	a.Error(err).Nil(buf)
}

func TestFilterDoc(t *testing.T) {
	a := assert.New(t)

	d := doctest.Get()
	o := &Options{}
	filterDoc(d, o)
	a.Equal(3, len(d.Tags))

	d = doctest.Get()
	o = &Options{
		Tags: []string{"t1"},
	}
	filterDoc(d, o)
	a.Equal(1, len(d.Tags)).
		Equal(2, len(d.Apis))

	d = doctest.Get()
	o = &Options{
		Tags: []string{"t1", "t2"},
	}
	filterDoc(d, o)
	a.Equal(2, len(d.Tags)).
		Equal(2, len(d.Apis))

	d = doctest.Get()
	o = &Options{
		Tags: []string{"tag1"},
	}
	filterDoc(d, o)
	a.Equal(1, len(d.Tags)).
		Equal(1, len(d.Apis))

	d = doctest.Get()
	o = &Options{
		Tags: []string{"not-exists"},
	}
	filterDoc(d, o)
	a.Equal(0, len(d.Tags)).
		Equal(0, len(d.Apis))
}
