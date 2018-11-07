// Copyright 2018 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package doc

// Callback 回调内容
type Callback struct {
	responses
	requests
	Method  string   `yaml:"method" json:"method"`
	Queries []*Param `yaml:"queries,omitempty" json:"queries,omitempty"` // 查询参数
	Params  []*Param `yaml:"params,omitempty" json:"params,omitempty"`   // URL 参数
}
