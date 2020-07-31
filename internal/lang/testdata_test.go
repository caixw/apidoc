// SPDX-License-Identifier: MIT

package lang

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/issue9/assert"
	"github.com/issue9/sliceutil"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/core/messagetest"
)

// 当前文件用于测试各个语言的解析是否都正常
//
// 测试数据在 ./testdata 目录，按语言 ID 进行分类，
// 每个目录下可以包含多个同种语言的多个源文件，
// 各个源文件中只要包含 blocks 中定义的部分内容作为注释代码即可。

var blocks = []string{
	"   SPDX-License-Identifier: MIT\n",
	"    line1\n",
	"   \n   line1\n   line2\n   line3\n   ",
}

func TestLanguages(t *testing.T) {
	a := assert.New(t)

	for _, l := range Langs() {
		dir := filepath.Join("testdata", l.ID)
		if _, err := os.Stat(dir); err != nil && os.IsNotExist(err) {
			t.Errorf("未实现 %s 的测试", l.ID)
			continue
		}

		parseLang(a, l.ID, dir)
	}
}

func parseLang(a *assert.Assertion, id, dir string) {
	count := 0
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		blks := make(chan core.Block, 10)
		rslt := messagetest.NewMessageHandler()

		Parse(rslt.Handler, id, core.Block{
			Data:     data,
			Location: core.Location{URI: core.URI(path)},
		}, blks)

		rslt.Handler.Stop()
		a.Empty(rslt.Errors)

		exit := make(chan bool)
		close(blks)
		go func() {
			for block := range blks {
				cnt := sliceutil.Count(blocks, func(i int) bool {
					// 忽略首层空格的差别，部分语言没有多行注释，
					// 单行注释生成的代码块会比多行注释生成的代码多一个换行符。
					return strings.TrimSpace(blocks[i]) == strings.TrimSpace(string(block.Data))
				})

				if cnt <= 0 {
					a.True(false, "%s:%s 未找到对应的数据", id, string(block.Data))
				}
			}
			exit <- true
		}()
		<-exit

		count++
		return nil
	})
	a.NotError(err)

	if count == 0 {
		a.TB().Errorf("未实现 %s 的测试", id)
	}
}
