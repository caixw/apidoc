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
