// SPDX-License-Identifier: MIT

// <apidoc version="1.1.1">
//     <title>示例文档</title>
//     <server name="admin" url="https://api.example.com/admin">
//     <description type="html"><![CDATA[
//     后台管理接口，<br /><br /><br /><br /><p style="color:red">admin</p>
//     ]]></description>
//     </server>
//     <server name="old-client" url="https://api.example.com/client" deprecated="1.1.1" summary="客户端接口" />
//     <server name="client" url="https://api.example.com" summary="客户端接口" />
//     <mimetype>application/xml</mimetype>
//     <mimetype>application/json</mimetype>
//     <tag name="t1" title="标签1" />
//     <tag name="t2" title="标签2" />
//     <license url="https://opensource.org/licenses/MIT" text="MIT" />
//     <contact name="name">
//         <url>https://example.com</url>
//         <email>example@example.com</email>
//     </contact>
//
//     <response status="400" type="object" name="result">
//         <param name="code" type="number" xml-attr="true" summary="状态码" optional="false" />
//         <param name="message" type="string" summary="错误信息" optional="false" />
//         <param name="detail" type="object" array="true" summary="错误明细">
//             <param name="id" type="string" summary="id" />
//             <param name="message" type="string" summary="message" />
//         </param>
//     </response>
//     <response status="404" summary="not found" />
//
//     <description type="html">
// <![CDATA[
//     <p>这是一个用于测试的文档用例</p>
//     状态码：
//     <ul>
//         <li>40300:xxx</li>
//         <li>40301:xxx</li>
//     </ul>
// ]]>
//     </description>
// </apidoc>
