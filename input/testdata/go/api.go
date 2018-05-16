// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package testdata

import "net/http"

// @api POST /users/login 登录
// group users
// tags: [t1,t2]
//
// request:
//   description: request body
//   content:
//     application/json:
//       schema:
//         type: object
//         properties:
//           username:
//             type: string
//             description: 登录账号
//           password:
//             type: string
//             description: 密码
// responses:
//   200
//     description:
//     headers:
//        header1: header1
//        header2: header2
//     content:
//       application/json:
//         schema:
//           type: object
//           properties:
//             token:
//               type: string
//               description: 登录账号
//             password:
//               type: string
//               description: 密码
func login(w http.ResponseWriter, r *http.Request) {
	println("/**********", "login")
	// TODO
}

// 123
// 123
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
func logout(w http.ResponseWriter, r *http.Request) {
	println("logout", "**********/")
	// TODO
}
