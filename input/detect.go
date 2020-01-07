// SPDX-License-Identifier: MIT

package input

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/caixw/apidoc/v6/internal/lang"
)

// Detect 检测指定目录下的内容，并为其生成一个合适的 input.Options 实例。
//
// 检测依据为根据扩展名来做统计，数量最大且被支持的获胜。
func Detect(dir string, recursive bool) ([]*Options, error) {
	exts, err := detectExts(dir, recursive)
	if err != nil {
		return nil, err
	}

	langs := detectLanguage(exts)

	opts := make([]*Options, 0, len(langs))
	for _, l := range langs {
		opts = append(opts, &Options{
			Lang:      l.Name,
			Dir:       dir,
			Exts:      l.Exts,
			Recursive: recursive,
		})
	}

	return opts, nil
}

type language struct {
	lang.Language
	count int
}

// 根据 exts 计算每个语言对应的文件数量，并按倒序返回
//
// exts 参数为从 detectExts 中获取的返回值
func detectLanguage(exts map[string]int) []*language {
	langs := make([]*language, 0, len(exts))

	for ext, count := range exts {
		l := lang.GetByExt(ext)
		if l == nil {
			continue
		}

		found := false
		for _, item := range langs {
			if item.Name == l.Name {
				item.count += count
				found = true
				break
			}
		}
		if !found {
			langs = append(langs, &language{
				count:    count,
				Language: *l,
			})
		}
	} // end for

	sort.SliceStable(langs, func(i, j int) bool {
		return langs[i].count > langs[j].count
	})

	return langs
}

// 返回 dir 目录下文件类型及对应的文件数量的一个集合。
// recursive 表示是否查找子目录。
func detectExts(dir string, recursive bool) (map[string]int, error) {
	exts := map[string]int{}

	walk := func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if fi.IsDir() {
			if !recursive && dir != path {
				return filepath.SkipDir
			}
		} else {
			ext := strings.ToLower(filepath.Ext(path))
			if len(ext) > 0 {
				exts[ext]++
			}
		}

		return nil
	}

	if err := filepath.Walk(dir, walk); err != nil {
		return nil, err
	}

	return exts, nil
}
