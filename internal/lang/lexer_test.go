// SPDX-License-Identifier: MIT

package lang

import (
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v6/core"
	"github.com/caixw/apidoc/v6/core/messagetest"
)

func TestLexer_block(t *testing.T) {
	a := assert.New(t)

	blocks := []Blocker{
		newCStyleSingleComment(),
		newCStyleMultipleComment(),
		newRubyMultipleComment("=pod", "=cut", ""),
		newCStyleString(),
	}

	data := []byte(`// scomment1
  // scomment2
func(){}
"/*中文1"
行首代码"//中文2"
 /*
mcomment1
mcomment2
*/

中文// scomment3
// scomment4
=pod
 mcomment3
 mcomment4
=cut
`)
	l, err := NewLexer(data, blocks)
	a.NotError(err).NotNil(l)

	b, pos := l.block() // scomment1
	a.NotNil(b)
	a.Equal(pos, core.Position{Line: 0, Character: 0})
	_, ok := b.(*singleComment)
	a.True(ok)
	raw, data, ok := b.EndFunc(l)
	a.True(ok).
		Equal(string(raw), "// scomment1\n  // scomment2\n").
		Equal(string(data), "   scomment1\n     scomment2\n")

	b, pos = l.block() // 中文1
	a.NotNil(b)
	a.Equal(pos, core.Position{Line: 3, Character: 0})
	_, ok = b.(*stringBlock)
	a.True(ok)
	_, _, ok = b.EndFunc(l)
	a.True(ok)

	b, pos = l.block() // 中文2
	a.NotNil(b)
	a.Equal(pos, core.Position{Line: 4, Character: 4})
	_, ok = b.(*stringBlock)
	a.True(ok)
	_, _, ok = b.EndFunc(l)
	a.True(ok)

	b, pos = l.block()
	a.NotNil(b)
	a.Equal(pos, core.Position{Line: 5, Character: 1})
	_, ok = b.(*multipleComment)
	a.True(ok)
	raw, data, ok = b.EndFunc(l)
	a.NotError(ok).
		Equal(string(raw), "/*\nmcomment1\nmcomment2\n*/").
		Equal(string(data), "  \nmcomment1\nmcomment2\n  ")

	/* 测试一段单行注释后紧跟 \n=pod 形式的多行注释，是否会出错 */

	b, pos = l.block() // scomment3,scomment4
	a.NotNil(b)
	a.Equal(pos, core.Position{Line: 10, Character: 2})
	_, ok = b.(*singleComment)
	a.True(ok)
	raw, data, ok = b.EndFunc(l)
	a.True(ok).
		Equal(string(raw), "// scomment3\n// scomment4\n").
		Equal(string(data), "   scomment3\n   scomment4\n")

	b, pos = l.block() // mcomment3,mcomment4
	a.NotNil(b)
	a.Equal(pos, core.Position{Line: 12, Character: 0})
	_, ok = b.(*rubyMultipleComment)
	a.True(ok)
	raw, data, ok = b.EndFunc(l)
	a.True(ok).
		Equal(string(raw), "=pod\n mcomment3\n mcomment4\n=cut\n").
		Equal(string(data), "      mcomment3\n mcomment4\n     ")
}

func TestLexer_Parse(t *testing.T) {
	a := assert.New(t)

	raw := `// <api method="GET">
// <path path="/apis/gbk" />
// <description>1223 中文 45 </description>
// <server>test</server>
// </api>
`

	blocks := make(chan core.Block, 100)
	erro, _, h := messagetest.MessageHandler()
	l, err := NewLexer([]byte(raw), cStyle)
	a.NotError(err).NotNil(l)
	l.Parse(blocks, h, core.URI("./testdata/gbk.php"))
	h.Stop()
	close(blocks)
	a.Equal(1, len(blocks))
	blk := <-blocks
	a.Equal(string(blk.Data), `   <api method="GET">
   <path path="/apis/gbk" />
   <description>1223 中文 45 </description>
   <server>test</server>
   </api>
`).
		Equal(string(blk.Raw), raw).
		Equal(blk.Location.Range, core.Range{
			Start: core.Position{Line: 0, Character: 0},
			End:   core.Position{Line: 5, Character: 0},
		})
	a.Empty(erro.String())

	// 没有正确的结束符号
	raw = `/* <api method="GET">
// <path path="/apis/gbk" />
// <description>1223 中文 45 </description>
// <server>test</server>
// </api>
`
	blocks = make(chan core.Block, 100)
	erro, _, h = messagetest.MessageHandler()
	l, err = NewLexer([]byte(raw), cStyle)
	a.NotError(err).NotNil(l)
	l.Parse(blocks, h, core.URI("./testdata/gbk.php"))
	h.Stop()
	close(blocks)
	a.NotEmpty(erro.String())
}
