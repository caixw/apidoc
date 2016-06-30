// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package doc

import (
	"testing"

	"github.com/issue9/assert"
)

// go1.6 BenchmarkDoc_Scan-4    	   50000	     33875 ns/op
func BenchmarkDoc_Scan(b *testing.B) {
	code := `
@api get /baseurl/api/login api summary
api description 1
api description 2
@apiGroup users
@apiQuery q1 int q1 summary
@apiQuery q2 int q2 summary
@apiParam p1 int p1 summary
@apiParam p2 int p2 summary
@apiSuccess 200 json
@apiHeader h1 v1
@apiHeader h2 v2
@apiParam p1 int optional p1 summary
@apiParam p2 int p2 summary
@apiExample json
{
    p1:v1,
    p2:v2
}
@apiExample xml
<root>
    <p1>v1</p1>
    <p2>v2</p2>
</root>
@apiError 200 json
@apiHeader h1 v1
@apiHeader h2 v2
`

	doc := New()
	for i := 0; i < b.N; i++ {
		err := doc.Scan([]rune(code))
		if err != nil {
			b.Error("BenchmarkLexer_scan:error")
		}
	}
}

// go1.6 BenchmarkTag_readWord-4	10000000	       131 ns/op
func BenchmarkTag_readWord(b *testing.B) {
	a := assert.New(b)
	t := &tag{data: []rune("line1\n @delimiter line2 \n")}
	a.NotNil(t)

	for i := 0; i < b.N; i++ {
		_ = t.readWord()
		t.pos = 0
	}
}

// go1.6 BenchmarkTag_readLine-4	10000000	       109 ns/op
func BenchmarkTag_readLine(b *testing.B) {
	a := assert.New(b)
	t := &tag{data: []rune("line1\n @delimiter line2 \n")}
	a.NotNil(t)

	for i := 0; i < b.N; i++ {
		_ = t.readLine()
		t.pos = 0
	}
}

// go1.6 BenchmarkTag_readEnd-4 	10000000	       181 ns/op
func BenchmarkTag_readEnd(b *testing.B) {
	a := assert.New(b)
	t := &tag{data: []rune("line1\n line2 \n")}
	a.NotNil(t)

	for i := 0; i < b.N; i++ {
		_ = t.readEnd()
		t.pos = 0
	}
}

// go1.6 BenchmarkNewLexer-4    	300000000	         5.63 ns/op
func BenchmarkNewLexer(b *testing.B) {
	data := []rune("line")
	for i := 0; i < b.N; i++ {
		l := newLexer(data)
		if l.atEOF() {
		}
	}
}
