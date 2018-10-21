// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

// Package encoding 为输入的文件提供编码支持。
package encoding

import (
	"bytes"
	"io/ioutil"
	"sort"
	"strings"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"

	"github.com/caixw/apidoc/locale"
)

// DefaultEncoding 默认的编码名称，只能是 utf-8。
//
// 这里给出一个常量，是方便给其它包引用时有一个统一名称。
const DefaultEncoding = "utf8"

// 编码名称与解码器的关联。若解码器为空，表示不需要解码。
//
// TODO: 采用 htmlindex 支持大部分的编码方式？
var encodings = map[string]encoding.Encoding{
	DefaultEncoding: nil,
	"utf-8":         nil,

	"gbk":     simplifiedchinese.GBK,
	"gb18030": simplifiedchinese.GB18030,
	"gb2312":  simplifiedchinese.HZGB2312,

	"big5": traditionalchinese.Big5,
}

// Encodings 返回支持的编码列表
func Encodings() []string {
	ret := make([]string, 0, len(encodings))
	for name := range encodings {
		ret = append(ret, name)
	}

	sort.SliceStable(ret, func(i, j int) bool {
		return ret[i] > ret[j]
	})

	return ret
}

// Transform 将 path 指向的文件内容，按 encoding 编码进行加载，
// 并转换成 utf-8 之后返回其内容。
func Transform(path, encoding string) ([]byte, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	enc, found := encodings[strings.ToLower(encoding)]
	if !found {
		return nil, locale.Errorf(locale.ErrUnsupportedEncoding, encoding)
	}

	if enc == nil {
		return data, nil
	}

	reader := transform.NewReader(bytes.NewReader(data), enc.NewDecoder())
	return ioutil.ReadAll(reader)
}
