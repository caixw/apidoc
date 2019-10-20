// SPDX-License-Identifier: MIT

// <api method="GET" summary="获取用户">
//     <description><![CDATA[
// 这是关于接口的详细说明文档
// 可以是一个 HTML 内容
// ]]></description>
//     <path path="/users">
//         <query name="page" type="number" default="0" summary="页码" />
//         <query name="size" type="number" default="20">数量</query>
//     </path>
//
//     <request type="none" mimetype="json">
//         <header name="name" type="string">desc</header>
//         <header name="name1" type="string">desc1</header>
//     </request>
//
//     <response status="200" array="true" type="object" mimetype="json">
//         <param name="count" type="number" optional="false" summary="summary" />
//         <param name="list" type="object" array="true" summary="list">
//             <param name="id" type="number" summary="用户 ID" />
//             <param name="name" type="string" summary="用户名" />
//             <param name="groups" type="string" array="true" optional="true" summary="用户所在的权限组">
//                 <param name="id" type="string" summary="权限组 ID" />
//                 <param name="name" type="string" summary="权限组名称" />
//             </param>
//         </param>
//     </response>
// </api>
fn getList() {}

// <api method="POST" summary="添加用户">
//     <description><![CDATA[
// 这是关于接口的详细说明文档
// 可以是一个 HTML 内容
// ]]></description>
//     <path path="/users" />
//
//     <request type="object" mimetype="json">
//         <param name="count" type="number" optional="false" summary="summary" />
//         <param name="list" type="object" array="true" summary="list">
//             <param name="id" type="number" summary="用户 ID" />
//             <param name="name" type="string" summary="用户名" />
//             <param name="groups" type="string" array="true" optional="true" summary="用户所在的权限组">
//                 <param name="id" type="string" summary="权限组 ID" />
//                 <param name="name" type="string" summary="权限组名称" />
//             </param>
//         </param>
//     </request>
//
//     <response status="200" array="true" type="none" mimetype="json">
//     </response>
// </api>
fn post() {}

// <api method="DELETE" summary="删除用户">
//     <description><![CDATA[
// 这是关于接口的详细说明文档
// 可以是一个 HTML 内容
// ]]></description>
//     <path path="/users" />
// </api>
fn delete() {}

// <api method="GET" summary="获取用户详情">
//     <description><![CDATA[
// 这是关于接口的详细说明文档
// 可以是一个 HTML 内容
// ]]></description>
//     <path path="/users/{id}">
//         <param name="id" type="number" summary="用户 ID" />
//         <query name="state" type="string" array="true" default="[normal,lock]" summary="state">
//              <enum value="normal">正常</enum>
//              <enum value="lock">锁定</enum>
//         </query>
//     </path>
//
//     <response status="200" array="true" type="object" mimetype="json">
//         <param name="id" type="number" summary="用户 ID" />
//         <param name="name" type="string" summary="用户名" />
//         <param name="groups" type="string" optional="true" summary="用户所在的权限组">
//             <param name="id" type="string" summary="权限组 ID" />
//             <param name="name" type="string" summary="权限组名称" />
//         </param>
//     </response>
// </api>
fn get() {}
