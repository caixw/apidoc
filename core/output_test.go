// Copyright 2015 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package core

import (
	"os"
	"testing"

	"github.com/issue9/assert"
)

var (
	block1 = `
@api 获取所有用户信息
获取所有用户的详细信息，包括用户所属的权限组，昵称等。
若没有权限，则返回空对象。
@apiURL /api/users
@apiMethods get
@apiGroup users
@apiQuery status string optional 只获取指定状态的数据，可用值为normal, locked
@apiStatus 200 json 成功获取用户信息
@apiParam users object 表示所有的用户列表
@apiExample json
{"users":[
	{"id":1, "name": "n1", "group": 1},
	{"id":2, "name": "n2", "group": 1},
	{"id":3, "name": "n3", "group": 1},
]}
@apiStatus 401 none 权限不足
`
	block2 = `
@api 获取指定用户的详细信息
@apiURL /api/users/{id}
@apiMethods get
@apiParam id int 用户的ID值
@apiGroup users
@apiStatus 200 json 成功获取信息
@apiParam id int 用户的ID
@apiParam name string 用户名称
@apiParam group int 用户所在的权限组ID
@apiExample json
{"id":1, "name": "n1", "group": 1}
`
	block3 = `
@api 请求登录用户
@apiURL /api/auth/login
@apiMethods post
@apiGroup auth
@apiRequest json
@apiParam username string 登录用户名
@apiParam password string 登录密码
@apiExample json
{"username": "admin", "password": "admin"}
@apiStatus 200 成功登录
@apiHeader token xxx
@apiStatus 401 none 权限不足
`
)

func TestLexer_OutputHtml(t *testing.T) {
	testdir := "./testdir"
	a := assert.New(t)

	// 创建测试目录
	a.NotError(os.MkdirAll(testdir, os.ModePerm), "无法创建测试目录")

	tree := NewTree()
	a.NotNil(tree)

	a.NotError(tree.scan([]byte(block1), 1, "test1.go"))
	a.NotError(tree.scan([]byte(block2), 100, "test1.go"))
	a.NotError(tree.scan([]byte(block3), 1, "test2.go"))

	a.NotError(tree.OutputHtml(testdir))
}
