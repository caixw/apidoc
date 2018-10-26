// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package schema

import (
	"bufio"
	"bytes"
	"strconv"
	"unicode"

	"github.com/caixw/apidoc/doc/lexer"
)

// 分析类型的内容。值可以有以下格式：
//  - type 单一类型
//  - type.subtype 集合类型，subtype 表示集全元素的类型，一般用于数组。
func parseType(tag *lexer.Tag, typ []byte) (t1, t2 string, err error) {
	types := bytes.SplitN(typ, seqaratorDot, 2)
	if len(types) == 0 {
		return "", "", tag.ErrInvalidFormat()
	}

	type0 := string(types[0])
	if type0 != Array {
		return type0, "", nil
	}

	if len(types) == 1 {
		return "", "", tag.ErrInvalidFormat()
	}

	return type0, string(types[1]), nil
}

// 分析可选类型，格式如下
//  optional
//  required
//  optional.defaultvalue
func parseOptional(typ, subtype string, optional []byte) (opt bool, val interface{}, err error) {
	index := bytes.IndexByte(optional, '.')
	if index < 0 {
		return isOptional(optional), nil, nil
	}

	opt = isOptional(optional[:index])

	optional = optional[index+1:]
	if typ == Array {
		fn := getConvertFunc(subtype)
		data := parseArray(optional)
		vals := make([]interface{}, 0, len(data))
		for _, item := range data {
			v, err := fn(string(item))
			if err != nil {
				return false, nil, err
			}
			vals = append(vals, v)
		}
		val = vals
	} else {
		fn := getConvertFunc(typ)
		val, err = fn(string(optional))
		if err != nil {
			return false, nil, err
		}
	}

	return opt, val, nil
}

var requiredBytes = []byte("required")

func isOptional(optional []byte) bool {
	return !bytes.Equal(bytes.ToLower(optional), requiredBytes)
}

// 解析数组
//  "[a1,a2,a3]" ==> {"a1","a2","a3"}
func parseArray(optional []byte) [][]byte {
	optional = bytes.TrimFunc(optional, func(r rune) bool {
		return r == '[' || r == ']'
	})

	ret := make([][]byte, 0, bytes.Count(optional, []byte{','}))
LOOP:
	for {
		index := bytes.IndexByte(optional, ',')
		switch {
		case index < 0:
			ret = append(ret, bytes.TrimSpace(optional))
			break LOOP
		case index > 0:
			ret = append(ret, bytes.TrimSpace(optional[:index]))
			optional = optional[index+1:]
		case index == 0:
			ret = append(ret, []byte{})
			optional = optional[1:]
		}
	}

	return ret
}

// 分析枚举内容
//
// 支持类似的的格式：
//  - s1 xxxx
//  - s2 xx
//  * s3 xxxx
// 将返回 s1,s2,s3
func parseEnum(data []byte) []string {
	enum := make([]string, 0, 5)

	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Split(bufio.ScanLines)

LOOP:
	for scanner.Scan() {
		line := scanner.Bytes()
		line = bytes.TrimSpace(line)

		// 过滤非 - 和 * 开头的行
		if (len(line) == 0) || (line[0] != '*' && line[0] != '-') {
			continue
		}

		// 去掉 * - 和空格
		line = bytes.TrimLeftFunc(line, func(r rune) bool {
			return r == '-' || r == '*' || unicode.IsSpace(r)
		})

		// 拿到首单词
		str := string(line)
		for index, c := range str {
			if unicode.IsSpace(c) {
				enum = append(enum, str[:index])
				continue LOOP
			}
		}
	}

	return enum
}

func convertEnumType(enum []string, typ string) ([]interface{}, error) {
	fn := getConvertFunc(typ)

	ret := make([]interface{}, 0, len(enum))
	for _, elem := range enum {
		v, err := fn(elem)
		if err != nil {
			return nil, err
		}
		ret = append(ret, v)
	}

	return ret, nil
}

func getConvertFunc(typ string) convertFunc {
	fn, found := converters[typ]
	if !found {
		return stringConvert
	}
	return fn
}

// 类型转换函数定义
type convertFunc func(val string) (interface{}, error)

var (
	numberConvert = func(v string) (interface{}, error) {
		return strconv.ParseInt(v, 10, 64)
	}
	stringConvert = func(v string) (interface{}, error) {
		return v, nil
	}
	boolConvert = func(v string) (interface{}, error) {
		return strconv.ParseBool(v)
	}
	converters = map[string]convertFunc{
		Number:  numberConvert,
		Integer: numberConvert,
		String:  stringConvert,
		Bool:    boolConvert,
	}
)
