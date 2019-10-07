// SPDX-License-Identifier: MIT

// <apidoc version="1.1.1">
//     <title>title of doc</title>
//     <server name="admin" url="https://api.example.com/admin">后台管理接口</server>
//     <server name="client" url="https://api.example.com">客户端接口</server>
//     <tag name="t1" title="标签1" />
//     <tag name="t2" title="标签2" />
//     <license url="https://opensource.org/licenses/MIT">MIT</license>
//     <contact name="name">
//         <url>https://example.com</url>
//         <email>example@example.com</email>
//     </contact>
//
//     <response status="404" type="object" mimetype="json">
//         <header name="authorization">token</header>
//         <param name="code" type="number" summary="状态码" required="true" />
//         <param name="message" type="string" summary="错误信息" required="true" />
//         <param name="detail" type="object" array="true" summary="错误明细">
//             <param name="id" type="string" summary="id" />
//             <param name="message" type="string" summary="message" />
//         </param>
//     </response>
//
//     <content>
//     <![CDATA[
//      这里可以是 markdown 文档
//     ]]>
//     </content>
// </apidoc>
