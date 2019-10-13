// SPDX-License-Identifier: MIT

// <api method="GET" summary="获取用户日志">
//     <description><![CDATA[
// 这是关于接口的详细说明文档
// ===
// 可以是一个 *markdown* 内容]]>
//     </description>
//     <path path="/users/{id}/logs/{lid}">
//         <param name="id" type="number">用户 ID</param>
//         <param name="lid" type="number">日志ID</param>
//         <query name="page" type="number" default="0" summary="页码" />
//         <query name="size" type="number" default="20">数量</query>
//         <query name="state" type="string" array="true" default="[normal,lock]" summary="state">
//              <enum value="normal">正常</enum>
//              <enum value="lock">锁定</enum>
//         </query>
//     </path>
//
//     <request type="object" mimetype="json">
//         <header name="name">desc</header>
//         <header name="name1">desc1</header>
//         <param name="count" type="number" optional="true" summary="summary" />
//         <param name="list" type="object" array="true" summary="list">
//             <param name="id" type="number" optional="true" summary="desc" />
//             <param name="name" type="string" optional="true" summary="desc" />
//             <param name="groups" type="string" array="true" optional="true" summary="desc">
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
//         ]]></example>
//     </request>
//
//     <response status="200" array="true" type="object" mimetype="json">
//         <param name="id" type="number" summary="desc" />
//         <param name="name" type="string" summary="desc" />
//         <param name="group" type="object" summary="object">
//             <param name="id" type="number" summary="id" />
//         </param>
//     </response>
//
//     <response status="404" type="object" mimetype="json">
//         <header name="authorization">token</header>
//         <param name="code" type="number" summary="状态码" optional="false" />
//         <param name="message" type="string" summary="错误信息" optional="false" />
//         <param name="detail" type="object" array="true" summary="错误明细">
//             <param name="id" type="string" summary="id" />
//             <param name="message" type="string" summary="message" />
//         </param>
//     </response>
//
//     <callback schema="https" summary="回调函数" method="POST">
//         <request mimetype="json" type="object">
//             <query name="size">size</query>
//
//             <param name="id" type="number" default="1" />
//             <param name="age" type="number" />
//             <example mimetype="json">
//             <![CDATA[
//             {
//                 id:1,
//                 sex: male,
//             }
//             ]]>
//             </example>
//         </request>
//
//         <response status="200" mimetype="text" type="string" />
//     </callback>
// </api>
void logs() {}
