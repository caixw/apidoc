// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package example

import "net/http"

// @api post /login 用户登录
//
// @apiRequest json
// @apiParam username string 用户名
// @apiParam password string 密码
// @apiExample json
// {
//     "username": "admin",
//     "password": "123456"
// }
//
// @apiSuccess 201 成功登录
// @apiParam expries int 过期的时间，单位秒
// @apiParam token string 存储 token
//
// @apiError 401 账号密码验证错误
func login(w http.ResponseWriter, r *http.Request) {
	// TODO
}

// @api delete /login 注销用户
//
// @apiRequest json
// @apiHeader Authorization 当前登录用户的 token
//
// @apiSuccess 204 注销成功
func logout(w http.ResponseWriter, r *http.Request) {
	// TODO
}
