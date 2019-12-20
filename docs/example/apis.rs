// SPDX-License-Identifier: MIT

// <api method="GET" summary="获取用户" deprecated="1.1.11">
//     <description>
// <![CDATA[
// <p>这是关于接口的详细说明文档</p><br />
// 可以是一个 HTML 内容
// ]]>
// </description>
//     <path path="/users">
//         <query name="page" type="number" default="0" summary="页码" />
//         <query name="size" type="number" default="20"><description><![CDATA[数量]]></description></query>
//     </path>
//
//     <tag>t1</tag>
//     <tag>t2</tag>
//     <server>admin</server>
//     <header name="name" type="string"><description><![CDATA[desc]]></description></header>
//     <header name="name1" type="string" summary="name1 desc" />
//
//     <response status="200" array="true" type="object" mimetype="application/json">
//         <param name="count" type="number" optional="false" summary="summary" />
//         <param name="list" type="object" array="true" summary="list">
//             <description type="html"><![CDATA[<span style="color:red">list description</span>]]></description>
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
//     <description>
// <![CDATA[
// 这是关于接口的详细说明文档<br />
// 可以是一个 HTML 内容
// ]]>
// </description>
//     <path path="/users" />
//
//     <tag>t2</tag>
//     <server>admin</server>
//     <server>old-client</server>
//
//     <request type="object" mimetype="application/json">
//         <header name="content-type" summary="application/json" type="string" />
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
//     <request type="object" mimetype="application/xml" name="users">
//         <param name="count" type="number" optional="false" summary="summary" />
//         <param name="list" type="object" array="true" summary="list">
//             <param name="id" type="number" summary="用户 ID" />
//             <param name="name" type="string" summary="用户名" />
//             <param name="groups" type="string" array="true" optional="true" summary="用户所在的权限组">
//                 <param name="id" type="string" summary="权限组 ID" />
//                 <param name="name" type="string" summary="权限组名称" />
//             </param>
//         </param>
//         <example mimetype="application/xml">
//         <![CDATA[
//             <users count="20">
//                 <user id="20" name="xx"></user>
//                 <user id="21" name="xx"></user>
//             </users>
//         ]]>
//         </example>
//     </request>
//
//     <response status="200" array="true" mimetype="application/json">
//     </response>
// </api>
fn post() {}

// <api method="DELETE" summary="删除用户">
//     <server>admin</server>
//     <description>
// <![CDATA[
// 这是关于接口的详细说明文档<br />
// 可以是一个 HTML 内容<br />
// ]]>
//     </description>
//     <path path="/users/{id}">
//         <param name="id" type="number" summary="用户 ID" />
//     </path>
// </api>
fn delete() {}

// <api method="GET" summary="获取用户详情">
//     <server>old-client</server>
//     <description>
// <![CDATA[
// 这是关于接口的详细说明文档
// 可以是一个 HTML 内容
// ]]>
//     </description>
//     <path path="/users/{id}">
//         <param name="id" type="number" summary="用户 ID" />
//         <query name="state" type="string" array="true" default="[normal,lock]" summary="state">
//              <enum value="normal" summary="正常" />
//              <enum value="lock"><description type="html"><![CDATA[<span style="color:red">锁定</span>]]></description></enum>
//         </query>
//     </path>
//
//     <response status="200" array="true" type="object" mimetype="application/json">
//         <param name="id" type="number" summary="用户 ID" />
//         <param name="name" type="string" summary="用户名" />
//         <param name="groups" type="string" optional="true" summary="用户所在的权限组">
//             <param name="id" type="string" summary="权限组 ID" />
//             <param name="name" type="string" summary="权限组名称" />
//         </param>
//     </response>
// </api>
fn get() {}
