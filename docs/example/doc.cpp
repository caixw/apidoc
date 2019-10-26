// SPDX-License-Identifier: MIT

// <apidoc version="1.1.1">
//     <title>示例文档</title>
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
//     <response status="400" type="object" mimetype="application/json">
//         <param name="code" type="number" summary="状态码" optional="false" />
//         <param name="message" type="string" summary="错误信息" optional="false" />
//         <param name="detail" type="object" array="true" summary="错误明细">
//             <param name="id" type="string" summary="id" />
//             <param name="message" type="string" summary="message" />
//         </param>
//     </response>
//     <response status="400" type="object" mimetype="application/xml" name="result">
//         <param name="code" type="number" summary="状态码" optional="false" />
//         <param name="message" type="string" summary="错误信息" optional="false" />
//         <param name="detail" type="object" array="true" summary="错误明细列表">
//             <param name="msg" type="object" summary="message">
//                 <param name="@id" type="string" summary="id" />
//             </param>
//         </param>
//         <example mimetype="application/xml">
//         <![CDATA[
//         <result code="40001">
//             <message>错误信息内容</message>
//             <detail>
//                 <msg id="name">不能为空</msg>
//                 <msg id="id">不能为空</msg>
//             </detail>
//         </result>
//         ]]>
//         </example>
//     </response>
//     <response status="404" type="none" mimetype="application/json application/xml" summary="not found" />
//
//     <description>
//     <![CDATA[
//     <p>这是一个用于测试的文档用例</p>
//     状态码：
//     <ul>
//         <li>40300:xxx</li>
//         <li>40301:xxx</li>
//     </ul>
//     ]]></description>
// </apidoc>
