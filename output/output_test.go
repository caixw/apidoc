// SPDX-License-Identifier: MIT

package output

import (
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v5/doc"
)

func getTestDoc() *doc.Doc {
	return &doc.Doc{
		Tags: []*doc.Tag{{Name: "t1"}, {Name: "t2"}},
		Apis: []*doc.API{
			{Tags: []string{"t1", "tag1"}},
			{Tags: []string{"t2", "tag2"}},
		},
	}
}

func TestRender(t *testing.T) {
	a := assert.New(t)
	doc := getTestDoc()
	o := &Options{
		Path: "./apidoc.xml",
	}

	a.NotError(Render(doc, o))
}

func TestBuffer(t *testing.T) {
	a := assert.New(t)
	doc := getTestDoc()

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

	d := getTestDoc()
	o := &Options{}
	filterDoc(d, o)
	a.Equal(2, len(d.Tags))

	d = getTestDoc()
	o = &Options{
		Tags: []string{"t1"},
	}
	filterDoc(d, o)
	a.Equal(1, len(d.Tags)).
		Equal(1, len(d.Apis))

	d = getTestDoc()
	o = &Options{
		Tags: []string{"t1", "t2"},
	}
	filterDoc(d, o)
	a.Equal(2, len(d.Tags)).
		Equal(2, len(d.Apis))

	d = getTestDoc()
	o = &Options{
		Tags: []string{"tag1"},
	}
	filterDoc(d, o)
	a.Equal(0, len(d.Tags)).
		Equal(1, len(d.Apis))

	d = getTestDoc()
	o = &Options{
		Tags: []string{"not-exists"},
	}
	filterDoc(d, o)
	a.Equal(0, len(d.Tags)).
		Equal(0, len(d.Apis))
}
