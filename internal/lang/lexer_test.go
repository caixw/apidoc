// SPDX-License-Identifier: MIT

package lang

import (
	"testing"

	"github.com/issue9/assert"

	"github.com/caixw/apidoc/v6/core"
	"github.com/caixw/apidoc/v6/core/messagetest"
	"github.com/caixw/apidoc/v6/spec"
)

func TestLexer_Position(t *testing.T) {
	a := assert.New(t)

	l, err := NewLexer([]byte("l0\nl1\nl2\nl3\n"), nil)
	a.NotError(err).NotNil(l)

	l.next(3)
	a.Equal(l.current, position{
		Position: core.Position{Line: 1, Character: 0},
		Offset:   3,
	})

	l.next(4)
	a.Equal(l.current, position{
		Position: core.Position{Line: 2, Character: 1},
		Offset:   7,
	})

	l.next(3)
	l.next(3)
	a.Equal(l.current, position{
		Position: core.Position{Line: 4, Character: 0},
		Offset:   12,
	})

	l, err = NewLexer([]byte("12中文ab"), nil)
	a.NotError(err).NotNil(l)
	l.next(2)
	a.Equal(l.current, position{
		Position: core.Position{Line: 0, Character: 2},
		Offset:   2,
	})

	l.next(2)
	a.Equal(l.current, position{
		Position: core.Position{Line: 0, Character: 4},
		Offset:   8,
	})
}

func TestLexer_match(t *testing.T) {
	a := assert.New(t)

	l := &Lexer{
		data: []byte("ab中\ncd"),
	}

	a.False(l.match("b")).Equal(0, l.current.Offset)
	a.True(l.match("ab")).Equal(2, l.current.Offset)
	a.False(l.match("ab")).Equal(2, l.current.Offset)
	a.True(l.match("中")).Equal(l.current, position{
		Position: core.Position{Line: 0, Character: 3},
		Offset:   5,
	})

	l.back()
	a.True(l.match("中")).Equal(l.current, position{
		Position: core.Position{Line: 0, Character: 3},
		Offset:   5,
	})

	l.next(len(l.data))
	a.False(l.match("ab")).Equal(l.current, position{
		Position: core.Position{Line: 1, Character: 2},
		Offset:   8,
	})
}

func TestLexer_block(t *testing.T) {
	a := assert.New(t)

	blocks := []Blocker{
		newCStyleSingleComment(),
		newCStyleMultipleComment(),
		newRubyMultipleComment("=pod", "=cut", ""),
		newCStyleString(),
	}

	l := &Lexer{
		data: []byte(`// scomment1
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
`),
		blocks: blocks,
	}

	b, pos := l.block() // scomment1
	a.NotNil(b)
	a.Equal(pos, core.Position{Line: 0, Character: 0})
	_, ok := b.(*singleComment)
	a.True(ok)
	raw, data, err := b.EndFunc(l)
	a.NotError(err).
		Equal(string(raw), "// scomment1\n  // scomment2\n").
		Equal(string(data), "   scomment1\n     scomment2\n")

	b, pos = l.block() // 中文1
	a.NotNil(b)
	a.Equal(pos, core.Position{Line: 3, Character: 0})
	_, ok = b.(*stringBlock)
	a.True(ok)
	_, _, err = b.EndFunc(l)
	a.NotError(err)

	b, pos = l.block() // 中文2
	a.NotNil(b)
	a.Equal(pos, core.Position{Line: 4, Character: 4})
	_, ok = b.(*stringBlock)
	a.True(ok)
	_, _, err = b.EndFunc(l)
	a.NotError(err)

	b, pos = l.block()
	a.NotNil(b)
	a.Equal(pos, core.Position{Line: 5, Character: 1})
	_, ok = b.(*multipleComment)
	a.True(ok)
	raw, data, err = b.EndFunc(l)
	a.NotError(err).
		Equal(string(raw), "/*\nmcomment1\nmcomment2\n*/").
		Equal(string(data), "  \nmcomment1\nmcomment2\n  ")

	/* 测试一段单行注释后紧跟 \n=pod 形式的多行注释，是否会出错 */

	b, pos = l.block() // scomment3,scomment4
	a.NotNil(b)
	a.Equal(pos, core.Position{Line: 10, Character: 2})
	_, ok = b.(*singleComment)
	a.True(ok)
	raw, data, err = b.EndFunc(l)
	a.NotError(err).
		Equal(string(raw), "// scomment3\n// scomment4\n").
		Equal(string(data), "   scomment3\n   scomment4\n")

	b, pos = l.block() // mcomment3,mcomment4
	a.NotNil(b)
	a.Equal(pos, core.Position{Line: 12, Character: 0})
	_, ok = b.(*rubyMultipleComment)
	a.True(ok)
	raw, data, err = b.EndFunc(l)
	a.NotError(err).
		Equal(string(raw), "=pod\n mcomment3\n mcomment4\n=cut\n").
		Equal(string(data), "      mcomment3\n mcomment4\n     ")
}

func TestLexer_skipSpace(t *testing.T) {
	a := assert.New(t)
	l, err := NewLexer([]byte("    0 \n  1 "), nil)
	a.NotError(err).NotNil(l)

	l.skipSpace()
	a.Equal(string(l.next(1)), "0")

	// 无法跳过换行符
	l.next(1)
	l.skipSpace()
	l.skipSpace()
	l.skipSpace()
	l.skipSpace()
	a.Empty(l.skipSpace())
	a.Equal(string(l.next(1)), "\n")

	l.next(1)
	a.Equal(1, len(l.skipSpace()))
	l.back()
	a.Equal(1, len(l.skipSpace()))
	a.Equal(string(l.next(1)), "1")

	l.next(1)
	l.skipSpace()
	l.skipSpace()
	a.Equal(l.current.Offset, len(l.data))
}

func TestLexer_delim(t *testing.T) {
	a := assert.New(t)

	l, err := NewLexer([]byte("123"), nil)
	a.NotError(err).NotNil(l)
	a.Nil(l.delim('\n'))

	l, err = NewLexer([]byte("123\n"), nil)
	a.NotError(err).NotNil(l)
	a.Equal(string(l.delim('\n')), "123\n").
		Equal(l.current.Offset, 4)

	l = &Lexer{data: []byte("123\n"), current: position{Offset: 1}}
	a.Equal(string(l.delim('\n')), "23\n").
		Equal(l.current.Offset, 4)
}

func TestLexer_Parse(t *testing.T) {
	a := assert.New(t)

	raw := `// <api method="GET">
// <path path="/apis/gbk" />
// <description>1223 中文 45 </description>
// <server>test</server>
// </api>
`

	blocks := make(chan spec.Block, 100)
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
	blocks = make(chan spec.Block, 100)
	erro, _, h = messagetest.MessageHandler()
	l, err = NewLexer([]byte(raw), cStyle)
	a.NotError(err).NotNil(l)
	l.Parse(blocks, h, core.URI("./testdata/gbk.php"))
	h.Stop()
	close(blocks)
	a.NotEmpty(erro.String())
}
