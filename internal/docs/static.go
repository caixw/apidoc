// 当前文件由工具自动生成，请勿手动修改！

package docs

var data = []*FileInfo{{
	Name:        "config.xml",
	ContentType: "application/xml; charset=utf-8",
	Content: []byte(`<?xml version="1.0" encoding="UTF-8"?>

<!-- 该文件由工具自动生成，请勿手动修改！-->

<config>
	<name>apidoc</name>
	<version>v6</version>
	<repo>https://github.com/caixw/apidoc</repo>
	<url>https://apidoc.tools</url>
	<languages>
		<language>C#</language>
		<language>C/C++</language>
		<language>D</language>
		<language>Erlang</language>
		<language>Go</language>
		<language>Groovy</language>
		<language>Java</language>
		<language>JavaScript</language>
		<language>Kotlin</language>
		<language>Pascal/Delphi</language>
		<language>Perl</language>
		<language>PHP</language>
		<language>Python</language>
		<language>Ruby</language>
		<language>Rust</language>
		<language>Scala</language>
		<language>Swift</language>
	</languages>
</config>
`),
},
	{
		Name:        "example/index.xml",
		ContentType: "application/xml; charset=utf-8",
		Content: []byte(`<?xml version="1.0" encoding="UTF-8"?>

<?xml-stylesheet type="text/xsl" href="../v6/apidoc.xsl"?>
<apidoc apidoc="6.0.0" created="2020-03-05T00:27:52+08:00" version="1.1.1">
	<title>示例文档</title>
	<description type="html"><![CDATA[
 
     <p>这是一个用于测试的文档用例</p>
     状态码：
     <ul>
         <li>40300:xxx</li>
         <li>40301:xxx</li>
     </ul>
 
     ]]></description>
	<contact name="name">
		<url>https://example.com</url>
		<email>example@example.com</email>
	</contact>
	<license text="MIT" url="https://opensource.org/licenses/MIT"></license>
	<tag name="t1" title="标签1"></tag>
	<tag name="t2" title="标签2"></tag>
	<server name="admin" url="https://api.example.com/admin">
		<description type="html"><![CDATA[
     后台管理接口，
     ]]></description>
	</server>
	<server name="old-client" url="https://api.example.com/client" deprecated="1.1.1" summary="客户端接口"></server>
	<server name="client" url="https://api.example.com" summary="客户端接口"></server>
	<api method="GET" summary="获取用户" deprecated="1.1.11">
		<path path="/users">
			<query name="page" type="number" default="0" summary="页码"></query>
			<query name="size" type="number" default="20">
				<description type="markdown"><![CDATA[数量]]></description>
			</query>
		</path>
		<description type="markdown"><![CDATA[
 
 <p>这是关于接口的详细说明文档</p><br />
 可以是一个 HTML 内容
 
 ]]></description>
		<response name="user" type="object" array="true" status="200">
			<param xml-attr="true" name="count" type="number" summary="summary"></param>
			<param name="list" type="object" array="true" summary="list">
				<param xml-attr="true" name="id" type="number" summary="用户 ID"></param>
				<param xml-attr="true" name="name" type="string" summary="用户名"></param>
				<param name="groups" type="string" optional="true" array="true" summary="用户所在的权限组">
					<param name="id" type="string" summary="权限组 ID"></param>
					<param name="name" type="string" summary="权限组名称"></param>
				</param>
				<description type="html"><![CDATA[<span style="color:red">list description</span>]]></description>
			</param>
		</response>
		<header name="name" type="string">
			<description type="markdown"><![CDATA[desc]]></description>
		</header>
		<header name="name1" type="string" summary="name1 desc"></header>
		<tag>t1</tag>
		<tag>t2</tag>
		<server>admin</server>
	</api>
	<api method="POST" summary="添加用户">
		<path path="/users"></path>
		<description type="markdown"><![CDATA[
 
 这是关于接口的详细说明文档<br />
 可以是一个 HTML 内容
 
 ]]></description>
		<request type="object">
			<param name="count" type="number" summary="summary"></param>
			<param name="list" type="object" array="true" summary="list">
				<param name="id" type="number" summary="用户 ID"></param>
				<param name="name" type="string" summary="用户名"></param>
				<param name="groups" type="string" optional="true" array="true" summary="用户所在的权限组">
					<param name="id" type="string" summary="权限组 ID"></param>
					<param name="name" type="string" summary="权限组名称"></param>
				</param>
			</param>
			<header name="content-type" type="string" summary="application/json"></header>
		</request>
		<request name="users" type="object" mimetype="application/xml">
			<param name="count" type="number" summary="summary"></param>
			<param name="list" type="object" array="true" summary="list">
				<param name="id" type="number" summary="用户 ID"></param>
				<param name="name" type="string" summary="用户名"></param>
				<param name="groups" type="string" optional="true" array="true" summary="用户所在的权限组">
					<param name="id" type="string" summary="权限组 ID"></param>
					<param name="name" type="string" summary="权限组名称"></param>
				</param>
			</param>
			<example mimetype="application/xml"><![CDATA[
         
             <users count="20">
                 <user id="20" name="xx"></user>
                 <user id="21" name="xx"></user>
             </users>
         
         ]]></example>
		</request>
		<response array="true" status="200" mimetype="application/json"></response>
		<header name="name" type="string">
			<description type="markdown"><![CDATA[desc]]></description>
		</header>
		<header name="name1" type="string" summary="name1 desc"></header>
		<tag>t2</tag>
		<server>admin</server>
		<server>old-client</server>
	</api>
	<api method="DELETE" summary="删除用户">
		<path path="/users/{id}">
			<param name="id" type="number" summary="用户 ID"></param>
		</path>
		<description type="markdown"><![CDATA[
 
 这是关于接口的详细说明文档<br />
 可以是一个 HTML 内容<br />
 
     ]]></description>
		<server>admin</server>
	</api>
	<api method="GET" summary="获取用户详情">
		<path path="/users/{id}">
			<param name="id" type="number" summary="用户 ID"></param>
			<query name="state" type="string" default="[normal,lock]" array="true" summary="state">
				<enum value="normal" summary="正常"></enum>
				<enum value="lock">
					<description type="html"><![CDATA[<span style="color:red">锁定</span>]]></description>
				</enum>
			</query>
		</path>
		<description type="markdown"><![CDATA[
 
 这是关于接口的详细说明文档
 可以是一个 HTML 内容
 
     ]]></description>
		<response type="object" array="true" status="200" mimetype="application/json">
			<param name="id" type="number" summary="用户 ID"></param>
			<param name="name" type="string" summary="用户名"></param>
			<param name="groups" type="string" optional="true" summary="用户所在的权限组">
				<param name="id" type="string" summary="权限组 ID"></param>
				<param name="name" type="string" summary="权限组名称"></param>
			</param>
		</response>
		<server>old-client</server>
	</api>
	<api method="GET" summary="获取用户日志">
		<path path="/users/{id}/logs">
			<param name="id" type="number">
				<description type="markdown"><![CDATA[用户 ID]]></description>
			</param>
			<query name="page" type="number" default="0" summary="页码"></query>
			<query name="size" type="number" default="20">
				<description type="markdown"><![CDATA[数量]]></description>
			</query>
		</path>
		<description type="html"><![CDATA[
 
 <p>这是关于接口的详细说明文档</p>
 <p style="color:red">可以是一个 HTML 内容</p>
 
     ]]></description>
		<response type="object" array="true" status="200" mimetype="application/json">
			<param name="count" type="number" optional="true" summary="summary"></param>
			<param name="list" type="object" array="true" summary="list">
				<param name="id" type="number" optional="true" summary="desc"></param>
				<param name="name" type="string" optional="true" summary="desc"></param>
				<param name="groups" type="string" optional="true" array="true" summary="desc">
					<enum value="xx1">
						<description type="markdown"><![CDATA[xx]]></description>
					</enum>
					<enum value="xx2" summary="xx"></enum>
				</param>
			</param>
			<example mimetype="application/json"><![CDATA[
 {
  count: 5,
  list: [
    {id:1, name: 'name1', 'groups': [1,2]},
    {id:2, name: 'name2', 'groups': [1,2]}
  ]
 }
         ]]></example>
			<header name="name" type="string" summary="desc"></header>
			<header name="name1" type="string" summary="desc1"></header>
		</response>
		<callback method="POST" summary="回调函数">
			<description type="html"><![CDATA[
 
         <p style="color:red">这是一个回调函数的详细说明</p>
         <p>为一个 html 文档</p>
 
     ]]></description>
			<response type="string" status="200" mimetype="text/plain"></response>
			<request type="object" mimetype="application/json">
				<param name="id" type="number" default="1" summary="id"></param>
				<param name="age" type="number" summary="age"></param>
				<example mimetype="application/json"><![CDATA[
             
             {
                 id:1,
                 sex: male,
             }
             
             ]]></example>
			</request>
		</callback>
		<server>client</server>
	</api>
	<response name="result" type="object" status="400">
		<param xml-attr="true" name="code" type="number" summary="状态码"></param>
		<param name="message" type="string" summary="错误信息"></param>
		<param name="detail" type="object" array="true" summary="错误明细">
			<param name="id" type="string" summary="id"></param>
			<param name="message" type="string" summary="message"></param>
		</param>
	</response>
	<response summary="not found" status="404"></response>
	<mimetype>application/xml</mimetype>
	<mimetype>application/json</mimetype>
</apidoc>`),
	},
	{
		Name:        "icon.svg",
		ContentType: "image/svg+xml; charset=utf-8",
		Content: []byte(`<?xml version="1.0" encoding="utf-8" ?>

<svg viewBox="0 0 512 512" xmlns="http://www.w3.org/2000/svg" stroke="#f60">
    <title>apidoc</title>
    <desc>https://apidoc.tools LOGO</desc>
    <desc>这是一只诞生于羊年的小工具</desc>

    <circle cx="256" cy="256" r="248" fill-opacity="0" stroke-width="16" />
    <circle cx="530" cy="385" r="300" fill-opacity="0" stroke-width="16" />
    <circle cx="-15" cy="385" r="300" fill-opacity="0" stroke-width="16" />
</svg>
`),
	},
	{
		Name:        "index.css",
		ContentType: "text/css; charset=utf-8",
		Content: []byte(`@charset "utf-8";

:root {
    --max-width: 1024px;
    --padding: 1rem;
    --article-padding: calc(var(--padding) / 2);

    --color: black;
    --accent-color: #0074d9;
    --background: white;
    --border-color: #e0e0e0;
}

@media (prefers-color-scheme: dark) {
    :root {
        --color: #b0b0b0;
        --accent-color: #0074d9;
        --background: black;
        --border-color: #303030;
    }
}

body {
    padding: 0;
    margin: 0 auto;
    color: var(--color);
    background: var(--background);
}

table {
    width: 100%;
    border: 1px solid var(--border-color);
    border-radius: 5px;
    padding: 0 var(--article-padding);
}

table thead tr {
    height: 2.5rem;
}

table th {
    white-space: nowrap;
    padding-right: .5rem;
}

table th, table td {
    text-align: left;
    border-bottom: 1px solid var(--border-color);
}

table tbody tr:last-of-type td,
table tbody tr:last-of-type th {
    border-bottom: none;
}

ul, ol {
    padding: 0;
    margin: 0;
    list-style-position: inside;
}

a {
    text-decoration: none;
    color: var(--accent-color);
}

a:hover {
    opacity: .7;
}

p {
    margin: 15px 0;
}

pre {
    border: 1px solid var(--border-color);
    padding: var(--article-padding);
    border-radius: 5px;
    white-space: pre-wrap;
}

/*************************** header ***********************/

header {
    top: 0;
    left: 0;
    position: sticky;
    z-index: 1000;
    box-sizing: border-box;
    background: var(--background);
    box-shadow: 2px 2px 2px var(--border-color);
}

header .wrap {
    margin: 0 auto;
    text-align: left;
    max-width: var(--max-width);
    padding: 0 var(--padding);
    display: flex;
    flex-wrap: wrap;
    justify-content: space-between;
    align-items: center;
}

header h1 {
    display: inline-block;
    margin: var(--article-padding) 0;
    text-transform: uppercase;
}

header h1 .version {
    font-size: 1rem;
}

header h1 img {
    height: 1.5rem;
    margin-right: .5rem;
}

header .menus {
    display: flex;
    align-items: baseline;
}

header .menu {
    margin-left: var(--padding);
    color: var(--color);
}

header .menu:hover {
    opacity: .7;
}

header .drop-menus {
    position: relative;
    cursor: pointer;
}

header .drop-menus ul {
    position: absolute;
    right: 0;
    display: none;
    list-style: none;
    background: var(--background);
    border: 1px solid var(--border-color);
    padding: var(--article-padding);
}

header .drop-menus a {
    color: var(--color);
}

header .drop-menus:hover ul {
    display: block;
}

/************************** main **************************/

main {
    margin: 0 auto;
    max-width: var(--max-width);
    padding: var(--padding);
}

main article .link {
    font-size: .8rem;
    display: none;
}

main article h2:hover .link,
main article h3:hover .link {
    display: inline-block;
}

/************************* footer **************************/

footer {
    border-top: 1px solid var(--border-color);
    padding: 0 var(--padding) var(--padding);
}

footer .wrap {
    max-width: var(--max-width);
    margin: 0 auto;
}

.goto-top {
    border: solid var(--color);
    border-width: 0 5px 5px 0;
    display: block;
    padding: 5px;
    transform: rotate(-135deg);
    position: fixed;
    bottom: 3rem;
    right: 3rem;
}
`),
	},
	{
		Name:        "index.js",
		ContentType: "application/javascript; charset=utf-8",
		Content: []byte(`'use strict';

window.onload = function() {
    initGotoTop();
};

function initGotoTop() {
    const top = document.querySelector('.goto-top');

    // 在最顶部时，隐藏按钮
    window.addEventListener('scroll', (e) => {
        const body = document.querySelector('html');
        if (body.scrollTop > 50) {
            top.style.display = 'block';
        } else {
            top.style.display = 'none';
        }
    });
}
`),
	},
	{
		Name:        "index.xml",
		ContentType: "application/xml; charset=utf-8",
		Content: []byte(`<?xml version="1.0" encoding="utf-8"?>
<?xml-stylesheet type="text/xsl" href="./index.xsl"?>

<!--
这是官网首页内容，同时也是简体中文的本地化内容。

其它语言的本化地内容，需要重新改写本文件中除注释外的所有内容。
-->

<docs lang="zh-hans">
    <title>apidoc | RESTful API 文档处理工具</title>
    <license url="https://creativecommons.org/licenses/by/4.0/deed.zh">署名 4.0 国际 (CC BY 4.0)</license>

    <!-- 类型描述中表格的相关本化地信息 -->
    <type-locale>
        <header>
            <name>名称</name>
            <type>类型</type>
            <required>必填</required>
            <description>描述</description>
        </header>
    </type-locale>

    <doc id="about" title="关于">
        <p>apidoc 是一个简单的 <abbr title="Representational State Transfer">RESTful</abbr> <abbr title="Application Programming Interface">API</abbr> 文档生成工具，它从代码注释中提取特定格式的内容生成文档。支持诸如 Go、Java、C++、Rust 等大部分开发语言，具体可使用 <code>apidoc lang</code> 命令行查看所有的支持列表。</p>

        <p>apidoc 拥有以下特点：</p>
        <ol>
            <li>跨平台，linux、windows、macOS 等都支持；</li>
            <li>支持语言广泛，即使是不支持，也很方便扩展；</li>
            <li>支持多个不同语言的多个项目生成一份文档；</li>
            <li>输出模板可自定义；</li>
            <li>根据文档生成 mock 数据；</li>
        </ol>

        <p>以下是一段简短的 C 语言风格代码下的示例：</p>
        <pre><code class="language-markup"><![CDATA[/**
 * <api method="GET" summary="获取所有的用户信息">
 *     <path path="/users">
 *         <query name="page" type="number" default="0">显示第几页的内容</query>
 *         <query name="size" type="number" default="20">每页显示的数量</query>
 *     </path>
 *     <tag>user</tag>
 *     <server>users</server>
 *     <response status="200" type="object" mimetype="application/json">
 *         <param name="count" type="int" optional="false" summary="符合条件的所有用户数量" />
 *         <param name="users" type="object" array="true" summary="用户列表">
 *             <param name="id" type="int" summary="唯一 ID" />
 *             <param name="name" type="string" summary="姓名" />
 *         </param>
 *     </response>
 *     <response status="500" mimetype="application/json" type="obj">
 *         <param name="code" type="int" summary="错误代码" />
 *         <param name="msg" type="string" summary="错误内容" />
 *     </response>
 * </api>
 */]]></code></pre>
        <p>多行注释中，每一行中以<code>空白字符+symbol+空白字符</code>开头的，这些字符将会被过滤，symbol 表示该注释块的起始字符中的任意字符。比如以上代码中，所有的 <var>*</var> 将被过滤。</p>
    </doc>

    <doc id="usage" title="使用" />

    <doc id="spec" title="文档格式">
        <p>文档采用 XML 格式。存在两个顶级标签：<code>apidoc</code> 和 <code>api</code>，用于描述整体内容和具体接口信息。</p>

        <p>文档被从注释中提取之后，最终会被合并成一个 XML 文件，在该文件中 <code>api</code> 作为 <code>apidoc</code> 的一个子元素存在，如果你的项目不想把文档写在注释中，也可以直接编写一个完整的 XML 文件，将 <code>api</code> 作为 <code>apidoc</code> 的一个子元素。</p>

        <p>具体可参考<a href="./example/index.xml">示例代码</a>。</p>

        <p>以下是对各个 XML 元素以及参数介绍，其中以 <code>@</code> 开头的表示 XML 属性；<code>.</code> 表示为当前元素的内容；其它表示子元素。</p>
    </doc>

    <doc id="install" title="安装" parent="usage">
        <p>可以直接从 <a href="https://github.com/caixw/apidoc/releases">https://github.com/caixw/apidoc/releases</a> 查找你需要的版本下载，放入 <code>PATH</code> 中即可使用。如果没有你需要的平台文件，则需要从源代码编译：</p>
        <ul>
            <li>下载 Go 编译工具</li>
            <li>下载源代码</li>
            <li>执行代码中 <code>build/build.sh</code> 或是 <code>build/build.cmd</code> 进行编译</li>
            <li>编译好的文件存放在 cmd/apidoc 下，可以将该文件放置在 PATH 目录</li>
        </ul>
    </doc>

    <doc id="env" title="环境变量" parent="usage">
        <p>apidoc 会读取 <var>LANG</var> 的环境变量作为其本地化的依据，若想指定其它语种，可以手动指定 <var>LANG</var> 环境变量：<samp>LANG=zh-Hant apidoc</samp>。在 windows 系统中，若不存在 <var>LANG</var> 环境变量，则会调用 <samp>GetUserDefaultLocaleName</samp> 函数来获取相应的语言信息。</p>
    </doc>

    <doc id="cli" title="命令行" parent="usage">
        <p>可以通过 <samp>apidoc help</samp> 查看命令行支持的子命令。主要包含了以下几个：</p>
        <table>
            <thead><tr><th>子命令</th><th>描述</th></tr></thead>
            <tbody>
                <tr><td>help</td><td>显示子命令的描述信息</td></tr>
                <tr><td>build</td><td>生成文档内容</td></tr>
                <tr><td>mock</td><td>根据文档提供 mock 服务</td></tr>
                <tr><td>static</td><td>提供查看文档的本地服务</td></tr>
                <tr><td>version</td><td>显示版本信息</td></tr>
                <tr><td>lang</td><td>列出当前支持的语言</td></tr>
                <tr><td>locale</td><td>列出当前支持的本地化内容</td></tr>
                <tr><td>detect</td><td>根据指定的目录生成配置文件</td></tr>
                <tr><td>test</td><td>检测语法是否准确</td></tr>
            </tbody>
        </table>
        <p>mock 子命令可以根据文档生成一些符合要求的随机数据。这些数据每次请求都不相同，包括数量、长度、数值大小等。</p>
    </doc>

    <!-- 配置文件的类型定义 -->
    <types parent="usage">
        <type name=".apidoc.yaml">
            <description>
                <p>配置文件名固定为 <code>.apidoc.yaml</code>，格式为 YAML，可参考 <a href="example/.apidoc.yaml">.apidoc.yaml</a>。文件可以通过命令 <code>apidoc detect</code> 生成。主要包含了以几个配置项：</p>
            </description>
            <item name="version" >产生此配置文件的 apidoc 版本</item>
            <item name="inputs">指定输入的数据，同一项目只能解析一种语言。</item>
            <item name="inputs.dir">需要解析的源文件所在目录</item>
            <item name="inputs.recursive">是否解析子目录下的源文件</item>
            <item name="inputs.encoding">编码，默认为 <code>utf-8</code>，值可以是 <a href="https://www.iana.org/assignments/character-sets/character-sets.xhtml">character-sets</a> 中的内容。</item>
            <item name="inputs.lang">源文件类型。具体支持的类型可通过 -l 参数进行查找</item>
            <item name="output">控制输出行为</item>
            <item name="output.path">指定输出的文件名，包含路径信息。</item>
            <item name="output.tags">只输出与这些标签相关联的文档，默认为全部。</item>
            <item name="output.style">为 XML 文件指定的 XSL 文件。</item>
        </type>
    </types>

    <types parent="spec">
        <type name="apidoc">
            <description><p>用于描述整个文档的相关内容，只能出现一次。</p></description>
            <item name="@version">文档的版本</item>
            <item name="@lang">内容的本地化 ID，比如 <samp><var>zh-hans</var></samp> 等。</item>
            <item name="@logo">图标，默认采用官网的 <var>https://apidoc.tools/icon.svg</var>，同时作用于 favicon 和 logo，只支持 SVG 格式。</item>
            <item name="@created">文档的生成时间</item>
            <item name="title">文档的标题</item>
            <item name="description">文档的整体介绍，可以是使用 HTML 内容。</item>
            <item name="contract">联系人信息</item>
            <item name="license">内容的版权信息</item>
            <item name="tag">可以用的标签列表</item>
            <item name="server">API 基地址列表，每个 API 最少应该有一个 server。</item>
            <item name="mimetype">接口所支持的 mimetype 类型</item>
            <item name="response">表示所有 API 都有可能返回的內容</item>
            <item name="api">API 文档内容</item>
        </type>

        <type name="link">
            <description><p>用于描述链接，一般转换为 HTML 的 a 标签。</p></description>
            <item name="@url">链接指向的 URL</item>
            <item name="@text">链接的文本内容</item>
        </type>

        <type name="contact">
            <description><p>用于描述联系方式</p></description>
            <item name="@url">链接的 URL，与邮箱必须二选一必填</item>
            <item name="@email">邮件地址，与 url 必须二选一必填</item>
            <item name=".">联系人名称</item>
        </type>

        <type name="tag">
            <description><p>定义标签，标签相当于关键字，作用于 API，相当于启到分类的作用。</p></description>
            <item name="@name">标签的唯一 ID，推荐采用英文字母表示。</item>
            <item name="@title">标签名称</item>
            <item name="@deprecated">表示该标签在大于等于该版本号时不再启作用</item>
        </type>

        <type name="server">
            <description><p>定义服务器的相关信息，作用于 API，决定该 API 与哪个服务器相关联。</p></description>
            <item name="@name">唯一 ID，推荐采用英文字母表示。</item>
            <item name="@url">服务基地址</item>
            <item name="@deprecated">表示在大于等于该版本号时不再启作用</item>
            <item name="@summary">简要的描述内容，或者通过 <code>description</code> 提供一份富文本内容。</item>
            <item name="description">对该服务的具体描述，可以使用 HTML 内容</item>
        </type>

        <type name="api">
            <description><p>定义接口的具体内容</p></description>
            <item name="@version">表示此接口在该版本中添加</item>
            <item name="@method">请求方法</item>
            <item name="@summary">简要介绍</item>
            <item name="@deprecated">表示在大于等于该版本号时不再启作用</item>
            <item name="@id">唯一 ID</item>
            <item name="description">该接口的详细介绍，为 HTML 内容。</item>
            <item name="path">定义路径信息</item>
            <item name="request">定义可用的请求信息</item>
            <item name="response">定义可能的返回信息</item>
            <item name="callback">定义回调接口内容</item>
            <item name="tag">关联的标签</item>
            <item name="server">关联的服务</item>
            <item name="header">传递的报头内容，如果是某个 mimetype 专用的，可以放在 request 元素中。</item>
        </type>

        <type name="path">
            <description><p>用于定义请求时与路径相关的内容</p></description>
            <item name="@path">接口地址</item>
            <item name="param">地址中的参数</item>
            <item name="query">地址中的查询参数</item>
        </type>

        <type name="request">
            <description><p>定义了请求和返回的相关内容</p></description>
            <item name="@xml-ns">XML 标签的命名空间</item>
            <item name="@xml-ns-prefix">XML 标签的命名空间名称前缀</item>
            <item name="@xml-wrapped">如果当前元素的 <code>@array</code> 为 <var>true</var>，是否将其包含在 wrapped 指定的标签中。</item>
            <item name="@name">当 mimetype 为 <var>application/xml</var> 时，此值表示 XML 的顶层元素名称，否则无用。</item>
            <item name="@type">值的类型，可以是 <del title="建议使用空值代替"><var>none</var></del>、<var>string</var>、<var>number</var>、<var>bool</var>、<var>object</var> 和 空值；空值表示不输出任何内容。</item>
            <item name="@deprecated">表示在大于等于该版本号时不再启作用</item>
            <item name="@summary">简要介绍</item>
            <item name="@array">是否为数组</item>
            <item name="@status">状态码。在 request 中，该值不可用，否则为必填项。</item>
            <item name="@mimetype">媒体类型，比如 <var>application/json</var> 等。</item>
            <item name="description">详细介绍，为 HTML 内容。</item>
            <item name="enum">当前参数可用的枚举值</item>
            <item name="param">子类型，比如对象的子元素。</item>
            <item name="example">示例代码。</item>
            <item name="header">传递的报头内容</item>
        </type>

        <type name="param">
            <description><p>参数类型，基本上可以作为 <code>request</code> 的子集使用。</p></description>
            <item name="@xml-attr">是否作为父元素的属性，仅作用于 XML 元素。</item>
            <item name="@xml-extract">将当前元素的内容作为父元素的内容，要求父元素必须为 <var>object</var>。</item>
            <item name="@xml-ns">XML 标签的命名空间</item>
            <item name="@xml-ns-prefix">XML 标签的命名空间名称前缀</item>
            <item name="@xml-attr">是否作为父元素的属性，仅用于 XML 的请求。</item>
            <item name="@xml-wrapped">如果当前元素的 <code>@array</code> 为 <var>true</var>，是否将其包含在 wrapped 指定的标签中。</item>
            <item name="@name">值的名称</item>
            <item name="@type">值的类型，可以是 <var>string</var>、<var>number</var>、<var>bool</var> 和 <var>object</var></item>
            <item name="@deprecated">表示在大于等于该版本号时不再启作用</item>
            <item name="@default">默认值</item>
            <item name="@optional">是否为可选的参数</item>
            <item name="@summary">简要介绍</item>
            <item name="@array">是否为数组</item>
            <item name="@array-style">是否以数组的形式展示数据，默认采用 form 形式，仅在 <code>@array</code> 为 <var>true</var> 时有效。</item>
            <item name="description">详细介绍，为 HTML 内容。</item>
            <item name="enum">当前参数可用的枚举值</item>
            <item name="param">子类型，比如对象的子元素。</item>
        </type>

        <type name="enum">
            <description><p>定义枚举类型的数所的枚举值</p></description>
            <item name="@value">枚举值</item>
            <item name="@deprecated">表示在大于等于该版本号时不再启作用</item>
            <item name=".">该值的详细介绍</item>
        </type>

        <type name="example">
            <description><p>示例代码</p></description>
            <item name="@mimetype">代码的 mimetype 类型。</item>
            <item name=".">示例代码的内容，需要使用 CDATA 包含代码。</item>
        </type>

        <type name="header">
            <description><p>定义了请求和返回的报头结构</p></description>
            <item name="@name">报头的名称</item>
            <item name="@deprecated">表示在大于等于该版本号时不再启作用</item>
            <item name="@summary">对报头的描述</item>
            <item name="description">对报头的描述</item>
        </type>

        <type name="callback">
            <description><p>定义接口回调的相关内容</p></description>
            <item name="@method">请求方法</item>
            <item name="@summary">简要介绍</item>
            <item name="@deprecated">表示在大于等于该版本号时不再启作用</item>
            <item name="description">该接口的详细介绍</item>
            <item name="path">定义路径信息</item>
            <item name="request">定义可用的请求信息</item>
            <item name="response">定义可能的返回信息</item>
        </type>

        <type name="richtext">
            <description><p>富文本信息，可以以不同的格式展示数据。</p></description>
            <item name="@type">富文本的格式，目前可以是 <var>html</var> 或是 <var>markdown</var></item>
            <item name=".">实际的文本内容，根据 <code>@type</code> 属性确定渲染的方式。</item>
        </type>

        <type name="version">
            <description>
                <p>版本号格式，遵循 <a href="https://semver.org/lang/zh-CN/">semver</a> 的规则。比如 <samp>1.1.1</samp>、<samp>0.1.0</samp> 等。</p>
            </description>
        </type>

        <type name="date">
            <description>
                <p>采用 <a href="https://tools.ietf.org/html/rfc3339">RFC3339</a> 格式表示的时间，比如：<samp>2019-12-16T00:35:48+08:00</samp></p>
            </description>
        </type>
    </types>

    <footer>
        <license>
            <p>当前页面内容托管于 </p><p>，并采用</p><p>进行许可。</p>
        </license>
    </footer>
</docs>
`),
	},
	{
		Name:        "index.xsl",
		ContentType: "text/xsl; charset=utf-8",
		Content: []byte(`<?xml version="1.0" encoding="UTF-8"?>

<xsl:stylesheet version="1.0" xmlns:xsl="http://www.w3.org/1999/XSL/Transform">
<xsl:output method="html" encoding="utf-8" indent="yes" version="5.0" doctype-system="about:legacy-compat" />

<xsl:variable name="curr-lang">
    <xsl:value-of select="/docs/@lang" />
</xsl:variable>

<!-- 获取当前文档的语言名称，如果水存在，则直接采用 @lang 属性 -->
<xsl:variable name="curr-lang-title">
    <xsl:variable name="title">
        <xsl:value-of select="document('locales.xml')/locales/locale[@id=$curr-lang]/@title" />
    </xsl:variable>

    <xsl:choose>
        <xsl:when test="$title=''"><xsl:value-of select="$curr-lang" /></xsl:when>
        <xsl:otherwise><xsl:value-of select="$title" /></xsl:otherwise>
    </xsl:choose>
</xsl:variable>

<xsl:variable name="keywords">
    <xsl:for-each select="document('config.xml')/config/languages/language">
        <xsl:value-of select="." /><xsl:value-of select="','" />
    </xsl:for-each>
</xsl:variable>

<xsl:template match="/">
    <html>
        <head>
            <title><xsl:value-of select="docs/title" /></title>
            <meta charset="UTF-8" />
            <meta name="viewport" content="width=device-width,initial-scale=1,maximum-scale=1"/>
            <meta name="keywords" content="{$keywords}RESTful API,document,apidoc" />
            <link rel="icon" type="image/svg+xml" href="./icon.svg" />
            <link rel="mask-icon" type="image/svg+xml" href="./icon.svg" color="black" />
            <link rel="canonical" href="{document('config.xml')/config/url}" />
            <link rel="stylesheet" type="text/css" href="./index.css" />
            <link rel="license" href="{/docs/liense/@url}" />
            <link rel="stylesheet" type="text/css" href="https://cdnjs.cloudflare.com/ajax/libs/prism/1.17.1/themes/prism-tomorrow.min.css" />
            <script src="./index.js" />
        </head>

        <body>
            <xsl:call-template name="header" />

            <main>
                <xsl:for-each select="docs/doc[not(@parent)]"> <!-- 不存在 @parent，表示是顶级的 doc 元素 -->
                    <xsl:call-template name="article">
                        <xsl:with-param name="doc" select="." />
                    </xsl:call-template>
                </xsl:for-each>
            </main>

            <footer>
                <div class="wrap">
                <p>
                    <xsl:value-of select="docs/footer/license/p[1]" />
                    <a href="{document('config.xml')/config/repo}">Github</a>
                    <xsl:value-of select="docs/footer/license/p[2]" />
                    <a href="{docs/license/@url}"><xsl:value-of select="docs/license" /></a>
                    <xsl:value-of select="docs/footer/license/p[3]" />
                </p>
                </div>
                <a href="#" class="goto-top" />
            </footer>

            <script src="https://cdnjs.cloudflare.com/ajax/libs/prism/1.17.1/components/prism-core.min.js"></script>
            <script src="https://cdnjs.cloudflare.com/ajax/libs/prism/1.17.1/plugins/autoloader/prism-autoloader.min.js"></script>
        </body>
    </html>
</xsl:template>

<xsl:template name="header">
    <header>
        <div class="wrap">
            <h1>
                <img src="./icon.svg" />
                <xsl:value-of select="document('config.xml')/config/name" />
                <span class="version">&#160;(<xsl:value-of select="document('config.xml')/config/version" />)</span>
            </h1>

            <div class="menus" role="navigation">
                <xsl:for-each select="docs/doc[not(@parent)]">
                    <a class="menu" href="#{@id}"><xsl:value-of select="@title" /></a>
                </xsl:for-each>

                <span class="drop-menus" role="menu">
                    <a class="menu">
                        <xsl:value-of select="$curr-lang-title" />
                        <span aria-hiddren="true">&#160;&#x25bc;</span>
                    </a>
                    <ul>
                        <xsl:for-each select="document('locales.xml')/locales/locale">
                            <li><a href="{@href}"><xsl:value-of select="@title" /></a></li>
                        </xsl:for-each>
                    </ul>
                </span>
            </div>
        </div>
    </header>
</xsl:template>

<!-- 将 doc 元素转换成 article HTML -->
<xsl:template name="article">
    <xsl:param name="doc" />

    <xsl:variable name="id" select="$doc/@id" />

    <article id="{$id}">
        <xsl:choose> <!-- 根据是否存在 $parent，决定是用 h3 还是 h2 标签 -->
            <xsl:when test="$doc/@parent">
                <h3>
                    <xsl:value-of select="$doc/@title" />
                    <a class="link" href="#{$id}">&#160;&#160;&#128279;</a>
                </h3>
            </xsl:when>
            <xsl:otherwise>
                <h2>
                    <xsl:value-of select="$doc/@title" />
                    <a class="link" href="#{$doc/@id}">&#160;&#160;&#128279;</a>
                </h2>
            </xsl:otherwise>
        </xsl:choose>

        <xsl:copy-of select="$doc/node()" />

        <xsl:for-each select="/docs/doc[@parent=$id]">
            <xsl:call-template name="article">
                <xsl:with-param name="doc" select="." />
            </xsl:call-template>
        </xsl:for-each>

        <xsl:for-each select="document('types.xml')/types/types[@parent=$id]/type">
            <xsl:call-template name="type">
                <xsl:with-param name="type" select="." />
                <xsl:with-param name="parent" select="$id" />
            </xsl:call-template>
        </xsl:for-each>
    </article>
</xsl:template>

<!-- 以下两个变量仅用于 type 模板，在模板中无法直接使用 /docs 元素，所以使用变量引用 -->
<xsl:variable name="header-locale" select="/docs/type-locale/header" />
<xsl:variable name="doc-types" select="/docs/types" />

<!-- 将类型显示为一个 table -->
<xsl:template name="type">
    <xsl:param name="type" />
    <xsl:param name="parent" />

    <xsl:variable name="type-locale" select="$doc-types[@parent=$parent]/type[@name=$type/@name]" />

    <article id="type_{$type/@name}">
        <h3>
            <xsl:value-of select="$type/@name" />
            <a class="link" href="#type_{$type/@name}">&#160;&#160;&#128279;</a>
        </h3>
        <xsl:copy-of select="$type-locale/description/node()" />
        <xsl:if test="item">
            <table>
                <thead>
                    <tr>
                        <th><xsl:copy-of select="$header-locale/name" /></th>
                        <th><xsl:copy-of select="$header-locale/type" /></th>
                        <th><xsl:copy-of select="$header-locale/required" /></th>
                        <th><xsl:copy-of select="$header-locale/description" /></th>
                    </tr>
                </thead>

                <tbody>
                    <xsl:for-each select="$type/item">
                    <xsl:variable name="name" select="@name" />
                    <tr>
                        <th><xsl:value-of select="@name" /></th>
                        <td><xsl:value-of select="@type" /></td>
                        <td>
                            <xsl:call-template name="checkbox">
                                <xsl:with-param name="chk" select="@required" />
                            </xsl:call-template>
                        </td>
                        <td><xsl:copy-of select="$type-locale/item[@name=$name]/node()" /></td>
                    </tr>
                    </xsl:for-each>
                </tbody>
            </table>
        </xsl:if>
    </article>
</xsl:template>

<xsl:template name="checkbox">
    <xsl:param name="chk" />
    <xsl:choose>
        <xsl:when test="$chk='true'">
            <input type="checkbox" checked="true" disabled="true" />
        </xsl:when>
        <xsl:otherwise>
            <input type="checkbox" disabled="true" />
        </xsl:otherwise>
    </xsl:choose>
</xsl:template>

</xsl:stylesheet>
`),
	},
	{
		Name:        "index.zh-hant.xml",
		ContentType: "application/xml; charset=utf-8",
		Content: []byte(`<?xml version="1.0" encoding="utf-8"?>
<?xml-stylesheet type="text/xsl" href="./index.xsl"?>

<!--
這是官網首頁內容，同時也是簡體中文的本地化內容。

其它語言的本化地內容，需要重新改寫本文件中除註釋外的所有內容。
-->

<docs lang="zh-hant">
    <title>apidoc | RESTful API 文檔處理工具</title>
    <license url="https://creativecommons.org/licenses/by/4.0/deed.zh">署名 4.0 國際 (CC BY 4.0)</license>

    <!-- 類型描述中表格的相關本化地信息 -->
    <type-locale>
        <header>
            <name>名稱</name>
            <type>類型</type>
            <required>必填</required>
            <description>描述</description>
        </header>
    </type-locale>

    <doc id="about" title="關於">
        <p>apidoc 是壹個簡單的 <abbr title="Representational State Transfer">RESTful</abbr> <abbr title="Application Programming Interface">API</abbr> 文檔生成工具，它從代碼註釋中提取特定格式的內容生成文檔。支持諸如 Go、Java、C++、Rust 等大部分開發語言，具體可使用 <code>apidoc lang</code> 命令行查看所有的支持列表。</p>

        <p>apidoc 擁有以下特點：</p>
        <ol>
            <li>跨平臺，linux、windows、macOS 等都支持；</li>
            <li>支持語言廣泛，即使是不支持，也很方便擴展；</li>
            <li>支持多個不同語言的多個項目生成壹份文檔；</li>
            <li>輸出模板可自定義；</li>
            <li>根據文檔生成 mock 數據；</li>
        </ol>

        <p>以下是壹段簡短的 C 語言風格代碼下的示例：</p>
        <pre><code class="language-markup"><![CDATA[/**
 * <api method="GET" summary="獲取所有的用戶信息">
 *     <path path="/users">
 *         <query name="page" type="number" default="0">顯示第幾頁的內容</query>
 *         <query name="size" type="number" default="20">每頁顯示的數量</query>
 *     </path>
 *     <tag>user</tag>
 *     <server>users</server>
 *     <response status="200" type="object" mimetype="application/json">
 *         <param name="count" type="int" optional="false" summary="符合條件的所有用戶數量" />
 *         <param name="users" type="object" array="true" summary="用戶列表">
 *             <param name="id" type="int" summary="唯壹 ID" />
 *             <param name="name" type="string" summary="姓名" />
 *         </param>
 *     </response>
 *     <response status="500" mimetype="application/json" type="obj">
 *         <param name="code" type="int" summary="錯誤代碼" />
 *         <param name="msg" type="string" summary="錯誤內容" />
 *     </response>
 * </api>
 */]]></code></pre>
        <p>多行註釋中，每壹行中以<code>空白字符+symbol+空白字符</code>開頭的，這些字符將會被過濾，symbol 表示該註釋塊的起始字符中的任意字符。比如以上代碼中，所有的 <var>*</var> 將被過濾。</p>
    </doc>

    <doc id="usage" title="使用" />

    <doc id="spec" title="文檔格式">
        <p>文檔采用 XML 格式。存在兩個頂級標簽：<code>apidoc</code> 和 <code>api</code>，用於描述整體內容和具體接口信息。</p>

        <p>文檔被從註釋中提取之後，最終會被合並成壹個 XML 文件，在該文件中 <code>api</code> 作為 <code>apidoc</code> 的壹個子元素存在，如果妳的項目不想把文檔寫在註釋中，也可以直接編寫壹個完整的 XML 文件，將 <code>api</code> 作為 <code>apidoc</code> 的壹個子元素。</p>

        <p>具體可參考<a href="./example/index.xml">示例代碼。</a></p>

        <p>以下是對各個 XML 元素以及參數介紹，其中以 <code>@</code> 開頭的表示 XML 屬性；<code>.</code> 表示為當前元素的內容；其它表示子元素。</p>
    </doc>

    <doc id="install" title="安裝" parent="usage">
          <p>可以直接從 <a href="https://github.com/caixw/apidoc/releases">https://github.com/caixw/apidoc/releases</a> 查找妳需要的版本下載，放入 <code>PATH</code> 中即可使用。如果沒有妳需要的平臺文件，則需要從源代碼編譯：</p>
        <ul>
            <li>下載 Go 編譯工具</li>
            <li>下載源代碼</li>
            <li>執行代碼中 <code>build/build.sh</code> 或是 <code>build/build.cmd</code> 進行編譯</li>
            <li>編譯好的文件存放在 cmd/apidoc 下，可以將該文件放置在 PATH 目錄</li>
        </ul>
    </doc>

    <doc id="env" title="環境變量" parent="usage">
        <p>apidoc 會讀取 <var>LANG</var> 的環境變量作為其本地化的依據，若想指定其它語種，可以手動指定 <var>LANG</var> 環境變量：<samp>LANG=zh-Hant apidoc</samp>。在 windows 系統中，若不存在 <var>LANG</var> 環境變量，則會調用 <samp>GetUserDefaultLocaleName</samp> 函數來獲取相應的語言信息。</p>
    </doc>

    <doc id="cli" title="命令行" parent="usage">
        <p>可以通過 <samp>apidoc help</samp> 查看命令行支持的子命令。主要包含了以下幾個：</p>
        <table>
            <thead><tr><th>子命令</th><th>描述</th></tr></thead>
            <tbody>
                <tr><td>help</td><td>顯示子命令的描述信息</td></tr>
                <tr><td>build</td><td>生成文檔內容</td></tr>
                <tr><td>mock</td><td>根據文檔提供 mock 服務</td></tr>
                <tr><td>static</td><td>提供查看文檔的本地服務</td></tr>
                <tr><td>version</td><td>顯示版本信息</td></tr>
                <tr><td>lang</td><td>列出當前支持的語言</td></tr>
                <tr><td>locale</td><td>列出當前支持的本地化內容</td></tr>
                <tr><td>detect</td><td>根據指定的目錄生成配置文件</td></tr>
                <tr><td>test</td><td>檢測語法是否準確</td></tr>
            </tbody>
        </table>
        <p>mock 子命令可以根據文檔生成壹些符合要求的隨機數據。這些數據每次請求都不相同，包括數量、長度、數值大小等。</p>
    </doc>

    <!-- 配置文件的類型定義 -->
    <types parent="usage">
        <type name=".apidoc.yaml">
            <description>
                <p>配置文件名固定為 <code>.apidoc.yaml</code>，格式為 YAML，可參考 <a href="example/.apidoc.yaml">.apidoc.yaml</a>。文件可以通過命令 <code>apidoc detect</code> 生成。主要包含了以幾個配置項：</p>
            </description>
            <item name="version" >產生此配置文件的 apidoc 版本</item>
            <item name="inputs">指定輸入的數據，同壹項目只能解析壹種語言。</item>
            <item name="inputs.dir">需要解析的源文件所在目錄</item>
            <item name="inputs.recursive">是否解析子目錄下的源文件</item>
            <item name="inputs.encoding">編碼，默認為 <code>utf-8</code>，值可以是 <a href="https://www.iana.org/assignments/character-sets/character-sets.xhtml">character-sets</a> 中的內容。</item>
            <item name="inputs.lang">源文件類型。具體支持的類型可通過 -l 參數進行查找</item>
            <item name="output">控制輸出行為</item>
            <item name="output.path">指定輸出的文件名，包含路徑信息。</item>
            <item name="output.type">指定輸出的文件格式，值可以是 <var>apidoc+xml</var>、<var>openapi+json</var> 和 <var>openapi+yaml</var>，其中 <var>apidoc+xml</var> 為默認值。</item>
            <item name="output.tags">只輸出與這些標簽相關聯的文檔，默認為全部。</item>
            <item name="output.style">為 XML 文件指定的 XSL 文件。</item>
        </type>
    </types>

    <types parent="spec">
        <type name="apidoc">
            <description><p>用於描述整個文檔的相關內容，只能出現壹次。</p></description>
            <item name="@version">文檔的版本</item>
            <item name="@lang">內容的本地化 ID，比如 <samp><var>zh-hans</var></samp> 等。</item>
            <item name="@logo">圖標，默認采用官網的 <var>https://apidoc.tools/icon.svg</var>，同時作用於 favicon 和 logo，只支持 SVG 格式。</item>
            <item name="@created">文檔的生成時間</item>
            <item name="title">文檔的標題</item>
            <item name="description">文檔的整體介紹，可以是使用 HTML 內容。</item>
            <item name="contract">聯系人信息</item>
            <item name="license">內容的版權信息</item>
            <item name="tag">可以用的標簽列表</item>
            <item name="server">API 基地址列表，每個 API 最少應該有壹個 server。</item>
            <item name="mimetype">接口所支持的 mimetype 類型</item>
            <item name="response">表示所有 API 都有可能返回的內容</item>
            <item name="api">API 文檔內容</item>
        </type>

        <type name="link">
            <description><p>用於描述鏈接，壹般轉換為 HTML 的 a 標簽。</p></description>
            <item name="@url">鏈接指向的 URL</item>
            <item name="@text">鏈接的文本內容</item>
        </type>

        <type name="contact">
            <description><p>用於描述聯系方式</p></description>
            <item name="@url">鏈接的 URL，與郵箱必須二選壹必填</item>
            <item name="@email">郵件地址，與 url 必須二選壹必填</item>
            <item name=".">聯系人名稱</item>
        </type>

        <type name="tag">
            <description><p>定義標簽，標簽相當於關鍵字，作用於 API，相當於啟到分類的作用。</p></description>
            <item name="@name">標簽的唯壹 ID，推薦采用英文字母表示。</item>
            <item name="@title">標簽名稱</item>
            <item name="@deprecated">表示該標簽在大於等於該版本號時不再啟作用</item>
        </type>

        <type name="server">
            <description><p>定義服務器的相關信息，作用於 API，決定該 API 與哪個服務器相關聯。</p></description>
            <item name="@name">唯壹 ID，推薦采用英文字母表示。</item>
            <item name="@url">服務基地址</item>
            <item name="@deprecated">表示在大於等於該版本號時不再啟作用</item>
            <item name="@summary">簡要的描述內容，或者通過 <code>description</code> 提供壹份富文本內容。</item>
            <item name="description">對該服務的具體描述，可以使用 HTML 內容</item>
        </type>

        <type name="api">
            <description><p>定義接口的具體內容</p></description>
            <item name="@version">表示此接口在該版本中添加</item>
            <item name="@method">請求方法</item>
            <item name="@summary">簡要介紹</item>
            <item name="@deprecated">表示在大於等於該版本號時不再啟作用</item>
            <item name="@id">唯壹 ID</item>
            <item name="description">該接口的詳細介紹，為 HTML 內容。</item>
            <item name="path">定義路徑信息</item>
            <item name="request">定義可用的請求信息</item>
            <item name="response">定義可能的返回信息</item>
            <item name="callback">定義回調接口內容</item>
            <item name="tag">關聯的標簽</item>
            <item name="server">關聯的服務</item>
            <item name="header">傳遞的報頭內容，如果是某個 mimetype 專用的，可以放在 request 元素中。</item>
        </type>

        <type name="path">
            <description><p>用於定義請求時與路徑相關的內容</p></description>
            <item name="@path">接口地址</item>
            <item name="param">地址中的參數</item>
            <item name="query">地址中的查詢參數</item>
        </type>

        <type name="request">
            <description><p>定義了請求和返回的相關內容</p></description>
            <item name="@xml-ns">XML 標簽的命名空間</item>
            <item name="@xml-ns-prefix">XML 標簽的命名空間名稱前綴</item>
            <item name="@xml-wrapped">如果当前元素的 <code>@array</code> 为 <var>true</var>，是否将其包含在 wrapped 指定的标签中。</item>
            <item name="@name">當 mimetype 為 <var>application/xml</var> 時，此值表示 XML 的頂層元素名稱，否則無用。</item>
            <item name="@type">值的類型，可以是 <del title="建議使用空值代替"><var>none</var></del>、<var>string</var>、<var>number</var>、<var>bool</var>、<var>object</var> 和 空值；空值表示不輸出任何內容。</item>
            <item name="@deprecated">表示在大於等於該版本號時不再啟作用</item>
            <item name="@summary">簡要介紹</item>
            <item name="@array">是否為數組</item>
            <item name="@status">狀態碼。在 request 中，該值不可用，否則為必填項。</item>
            <item name="@mimetype">媒體類型，比如 <var>application/json</var> 等。</item>
            <item name="description">詳細介紹，為 HTML 內容。</item>
            <item name="enum">當前參數可用的枚舉值</item>
            <item name="param">子類型，比如對象的子元素。</item>
            <item name="example">示例代碼。</item>
            <item name="header">傳遞的報頭內容</item>
        </type>

        <type name="param">
            <description><p>參數類型，基本上可以作為 <code>request</code> 的子集使用。</p></description>
            <item name="@xml-attr">是否作為父元素的屬性，僅作用於 XML 元素。</item>
            <item name="@xml-extract">將當前元素的內容作為父元素的內容，要求父元素必須為 <var>object</var>。</item>
            <item name="@xml-ns">XML 標簽的命名空間</item>
            <item name="@xml-ns-prefix">XML 標簽的命名空間名稱前綴</item>
            <item name="@xml-attr">是否作為父元素的屬性，僅用於 XML 的請求。</item>
            <item name="@xml-wrapped">如果当前元素的 <code>@array</code> 为 <var>true</var>，是否将其包含在 wrapped 指定的标签中。</item>
            <item name="@name">值的名稱</item>
            <item name="@type">值的類型，可以是 <var>string</var>、<var>number</var>、<var>bool</var> 和 <var>object</var></item>
            <item name="@deprecated">表示在大於等於該版本號時不再啟作用</item>
            <item name="@default">默認值</item>
            <item name="@optional">是否為可選的參數</item>
            <item name="@summary">簡要介紹</item>
            <item name="@array">是否為數組</item>
            <item name="@array-style">是否以數組的形式展示數據，默認采用 form 形式，僅在 <code>@array</code> 為 <var>true</var> 時有效。</item>
            <item name="description">詳細介紹，為 HTML 內容。</item>
            <item name="enum">當前參數可用的枚舉值</item>
            <item name="param">子類型，比如對象的子元素。</item>
        </type>

        <type name="enum">
            <description><p>定義枚舉類型的數所的枚舉值</p></description>
            <item name="@value">枚舉值</item>
            <item name="@deprecated">表示在大於等於該版本號時不再啟作用</item>
            <item name=".">該值的詳細介紹</item>
        </type>

        <type name="example">
            <description><p>示例代碼</p></description>
            <item name="@mimetype">代碼的 mimetype 類型。</item>
            <item name=".">示例代碼的內容，需要使用 CDATA 包含代碼。</item>
        </type>

        <type name="header">
            <description><p>定義了請求和返回的報頭結構</p></description>
            <item name="@name">報頭的名稱</item>
            <item name="@deprecated">表示在大於等於該版本號時不再啟作用</item>
            <item name="@summary">對報頭的描述</item>
            <item name="description">對報頭的描述</item>
        </type>

        <type name="callback">
            <description><p>定義接口回調的相關內容</p></description>
            <item name="@method">請求方法</item>
            <item name="@summary">簡要介紹</item>
            <item name="@deprecated">表示在大於等於該版本號時不再啟作用</item>
            <item name="description">該接口的詳細介紹</item>
            <item name="path">定義路徑信息</item>
            <item name="request">定義可用的請求信息</item>
            <item name="response">定義可能的返回信息</item>
        </type>

        <type name="richtext">
            <description><p>富文本信息，可以以不同的格式展示數據。</p></description>
            <item name="@type">富文本的格式，目前可以是 <var>html</var> 或是 <var>markdown</var></item>
            <item name=".">實際的文本內容，根據 <code>@type</code> 屬性確定渲染的方式。</item>
        </type>

        <type name="version">
            <description>
                <p>版本號格式，遵循 <a href="https://semver.org/lang/zh-CN/">semver</a> 的規則。比如 <samp>1.1.1</samp>、<samp>0.1.0</samp> 等。</p>
            </description>
        </type>

        <type name="date">
            <description>
                <p>采用 <a href="https://tools.ietf.org/html/rfc3339">RFC3339</a> 格式表示的時間，比如：<samp>2019-12-16T00:35:48+08:00</samp></p>
            </description>
        </type>
    </types>

    <footer>
        <license>
            <p>當前頁面內容托管於 </p><p>，並采用</p><p>進行許可。</p>
        </license>
    </footer>
</docs>
`),
	},
	{
		Name:        "locales.xml",
		ContentType: "application/xml; charset=utf-8",
		Content: []byte(`<?xml version="1.0" encoding="utf-8"?>

<locales>
    <locale id="zh-hans" href="index.xml" title="简体中文" />
    <locale id="zh-hant" href="index.zh-hant.xml" title="繁体中文" />
</locales>
`),
	},
	{
		Name:        "types.xml",
		ContentType: "application/xml; charset=utf-8",
		Content: []byte(`<?xml version="1.0" encoding="utf-8"?>

<types>

    <!-- 配置文件的类型定义 -->
    <types parent="usage">
        <type name=".apidoc.yaml">
            <item name="version" type="version" required="true" />
            <item name="inputs" type="object[]" required="true" />
            <item name="inputs.dir" type="string" required="true" />
            <item name="inputs.recursive" type="bool" required="false" />
            <item name="inputs.encoding" type="string" required="false" />
            <item name="inputs.lang" type="string" required="true" />
            <item name="output" type="object" required="true" />
            <item name="output.path" type="string" required="true" />
            <item name="output.type" type="string" required="false" />
            <item name="output.tags" type="string[]" required="false" />
            <item name="output.style" type="string" required="false" />
        </type>
    </types>

    <types parent="spec">
        <type name="apidoc">
            <item name="@version" type="version" required="true" />
            <item name="@lang" type="string" required="false" />
            <item name="@logo" type="string" required="false" />
            <item name="@created" type="date" required="false" />
            <item name="title" type="string" required="true" />
            <item name="description" type="richtext" required="false" />
            <item name="contract" type="contract" required="false" />
            <item name="license" type="link" required="false" />
            <item name="tag" type="tag[]" required="false" />
            <item name="server" type="server[]" required="true" />
            <item name="mimetype" type="string[]" required="true" />
            <item name="response" type="request[]" required="false" />
            <item name="api" type="api[]" required="false" />
        </type>

        <type name="link">
            <item name="@url" type="string" required="true" />
            <item name="@text" type="string" required="true" />
        </type>

        <type name="contact">
            <item name="@url" type="string" required="true" />
            <item name="@email" type="string" required="true" />
            <item name="." type="string" required="true" />
        </type>

        <type name="tag">
            <item name="@name" type="string" required="true" />
            <item name="@title" type="string" required="true" />
            <item name="@deprecated" type="version" required="false" />
        </type>

        <type name="server">
            <item name="@name" type="string" required="true" />
            <item name="@url" type="string" required="true" />
            <item name="@summary" type="string" required="false" />
            <item name="@deprecated" type="version" required="false" />
            <item name="description" type="richtext" required="false" />
        </type>

        <type name="api">
            <item name="@version" type="version" required="false" />
            <item name="@method" type="string" required="true" />
            <item name="@summary" type="string" required="true" />
            <item name="@deprecated" type="version" required="false" />
            <item name="@id" type="string" required="false" />
            <item name="description" type="richtext" required="false" />
            <item name="path" type="path" required="false" />
            <item name="request" type="request[]" required="false" />
            <item name="response" type="request[]" required="false" />
            <item name="callback" type="callback" required="false" />
            <item name="tag" type="string[]" required="false" />
            <item name="server" type="string[]" required="false" />
            <item name="header" type="header[]" required="false" />
        </type>

        <type name="path">
            <item name="@path" type="string" required="true" />
            <item name="param" type="param[]" required="false" />
            <item name="query" type="param[]" required="false" />
        </type>

        <type name="request">
            <item name="@xml-ns" type="bool" required="false" />
            <item name="@xml-ns-prefix" type="bool" required="false" />
            <item name="@xml-wrapped" type="string" required="false" />
            <item name="@name" type="string" required="true" />
            <item name="@type" type="string" required="false" />
            <item name="@deprecated" type="version" required="false" />
            <item name="@summary" type="string" required="true" />
            <item name="@array" type="bool" required="false" />
            <item name="@status" type="number" required="true" />
            <item name="@mimetype" type="string" required="false" />
            <item name="description" type="richtext" required="false" />
            <item name="enum" type="enum[]" required="false" />
            <item name="param" type="param[]" required="false" />
            <item name="example" type="example[]" required="false" />
            <item name="header" type="header[]" required="false" />
        </type>

        <type name="param">
            <item name="@xml-attr" type="bool" required="false" />
            <item name="@xml-extract" type="bool" required="false" />
            <item name="@xml-ns" type="bool" required="false" />
            <item name="@xml-ns-prefix" type="bool" required="false" />
            <item name="@xml-wrapped" type="string" required="false" />
            <item name="@name" type="string" required="true" />
            <item name="@type" type="string" required="true" />
            <item name="@deprecated" type="version" required="false" />
            <item name="@default" type="string" required="false" />
            <item name="@optional" type="bool" required="false" />
            <item name="@summary" type="string" required="true" />
            <item name="@array" type="bool" required="false" />
            <item name="description" type="richtext" required="false" />
            <item name="@array-style" type="bool" required="false" />
            <item name="enum" type="enum[]" required="false" />
            <item name="param" type="param[]" required="false" />
        </type>

        <type name="enum">
            <item name="@value" type="string" required="true" />
            <item name="@deprecated" type="version" required="false" />
            <item name="." type="string" required="true" />
        </type>

        <type name="example">
            <item name="@mimetype" type="string" required="true" />
            <item name="." type="string" required="true" />
        </type>

        <type name="header">
            <item name="@name" type="string" required="true" />
            <item name="@deprecated" type="version" required="false" />
            <item name="@summary" type="string" required="true" />
            <item name="description" type="richtext" required="true" />
        </type>

        <type name="callback">
            <item name="@method" type="string" required="true" />
            <item name="@summary" type="string" required="true" />
            <item name="@deprecated" type="version" required="false" />
            <item name="description" type="richtext" required="false" />
            <item name="path" type="path" required="true" />
            <item name="request" type="request[]" required="true" />
            <item name="response" type="request[]" required="true" />
        </type>

        <type name="richtext">
            <item name="@type" type="string" required="true" />
            <item name="." type="string" required="true" />
        </type>

        <type name="version" />

        <type name="date" />
    </types>
</types>
`),
	},
	{
		Name:        "v5/apidoc.css",
		ContentType: "text/css; charset=utf-8",
		Content: []byte(`@charset "utf-8";

:root {
    --max-width: 2048px;
    --min-width: 200px; /* 当列的宽度小于此值时，部分行内容会被从横向改变纵向排列。 */
    --padding: 1rem;
    --article-padding: calc(var(--padding) / 2);

    --color: black;
    --accent-color: #0074d9;
    --background: white;
    --border-color: #e0e0e0;
    --delete-color: red;

    /* method */
    --method-get-color: green;
    --method-options-color: green;
    --method-post-color: darkorange;
    --method-put-color: darkorange;
    --method-patch-color: darkorange;
    --method-delete-color: red;
}

@media (prefers-color-scheme: dark) {
    :root {
        --color: #b0b0b0;
        --accent-color: #0074d9;
        --background: black;
        --border-color: #303030;
        --delete-color: red;

        /* method */
        --method-get-color: green;
        --method-options-color: green;
        --method-post-color: darkorange;
        --method-put-color: darkorange;
        --method-patch-color: darkorange;
        --method-delete-color: red;
    }
}

body {
    padding: 0;
    margin: 0;
    color: var(--color);
    background: var(--background);
    text-align: center;
}

table {
    width: 100%;
}

table th, table td {
    font-weight: normal;
    text-align: left;
    border-bottom: 1px solid transparent;
}

table tr:hover th,
table tr:hover td {
    border-bottom: 1px solid var(--border-color);
}

ul, ol, ul li, ol li {
    padding: 0;
    margin: 0;
    list-style-position: inside;
}

p {
    margin: 0;
}

summary, input {
    outline: none;
}

a {
    text-decoration: none;
    color: var(--accent-color);
}

a:hover {
    opacity: .7;
}

.del {
    text-decoration: line-through;
    text-decoration-color: var(--delete-color);
}

.hidden {
    display: none;
}

/*************************** header ***********************/

header {
    position: sticky;
    top: 0;
    display: block;
    z-index: 1000;
    box-sizing: border-box;
    background: var(--background);
    box-shadow: 2px 2px 2px var(--border-color);
}

header .wrap {
    margin: 0 auto;
    text-align: left;
    max-width: var(--max-width);
    padding: 0 var(--padding);
    display: flex;
    flex-wrap: wrap;
    justify-content: space-between;
    align-items: center;
}

header h1 {
    margin: var(--padding) 0 var(--article-padding);
    display: inline-block;
}

header h1 .version {
    font-size: 1rem;
}

header h1 img {
    height: 1.5rem;
    margin-right: .5rem;
}


header .menus {
    display: flex;
}

header .menu {
    cursor: pointer;
    position: relative;
    margin-left: var(--padding);
    display: none; /* 默认不可见，只有在有 JS 的情况下，通过 JS 控制其可见性 */
}

header .menu:hover span {
    opacity: .7;
}

header .menu ul {
    position: absolute;
    min-width: 4rem;
    right: 0;
    display: none;
    list-style: none;
    background: var(--background);
    border: 1px solid var(--border-color);
    padding: var(--article-padding);
}

header .menu ul li {
    padding-bottom: var(--article-padding);
}

/* 可以保证 label 的内容在同一行 */
header .menu ul li label {
    display: inline-flex;
    align-items: baseline;
    word-break: keep-all;
    white-space: nowrap;
}

header .menu ul li:last-of-type {
    padding-bottom: 0;
}

header .menu:hover ul {
    display: block;
}

/*************************** main ***********************/

main {
    padding: 0 var(--padding);
    margin: 0 auto;
    max-width: var(--max-width);
    text-align: left;
}

main .content {
    margin: var(--padding) 0;
}

/****************** .servers *******************/

main .servers {
    display: flex;
    flex-flow: wrap;
    justify-content: space-between;
    margin-bottom: var(--padding);
}

main .servers .server {
    flex-grow: 1;
    min-width: var(--min-width);
    box-sizing: border-box;
    border: 1px solid var(--border-color);
    padding: var(--padding) var(--article-padding);
}

main .servers .server:hover {
    border: 1px solid var(--accent-color);
}

main .servers .server h4 {
    margin: 0 0 var(--padding);
}

/********************** .api **********************/

main .api {
    margin-bottom: var(--article-padding);
    border: 1px solid var(--border-color);
}

main .api summary {
    margin: 0;
    padding: var(--article-padding);
    border-bottom: 1px solid var(--border-color);
    cursor: pointer;
    line-height: 1;
}


main details.api:not([open]) summary {
    border: none;
}

main .api summary .action {
    min-width: 4rem;
    font-weight: bold;
    display: inline-block;
    margin-right: 1rem;
}

main .api summary .link {
    margin-right: 10px;
    text-decoration: none;
}

main .api .callback .summary,
main .api summary .summary {
    float: right;
    font-weight: 400;
    opacity: .5;
}

main .api .description {
    padding: var(--article-padding);
    margin: 0;
    border-bottom: 1px solid var(--border-color);
}

main .api[data-method='GET,'][open],
main .api[data-method='GET,']:hover,
main .callback[data-method='GET,'] h3 {
    border: 1px solid var(--method-get-color);
}
main .api[data-method='GET,'] summary {
    border-bottom: 1px solid var(--method-get-color);
}

main .api[data-method='POST,'][open],
main .api[data-method='POST,']:hover,
main .callback[data-method='POST,'] h3 {
    border: 1px solid var(--method-post-color);
}
main .api[data-method='POST,'] summary {
    border-bottom: 1px solid var(--method-post-color);
}

main .api[data-method='PUT,'][open],
main .api[data-method='PUT,']:hover,
main .callback[data-method='PUT,'] h3 {
    border: 1px solid var(--method-put-color);
}
main .api[data-method='PUT,'] summary {
    border-bottom: 1px solid var(--method-put-color);
}

main .api[data-method='PATCH,'][open],
main .api[data-method='PATCH,']:hover,
main .callback[data-method='PATCH,'] h3 {
    border: 1px solid var(--method-patch-color);
}
main .api[data-method='PATCH,'] summary {
    border-bottom: 1px solid var(--method-patch-color);
}

main .api[data-method='DELETE,'][open],
main .api[data-method='DELETE,']:hover,
main .callback[data-method='DELETE,'] h3 {
    border: 1px solid var(--method-delete-color);
}
main .api[data-method='DELETE,'] summary {
    border-bottom: 1px solid var(--method-delete-color);
}

main .api[data-method='OPTIONS,'][open],
main .api[data-method='OPTIONS,']:hover,
main .callback[data-method='OPTIONS,'] h3 {
    border: 1px solid var(--method-options-color);
}
main .api[data-method='OPTIONS,'] summary {
    border-bottom: 1px solid var(--method-options-color);
}

main .callback h3 {
    padding: var(--article-padding) var(--padding);
    margin: 0;
    border-left: none !important;
    border-right: none !important;
    cursor: pointer;
    line-height: 1;
}

main .api .body {
    display: flex;
    flex-flow: wrap;
}

main .api .body .requests,
main .api .body .responses {
    flex: 1 1 50%;
    box-sizing: border-box;
    min-width: var(--min-width);
    padding: var(--article-padding);
}
main .api .body .requests {
    border-right: 1px dotted var(--border-color);
}

main .api .body .requests .title,
main .api .body .responses .title {
    margin: 0;
    opacity: .5;
}

main .api .param {
    margin-top: var(--padding);
}

main .api .param .title,
main .api .param .title {
    margin: 0;
    opacity: 1 !important;
    font-weight: normal;
}

main .api .param .example {
    display: none;
    margin: 0;
}

main .api .param .toggle-example {
    cursor: pointer;
}

main .api .body .responses .status,
main .api .body .requests .status {
    margin: calc(var(--padding) + var(--article-padding)) 0 var(--article-padding);
    border-bottom: 1px solid var(--border-color);
    padding-bottom: 3px;
}

main .api .param .parent-type,
main .api .status .mimetype {
    opacity: .5;
}

/*************************** footer ***********************/

footer {
    border-top: 1px solid var(--border-color);
    padding: var(--padding) var(--padding);
    text-align: left;
    margin-top: var(--padding);
}

footer .wrap {
    margin: 0 auto;
    max-width: var(--max-width);
    display: flex;
    flex-flow: wrap;
    justify-content: space-between;
}
`),
	},
	{
		Name:        "v5/apidoc.js",
		ContentType: "application/javascript; charset=utf-8",
		Content: []byte(`'use strict';

window.onload = function () {
    registerFilter('method');
    registerFilter('server');
    registerFilter('tag');
    registerExpand();
    registerLanguageFilter();

    initExample();

    prettyDescription();
};

function registerFilter(type) {
    const menu = document.querySelector('.' + type + '-selector');
    if (menu === null) { // 可能为空，表示不存在该过滤项
        return;
    }

    menu.style.display = 'block'; // 有 JS 的情况下，展示过滤菜单

    menu.querySelectorAll('li input').forEach((val) => {
        val.addEventListener('change', (event) => {
            if (event.target === null) {
                return;
            }

            const chk = event.target.checked;
            const tag = event.target.parentNode.parentNode.getAttribute('data-' + type);
            document.querySelectorAll('.api').forEach((api) => {
                if (!api.getAttribute('data-' + type).includes(tag + ',')) {
                    return;
                }

                api.setAttribute("data-hidden-" + type, chk ? "" : "true");

                const hidden = api.getAttribute('data-hidden-tag') === 'true' ||
                    api.getAttribute('data-hidden-server') === 'true' ||
                    api.getAttribute('data-hidden-method') === 'true';
                api.style.display = hidden ? 'none' : 'block';
            });
        }); // end addEventListener('change')
    }); // end forEach('li input')
}

function registerLanguageFilter() {
    const menu = document.querySelector('.languages-selector');

    menu.style.display = 'block';

    menu.querySelectorAll('li input').forEach((val) => {
        val.addEventListener('change', (event) => {
            if (event.target === null || !event.target.checked) {
                return;
            }

            const lang = event.target.parentNode.parentNode.getAttribute('lang');
            document.querySelectorAll('[data-locale]').forEach((elem) => {
                elem.className = elem.getAttribute('lang') === lang ? '' : 'hidden';
            });
        }); // end addEventListener('change')
    }); // end forEach('li input')
}

function registerExpand() {
    const expand = document.querySelector('.expand-selector');
    if (expand === null) {
        return;
    }

    expand.style.display = 'block';

    expand.querySelector('input').addEventListener('change', (event) => {
        const chk = event.target.checked;
        document.querySelectorAll('details').forEach((elem) => {
            elem.open = chk;
        });
    });
}

function initExample() {
    document.querySelectorAll('.toggle-example').forEach((btn) => {
        btn.addEventListener('click', (event) => {
            if (event.target === null) {
                return;
            }

            const parent = event.target.parentNode.parentNode.parentNode;
            const table = parent.querySelector('table');
            const pre = parent.querySelector('pre');

            if (table.getAttribute('data-visible') === 'true') {
                table.setAttribute('data-visible', 'false');
                table.style.display = 'none';
            } else {
                table.setAttribute('data-visible', 'true');
                table.style.display = 'table';
            }

            if (pre.getAttribute('data-visible') === 'true') {
                pre.setAttribute('data-visible', 'false');
                pre.style.display = 'none';
            } else {
                pre.setAttribute('data-visible', 'true');
                pre.style.display = 'block';
            }
        });
    });
}

// 美化描述内容
//
// 即将 html 内容转换成真的 HTML 格式，而 markdown 则依然是 pre 显示。
function prettyDescription() {
    document.querySelectorAll('[data-type]').forEach((elem) => {
        const type = elem.getAttribute('data-type');
        if (type !== 'html') {
            return;
        }

        elem.innerHTML = elem.getElementsByTagName('pre')[0].innerText;
    });
}
`),
	},
	{
		Name:        "v5/apidoc.xsl",
		ContentType: "text/xsl; charset=utf-8",
		Content: []byte(`<?xml version="1.0" encoding="UTF-8"?>

<xsl:stylesheet version="1.0" xmlns:xsl="http://www.w3.org/1999/XSL/Transform">
<xsl:import href="./locales.xsl" />
<xsl:output
    method="html"
    encoding="utf-8"
    indent="yes"
    version="5.0"
    doctype-system="about:legacy-compat"
/>

<xsl:template match="/">
    <html lang="{$curr-lang}">
        <head>
            <title><xsl:value-of select="apidoc/title" /></title>
            <meta charset="UTF-8" />
            <meta name="viewport" content="width=device-width,initial-scale=1,maximum-scale=1" />
            <meta name="generator" content="apidoc" />
            <link rel="icon" type="image/svg+xml" href="{$icon}" />
            <link rel="mask-icon" type="image/svg+xml" href="{$icon}" color="black" />
            <xsl:if test="apidoc/license"><link rel="license" href="{apidoc/license/@url}" /></xsl:if>
            <link rel="stylesheet" type="text/css" href="{$base-url}apidoc.css" />
            <script src="{$base-url}apidoc.js"></script>
        </head>
        <body>
            <xsl:call-template name="header" />

            <main>
                <div class="content" data-type="{description/@type}">
                    <pre><xsl:copy-of select="apidoc/description/node()" /></pre>
                </div>
                <div class="servers"><xsl:apply-templates select="apidoc/server" /></div>
                <xsl:apply-templates select="apidoc/api" />
            </main>

            <footer>
            <div class="wrap">
                <xsl:if test="apidoc/license"><xsl:copy-of select="$locale-license" /></xsl:if>
                <xsl:copy-of select="$locale-generator" />
            </div>
            </footer>
        </body>
    </html>
</xsl:template>

<xsl:template match="/apidoc/server">
    <div class="server">
        <h4>
            <xsl:call-template name="deprecated">
                <xsl:with-param name="deprecated" select="@deprecated" />
            </xsl:call-template>
            <xsl:value-of select="@name" />
        </h4>

        <p><xsl:value-of select="@url" /></p>
        <div>
            <xsl:choose>
                <xsl:when test="description">
                    <xsl:attribute name="data-type">
                        <xsl:value-of select="description/@type" />
                    </xsl:attribute>
                    <pre><xsl:copy-of select="description/node()" /></pre>
                </xsl:when>
                <xsl:otherwise><xsl:value-of select="@summary" /></xsl:otherwise>
            </xsl:choose>
        </div>
    </div>
</xsl:template>

<!-- header 界面元素 -->
<xsl:template name="header">
<header>
<div class="wrap">
    <h1>
        <img src="{$icon}" />
        <xsl:value-of select="/apidoc/title" />
        <span class="version">&#160;(<xsl:value-of select="/apidoc/@version" />)</span>
    </h1>

    <div class="menus">
        <!-- expand -->
        <label class="menu expand-selector" role="checkbox">
            <input type="checkbox" /><xsl:copy-of select="$locale-expand" />
        </label>

        <!-- server -->
        <xsl:if test="apidoc/server">
        <div class="menu server-selector" role="menu" aria-haspopup="true">
            <xsl:copy-of select="$locale-server" />
            <span aria-hiddren="true">&#160;&#x25bc;</span>
            <ul role="menu" aria-hiddren="true">
                <xsl:for-each select="apidoc/server">
                <li data-server="{@name}" role="menuitemcheckbox">
                    <label>
                        <input type="checkbox" checked="checked" />&#160;<xsl:value-of select="@name" />
                    </label>
                </li>
                </xsl:for-each>
            </ul>
        </div>
        </xsl:if>

        <!-- tag -->
        <xsl:if test="apidoc/tag">
        <div class="menu tag-selector" role="menu" aria-haspopup="true">
            <xsl:copy-of select="$locale-tag" />
            <span aria-hiddren="true">&#160;&#x25bc;</span>
            <ul role="menu" aria-hiddren="true">
                <xsl:for-each select="apidoc/tag">
                <li data-tag="{@name}" role="menuitemcheckbox">
                    <label>
                        <input type="checkbox" checked="checked" />&#160;<xsl:value-of select="@title" />
                    </label>
                </li>
                </xsl:for-each>
            </ul>
        </div>
        </xsl:if>

        <!-- method -->
        <div class="menu method-selector" role="menu" aria-haspopup="true">
            <xsl:copy-of select="$locale-method" />
            <span aria-hiddren="true">&#160;&#x25bc;</span>
            <ul role="menu" aria-hiddren="true">
                <!-- 浏览器好像都不支持 xpath 2.0，所以无法使用 distinct-values -->
                <!-- xsl:for-each select="distinct-values(/apidoc/api/@method)" -->
                <xsl:for-each select="/apidoc/api/@method[not(../preceding-sibling::api/@method = .)]">
                <li data-method="{.}" role="menuitemcheckbox">
                    <label><input type="checkbox" checked="checked" />&#160;<xsl:value-of select="." /></label>
                </li>
                </xsl:for-each>
            </ul>
        </div>

        <!-- language -->
        <div class="menu languages-selector" role="menu" aria-haspopup="true">
            <xsl:copy-of select="$locale-language" />
            <span aria-hiddren="true">&#160;&#x25bc;</span>
            <ul role="menu" aria-hiddren="true"><xsl:call-template name="languages" /></ul>
        </div>
    </div> <!-- end .menus -->
</div> <!-- end .wrap -->
</header>
</xsl:template>

<!-- api 界面元素 -->
<xsl:template match="/apidoc/api">
    <xsl:variable name="id" select="concat(@method, translate(path/@path, $id-from, $id-to))" />

    <details id="{$id}" class="api" data-method="{@method},">
    <xsl:attribute name="data-tag">
        <xsl:for-each select="tag"><xsl:value-of select="concat(., ',')" /></xsl:for-each>
    </xsl:attribute>
    <xsl:attribute name="data-server">
        <xsl:for-each select="server"><xsl:value-of select="concat(., ',')" /></xsl:for-each>
    </xsl:attribute>

        <summary>
            <a class="link" href="#{$id}">&#128279;</a> <!-- 链接符号 -->

            <span class="action"><xsl:value-of select="@method" /></span>
            <span>
                <xsl:call-template name="deprecated">
                    <xsl:with-param name="deprecated" select="@deprecated" />
                </xsl:call-template>

                <xsl:value-of select="path/@path" />
            </span>

            <span class="summary"><xsl:value-of select="@summary" /></span>
        </summary>

        <xsl:if test="description">
            <div class="description" data-type="{description/@type}">
                <pre><xsl:copy-of select="description/node()" /></pre>
            </div>
        </xsl:if>

        <div class="body">
            <div class="requests">
                <h4 class="title"><xsl:copy-of select="$locale-request" /></h4>
                <xsl:call-template name="requests">
                    <xsl:with-param name="requests" select="request" />
                    <xsl:with-param name="path" select="path" />
                    <xsl:with-param name="headers" select="header" />
                </xsl:call-template>
            </div>
            <div class="responses">
                <h4 class="title"><xsl:copy-of select="$locale-response" /></h4>

                <xsl:for-each select="response">
                    <xsl:call-template name="response">
                        <xsl:with-param name="response" select="." />
                    </xsl:call-template>
                </xsl:for-each>

                <xsl:for-each select="/apidoc/response"><!-- 公有的 response -->
                    <xsl:call-template name="response">
                        <xsl:with-param name="response" select="." />
                    </xsl:call-template>
                </xsl:for-each>
            </div>
        </div>

        <xsl:if test="./callback"><xsl:apply-templates select="./callback" /></xsl:if>
    </details>
</xsl:template>

<!-- 回调内容 -->
<xsl:template match="/apidoc/api/callback">
    <div class="callback" data-method="{./@method},">
        <h3>
            <xsl:copy-of select="$locale-callback" />
            <span class="summary"><xsl:value-of select="@summary" /></span>
        </h3>
        <xsl:if test="description">
            <div class="description" data-type="{description/@type}">
                <pre><xsl:copy-of select="description/node()" /></pre>
            </div>
        </xsl:if>

        <div class="body">
            <div class="requests">
                <h4 class="title"><xsl:copy-of select="$locale-request" /></h4>
                <xsl:call-template name="requests">
                    <xsl:with-param name="requests" select="request" />
                    <xsl:with-param name="path" select="path" />
                    <xsl:with-param name="headers" select="header" />
                </xsl:call-template>
            </div>

            <xsl:if test="response">
                <div class="responses">
                    <h4 class="title"><xsl:copy-of select="$locale-response" /></h4>
                    <xsl:for-each select="response">
                        <xsl:call-template name="response">
                            <xsl:with-param name="response" select="." />
                        </xsl:call-template>
                    </xsl:for-each>
                </div>
            </xsl:if>
        </div> <!-- end .body -->
    </div> <!-- end .callback -->
</xsl:template>

<!-- api/request 的界面元素 -->
<xsl:template name="requests">
<xsl:param name="requests" />
<xsl:param name="path" />
<xsl:param name="headers" /> <!-- 公用的报头 -->

<div class="request">
    <xsl:if test="$path/param">
        <xsl:call-template name="param">
            <xsl:with-param name="title">
                <xsl:copy-of select="$locale-path-param" />
            </xsl:with-param>
            <xsl:with-param name="param" select="$path/param" />
            <xsl:with-param name="simple" select="'true'" />
        </xsl:call-template>
    </xsl:if>

    <xsl:if test="$path/query">
        <xsl:call-template name="param">
            <xsl:with-param name="title">
                <xsl:copy-of select="$locale-query" />
            </xsl:with-param>
            <xsl:with-param name="param" select="$path/query" />
            <xsl:with-param name="simple" select="'true'" />
        </xsl:call-template>
    </xsl:if>

    <xsl:if test="$headers">
        <xsl:call-template name="param">
            <xsl:with-param name="title">
                <xsl:copy-of select="$locale-header" />
            </xsl:with-param>
            <xsl:with-param name="param" select="$headers" />
            <xsl:with-param name="simple" select="'true'" />
        </xsl:call-template>
    </xsl:if>

    <xsl:for-each select="$requests">
        <h5 class="status"><xsl:value-of select="@mimetype" /></h5>

        <xsl:if test="header">
            <xsl:call-template name="param">
                <xsl:with-param name="title">
                    <xsl:copy-of select="$locale-header" />
                </xsl:with-param>
                <xsl:with-param name="param" select="header" />
                <xsl:with-param name="simple" select="'true'" />
            </xsl:call-template>
        </xsl:if>

        <xsl:call-template name="param">
            <xsl:with-param name="title"><xsl:copy-of select="$locale-body" /></xsl:with-param>
            <xsl:with-param name="param" select="." />
        </xsl:call-template>
    </xsl:for-each>
</div>
</xsl:template>

<!-- api/response 的界面 -->
<xsl:template name="response">
    <xsl:param name="response" />

    <h5 class="status">
        <xsl:value-of select="$response/@status" />
        <span class="mimetype">&#160;(<xsl:value-of select="$response/@mimetype" />)</span>
    </h5>

    <xsl:if test="$response/header">
        <xsl:call-template name="param">
            <xsl:with-param name="title">
                <xsl:copy-of select="$locale-header" />
            </xsl:with-param>
            <xsl:with-param name="param" select="$response/header" />
            <xsl:with-param name="simple" select="'true'" />
        </xsl:call-template>
    </xsl:if>

    <xsl:call-template name="param">
        <xsl:with-param name="title"><xsl:copy-of select="$locale-body" /></xsl:with-param>
        <xsl:with-param name="param" select="$response" />
    </xsl:call-template>
</xsl:template>

<!-- path param, path query, header 等的界面 -->
<xsl:template name="param">
    <xsl:param name="title" />
    <xsl:param name="param" />
    <xsl:param name="simple" select="'false'" /> <!-- 简单的类型，不存在嵌套类型，也不会有示例代码 -->

    <xsl:if test="not($param/@type='none')">
        <div class="param">
            <h4 class="title">
                &#x27a4;&#160;<xsl:copy-of select="$title" />
                <xsl:if test="$param/example">
                    &#160;(<a class="toggle-example"><xsl:copy-of select="$locale-example" /></a>)
                </xsl:if>
            </h4>

            <table class="param-list" data-visible="true">
                <thead>
                    <tr>
                        <th><xsl:copy-of select="$locale-var" /></th>
                        <th><xsl:copy-of select="$locale-type" /></th>
                        <th><xsl:copy-of select="$locale-value" /></th>
                        <th><xsl:copy-of select="$locale-description" /></th>
                    </tr>
                </thead>
                <tbody>
                    <xsl:choose>
                        <xsl:when test="$simple='true'">
                            <xsl:call-template name="simple-param-list">
                                <xsl:with-param name="param" select="$param" />
                            </xsl:call-template>
                        </xsl:when>
                        <xsl:otherwise>
                            <xsl:call-template name="param-list">
                                <xsl:with-param name="param" select="$param" />
                            </xsl:call-template>
                        </xsl:otherwise>
                    </xsl:choose>
                </tbody>
            </table>

            <xsl:if test="$param/example">
            <pre class="example" data-visible="false" data-mimetype="{$param/example/@mimetype}"><xsl:copy-of select="$param/example/node()" /></pre>
            </xsl:if>
        </div>
    </xsl:if>
</xsl:template>

<!-- 列顺序必须要与 param 中的相同 -->
<xsl:template name="simple-param-list">
    <xsl:param name="param" />

    <xsl:for-each select="$param">
        <xsl:call-template name="param-list-tr">
            <xsl:with-param name="param" select="." />
        </xsl:call-template>
    </xsl:for-each>
</xsl:template>

<!-- 列顺序必须要与 param 中的相同 -->
<xsl:template name="param-list">
    <xsl:param name="param" />
    <xsl:param name="parent" select="''" /> <!-- 上一级的名称，嵌套对象时可用 -->

    <xsl:for-each select="$param">
        <xsl:call-template name="param-list-tr">
            <xsl:with-param name="param" select="." />
            <xsl:with-param name="parent" select="$parent" />
        </xsl:call-template>

        <xsl:if test="param">
            <xsl:variable name="p">
                    <xsl:value-of select="concat($parent, @name)" />
                    <xsl:if test="@name"><xsl:value-of select="'.'" /></xsl:if>
            </xsl:variable>

            <xsl:call-template name="param-list">
                <xsl:with-param name="param" select="param" />
                <xsl:with-param name="parent" select="$p" />
            </xsl:call-template>
        </xsl:if>
    </xsl:for-each>
</xsl:template>

<!-- 显示第一行的参数数据 -->
<xsl:template name="param-list-tr">
    <xsl:param name="param" />
    <xsl:param name="parent" select="''" />

    <tr>
        <xsl:call-template name="deprecated">
            <xsl:with-param name="deprecated" select="$param/@deprecated" />
        </xsl:call-template>
        <th>
            <span class="parent-type"><xsl:value-of select="$parent" /></span>
            <xsl:value-of select="$param/@name" />
        </th>

        <td>
            <xsl:value-of select="$param/@type" />
            <xsl:if test="$param/@array='true'"><xsl:value-of select="'[]'" /></xsl:if>
        </td>

        <td>
            <xsl:choose>
                <xsl:when test="$param/@optional='true'"><xsl:value-of select="'O'" /></xsl:when>
                <xsl:otherwise><xsl:value-of select="'R'" /></xsl:otherwise>
            </xsl:choose>
            <xsl:value-of select="concat(' ', $param/@default)" />
        </td>

        <td>
            <xsl:choose>
                <xsl:when test="description">
                    <xsl:attribute name="data-type">
                        <xsl:value-of select="description/@type" />
                    </xsl:attribute>
                    <pre><xsl:copy-of select="description/node()" /></pre>
                </xsl:when>
                <xsl:otherwise><xsl:value-of select="@summary" /></xsl:otherwise>
            </xsl:choose>
            <xsl:call-template name="enum">
                <xsl:with-param name="enum" select="$param/enum"/>
            </xsl:call-template>
        </td>
    </tr>
</xsl:template>

<!-- 显示枚举类型的内容 -->
<xsl:template name="enum">
    <xsl:param name="enum" />

    <xsl:if test="$enum">
        <p><xsl:copy-of select="$locale-enum" /></p>
        <ul>
        <xsl:for-each select="$enum">
            <li>
            <xsl:call-template name="deprecated">
                <xsl:with-param name="deprecated" select="@deprecated" />
            </xsl:call-template>

            <xsl:value-of select="@value" />:
            <xsl:choose>
                <xsl:when test="description">
                    <div data-type="{description/@type}">
                        <pre><xsl:copy-of select="description/node()" /></pre>
                    </div>
                </xsl:when>
                <xsl:otherwise><xsl:value-of select="summary" /></xsl:otherwise>
            </xsl:choose>
            </li>
        </xsl:for-each>
        </ul>
    </xsl:if>
</xsl:template>

<!--
给指定的元素添加已弃用的标记

该模板会给父元素添加 class 和 title 属性，
所以必须要在父元素的任何子元素之前，否则 chrome 和 safari 可能无法正常解析。
-->
<xsl:template name="deprecated">
    <xsl:param name="deprecated" />

    <xsl:if test="$deprecated">
        <xsl:attribute name="class"><xsl:value-of select="'del'" /></xsl:attribute>
        <xsl:attribute name="title">
            <xsl:value-of select="$deprecated" />
        </xsl:attribute>
    </xsl:if>
</xsl:template>

<!-- 用于将 API 地址转换成合法的 ID 标记 -->
<xsl:variable name="id-from" select="'{}/'" />
<xsl:variable name="id-to" select="'__-'" />

<!-- 根据情况获取相应的图标 -->
<xsl:variable name="icon">
    <xsl:choose>
        <xsl:when test="/apidoc/@logo">
            <xsl:value-of select="/apidoc/@logo" />
        </xsl:when>
        <xsl:otherwise>
            <xsl:value-of select="concat($base-url, '../icon.svg')" />
        </xsl:otherwise>
    </xsl:choose>
</xsl:variable>

<!--
获取相对于当前 xsl 文件的基地址
xsl 2.0 可以直接采用 base-uri(document(''))
-->
<xsl:variable name="base-url">
    <xsl:apply-templates select="processing-instruction('xml-stylesheet')" />
</xsl:variable>

<xsl:template match="processing-instruction('xml-stylesheet')">
    <xsl:variable name="v1" select="substring-after(., 'href=&quot;')" />
    <!-- NOTE: 此处假定当前文件叫作 apidoc.xsl，如果不是的话，需要另外处理此代码 -->
    <xsl:variable name="v2" select="substring-before($v1, 'apidoc.xsl&quot;')" />
    <xsl:value-of select="$v2" />
</xsl:template>

</xsl:stylesheet>
`),
	},
	{
		Name:        "v5/locales.xsl",
		ContentType: "text/xsl; charset=utf-8",
		Content: []byte(`<?xml version="1.0" encoding="UTF-8"?>

<!-- 当前文件实现了简单的翻译功能 -->

<xsl:stylesheet
version="1.0"
xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
xmlns:l="urn:locale"
exclude-result-prefixes="l">

<!-- 当前支持的本地化列表，其中第一个会被当作默认值。 -->
<l:locales>
    <locale id="zh-hans">简体中文</locale>
    <locale id="zh-hant">繁體中文</locale>
</l:locales>

<xsl:template name="languages">
    <xsl:for-each select="document('')/xsl:stylesheet/l:locales/locale">
    <li lang="{@id}" role="menuitemradio">
        <label><input type="radio" name="lang" checked="{$curr-lang=@id}" />&#160;<xsl:value-of select="." /></label>
    </li>
    </xsl:for-each>
</xsl:template>

<!-- language -->
<xsl:variable name="locale-language">
    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hans'" />
        <xsl:with-param name="text" select="document('')/xsl:stylesheet/l:locales/locale[@id='zh-hans']" />
    </xsl:call-template>

    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hant'" />
        <xsl:with-param name="text" select="document('')/xsl:stylesheet/l:locales/locale[@id='zh-hant']" />
    </xsl:call-template>
</xsl:variable>

<!-- server -->
<xsl:variable name="locale-server">
    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hans'" />
        <xsl:with-param name="text" select="'服务'" />
    </xsl:call-template>

    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hant'" />
        <xsl:with-param name="text" select="'服務'" />
    </xsl:call-template>
</xsl:variable>

<!-- tag -->
<xsl:variable name="locale-tag">
    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hans'" />
        <xsl:with-param name="text" select="'标签'" />
    </xsl:call-template>

    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hant'" />
        <xsl:with-param name="text" select="'標簽'" />
    </xsl:call-template>
</xsl:variable>

<!-- expand -->
<xsl:variable name="locale-expand">
    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hans'" />
        <xsl:with-param name="text" select="'展开'" />
    </xsl:call-template>

    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hant'" />
        <xsl:with-param name="text" select="'展開'" />
    </xsl:call-template>
</xsl:variable>

<!-- method -->
<xsl:variable name="locale-method">
    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hans'" />
        <xsl:with-param name="text" select="'请求方法'" />
    </xsl:call-template>

    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hant'" />
        <xsl:with-param name="text" select="'請求方法'" />
    </xsl:call-template>
</xsl:variable>

<!-- request -->
<xsl:variable name="locale-request">
    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hans'" />
        <xsl:with-param name="text" select="'请求'" />
    </xsl:call-template>

    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hant'" />
        <xsl:with-param name="text" select="'請求'" />
    </xsl:call-template>
</xsl:variable>

<!-- response -->
<xsl:variable name="locale-response">
    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hans'" />
        <xsl:with-param name="text" select="'返回'" />
    </xsl:call-template>

    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hant'" />
        <xsl:with-param name="text" select="'返回'" />
    </xsl:call-template>
</xsl:variable>

<!-- callback -->
<xsl:variable name="locale-callback">
    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hans'" />
        <xsl:with-param name="text" select="'回调'" />
    </xsl:call-template>

    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hant'" />
        <xsl:with-param name="text" select="'回調'" />
    </xsl:call-template>
</xsl:variable>

<!-- path param -->
<xsl:variable name="locale-path-param">
    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hans'" />
        <xsl:with-param name="text" select="'路径参数'" />
    </xsl:call-template>

    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hant'" />
        <xsl:with-param name="text" select="'路徑參數'" />
    </xsl:call-template>
</xsl:variable>

<!-- query -->
<xsl:variable name="locale-query">
    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hans'" />
        <xsl:with-param name="text" select="'查询参数'" />
    </xsl:call-template>

    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hant'" />
        <xsl:with-param name="text" select="'查詢參數'" />
    </xsl:call-template>
</xsl:variable>

<!-- header -->
<xsl:variable name="locale-header">
    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hans'" />
        <xsl:with-param name="text" select="'报头'" />
    </xsl:call-template>

    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hant'" />
        <xsl:with-param name="text" select="'報頭'" />
    </xsl:call-template>
</xsl:variable>

<!-- body -->
<xsl:variable name="locale-body">
    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hans'" />
        <xsl:with-param name="text" select="'报文'" />
    </xsl:call-template>

    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hant'" />
        <xsl:with-param name="text" select="'報文'" />
    </xsl:call-template>
</xsl:variable>

<!-- example -->
<xsl:variable name="locale-example">
    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hans'" />
        <xsl:with-param name="text" select="'示例代码'" />
    </xsl:call-template>

    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hant'" />
        <xsl:with-param name="text" select="'示例代碼'" />
    </xsl:call-template>
</xsl:variable>

<!-- var -->
<xsl:variable name="locale-var">
    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hans'" />
        <xsl:with-param name="text" select="'变量'" />
    </xsl:call-template>

    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hant'" />
        <xsl:with-param name="text" select="'變量'" />
    </xsl:call-template>
</xsl:variable>

<!-- type -->
<xsl:variable name="locale-type">
    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hans'" />
        <xsl:with-param name="text" select="'类型'" />
    </xsl:call-template>

    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hant'" />
        <xsl:with-param name="text" select="'類型'" />
    </xsl:call-template>
</xsl:variable>

<!-- value -->
<xsl:variable name="locale-value">
    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hans'" />
        <xsl:with-param name="text" select="'值'" />
    </xsl:call-template>

    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hant'" />
        <xsl:with-param name="text" select="'值'" />
    </xsl:call-template>
</xsl:variable>

<!-- description -->
<xsl:variable name="locale-description">
    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hans'" />
        <xsl:with-param name="text" select="'描述'" />
    </xsl:call-template>

    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hant'" />
        <xsl:with-param name="text" select="'描述'" />
    </xsl:call-template>
</xsl:variable>

<!-- enum -->
<xsl:variable name="locale-enum">
    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hans'" />
        <xsl:with-param name="text" select="'枚举'" />
    </xsl:call-template>

    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hant'" />
        <xsl:with-param name="text" select="'枚舉'" />
    </xsl:call-template>
</xsl:variable>

<!-- license -->
<xsl:variable name="locale-license">
    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hans'" />
        <xsl:with-param name="text">
            文档版权为 <a href="{apidoc/license/@url}"><xsl:value-of select="apidoc/license/@text" /></a>
        </xsl:with-param>
    </xsl:call-template>

    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hant'" />
        <xsl:with-param name="text">
            文檔版權為 <a href="{apidoc/license/@url}"><xsl:value-of select="apidoc/license/@text" /></a>
        </xsl:with-param>
    </xsl:call-template>
</xsl:variable>

<!-- generator -->
<xsl:variable name="locale-generator">
    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hans'" />
        <xsl:with-param name="text">
            由 <a href="https://apidoc.tools">apidoc</a> 生成于 <time><xsl:value-of select="apidoc/@created" /></time>
        </xsl:with-param>
    </xsl:call-template>

    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hant'" />
        <xsl:with-param name="text">
            由 <a href="https://apidoc.tools">apidoc</a> 生成於 <time><xsl:value-of select="apidoc/@created" /></time>
        </xsl:with-param>
    </xsl:call-template>
</xsl:variable>

<xsl:template name="build-locale">
    <xsl:param name="lang" />
    <xsl:param name="text" />

    <xsl:variable name="class">
        <xsl:choose>
            <xsl:when test="$curr-lang=translate($lang, $uppercase, $lowercase)">
                <xsl:value-of select="''" />
            </xsl:when>
            <xsl:otherwise>
                <xsl:value-of select="'hidden'" />
            </xsl:otherwise>
        </xsl:choose>
    </xsl:variable>

    <!-- data-locale 属性表示该元素是一个本地化信息元素，JS 代码通过该标记切换语言。 -->
    <span data-locale="true" lang="{$lang}" class="{$class}"><xsl:copy-of select="$text" /></span>
</xsl:template>

<!--
返回当前文档的语言，会转换为小写，_ 也会被转换成 -
如果文档指定的语言不存在，则会采取 l:locales 中的第一个元素作为默认语言。
-->
<xsl:variable name="curr-lang">
    <xsl:variable name="curr" select="translate(/apidoc/@lang, $uppercase, $lowercase)" />

    <xsl:variable name="r1">
        <xsl:for-each select="document('')/xsl:stylesheet/l:locales/locale">
            <xsl:if test="@id=$curr"><xsl:value-of select="$curr" /></xsl:if>
        </xsl:for-each>
    </xsl:variable>

    <xsl:variable name="r2">
    <xsl:choose>
        <xsl:when test="$r1 and not($r1='')"> 
            <xsl:value-of select="$r1" />
        </xsl:when>
        <xsl:otherwise>
            <xsl:value-of select="document('')/xsl:stylesheet/l:locales/locale[1]/@id" />
        </xsl:otherwise>
    </xsl:choose>
    </xsl:variable>

    <xsl:value-of select="$r2" />
</xsl:variable>

<!-- 用于实现 lower-case 和 upper-case，如果将来某天浏览器支持 xsl 2.0 了，可以直接采用相关函数 -->
<xsl:variable name="lowercase" select="'abcdefghijklmnopqrstuvwxyz-'" />
<xsl:variable name="uppercase" select="'ABCDEFGHIJKLMNOPQRSTUVWXYZ_'" />

</xsl:stylesheet>
`),
	},
	{
		Name:        "v5/view.html",
		ContentType: "text/html; charset=utf-8",
		Content: []byte(`<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>Document</title>

    <style>
    :root {
        --background: white;
    }
    @media (prefers-color-scheme: dark) {
        :root {
            --background: black;
        }
    }

    html,
    body {
        margin: 0;
        padding: 0;
        background: var(--background);
    }

    iframe {
        border: none;
        width: 100%;
        background: var(--background);
    }
    </style>
</head>

<body>
    <!--
    NOTE: xslt 中引用的 js，只有在 iframe 中才能执行，
    直接用 document.replace() 替换当前内容，不会执行引用的 JS 代码。
    -->
    <iframe id="apidoc"></iframe>
    <script>
    async function loadXML(path) {
        const obj = await fetch(path);
        return (new DOMParser()).parseFromString(await obj.text(), "text/xml");
    }

    function changeFrameHeight() {
        const iframe = document.getElementById("apidoc");
        iframe.height = document.documentElement.clientHeight;
    }

    window.onresize = function () {
        changeFrameHeight();
    } 

    const queries = new URLSearchParams(window.location.search);
    const url = queries.get('url');

    async function init() {
        const processor = new XSLTProcessor();
        processor.importStylesheet(await loadXML('./apidoc.xsl'));

        const xml = await loadXML(url);
        const doc = processor.transformToDocument(xml);
        const html = (new XMLSerializer()).serializeToString(doc);
        const blob = new Blob([html], { type: 'text/html' })
        const obj = URL.createObjectURL(blob)

        const iframe = document.getElementById('apidoc');
        iframe.src = obj;
        iframe.addEventListener('load', (e) => {
            changeFrameHeight();
        });
    }

    try{
        init()
    } catch (e) {
        console.error(e)
    }
    </script>
</body>

</html>
`),
	},
	{
		Name:        "v6/apidoc.css",
		ContentType: "text/css; charset=utf-8",
		Content: []byte(`@charset "utf-8";

:root {
    --max-width: 2048px;
    --min-width: 200px; /* 当列的宽度小于此值时，部分行内容会被从横向改变纵向排列。 */
    --padding: 1rem;
    --article-padding: calc(var(--padding) / 2);

    --color: black;
    --accent-color: #0074d9;
    --background: white;
    --border-color: #e0e0e0;
    --delete-color: red;

    /* method */
    --method-get-color: green;
    --method-options-color: green;
    --method-post-color: darkorange;
    --method-put-color: darkorange;
    --method-patch-color: darkorange;
    --method-delete-color: red;
}

@media (prefers-color-scheme: dark) {
    :root {
        --color: #b0b0b0;
        --accent-color: #0074d9;
        --background: black;
        --border-color: #303030;
        --delete-color: red;

        /* method */
        --method-get-color: green;
        --method-options-color: green;
        --method-post-color: darkorange;
        --method-put-color: darkorange;
        --method-patch-color: darkorange;
        --method-delete-color: red;
    }
}

body {
    padding: 0;
    margin: 0;
    color: var(--color);
    background: var(--background);
    text-align: center;
}

table {
    width: 100%;
}

table th, table td {
    font-weight: normal;
    text-align: left;
    border-bottom: 1px solid transparent;
}

table tr:hover th,
table tr:hover td {
    border-bottom: 1px solid var(--border-color);
}

ul, ol, ul li, ol li {
    padding: 0;
    margin: 0;
    list-style-position: inside;
}

p {
    margin: 0;
}

summary, input {
    outline: none;
}

a {
    text-decoration: none;
    color: var(--accent-color);
}

a:hover {
    opacity: .7;
}

.del {
    text-decoration: line-through;
    text-decoration-color: var(--delete-color);
}

.hidden {
    display: none;
}

/*************************** header ***********************/

header {
    position: sticky;
    top: 0;
    display: block;
    z-index: 1000;
    box-sizing: border-box;
    background: var(--background);
    box-shadow: 2px 2px 2px var(--border-color);
}

header .wrap {
    margin: 0 auto;
    text-align: left;
    max-width: var(--max-width);
    padding: 0 var(--padding);
    display: flex;
    flex-wrap: wrap;
    justify-content: space-between;
    align-items: center;
}

header h1 {
    margin: var(--padding) 0 var(--article-padding);
    display: inline-block;
}

header h1 .version {
    font-size: 1rem;
}

header h1 img {
    height: 1.5rem;
    margin-right: .5rem;
}


header .menus {
    display: flex;
}

header .menu {
    cursor: pointer;
    position: relative;
    margin-left: var(--padding);
    display: none; /* 默认不可见，只有在有 JS 的情况下，通过 JS 控制其可见性 */
}

header .menu:hover span {
    opacity: .7;
}

header .menu ul {
    position: absolute;
    min-width: 4rem;
    right: 0;
    display: none;
    list-style: none;
    background: var(--background);
    border: 1px solid var(--border-color);
    padding: var(--article-padding);
}

header .menu ul li {
    padding-bottom: var(--article-padding);
}

/* 可以保证 label 的内容在同一行 */
header .menu ul li label {
    display: inline-flex;
    align-items: baseline;
    word-break: keep-all;
    white-space: nowrap;
}

header .menu ul li:last-of-type {
    padding-bottom: 0;
}

header .menu:hover ul {
    display: block;
}

/*************************** main ***********************/

main {
    padding: 0 var(--padding);
    margin: 0 auto;
    max-width: var(--max-width);
    text-align: left;
}

main .content {
    margin: var(--padding) 0;
}

/****************** .servers *******************/

main .servers {
    display: flex;
    flex-flow: wrap;
    justify-content: space-between;
    margin-bottom: var(--padding);
}

main .servers .server {
    flex-grow: 1;
    min-width: var(--min-width);
    box-sizing: border-box;
    border: 1px solid var(--border-color);
    padding: var(--padding) var(--article-padding);
}

main .servers .server:hover {
    border: 1px solid var(--accent-color);
}

main .servers .server h4 {
    margin: 0 0 var(--padding);
}

/********************** .api **********************/

main .api {
    margin-bottom: var(--article-padding);
    border: 1px solid var(--border-color);
}

main .api>summary {
    margin: 0;
    padding: var(--article-padding);
    border-bottom: 1px solid var(--border-color);
    cursor: pointer;
    line-height: 1;
}


main details.api:not([open])>summary {
    border: none;
}

main .api>summary .action {
    min-width: 4rem;
    font-weight: bold;
    display: inline-block;
    margin-right: 1rem;
}

main .api>summary .link {
    margin-right: 10px;
    text-decoration: none;
}

main .api .callback .summary,
main .api>summary .summary {
    float: right;
    font-weight: 400;
    opacity: .5;
}

main .api .description {
    padding: var(--article-padding);
    margin: 0;
    border-bottom: 1px solid var(--border-color);
}

main .api[data-method='GET,'][open],
main .api[data-method='GET,']:hover,
main .callback[data-method='GET,'] h3 {
    border: 1px solid var(--method-get-color);
}
main .api[data-method='GET,']>summary {
    border-bottom: 1px solid var(--method-get-color);
}

main .api[data-method='POST,'][open],
main .api[data-method='POST,']:hover,
main .callback[data-method='POST,'] h3 {
    border: 1px solid var(--method-post-color);
}
main .api[data-method='POST,']>summary {
    border-bottom: 1px solid var(--method-post-color);
}

main .api[data-method='PUT,'][open],
main .api[data-method='PUT,']:hover,
main .callback[data-method='PUT,'] h3 {
    border: 1px solid var(--method-put-color);
}
main .api[data-method='PUT,']>summary {
    border-bottom: 1px solid var(--method-put-color);
}

main .api[data-method='PATCH,'][open],
main .api[data-method='PATCH,']:hover,
main .callback[data-method='PATCH,'] h3 {
    border: 1px solid var(--method-patch-color);
}
main .api[data-method='PATCH,']>summary {
    border-bottom: 1px solid var(--method-patch-color);
}

main .api[data-method='DELETE,'][open],
main .api[data-method='DELETE,']:hover,
main .callback[data-method='DELETE,'] h3 {
    border: 1px solid var(--method-delete-color);
}
main .api[data-method='DELETE,']>summary {
    border-bottom: 1px solid var(--method-delete-color);
}

main .api[data-method='OPTIONS,'][open],
main .api[data-method='OPTIONS,']:hover,
main .callback[data-method='OPTIONS,'] h3 {
    border: 1px solid var(--method-options-color);
}
main .api[data-method='OPTIONS,']>summary {
    border-bottom: 1px solid var(--method-options-color);
}

main .callback h3 {
    padding: var(--article-padding) var(--padding);
    margin: 0;
    border-left: none !important;
    border-right: none !important;
    cursor: pointer;
    line-height: 1;
}

main .api .body {
    display: flex;
    flex-flow: wrap;
}

main .api .title,
main .api .title {
    margin: var(--article-padding) 0 0;
    font-weight: normal;
}

main .api .body .responses .status {
    margin: var(--padding) 0 var(--article-padding);
    border-bottom: 1px solid var(--border-color);
    padding-bottom: 3px;
}

main .api .param .parent-type {
    opacity: .5;
}

main .api .body .requests,
main .api .body .responses {
    flex: 1 1 50%;
    box-sizing: border-box;
    min-width: var(--min-width);
}
main .api .body .requests {
    border-right: 1px dotted var(--border-color);
}

main .api .body .requests .header,
main .api .body .responses .header {
    margin: 0;
    opacity: .5;
    padding: var(--article-padding);
}

main .api .body .requests details:not([open]),
main .api .body .responses details:not([open]) {
    margin-bottom: var(--article-padding);
}

main .api .body .requests>.param,
main .api .body .requests details summary ~ *,
main .api .body .responses details summary ~ * {
    padding: 0 var(--article-padding);
}

main .api .responses details summary,
main .api .requests details summary {
    cursor: pointer;
    padding: var(--article-padding);
}

main .api .responses details[open] summary,
main .api .requests details[open] summary,
main .api .responses details summary:hover,
main .api .requests details summary:hover {
    cursor: pointer;
    background-color: var(--border-color);
    padding: var(--article-padding);
}

/*************************** footer ***********************/

footer {
    border-top: 1px solid var(--border-color);
    padding: var(--padding) var(--padding);
    text-align: left;
    margin-top: var(--padding);
}

footer .wrap {
    margin: 0 auto;
    max-width: var(--max-width);
    display: flex;
    flex-flow: wrap;
    justify-content: space-between;
}
`),
	},
	{
		Name:        "v6/apidoc.js",
		ContentType: "application/javascript; charset=utf-8",
		Content: []byte(`'use strict';

window.onload = function () {
    registerFilter('method');
    registerFilter('server');
    registerFilter('tag');
    registerExpand();
    registerLanguageFilter();
    prettyDescription();
};

function registerFilter(type) {
    const menu = document.querySelector('.' + type + '-selector');
    if (menu === null) { // 可能为空，表示不存在该过滤项
        return;
    }

    menu.style.display = 'block'; // 有 JS 的情况下，展示过滤菜单

    menu.querySelectorAll('li input').forEach((val) => {
        val.addEventListener('change', (event) => {
            if (event.target === null) {
                return;
            }

            const chk = event.target.checked;
            const tag = event.target.parentNode.parentNode.getAttribute('data-' + type);
            document.querySelectorAll('.api').forEach((api) => {
                if (!api.getAttribute('data-' + type).includes(tag + ',')) {
                    return;
                }

                api.setAttribute("data-hidden-" + type, chk ? "" : "true");

                const hidden = api.getAttribute('data-hidden-tag') === 'true' ||
                    api.getAttribute('data-hidden-server') === 'true' ||
                    api.getAttribute('data-hidden-method') === 'true';
                api.style.display = hidden ? 'none' : 'block';
            });
        }); // end addEventListener('change')
    }); // end forEach('li input')
}

function registerLanguageFilter() {
    const menu = document.querySelector('.languages-selector');

    menu.style.display = 'block';

    menu.querySelectorAll('li input').forEach((val) => {
        val.addEventListener('change', (event) => {
            if (event.target === null || !event.target.checked) {
                return;
            }

            const lang = event.target.parentNode.parentNode.getAttribute('lang');
            document.querySelectorAll('[data-locale]').forEach((elem) => {
                elem.className = elem.getAttribute('lang') === lang ? '' : 'hidden';
            });
        }); // end addEventListener('change')
    }); // end forEach('li input')
}

function registerExpand() {
    const expand = document.querySelector('.expand-selector');
    if (expand === null) {
        return;
    }

    expand.style.display = 'block';

    expand.querySelector('input').addEventListener('change', (event) => {
        const chk = event.target.checked;
        document.querySelectorAll('details.api').forEach((elem) => {
            elem.open = chk;
        });
    });
}

// 美化描述内容
//
// 即将 html 内容转换成真的 HTML 格式，而 markdown 则依然是 pre 显示。
function prettyDescription() {
    document.querySelectorAll('[data-type]').forEach((elem) => {
        const type = elem.getAttribute('data-type');
        if (type !== 'html') {
            return;
        }

        elem.innerHTML = elem.getElementsByTagName('pre')[0].innerText;
    });
}
`),
	},
	{
		Name:        "v6/apidoc.xsl",
		ContentType: "text/xsl; charset=utf-8",
		Content: []byte(`<?xml version="1.0" encoding="UTF-8"?>

<xsl:stylesheet version="1.0" xmlns:xsl="http://www.w3.org/1999/XSL/Transform">
<xsl:import href="./locales.xsl" />
<xsl:output
    method="html"
    encoding="utf-8"
    indent="yes"
    version="5.0"
    doctype-system="about:legacy-compat"
/>

<xsl:template match="/">
<html lang="{$curr-lang}">
    <head>
        <title><xsl:value-of select="apidoc/title" /></title>
        <meta charset="UTF-8" />
        <meta name="viewport" content="width=device-width,initial-scale=1,maximum-scale=1" />
        <meta name="generator" content="apidoc" />
        <link rel="icon" type="image/svg+xml" href="{$icon}" />
        <link rel="mask-icon" type="image/svg+xml" href="{$icon}" color="black" />
        <xsl:if test="apidoc/license"><link rel="license" href="{apidoc/license/@url}" /></xsl:if>
        <link rel="stylesheet" type="text/css" href="{$base-url}apidoc.css" />
        <script src="{$base-url}apidoc.js" />
    </head>
    <body>
        <xsl:call-template name="header" />

        <main>
            <div class="content" data-type="{description/@type}">
                <pre><xsl:copy-of select="apidoc/description/node()" /></pre>
            </div>
            <div class="servers"><xsl:apply-templates select="apidoc/server" /></div>
            <xsl:apply-templates select="apidoc/api" />
        </main>

        <footer>
        <div class="wrap">
            <xsl:if test="apidoc/license"><xsl:copy-of select="$locale-license" /></xsl:if>
            <xsl:copy-of select="$locale-generator" />
        </div>
        </footer>
    </body>
</html>
</xsl:template>


<xsl:template match="/apidoc/server">
<div class="server">
    <h4>
        <xsl:call-template name="deprecated">
            <xsl:with-param name="deprecated" select="@deprecated" />
        </xsl:call-template>
        <xsl:value-of select="@name" />
    </h4>

    <p><xsl:value-of select="@url" /></p>
    <div>
        <xsl:choose>
            <xsl:when test="description">
                <xsl:attribute name="data-type">
                    <xsl:value-of select="description/@type" />
                </xsl:attribute>
                <pre><xsl:copy-of select="description/node()" /></pre>
            </xsl:when>
            <xsl:otherwise><xsl:value-of select="@summary" /></xsl:otherwise>
        </xsl:choose>
    </div>
</div>
</xsl:template>


<!-- header 界面元素 -->
<xsl:template name="header">
<header>
<div class="wrap">
    <h1>
        <img src="{$icon}" />
        <xsl:value-of select="/apidoc/title" />
        <span class="version">&#160;(<xsl:value-of select="/apidoc/@version" />)</span>
    </h1>

    <div class="menus">
        <!-- expand -->
        <label class="menu expand-selector" role="checkbox">
            <input type="checkbox" /><xsl:copy-of select="$locale-expand" />
        </label>

        <!-- server -->
        <xsl:if test="apidoc/server">
        <div class="menu server-selector" role="menu" aria-haspopup="true">
            <xsl:copy-of select="$locale-server" />
            <span aria-hiddren="true">&#160;&#x25bc;</span>
            <ul role="menu" aria-hiddren="true">
                <xsl:for-each select="apidoc/server">
                <li data-server="{@name}" role="menuitemcheckbox">
                    <label>
                        <input type="checkbox" checked="checked" />&#160;<xsl:value-of select="@name" />
                    </label>
                </li>
                </xsl:for-each>
            </ul>
        </div>
        </xsl:if>

        <!-- tag -->
        <xsl:if test="apidoc/tag">
        <div class="menu tag-selector" role="menu" aria-haspopup="true">
            <xsl:copy-of select="$locale-tag" />
            <span aria-hiddren="true">&#160;&#x25bc;</span>
            <ul role="menu" aria-hiddren="true">
                <xsl:for-each select="apidoc/tag">
                <li data-tag="{@name}" role="menuitemcheckbox">
                    <label>
                        <input type="checkbox" checked="checked" />&#160;<xsl:value-of select="@title" />
                    </label>
                </li>
                </xsl:for-each>
            </ul>
        </div>
        </xsl:if>

        <!-- method -->
        <div class="menu method-selector" role="menu" aria-haspopup="true">
            <xsl:copy-of select="$locale-method" />
            <span aria-hiddren="true">&#160;&#x25bc;</span>
            <ul role="menu" aria-hiddren="true">
                <!-- 浏览器好像都不支持 xpath 2.0，所以无法使用 distinct-values -->
                <!-- xsl:for-each select="distinct-values(/apidoc/api/@method)" -->
                <xsl:for-each select="/apidoc/api/@method[not(../preceding-sibling::api/@method = .)]">
                <li data-method="{.}" role="menuitemcheckbox">
                    <label><input type="checkbox" checked="checked" />&#160;<xsl:value-of select="." /></label>
                </li>
                </xsl:for-each>
            </ul>
        </div>

        <!-- language -->
        <div class="menu languages-selector" role="menu" aria-haspopup="true">
            <xsl:copy-of select="$locale-language" />
            <span aria-hiddren="true">&#160;&#x25bc;</span>
            <ul role="menu" aria-hiddren="true"><xsl:call-template name="languages" /></ul>
        </div>
    </div> <!-- end .menus -->
</div> <!-- end .wrap -->
</header>
</xsl:template>


<!-- api 界面元素 -->
<xsl:template match="/apidoc/api">
<xsl:variable name="id" select="concat(@method, translate(path/@path, $id-from, $id-to))" />

<details id="{$id}" class="api" data-method="{@method},">
<xsl:attribute name="data-tag">
    <xsl:for-each select="tag"><xsl:value-of select="concat(., ',')" /></xsl:for-each>
</xsl:attribute>
<xsl:attribute name="data-server">
    <xsl:for-each select="server"><xsl:value-of select="concat(., ',')" /></xsl:for-each>
</xsl:attribute>

    <summary>
        <a class="link" href="#{$id}">&#128279;</a> <!-- 链接符号 -->

        <span class="action"><xsl:value-of select="@method" /></span>
        <span>
            <xsl:call-template name="deprecated">
                <xsl:with-param name="deprecated" select="@deprecated" />
            </xsl:call-template>

            <xsl:value-of select="path/@path" />
        </span>

        <span class="summary"><xsl:value-of select="@summary" /></span>
    </summary>

    <xsl:if test="description">
        <div class="description" data-type="{description/@type}">
            <pre><xsl:copy-of select="description/node()" /></pre>
        </div>
    </xsl:if>

    <div class="body">
        <div class="requests">
            <h4 class="header"><xsl:copy-of select="$locale-request" /></h4>
            <xsl:call-template name="requests">
                <xsl:with-param name="requests" select="request" />
                <xsl:with-param name="path" select="path" />
                <xsl:with-param name="headers" select="header" />
            </xsl:call-template>
        </div>
        <div class="responses">
            <h4 class="header"><xsl:copy-of select="$locale-response" /></h4>
            <xsl:call-template name="responses">
                <xsl:with-param name="responses" select="response | /apidoc/response" />
            </xsl:call-template>
        </div>
    </div>

    <xsl:if test="./callback"><xsl:apply-templates select="./callback" /></xsl:if>
</details>
</xsl:template>


<!-- 回调内容 -->
<xsl:template match="/apidoc/api/callback">
<div class="callback" data-method="{./@method},">
    <h3>
        <xsl:copy-of select="$locale-callback" />
        <span class="summary"><xsl:value-of select="@summary" /></span>
    </h3>
    <xsl:if test="description">
        <div class="description" data-type="{description/@type}">
            <pre><xsl:copy-of select="description/node()" /></pre>
        </div>
    </xsl:if>

    <div class="body">
        <div class="requests">
            <h4 class="header"><xsl:copy-of select="$locale-request" /></h4>
            <xsl:call-template name="requests">
                <xsl:with-param name="requests" select="request" />
                <xsl:with-param name="path" select="path" />
                <xsl:with-param name="headers" select="header" />
            </xsl:call-template>
        </div>

        <xsl:if test="response">
            <div class="responses">
                <h4 class="header"><xsl:copy-of select="$locale-response" /></h4>
                <xsl:call-template name="responses">
                    <xsl:with-param name="responses" select="response" />
                </xsl:call-template>
            </div>
        </xsl:if>
    </div> <!-- end .body -->
</div> <!-- end .callback -->
</xsl:template>


<!-- api/request 的界面元素 -->
<xsl:template name="requests">
<xsl:param name="requests" />
<xsl:param name="path" />
<xsl:param name="headers" /> <!-- 公用的报头 -->
<xsl:if test="$path/param">
    <xsl:call-template name="param">
        <xsl:with-param name="title">
            <xsl:copy-of select="$locale-path-param" />
        </xsl:with-param>
        <xsl:with-param name="param" select="$path/param" />
        <xsl:with-param name="simple" select="'true'" />
    </xsl:call-template>
</xsl:if>

<xsl:if test="$path/query">
    <xsl:call-template name="param">
        <xsl:with-param name="title">
            <xsl:copy-of select="$locale-query" />
        </xsl:with-param>
        <xsl:with-param name="param" select="$path/query" />
        <xsl:with-param name="simple" select="'true'" />
    </xsl:call-template>
</xsl:if>

<xsl:if test="$headers">
    <xsl:call-template name="param">
        <xsl:with-param name="title">
            <xsl:copy-of select="$locale-header" />
        </xsl:with-param>
        <xsl:with-param name="param" select="$headers" />
        <xsl:with-param name="simple" select="'true'" />
    </xsl:call-template>
</xsl:if>

<xsl:variable name="request-any" select="$requests[not(@mimetype)]" />
<xsl:for-each select="/apidoc/mimetype | $requests/@mimetype[not(/apidoc/mimetype=.)]">
    <xsl:variable name="mimetype" select="." />
    <xsl:variable name="request" select="$requests[@mimetype=$mimetype]" />
    <xsl:if test="$request">
        <xsl:call-template name="request-body">
            <xsl:with-param name="mimetype" select="$mimetype" />
            <xsl:with-param name="request" select="$request" />
        </xsl:call-template>
    </xsl:if>
    <xsl:if test="not($request) and $request-any">
        <xsl:call-template name="request-body">
            <xsl:with-param name="mimetype" select="$mimetype" />
            <xsl:with-param name="request" select="$request-any" />
        </xsl:call-template>
    </xsl:if>
</xsl:for-each>
</xsl:template>


<xsl:template name="request-body">
<xsl:param name="mimetype" />
<xsl:param name="request" />
<details>
    <summary><xsl:value-of select="$mimetype" /></summary>
    <xsl:call-template name="top-param">
        <xsl:with-param name="mimetype" select="$mimetype" />
        <xsl:with-param name="param" select="$request" />
    </xsl:call-template>
</details>
</xsl:template>

<xsl:template name="responses">
<xsl:param name="responses" />
<xsl:for-each select="/apidoc/mimetype | $responses/@mimetype[not(/apidoc/mimetype=.)]">
    <xsl:variable name="mimetype" select="." />
    <xsl:if test="$responses[@mimetype=$mimetype] | $responses[not(@mimetype)]">
        <details>
        <summary><xsl:value-of select="$mimetype" /></summary>
        <xsl:for-each select="$responses">
            <xsl:variable name="resp" select="current()[@mimetype=$mimetype]" />
            <xsl:variable name="resp-any" select="current()[not(@mimetype)]" />
            <xsl:if test="$resp"><!-- 可能同时存在符合 resp 和 resp-any 的数据，优先取 resp -->
                <xsl:call-template name="response">
                    <xsl:with-param name="response" select="$resp" />
                    <xsl:with-param name="mimetype" select="$mimetype" />
                </xsl:call-template>
            </xsl:if>
            <xsl:if test="not($resp) and $resp-any">
                <xsl:call-template name="response">
                    <xsl:with-param name="response" select="$resp-any" />
                    <xsl:with-param name="mimetype" select="$mimetype" />
                </xsl:call-template>
            </xsl:if>
        </xsl:for-each>
        </details>
    </xsl:if>
</xsl:for-each>
</xsl:template>


<!-- api/response 的界面 -->
<xsl:template name="response">
<xsl:param name="response" />
<xsl:param name="mimetype" />

<h5 class="status"><xsl:value-of select="$response/@status" /></h5>
<xsl:call-template name="top-param">
    <xsl:with-param name="mimetype" select="$mimetype" />
    <xsl:with-param name="param" select="$response" />
</xsl:call-template>
</xsl:template>

<!-- request 和 response 参数的顶层元素调用模板 -->
<xsl:template name="top-param">
<xsl:param name="mimetype" />
<xsl:param name="param" />
<xsl:if test="$param/header">
    <xsl:call-template name="param">
        <xsl:with-param name="title">
            <xsl:copy-of select="$locale-header" />
        </xsl:with-param>
        <xsl:with-param name="param" select="$param/header" />
        <xsl:with-param name="simple" select="'true'" />
    </xsl:call-template>
</xsl:if>


<xsl:call-template name="param">
    <xsl:with-param name="title"><xsl:copy-of select="$locale-body" /></xsl:with-param>
    <xsl:with-param name="param" select="$param" />
</xsl:call-template>

<xsl:if test="$param/example">
    <h4 class="title">&#x27a4;&#160;<xsl:copy-of select="$locale-example" /></h4>
    <xsl:for-each select="$param/example">
        <xsl:if test="not(@mimetype) or @mimetype=$mimetype">
            <pre class="example"><xsl:copy-of select="node()" /></pre>
        </xsl:if>
    </xsl:for-each>
</xsl:if>
</xsl:template>


<!-- path param, path query, header 等的界面 -->
<xsl:template name="param">
<xsl:param name="title" />
<xsl:param name="param" />
<xsl:param name="simple" select="'false'" /> <!-- 简单的类型，不存在嵌套类型，也不会有示例代码 -->

<xsl:if test="$param/@type">
    <div class="param">
        <h4 class="title">&#x27a4;&#160;<xsl:copy-of select="$title" /></h4>

        <table class="param-list">
            <thead>
                <tr>
                    <th><xsl:copy-of select="$locale-var" /></th>
                    <th><xsl:copy-of select="$locale-type" /></th>
                    <th><xsl:copy-of select="$locale-value" /></th>
                    <th><xsl:copy-of select="$locale-description" /></th>
                </tr>
            </thead>
            <tbody>
                <xsl:choose>
                    <xsl:when test="$simple='true'">
                        <xsl:call-template name="simple-param-list">
                            <xsl:with-param name="param" select="$param" />
                        </xsl:call-template>
                    </xsl:when>
                    <xsl:otherwise>
                        <xsl:call-template name="param-list">
                            <xsl:with-param name="param" select="$param" />
                        </xsl:call-template>
                    </xsl:otherwise>
                </xsl:choose>
            </tbody>
        </table>
    </div>
</xsl:if>
</xsl:template>


<!-- 列顺序必须要与 param 中的相同 -->
<xsl:template name="simple-param-list">
<xsl:param name="param" />

<xsl:for-each select="$param">
    <xsl:call-template name="param-list-tr">
        <xsl:with-param name="param" select="." />
    </xsl:call-template>
</xsl:for-each>
</xsl:template>

<!-- 列顺序必须要与 param 中的相同 -->
<xsl:template name="param-list">
<xsl:param name="param" />
<xsl:param name="parent" select="''" /> <!-- 上一级的名称，嵌套对象时可用 -->

<xsl:for-each select="$param">
    <xsl:call-template name="param-list-tr">
        <xsl:with-param name="param" select="." />
        <xsl:with-param name="parent" select="$parent" />
    </xsl:call-template>

    <xsl:if test="param">
        <xsl:variable name="p">
                <xsl:value-of select="concat($parent, @name)" />
                <xsl:if test="@name"><xsl:value-of select="'.'" /></xsl:if>
        </xsl:variable>

        <xsl:call-template name="param-list">
            <xsl:with-param name="param" select="param" />
            <xsl:with-param name="parent" select="$p" />
        </xsl:call-template>
    </xsl:if>
</xsl:for-each>
</xsl:template>


<!-- 显示一行参数数据 -->
<xsl:template name="param-list-tr">
<xsl:param name="param" />
<xsl:param name="parent" select="''" />
<tr>
    <xsl:call-template name="deprecated">
        <xsl:with-param name="deprecated" select="$param/@deprecated" />
    </xsl:call-template>
    <th>
        <span class="parent-type"><xsl:value-of select="$parent" /></span>
        <xsl:value-of select="$param/@name" />
    </th>

    <td>
        <xsl:value-of select="$param/@type" />
        <xsl:if test="$param/@array='true'"><xsl:value-of select="'[]'" /></xsl:if>
    </td>

    <td>
        <xsl:choose>
            <xsl:when test="$param/@optional='true'"><xsl:value-of select="'O'" /></xsl:when>
            <xsl:otherwise><xsl:value-of select="'R'" /></xsl:otherwise>
        </xsl:choose>
        <xsl:value-of select="concat(' ', $param/@default)" />
    </td>

    <td>
        <xsl:choose>
            <xsl:when test="description">
                <xsl:attribute name="data-type">
                    <xsl:value-of select="description/@type" />
                </xsl:attribute>
                <pre><xsl:copy-of select="description/node()" /></pre>
            </xsl:when>
            <xsl:otherwise><xsl:value-of select="@summary" /></xsl:otherwise>
        </xsl:choose>
        <xsl:call-template name="enum">
            <xsl:with-param name="enum" select="$param/enum"/>
        </xsl:call-template>
    </td>
</tr>
</xsl:template>


<!-- 显示枚举类型的内容 -->
<xsl:template name="enum">
<xsl:param name="enum" />
<xsl:if test="$enum">
    <p><xsl:copy-of select="$locale-enum" /></p>
    <ul>
    <xsl:for-each select="$enum">
        <li>
        <xsl:call-template name="deprecated">
            <xsl:with-param name="deprecated" select="@deprecated" />
        </xsl:call-template>

        <xsl:value-of select="@value" />:
        <xsl:choose>
            <xsl:when test="description">
                <div data-type="{description/@type}">
                    <pre><xsl:copy-of select="description/node()" /></pre>
                </div>
            </xsl:when>
            <xsl:otherwise><xsl:value-of select="summary" /></xsl:otherwise>
        </xsl:choose>
        </li>
    </xsl:for-each>
    </ul>
</xsl:if>
</xsl:template>


<!--
给指定的元素添加已弃用的标记

该模板会给父元素添加 class 和 title 属性，
所以必须要在父元素的任何子元素之前，否则 chrome 和 safari 可能无法正常解析。
-->
<xsl:template name="deprecated">
<xsl:param name="deprecated" />
<xsl:if test="$deprecated">
    <xsl:attribute name="class"><xsl:value-of select="'del'" /></xsl:attribute>
    <xsl:attribute name="title">
        <xsl:value-of select="$deprecated" />
    </xsl:attribute>
</xsl:if>
</xsl:template>

<!-- 用于将 API 地址转换成合法的 ID 标记 -->
<xsl:variable name="id-from" select="'{}/'" />
<xsl:variable name="id-to" select="'__-'" />

<!-- 根据情况获取相应的图标 -->
<xsl:variable name="icon">
    <xsl:choose>
        <xsl:when test="/apidoc/@logo">
            <xsl:value-of select="/apidoc/@logo" />
        </xsl:when>
        <xsl:otherwise>
            <xsl:value-of select="concat($base-url, '../icon.svg')" />
        </xsl:otherwise>
    </xsl:choose>
</xsl:variable>

<!--
获取相对于当前 xsl 文件的基地址
xsl 2.0 可以直接采用 base-uri(document(''))
-->
<xsl:variable name="base-url">
    <xsl:apply-templates select="processing-instruction('xml-stylesheet')" />
</xsl:variable>

<xsl:template match="processing-instruction('xml-stylesheet')">
    <xsl:variable name="v1" select="substring-after(., 'href=&quot;')" />
    <!-- NOTE: 此处假定当前文件叫作 apidoc.xsl，如果不是的话，需要另外处理此代码 -->
    <xsl:variable name="v2" select="substring-before($v1, 'apidoc.xsl&quot;')" />
    <xsl:value-of select="$v2" />
</xsl:template>

</xsl:stylesheet>
`),
	},
	{
		Name:        "v6/locales.xsl",
		ContentType: "text/xsl; charset=utf-8",
		Content: []byte(`<?xml version="1.0" encoding="UTF-8"?>

<!-- 当前文件实现了简单的翻译功能 -->

<xsl:stylesheet
version="1.0"
xmlns:xsl="http://www.w3.org/1999/XSL/Transform"
xmlns:l="urn:locale"
exclude-result-prefixes="l">

<!-- 当前支持的本地化列表，其中第一个会被当作默认值。 -->
<l:locales>
    <locale id="zh-hans">简体中文</locale>
    <locale id="zh-hant">繁體中文</locale>
</l:locales>

<xsl:template name="languages">
    <xsl:for-each select="document('')/xsl:stylesheet/l:locales/locale">
    <li lang="{@id}" role="menuitemradio">
        <label><input type="radio" name="lang" checked="{$curr-lang=@id}" />&#160;<xsl:value-of select="." /></label>
    </li>
    </xsl:for-each>
</xsl:template>

<!-- language -->
<xsl:variable name="locale-language">
    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hans'" />
        <xsl:with-param name="text" select="document('')/xsl:stylesheet/l:locales/locale[@id='zh-hans']" />
    </xsl:call-template>

    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hant'" />
        <xsl:with-param name="text" select="document('')/xsl:stylesheet/l:locales/locale[@id='zh-hant']" />
    </xsl:call-template>
</xsl:variable>

<!-- server -->
<xsl:variable name="locale-server">
    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hans'" />
        <xsl:with-param name="text" select="'服务'" />
    </xsl:call-template>

    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hant'" />
        <xsl:with-param name="text" select="'服務'" />
    </xsl:call-template>
</xsl:variable>

<!-- tag -->
<xsl:variable name="locale-tag">
    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hans'" />
        <xsl:with-param name="text" select="'标签'" />
    </xsl:call-template>

    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hant'" />
        <xsl:with-param name="text" select="'標簽'" />
    </xsl:call-template>
</xsl:variable>

<!-- expand -->
<xsl:variable name="locale-expand">
    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hans'" />
        <xsl:with-param name="text" select="'展开'" />
    </xsl:call-template>

    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hant'" />
        <xsl:with-param name="text" select="'展開'" />
    </xsl:call-template>
</xsl:variable>

<!-- method -->
<xsl:variable name="locale-method">
    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hans'" />
        <xsl:with-param name="text" select="'请求方法'" />
    </xsl:call-template>

    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hant'" />
        <xsl:with-param name="text" select="'請求方法'" />
    </xsl:call-template>
</xsl:variable>

<!-- request -->
<xsl:variable name="locale-request">
    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hans'" />
        <xsl:with-param name="text" select="'请求'" />
    </xsl:call-template>

    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hant'" />
        <xsl:with-param name="text" select="'請求'" />
    </xsl:call-template>
</xsl:variable>

<!-- response -->
<xsl:variable name="locale-response">
    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hans'" />
        <xsl:with-param name="text" select="'返回'" />
    </xsl:call-template>

    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hant'" />
        <xsl:with-param name="text" select="'返回'" />
    </xsl:call-template>
</xsl:variable>

<!-- callback -->
<xsl:variable name="locale-callback">
    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hans'" />
        <xsl:with-param name="text" select="'回调'" />
    </xsl:call-template>

    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hant'" />
        <xsl:with-param name="text" select="'回調'" />
    </xsl:call-template>
</xsl:variable>

<!-- path param -->
<xsl:variable name="locale-path-param">
    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hans'" />
        <xsl:with-param name="text" select="'路径参数'" />
    </xsl:call-template>

    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hant'" />
        <xsl:with-param name="text" select="'路徑參數'" />
    </xsl:call-template>
</xsl:variable>

<!-- query -->
<xsl:variable name="locale-query">
    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hans'" />
        <xsl:with-param name="text" select="'查询参数'" />
    </xsl:call-template>

    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hant'" />
        <xsl:with-param name="text" select="'查詢參數'" />
    </xsl:call-template>
</xsl:variable>

<!-- header -->
<xsl:variable name="locale-header">
    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hans'" />
        <xsl:with-param name="text" select="'报头'" />
    </xsl:call-template>

    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hant'" />
        <xsl:with-param name="text" select="'報頭'" />
    </xsl:call-template>
</xsl:variable>

<!-- body -->
<xsl:variable name="locale-body">
    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hans'" />
        <xsl:with-param name="text" select="'报文'" />
    </xsl:call-template>

    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hant'" />
        <xsl:with-param name="text" select="'報文'" />
    </xsl:call-template>
</xsl:variable>

<!-- example -->
<xsl:variable name="locale-example">
    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hans'" />
        <xsl:with-param name="text" select="'示例代码'" />
    </xsl:call-template>

    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hant'" />
        <xsl:with-param name="text" select="'示例代碼'" />
    </xsl:call-template>
</xsl:variable>

<!-- var -->
<xsl:variable name="locale-var">
    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hans'" />
        <xsl:with-param name="text" select="'变量'" />
    </xsl:call-template>

    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hant'" />
        <xsl:with-param name="text" select="'變量'" />
    </xsl:call-template>
</xsl:variable>

<!-- type -->
<xsl:variable name="locale-type">
    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hans'" />
        <xsl:with-param name="text" select="'类型'" />
    </xsl:call-template>

    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hant'" />
        <xsl:with-param name="text" select="'類型'" />
    </xsl:call-template>
</xsl:variable>

<!-- value -->
<xsl:variable name="locale-value">
    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hans'" />
        <xsl:with-param name="text" select="'值'" />
    </xsl:call-template>

    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hant'" />
        <xsl:with-param name="text" select="'值'" />
    </xsl:call-template>
</xsl:variable>

<!-- description -->
<xsl:variable name="locale-description">
    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hans'" />
        <xsl:with-param name="text" select="'描述'" />
    </xsl:call-template>

    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hant'" />
        <xsl:with-param name="text" select="'描述'" />
    </xsl:call-template>
</xsl:variable>

<!-- enum -->
<xsl:variable name="locale-enum">
    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hans'" />
        <xsl:with-param name="text" select="'枚举'" />
    </xsl:call-template>

    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hant'" />
        <xsl:with-param name="text" select="'枚舉'" />
    </xsl:call-template>
</xsl:variable>

<!-- license -->
<xsl:variable name="locale-license">
    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hans'" />
        <xsl:with-param name="text">
            文档版权为 <a href="{apidoc/license/@url}"><xsl:value-of select="apidoc/license/@text" /></a>
        </xsl:with-param>
    </xsl:call-template>

    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hant'" />
        <xsl:with-param name="text">
            文檔版權為 <a href="{apidoc/license/@url}"><xsl:value-of select="apidoc/license/@text" /></a>
        </xsl:with-param>
    </xsl:call-template>
</xsl:variable>

<!-- generator -->
<xsl:variable name="locale-generator">
    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hans'" />
        <xsl:with-param name="text">
            由 <a href="https://apidoc.tools">apidoc</a> 生成于 <time><xsl:value-of select="apidoc/@created" /></time>
        </xsl:with-param>
    </xsl:call-template>

    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'zh-hant'" />
        <xsl:with-param name="text">
            由 <a href="https://apidoc.tools">apidoc</a> 生成於 <time><xsl:value-of select="apidoc/@created" /></time>
        </xsl:with-param>
    </xsl:call-template>
</xsl:variable>

<xsl:template name="build-locale">
    <xsl:param name="lang" />
    <xsl:param name="text" />

    <xsl:variable name="class">
        <xsl:choose>
            <xsl:when test="$curr-lang=translate($lang, $uppercase, $lowercase)">
                <xsl:value-of select="''" />
            </xsl:when>
            <xsl:otherwise>
                <xsl:value-of select="'hidden'" />
            </xsl:otherwise>
        </xsl:choose>
    </xsl:variable>

    <!-- data-locale 属性表示该元素是一个本地化信息元素，JS 代码通过该标记切换语言。 -->
    <span data-locale="true" lang="{$lang}" class="{$class}"><xsl:copy-of select="$text" /></span>
</xsl:template>

<!--
返回当前文档的语言，会转换为小写，_ 也会被转换成 -
如果文档指定的语言不存在，则会采取 l:locales 中的第一个元素作为默认语言。
-->
<xsl:variable name="curr-lang">
    <xsl:variable name="curr" select="translate(/apidoc/@lang, $uppercase, $lowercase)" />

    <xsl:variable name="r1">
        <xsl:for-each select="document('')/xsl:stylesheet/l:locales/locale">
            <xsl:if test="@id=$curr"><xsl:value-of select="$curr" /></xsl:if>
        </xsl:for-each>
    </xsl:variable>

    <xsl:variable name="r2">
    <xsl:choose>
        <xsl:when test="$r1 and not($r1='')"> 
            <xsl:value-of select="$r1" />
        </xsl:when>
        <xsl:otherwise>
            <xsl:value-of select="document('')/xsl:stylesheet/l:locales/locale[1]/@id" />
        </xsl:otherwise>
    </xsl:choose>
    </xsl:variable>

    <xsl:value-of select="$r2" />
</xsl:variable>

<!-- 用于实现 lower-case 和 upper-case，如果将来某天浏览器支持 xsl 2.0 了，可以直接采用相关函数 -->
<xsl:variable name="lowercase" select="'abcdefghijklmnopqrstuvwxyz-'" />
<xsl:variable name="uppercase" select="'ABCDEFGHIJKLMNOPQRSTUVWXYZ_'" />

</xsl:stylesheet>
`),
	},
	{
		Name:        "v6/view.html",
		ContentType: "text/html; charset=utf-8",
		Content: []byte(`<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>Document</title>

    <style>
    :root {
        --background: white;
    }
    @media (prefers-color-scheme: dark) {
        :root {
            --background: black;
        }
    }

    html,
    body {
        margin: 0;
        padding: 0;
        background: var(--background);
    }

    iframe {
        border: none;
        width: 100%;
        background: var(--background);
    }
    </style>
</head>

<body>
    <!--
    NOTE: xslt 中引用的 js，只有在 iframe 中才能执行，
    直接用 document.replace() 替换当前内容，不会执行引用的 JS 代码。
    -->
    <iframe id="apidoc"></iframe>
    <script>
    async function loadXML(path) {
        const obj = await fetch(path);
        return (new DOMParser()).parseFromString(await obj.text(), "text/xml");
    }

    function changeFrameHeight() {
        const iframe = document.getElementById("apidoc");
        iframe.height = document.documentElement.clientHeight;
    }

    window.onresize = function () {
        changeFrameHeight();
    } 

    const queries = new URLSearchParams(window.location.search);
    const url = queries.get('url');

    async function init() {
        const processor = new XSLTProcessor();
        processor.importStylesheet(await loadXML('./apidoc.xsl'));

        const xml = await loadXML(url);
        const doc = processor.transformToDocument(xml);
        const html = (new XMLSerializer()).serializeToString(doc);
        const blob = new Blob([html], { type: 'text/html' })
        const obj = URL.createObjectURL(blob)

        const iframe = document.getElementById('apidoc');
        iframe.src = obj;
        iframe.addEventListener('load', (e) => {
            changeFrameHeight();
        });
    }

    try{
        init()
    } catch (e) {
        console.error(e)
    }
    </script>
</body>

</html>
`),
	},
}
