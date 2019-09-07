// SPDX-License-Identifier: MIT

// <api method="GET" summary="获取用户日志">
//     <description><![CDATA[
// 这是关于接口的详细说明文档
// ===
// 可以是一个 *markdown* 内容]]>
//     </description>
//     <path path="/users/{id}/logs/{lid}">
//         <param name="id" type="int">用户 ID</param>
//         <param name="lid" type="int">日志ID</param>
//         <query name="page" type="int" default="0" summary="页码" />
//         <query name="size" type="int" default="20">数量</query>
//         <query name="state" type="string" array="true" default="[normal,lock]" summary="state">
//              <enum value="normal">正常</enum>
//              <enum value="lock">锁定</enum>
//         </query>
//     </path>
//
//     <request type="object">
//         <header name="name">desc</header>
//         <header name="name1">desc1</header>
//         <param name="count" type="int" required="true" summary="summary" />
//         <param name="list" type="object" array="true" summary="list">
//             <param name="id" type="int" required="true" summary="desc" />
//             <param name="name" type="string" required="true" summary="desc" />
//             <param name="groups" type="string" array="true" required="true" summary="desc">
//                 <enum value="xx1">xx</enum>
//                 <enum value="xx2">xx</enum>
//             </param>
//         </param>
//         <example mimetype="json"><![CDATA[
// {
//  count: 5,
//  list: [
//    {id:1, name: 'name1', 'groups': [1,2]},
//    {id:2, name: 'name2', 'groups': [1,2]}
//  ]
// }
//         ]]</example>
//     </request>
//
//     <response status="200" type="array.object">
//         <param name="id"
//     </response>
// </api>


//
// @apiresponse 200 array.object * 通用的返回内容定义
// @apiheader string required desc
// @apiparam id int reqiured desc
// @apiparam name string reqiured desc
// @apiparam group object reqiured desc
// @apiparam group.id int reqiured desc
//
// @apiresponse 404 object application/json 错误的返回内容
// @apiheader string required desc38G
// @apiparam code int reqiured desc
// @apiparam message string reqiured desc
// @apiparam detail array.object reqiured desc
// @apiparam detail.id string reqiured desc
// @apiparam detail.message string reqiured desc
//
// @apiCallback GET 回调内容
// @apirequest object application/xml 特定的请求主体
// @apiresponse 404 object application/json 错误的返回内容
void logs() {}
