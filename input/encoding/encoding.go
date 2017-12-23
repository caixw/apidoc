// Copyright 2017 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package encoding

import (
	"bytes"
	"errors"
	"io/ioutil"
	"strings"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
)

// DefaultEncoding 默认的编码名称，只能是 utf-8
// 这里给出一个常量，是方便给其它包引用，统一名称。
var DefaultEncoding = "utf8"

// 一个编码名称与解码器的关联。
// 若解码器为空，表示不需要解码。
var encodings = map[string]encoding.Encoding{
	DefaultEncoding: nil,
	"utf-8":         nil,

	"gbk":     simplifiedchinese.GBK,
	"gb18083": simplifiedchinese.GB18030,
	"gb2312":  simplifiedchinese.HZGB2312,

	"big5": traditionalchinese.Big5,
}

// Encodings 返回支持的编码列表
func Encodings() []string {
	ret := make([]string, 0, len(encodings))

	for name := range encodings {
		ret = append(ret, name)
	}

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
		return nil, errors.New("not found")
	}

	if enc == nil {
		return data, nil
	}

	reader := transform.NewReader(bytes.NewReader(data), enc.NewDecoder())
	ret, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return ret, nil
}
