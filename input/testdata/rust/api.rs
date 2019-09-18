// SPDX-License-Identifier: MIT

mod std{

    /// @api POST /users/login 登录
    /// group users
    /// tags: [t1,t2]
    ///
    /// request:
    ///   description: request body
    ///   content:
    ///     application/json:
    ///       schema:
    ///         type: object
    ///         properties:
    ///           username:
    ///             type: string
    ///             description: 登录账号
    ///           password:
    ///             type: string
    ///             description: 密码
    fn login(w: http::ResponseWriter, r: &http::Request) {
        println("/********* login");
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
    fn logout(w: http::ResponseWriter, r: &http::Request) {
        println("logout  *******/");
        // TODO
    }
}
