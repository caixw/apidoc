// SPDX-License-Identifier: MIT

package build

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/issue9/sliceutil"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/ianaindex"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/lang"
	"github.com/caixw/apidoc/v7/internal/locale"
)

// Input 指定输入内容的相关信息。
type Input struct {
	Lang      string   `yaml:"lang"`                // 输入的目标语言，值为 internal/lang 中的 Language.ID
	Dir       core.URI `yaml:"dir"`                 // 源代码目录
	Exts      []string `yaml:"exts,omitempty"`      // 需要扫描的文件扩展名，为空则表示采用默认规则。
	Recursive bool     `yaml:"recursive,omitempty"` // 是否查找 Dir 的子目录
	Encoding  string   `yaml:"encoding,omitempty"`  // 源文件的编码，默认为 UTF-8
	Ignores   []string `yaml:"ignores,omitempty"`   // 忽略的文件或目录，比如 node_modules 等可在此指定

	paths     []core.URI        // 根据 Dir、Exts、Ignores 和 Recursive 生成
	encoding  encoding.Encoding // 根据 Encoding 生成
	sanitized bool
}

func (o *Input) sanitize() error {
	if o.sanitized {
		return nil
	}

	if len(o.Dir) == 0 {
		return core.NewError(locale.ErrIsEmpty, "dir").WithField("dir")
	}

	exists, err := o.Dir.Exists()
	if err != nil {
		return core.WithError(err).WithField("dir")
	}
	if !exists {
		return core.NewError(locale.ErrDirNotExists).WithField("dir")
	}

	if scheme, _ := o.Dir.Parse(); scheme != "" && scheme != core.SchemeFile {
		return core.NewError(locale.ErrInvalidURIScheme, scheme).WithField("dir")
	}

	if len(o.Lang) == 0 {
		return core.NewError(locale.ErrIsEmpty, "lang").WithField("lang")
	}

	language := lang.Get(o.Lang)
	if language == nil {
		return core.NewError(locale.ErrInvalidValue).WithField("lang")
	}

	if len(o.Exts) > 0 {
		exts := make([]string, 0, len(o.Exts))
		for _, ext := range o.Exts {
			if len(ext) == 0 {
				continue
			}

			if ext[0] != '.' {
				ext = "." + ext
			}
			exts = append(exts, ext)
		}
		o.Exts = exts
	} else {
		o.Exts = language.Exts
	}

	if err = o.recursivePath(); err != nil {
		return err
	}

	if o.Encoding != "" {
		o.encoding, err = ianaindex.IANA.Encoding(o.Encoding)
		if err != nil {
			return core.WithError(err).WithField("encoding")
		}
	}

	o.sanitized = true
	return nil
}

// 按 Input 中的规则查找所有符合条件的文件列表并保存至 Input.paths
func (o *Input) recursivePath() error {
	local, err := o.Dir.File()
	if err != nil {
		return core.WithError(err).WithField("dir")
	}
	local = filepath.Clean(local)

	for i, pattern := range o.Ignores {
		o.Ignores[i] = filepath.FromSlash(pattern)
	}

	walk := func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fi.IsDir() && !o.Recursive && path != local {
			return filepath.SkipDir
		}

		ignore, err := o.isIgnore(local, path)
		if err != nil {
			return err
		}
		if !ignore {
			o.paths = append(o.paths, core.FileURI(path))
		}
		return nil
	}

	if err := filepath.Walk(local, walk); err != nil {
		return core.WithError(err).WithField("dir")
	}

	if len(o.paths) == 0 {
		return core.NewError(locale.ErrNoFiles).WithField("dir")
	}
	return nil
}

func (o *Input) isIgnore(root, path string) (bool, error) {
	ext := filepath.Ext(path)
	if sliceutil.Count(o.Exts, func(i string) bool { return i == ext }) == 0 {
		return true, nil
	}

	path = strings.TrimPrefix(filepath.FromSlash(path), filepath.FromSlash(root))
	if path == "" {
		return false, nil
	}

	if path[0] == '/' || path[0] == os.PathSeparator {
		path = path[1:]
	}

	for _, ignore := range o.Ignores {
		match, err := filepath.Match(ignore, path)
		if err != nil {
			return false, err
		}
		if match {
			return true, nil
		}
	}

	return false, nil
}

// ParseInputs 分析 opt 中所指定的内容并输出到 blocks
//
// 分析后的内容推送至 blocks 中。
func ParseInputs(blocks chan core.Block, h *core.MessageHandler, opt ...*Input) {
	wg := &sync.WaitGroup{}
	for _, i := range opt {
		for _, path := range i.paths {
			wg.Add(1)
			go func(path core.URI, i *Input) {
				i.ParseFile(blocks, h, path)
				wg.Done()
			}(path, i)
		}
	}
	wg.Wait()
}

// ParseFile 分析 uri 指向的文件并输出到 blocks
func (o *Input) ParseFile(blocks chan core.Block, h *core.MessageHandler, uri core.URI) {
	data, err := uri.ReadAll(o.encoding)
	if err != nil {
		h.Error((core.Location{URI: uri}).WithError(err))
		return
	}

	lang.Parse(h, o.Lang, core.Block{
		Data:     data,
		Location: core.Location{URI: uri},
	}, blocks)
}
