// SPDX-License-Identifier: MIT

// <api method="GET" summary="获取用户日志">
//     <server>client</server>
//     <description type="html">
// <![CDATA[
// <p>这是关于接口的详细说明文档</p>
// <p style="color:red">可以是一个 HTML 内容</p>
// ]]>
//     </description>
//     <path path="/users/{id}/logs">
//         <param name="id" type="number"><description><![CDATA[用户 ID]]></description></param>
//         <query name="page" type="number" default="0" summary="页码" />
//         <query name="size" type="number" default="20"><description><![CDATA[数量]]></description></query>
//     </path>
//
//     <response status="200" array="true" type="object" mimetype="application/json">
//         <header name="name" type="string" summary="desc" />
//         <header name="name1" type="string" summary="desc1" />
//         <param name="count" type="number" optional="true" summary="summary" />
//         <param name="list" type="object" array="true" summary="list">
//             <param name="id" type="number" optional="true" summary="desc" />
//             <param name="name" type="string" optional="true" summary="desc" />
//             <param name="groups" type="string" array="true" optional="true" summary="desc">
//                 <enum value="xx1"><description>xx</description></enum>
//                 <enum value="xx2" summary="xx" />
//             </param>
//         </param>
//         <example mimetype="application/json"><![CDATA[
// {
//  count: 5,
//  list: [
//    {id:1, name: 'name1', 'groups': [1,2]},
//    {id:2, name: 'name2', 'groups': [1,2]}
//  ]
// }
//         ]]></example>
//     </response>
//
//     <callback schema="https" summary="回调函数" method="POST">
//     <description type="html">
// <![CDATA[
//         <p style="color:red">这是一个回调函数的详细说明</p>
//         <p>为一个 html 文档</p>
// ]]>
//     </description>
//         <request mimetype="application/json" type="object">
//             <query name="size">size</query>
//
//             <param name="id" type="number" default="1" summary="id" />
//             <param name="age" type="number" summary="age" />
//             <example mimetype="application/json">
//             <![CDATA[
//             {
//                 id:1,
//                 sex: male,
//             }
//             ]]>
//             </example>
//         </request>
//
//         <response status="200" mimetype="text/plain" type="string" />
//     </callback>
// </api>
void logs() {}
