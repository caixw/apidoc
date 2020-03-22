// SPDX-License-Identifier: MIT

package spec

import (
	"bytes"
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v6/core"
)

func TestBegin(t *testing.T) {
	a := assert.New(t)

	a.True(bytes.HasPrefix([]byte("<apidoc "), apidocBegin))
	a.True(bytes.HasPrefix([]byte("<apidoc\t"), apidocBegin))
	a.True(bytes.HasPrefix([]byte("<apidoc\n"), apidocBegin))

	a.True(bytes.HasPrefix([]byte("<api "), apiBegin))
	a.True(bytes.HasPrefix([]byte("<api\t"), apiBegin))
	a.True(bytes.HasPrefix([]byte("<api\n"), apiBegin))
}

func TestDoc_fromXML(t *testing.T) {
	a := assert.New(t)
	doc := NewAPIDoc()
	a.NotNil(doc)

	data := []byte(`<apidoc version="x.0.1"></apidoc>`)
	loc := core.Location{
		Range: core.Range{Start: core.Position{Line: 11, Character: 12}},
	}
	err := doc.fromXML(&Block{Location: loc, Data: data})
	a.Equal(err.(*core.SyntaxError).Location.Range.Start, core.Position{Line: 11, Character: 12})

	data = []byte(`<apidoc
	
	version="x.1.1">
	</apidoc>`)
	loc = core.Location{
		Range: core.Range{Start: core.Position{Line: 12, Character: 21}},
	}
	err = doc.fromXML(&Block{Location: loc, Data: data})
	a.Equal(err.(*core.SyntaxError).Location.Range.Start, core.Position{Line: 14, Character: 1})
}

func TestDoc_appendAPI(t *testing.T) {
	a := assert.New(t)
	doc := loadDoc(a)
	loc := core.Location{
		URI: "file:///file.php",
		Range: core.Range{
			Start: core.Position{
				Line:      11,
				Character: 22,
			},
			End: core.Position{},
		},
	}

	data := []byte(`<api version="x.0.1"></api>`)
	err := doc.appendAPI(&Block{Location: loc, Data: data})
	a.Equal(err.(*core.SyntaxError).Location.Range.Start, core.Position{Line: 11, Character: 27})

	data = []byte(`<api version="0.1.1">

	    <callback method="not-exists" />
	</api>`)
	err = doc.appendAPI(&Block{Location: loc, Data: data})
	a.Equal(err.(*core.SyntaxError).Location.Range.Start, core.Position{Line: 13, Character: 5})
}
