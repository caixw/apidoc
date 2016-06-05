// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

"use strict"

// @api POST /users/login 登录
// @apiGroup users
//
// @apiRequest json
// @apiParam username string 登录账号
// @apiParam password string 密码
//
// @apiSuccess 201 OK
// @apiParam expires int 过期时间
// @apiParam token string 凭证
// @apiExample json
// {
//     "expires": 11111111,
//     "token": "adl;kfqwer;q;afd"
// }
//
// @apiError 401 账号或密码错误
function login(w, r) {
	println("/********** login");
	println(/login/.test('login'));
	// TODO
}

/* @api DELETE /users/login 注销登录
@apiGroup users

@apiRequest json
@apiHeader Authorization xxxx

@apiSuccess 201 OK
@apiParam expires int 过期时间
@apiParam token string 凭证
@apiExample json
{
    "expires": 11111111,
    "token": "adl;kfqwer;q;afd"
}
*/
function logout(w, r) {
	println('logout **********/')
	// TODO
}
