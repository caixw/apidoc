// Copyright 2016 by caixw, All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package testdata;

public class test1{

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
    public void login() {
        System.out.println("/********** login");
        // TODO
    }

    // 123
    // 123
    /* @api DELETE /users/login 注销登录
    group users
    tags: [t1,t2]

    request:
    description: request body
    content:
        application/json:
        schema:
            type: object
            properties:
            username:
                type: string
                description: 登录账号
            password:
                type: string
                description: 密码
    */
    public void logout() {
        System.out.println("logout **********/");
        // TODO
    }
}
