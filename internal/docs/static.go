// 该文件由工具自动生成，请勿手动修改！

package docs

var data = []*FileInfo{{
	Name:        "example/index.xml",
	ContentType: "application/xml; charset=utf-8",
	Content: []byte(`<?xml version="1.0" encoding="UTF-8"?>

<?xml-stylesheet type="text/xsl" href="../v6/apidoc.xsl"?>
<apidoc apidoc="6.1.0" created="2020-06-03T02:25:03+08:00" version="1.1.1">
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
       后台管理接口，<br /><br /><br /><br /><p style="color:red">admin</p>
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
			<param xml-attr="true" name="count" type="number" optional="false" summary="summary"></param>
			<param xml-wrapped="users" name="user" type="object" array="true" summary="list">
				<param xml-attr="true" name="id" type="number" summary="用户 ID"></param>
				<param xml-attr="true" name="name" type="string" summary="用户名"></param>
				<param name="groups" type="object" optional="true" array="true" summary="用户所在的权限组">
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
			<param name="count" type="number" optional="false" summary="summary"></param>
			<param name="list" type="object" array="true" summary="list">
				<param name="id" type="number" summary="用户 ID"></param>
				<param name="name" type="string" summary="用户名"></param>
				<param name="groups" type="object" optional="true" array="true" summary="用户所在的权限组">
					<param name="id" type="string" summary="权限组 ID"></param>
					<param name="name" type="string" summary="权限组名称"></param>
				</param>
			</param>
			<header name="content-type" type="string" summary="application/json"></header>
		</request>
		<request name="users" type="object" mimetype="application/xml">
			<param name="count" type="number" optional="false" summary="summary"></param>
			<param name="list" type="object" array="true" summary="list">
				<param name="id" type="number" summary="用户 ID"></param>
				<param name="name" type="string" summary="用户名"></param>
				<param name="groups" type="object" optional="true" array="true" summary="用户所在的权限组">
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
			<param name="groups" type="object" optional="true" summary="用户所在的权限组">
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
		<server>admin</server>
	</api>
	<header name="Authorization" type="string" summary="token 值"></header>
	<header name="Accept" type="string" summary="能接受的字符集"></header>
	<response name="result" type="object" status="400">
		<param xml-attr="true" name="code" type="number" optional="false" summary="状态码"></param>
		<param name="message" type="string" optional="false" summary="错误信息"></param>
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
		Name:        "index.cmn-Hant.xml",
		ContentType: "application/xml; charset=utf-8",
		Content: []byte(`<?xml version="1.0" encoding="utf-8"?>
<?xml-stylesheet type="text/xsl" href="./index.xsl"?>

<!--
這是官網首頁內容，同時也是簡體中文的本地化內容。

其它語言的本化地內容，需要重新改寫本文件中除註釋外的所有內容。
-->

<docs lang="cmn-Hant">
    <title>apidoc | RESTful API 文檔處理工具</title>
    <license url="https://creativecommons.org/licenses/by/4.0/deed.zh">署名 4.0 國際 (CC BY 4.0)</license>

    <!-- 壹些翻譯比較零散的翻譯內容 -->
    <locales>
        <header> <!-- 表格報頭 -->
            <name>名稱</name>
            <type>類型</type>
            <required>必填</required>
            <description>描述</description>
        </header>
        <goto-top>返回頂部</goto-top>
    </locales>

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
    </doc>

    <doc id="usage" title="使用" />

    <doc id="spec" title="文檔格式">
        <p>文檔采用 XML 格式。存在兩個頂級標簽：<code>apidoc</code> 和 <code>api</code>，用於描述整體內容和具體接口信息。</p>

        <p>文檔被從註釋中提取之後，最終會被合並成壹個 XML 文件，在該文件中 <code>api</code> 作為 <code>apidoc</code> 的壹個子元素存在，如果妳的項目不想把文檔寫在註釋中，也可以直接編寫壹個完整的 XML 文件，將 <code>api</code> 作為 <code>apidoc</code> 的壹個子元素。</p>

        <p>具體可參考<a href="./example/index.xml">示例代碼。</a></p>

        <p>以下是對各個 XML 元素以及參數介紹，其中以 <code>@</code> 開頭的表示 XML 屬性；<code>.</code> 表示為當前元素的內容；其它表示子元素。</p>
    </doc>

    <!--######################### 以下为文档内容的子项 ###########################-->

    <doc id="install" title="安裝" parent="usage">
          <p>可以直接從 <a href="https://github.com/caixw/apidoc/releases">https://github.com/caixw/apidoc/releases</a> 查找妳需要的版本下載，放入 <code>PATH</code> 中即可使用。如果沒有妳需要的平臺文件，則需要從源代碼編譯：</p>
        <ul>
            <li>下載 <a href="https://golang.org/dl/">Go</a>；</li>
            <li>下載源代碼，<samp>git clone github.com/caixw/apidoc</samp>；</li>
            <li>執行代碼中 <code>build/build.sh</code> 或是 <code>build/build.cmd</code> 進行編譯；</li>
            <li>編譯好的文件存放在 cmd/apidoc 下，可以將該文件放置在 PATH 目錄；</li>
        </ul>
    </doc>

    <doc id="env" title="環境變量" parent="usage">
        <p>apidoc 會讀取 <var>LANG</var> 的環境變量作為其本地化的依據，若想指定其它語種，可以手動指定 <var>LANG</var> 環境變量：<samp>LANG=zh-Hant apidoc</samp>。在 windows 系統中，若不存在 <var>LANG</var> 環境變量，則會調用 <samp>GetUserDefaultLocaleName</samp> 函數來獲取相應的語言信息。</p>
    </doc>

    <doc id="cli" title="命令行" parent="usage">
        <p>可以通過 <samp>apidoc help</samp> 查看命令行支持的子命令。包含了以下幾個：</p>
    </doc>

    <doc id="apidoc.yaml" title=".apidoc.yaml" parent="usage">
        <p>配置文件名固定為 <code>.apidoc.yaml</code>，格式為 YAML，可參考 <a href="example/.apidoc.yaml">.apidoc.yaml</a>。文件可以通過命令 <code>apidoc detect</code> 生成。主要包含了以幾個配置項：</p>
    </doc>

    <footer>
        <license>
            <p>當前頁面內容托管於 </p><p>，並采用</p><p>進行許可。</p>
        </license>
    </footer>
</docs>
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
    opacity: .7
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

.type {
    color: var(--color)
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
    right: calc(((100% - var(--max-width)) / 2) - 3rem);
}

/*
 * 用于计算 goto-top 按钮的位置始终保持在内容主体的右侧边上。
 * 1024px 为 --max-width 的值，但是 CSS 并不支持直接在 @media
 * 使用 var() 函数。
 */
@media (max-width: calc(1024px + 6rem)) {
    .goto-top {
        right: 3rem;
    }
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

<docs lang="cmn-Hans">
    <title>apidoc | RESTful API 文档处理工具</title>
    <license url="https://creativecommons.org/licenses/by/4.0/deed.zh">署名 4.0 国际 (CC BY 4.0)</license>

    <!-- 一些翻译比较零散的翻译内容 -->
    <locales>
        <header> <!-- 表格报头 -->
            <name>名称</name>
            <type>类型</type>
            <required>必填</required>
            <description>描述</description>
        </header>
        <goto-top>返回顶部</goto-top>
    </locales>

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
    </doc>

    <doc id="usage" title="使用" />

    <doc id="spec" title="文档格式">
        <p>文档采用 XML 格式。存在两个顶级标签：<code>apidoc</code> 和 <code>api</code>，用于描述整体内容和具体接口信息。</p>

        <p>文档被从注释中提取之后，最终会被合并成一个 XML 文件，在该文件中 <code>api</code> 作为 <code>apidoc</code> 的一个子元素存在，如果你的项目不想把文档写在注释中，也可以直接编写一个完整的 XML 文件，将 <code>api</code> 作为 <code>apidoc</code> 的一个子元素。</p>

        <p>具体可参考<a href="./example/index.xml">示例代码</a>。</p>

        <p>以下是对各个 XML 元素以及参数介绍，其中以 <code>@</code> 开头的表示 XML 属性；<code>.</code> 表示为当前元素的内容；其它表示子元素。</p>
    </doc>

    <!--######################### 以下为文档内容的子项 ###########################-->

    <doc id="install" title="安装" parent="usage">
        <p>可以直接从 <a href="https://github.com/caixw/apidoc/releases">https://github.com/caixw/apidoc/releases</a> 查找你需要的版本下载，放入 <code>PATH</code> 中即可使用。如果没有你需要的平台文件，则需要从源代码编译：</p>
        <ul>
            <li>下载 <a href="https://golang.org/dl/">Go</a>；</li>
            <li>下载源代码，<samp>git clone github.com/caixw/apidoc</samp>；</li>
            <li>执行代码中 <code>build/build.sh</code> 或是 <code>build/build.cmd</code> 进行编译；</li>
            <li>编译好的文件存放在 cmd/apidoc 下，可以将该文件放置在 PATH 目录；</li>
        </ul>
    </doc>

    <doc id="env" title="环境变量" parent="usage">
        <p>apidoc 会读取 <var>LANG</var> 的环境变量作为其本地化的依据，若想指定其它语种，可以手动指定 <var>LANG</var> 环境变量：<samp>LANG=zh-Hant apidoc</samp>。在 windows 系统中，若不存在 <var>LANG</var> 环境变量，则会调用 <samp>GetUserDefaultLocaleName</samp> 函数来获取相应的语言信息。</p>
    </doc>

    <doc id="cli" title="命令行" parent="usage">
        <p>可以通过 <samp>apidoc help</samp> 查看命令行支持的子命令。包含了以下几个：</p>
    </doc>

    <doc id="apidoc.yaml" title=".apidoc.yaml" parent="usage">
        <p>配置文件名固定为 <code>.apidoc.yaml</code>，格式为 YAML，可参考 <a href="example/.apidoc.yaml">.apidoc.yaml</a>。文件可以通过命令 <code>apidoc detect</code> 生成。主要包含了以几个配置项：</p>
    </doc>

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

<xsl:variable name="locale-file">
    <xsl:value-of select="document('site.xml')/site/locales/locale[@id=$curr-lang]/@doc" />
</xsl:variable>

<!-- 获取当前文档的语言名称，如果不存在，则直接采用 @lang 属性 -->
<xsl:variable name="curr-lang-title">
    <xsl:variable name="title">
        <xsl:value-of select="document('site.xml')/site/locales/locale[@id=$curr-lang]/@title" />
    </xsl:variable>

    <xsl:choose>
        <xsl:when test="$title=''"><xsl:value-of select="$curr-lang" /></xsl:when>
        <xsl:otherwise><xsl:value-of select="$title" /></xsl:otherwise>
    </xsl:choose>
</xsl:variable>

<xsl:variable name="keywords">
    <xsl:for-each select="document('site.xml')/site/languages/language">
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
            <link rel="canonical" href="{document('site.xml')/site/url}" />
            <link rel="stylesheet" type="text/css" href="./index.css" />
            <link rel="license" href="{/docs/liense/@url}" />
            <link rel="stylesheet" type="text/css" href="https://cdnjs.cloudflare.com/ajax/libs/prism/1.20.0/themes/prism-tomorrow.min.css" />
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
                    <a href="{document('site.xml')/site/repo}">Github</a>
                    <xsl:value-of select="docs/footer/license/p[2]" />
                    <a href="{docs/license/@url}"><xsl:value-of select="docs/license" /></a>
                    <xsl:value-of select="docs/footer/license/p[3]" />
                </p>
                </div>
                <a href="#" class="goto-top" title="{docs/locales/goto-top}" aria-label="{docs/locales/goto-top}" />
            </footer>

            <script src="https://cdnjs.cloudflare.com/ajax/libs/prism/1.20.0/components/prism-core.min.js"></script>
            <script src="https://cdnjs.cloudflare.com/ajax/libs/prism/1.20.0/plugins/autoloader/prism-autoloader.min.js"></script>
        </body>
    </html>
</xsl:template>

<xsl:template name="header">
    <header>
        <div class="wrap">
            <h1>
                <img alt="logo" src="./icon.svg" />
                <xsl:value-of select="document('site.xml')/site/name" />
                <span class="version">&#160;(<xsl:value-of select="document('site.xml')/site/version" />)</span>
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
                        <xsl:for-each select="document('site.xml')/site/locales/locale">
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

        <xsl:if test="$id='spec'">
            <xsl:for-each select="document($locale-file)/locale/spec/type">
                <xsl:call-template name="type">
                    <xsl:with-param name="type" select="." />
                </xsl:call-template>
            </xsl:for-each>
        </xsl:if>

        <xsl:if test="$id='cli'">
            <xsl:call-template name="commands" />
        </xsl:if>

        <xsl:if test="$id='apidoc.yaml'">
            <xsl:call-template name="apidocYAML" />
        </xsl:if>
    </article>
</xsl:template>

<!-- 以下两个变量仅用于 type 和 commands 模板，在模板中无法直接使用 /docs 元素，所以使用变量引用 -->
<xsl:variable name="header-locale" select="/docs/locales/header" />

<!-- 将子命令显示为一个 table -->
<xsl:template name="commands">
    <table>
        <thead>
            <tr>
                <th><xsl:copy-of select="$header-locale/name" /></th>
                <th><xsl:copy-of select="$header-locale/description" /></th>
            </tr>
        </thead>

        <tbody>
            <xsl:for-each select="document($locale-file)/locale/commands/command">
            <tr>
                <th><xsl:value-of select="@name" /></th>
                <td><xsl:copy-of select="node()" /></td>
            </tr>
            </xsl:for-each>
        </tbody>
    </table>
</xsl:template>

<!-- 将 apidoc.yaml 类型说明显示为 table -->
<xsl:template name="apidocYAML">
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
            <xsl:for-each select="document($locale-file)/locale/config/item">
                <tr>
                    <th><xsl:value-of select="@name" /></th>
                    <td>
                        <xsl:value-of select="@type" />
                        <xsl:if test="@array='true'"><xsl:value-of select="'[]'" /></xsl:if>
                    </td>
                    <td>
                        <xsl:call-template name="checkbox">
                            <xsl:with-param name="chk" select="@required" />
                        </xsl:call-template>
                    </td>
                    <td><xsl:copy-of select="node()" /></td>
                </tr>
            </xsl:for-each>
        </tbody>
    </table>
</xsl:template>


<!-- 将类型显示为一个 table -->
<xsl:template name="type">
    <xsl:param name="type" />

    <article id="type_{$type/@name}">
        <h3>
            <xsl:value-of select="$type/@name" />
            <a class="link" href="#type_{$type/@name}">&#160;&#160;&#128279;</a>
        </h3>
        <xsl:copy-of select="$type/usage/node()" />
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
                    <tr>
                        <th><xsl:value-of select="@name" /></th>
                        <td>
                            <a class="type" href="#type_{@type}"><xsl:value-of select="@type" /></a>
                            <xsl:if test="@array='true'"><xsl:value-of select="'[]'" /></xsl:if>
                        </td>
                        <td>
                            <xsl:call-template name="checkbox">
                                <xsl:with-param name="chk" select="@required" />
                            </xsl:call-template>
                        </td>
                        <td><xsl:copy-of select="node()" /></td>
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
            <label aria-label="true">
                <input type="checkbox" checked="true" disabled="true" aria-hiddren="true" />
            </label>
        </xsl:when>
        <xsl:otherwise>
            <label aria-label="false">
                <input type="checkbox" disabled="true" aria-hiddren="true" />
            </label>
        </xsl:otherwise>
    </xsl:choose>
</xsl:template>

</xsl:stylesheet>
`),
	},
	{
		Name:        "locale.cmn-Hans.xml",
		ContentType: "application/xml; charset=utf-8",
		Content: []byte(`<?xml version="1.0" encoding="UTF-8"?>

<!-- 该文件由工具自动生成，请勿手动修改！ -->

<locale>
	<spec>
		<type name="apidoc">
			<usage>用于描述整个文档的相关内容，只能出现一次。</usage>
			<item name="@apidoc" type="string" array="false" required="false">文档的版本要号</item>
			<item name="@lang" type="string" array="false" required="false">文档内容的本地化 ID，比如 <var>zh-Hans</var>、<var>en-US</var> 等。</item>
			<item name="@logo" type="string" array="false" required="false">文档的图标，仅可使用 SVG 格式图标。</item>
			<item name="@created" type="date" array="false" required="false">文档的创建时间</item>
			<item name="@version" type="version" array="false" required="false">文档的版本号</item>
			<item name="title" type="string" array="false" required="true">文档的标题</item>
			<item name="description" type="richtext" array="false" required="false">文档的整体描述内容</item>
			<item name="contact" type="contact" array="false" required="false">文档作者的联系方式</item>
			<item name="license" type="link" array="false" required="false">文档的版权信息</item>
			<item name="tag" type="tag" array="true" required="false">文档中定义的所有标签</item>
			<item name="server" type="server" array="true" required="false">API 基地址列表，每个 API 最少应该有一个 server。</item>
			<item name="api" type="api" array="true" required="false">文档中的 API 文档</item>
			<item name="header" type="param" array="true" required="false">文档中所有 API 都包含的公共报头</item>
			<item name="response" type="request" array="true" required="false">文档中所有 API 文档都需要支持的返回内容</item>
			<item name="mimetype" type="string" array="true" required="true">文档所支持的 mimetype</item>
		</type>
		<type name="richtext">
			<usage>富文本内容</usage>
			<item name="@type" type="string" array="false" required="true">指定富文本内容的格式，目前支持 <var>html</var> 和 <var>markdown</var>。</item>
			<item name="." type="string" array="false" required="true">富文本的实际内容</item>
		</type>
		<type name="contact">
			<usage>用于描述联系方式</usage>
			<item name="@name" type="string" array="false" required="true">联系人的名称</item>
			<item name="url" type="string" array="false" required="false">联系人的 URL</item>
			<item name="email" type="string" array="false" required="false">联系人的电子邮件</item>
		</type>
		<type name="link">
			<usage>用于描述链接信息，一般转换为 HTML 的 <code>a</code> 标签。</usage>
			<item name="@text" type="string" array="false" required="true">链接的字面文字</item>
			<item name="@url" type="string" array="false" required="true">链接指向的文本</item>
		</type>
		<type name="tag">
			<usage>用于对各个 API 进行分类</usage>
			<item name="@name" type="string" array="false" required="true">标签的唯一 ID</item>
			<item name="@title" type="string" array="false" required="true">标签的字面名称</item>
			<item name="@deprecated" type="version" array="false" required="false">该标签在大于该版本时被弃用</item>
		</type>
		<type name="server">
			<usage>用于指定各个 API 的服务器地址</usage>
			<item name="@name" type="string" array="false" required="true">服务唯一 ID</item>
			<item name="@url" type="string" array="false" required="true">服务的基地址，与该服务关联的 API，访问地址都是相对于此地址的。</item>
			<item name="@deprecated" type="version" array="false" required="false">服务在大于该版本时被弃用</item>
			<item name="@summary" type="string" array="false" required="false">服务的摘要信息</item>
			<item name="description" type="richtext" array="false" required="false">服务的详细描述</item>
		</type>
		<type name="api">
			<usage>用于定义单个 API 接口的具体内容</usage>
			<item name="@version" type="version" array="false" required="false">表示此接口在该版本中添加</item>
			<item name="@method" type="string" array="false" required="true">当前接口所支持的请求方法</item>
			<item name="@id" type="string" array="false" required="false">接口的唯一 ID</item>
			<item name="@summary" type="string" array="false" required="false">简要介绍</item>
			<item name="@deprecated" type="version" array="false" required="false">在此版本之后将会被弃用</item>
			<item name="path" type="path" array="false" required="true">定义路径信息</item>
			<item name="description" type="richtext" array="false" required="false">该接口的详细介绍，为 HTML 内容。</item>
			<item name="request" type="request" array="true" required="false">定义可用的请求信息</item>
			<item name="response" type="request" array="true" required="false">定义可能的返回信息</item>
			<item name="callback" type="callback" array="false" required="false">定义回调接口内容</item>
			<item name="header" type="param" array="true" required="false">传递的报头内容，如果是某个 mimetype 专用的，可以放在 request 元素中。</item>
			<item name="tag" type="string" array="true" required="false">关联的标签</item>
			<item name="server" type="string" array="true" required="false">关联的服务</item>
		</type>
		<type name="path">
			<usage>用于定义请求时与路径相关的内容</usage>
			<item name="@path" type="string" array="false" required="true">接口地址</item>
			<item name="param" type="param" array="true" required="false">地址中的参数</item>
			<item name="query" type="param" array="true" required="false">地址中的查询参数</item>
		</type>
		<type name="param">
			<usage>参数类型，基本上可以作为 request 的子集使用。</usage>
			<item name="@xml-attr" type="bool" array="false" required="false">是否作为父元素的属性，仅作用于 XML 元素。是否作为父元素的属性，仅用于 XML 的请求。</item>
			<item name="@xml-extract" type="bool" array="false" required="false">将当前元素的内容作为父元素的内容，要求父元素必须为 <var>object</var>。</item>
			<item name="@xml-cdata" type="bool" array="false" required="false">当前内容为 CDATA，與 xml-attr 互斥。</item>
			<item name="@xml-ns" type="string" array="false" required="false">XML 标签的命名空间</item>
			<item name="@xml-ns-prefix" type="string" array="false" required="false">XML 标签的命名空间名称前缀</item>
			<item name="@xml-wrapped" type="string" array="false" required="false">如果当前元素的 <code>@array</code> 为 <var>true</var>，是否将其包含在 wrapped 指定的标签中。</item>
			<item name="@name" type="string" array="false" required="true">值的名称</item>
			<item name="@type" type="string" array="false" required="true">值的类型，可以是 <var>string</var>、<var>number</var>、<var>bool</var> 和 <var>object</var></item>
			<item name="@deprecated" type="version" array="false" required="false">表示在大于等于该版本号时不再启作用</item>
			<item name="@default" type="string" array="false" required="false">默认值</item>
			<item name="@optional" type="bool" array="false" required="false">是否为可选的参数</item>
			<item name="@array" type="bool" array="false" required="false">是否为数组</item>
			<item name="@summary" type="string" array="false" required="false">简要介绍</item>
			<item name="@array-style" type="bool" array="false" required="false">以数组的方式展示数据</item>
			<item name="param" type="param" array="true" required="false">子类型，比如对象的子元素。</item>
			<item name="enum" type="enum" array="true" required="false">当前参数可用的枚举值</item>
			<item name="description" type="richtext" array="false" required="false">详细介绍，为 HTML 内容。</item>
		</type>
		<type name="enum">
			<usage>定义枚举类型的数所的枚举值</usage>
			<item name="@deprecated" type="version" array="false" required="false">该属性弃用的版本号</item>
			<item name="@value" type="string" array="false" required="true">枚举值</item>
			<item name="@summary" type="string" array="false" required="false">枚举值的说明</item>
			<item name="description" type="richtext" array="false" required="false">枚举值的详细说明</item>
		</type>
		<type name="request">
			<usage>定义了请求和返回的相关内容</usage>
			<item name="@xml-attr" type="bool" array="false" required="false">是否作为父元素的属性，仅作用于 XML 元素。是否作为父元素的属性，仅用于 XML 的请求。</item>
			<item name="@xml-extract" type="bool" array="false" required="false">将当前元素的内容作为父元素的内容，要求父元素必须为 <var>object</var>。</item>
			<item name="@xml-cdata" type="bool" array="false" required="false">当前内容为 CDATA，與 xml-attr 互斥。</item>
			<item name="@xml-ns" type="string" array="false" required="false">XML 标签的命名空间</item>
			<item name="@xml-ns-prefix" type="string" array="false" required="false">XML 标签的命名空间名称前缀</item>
			<item name="@xml-wrapped" type="string" array="false" required="false">如果当前元素的 <code>@array</code> 为 <var>true</var>，是否将其包含在 wrapped 指定的标签中。</item>
			<item name="@name" type="string" array="false" required="false">当 mimetype 为 <var>application/xml</var> 时，此值表示 XML 的顶层元素名称，否则无用。</item>
			<item name="@type" type="string" array="false" required="false">值的类型，可以是 <var>string</var>、<var>number</var>、<var>bool</var>、<var>object</var> 和空值；空值表示不输出任何内容。</item>
			<item name="@deprecated" type="version" array="false" required="false">表示在大于等于该版本号时不再启作用</item>
			<item name="@array" type="bool" array="false" required="false">是否为数组</item>
			<item name="@summary" type="string" array="false" required="false">简要介绍</item>
			<item name="@status" type="number" array="false" required="false">状态码。在 request 中，该值不可用，否则为必填项。</item>
			<item name="@mimetype" type="string" array="false" required="false">媒体类型，比如 <var>application/json</var> 等。</item>
			<item name="enum" type="enum" array="true" required="false">当前参数可用的枚举值</item>
			<item name="param" type="param" array="true" required="false">子类型，比如对象的子元素。</item>
			<item name="example" type="example" array="true" required="false">示例代码</item>
			<item name="header" type="param" array="true" required="false">传递的报头内容</item>
			<item name="description" type="richtext" array="false" required="false">详细介绍，为 HTML 内容。</item>
		</type>
		<type name="example">
			<usage>示例代码</usage>
			<item name="@mimetype" type="string" array="false" required="true">特定于类型的示例代码</item>
			<item name="@summary" type="string" array="false" required="false">示例代码的概要信息</item>
			<item name="." type="string" array="false" required="true">示例代码的内容，需要使用 CDATA 包含代码。</item>
		</type>
		<type name="callback">
			<usage>定义接口的回调内容</usage>
			<item name="@method" type="string" array="false" required="true">回调的请求方法</item>
			<item name="@summary" type="string" array="false" required="false">简要介绍</item>
			<item name="@deprecated" type="version" array="false" required="false">在此版本之后将会被弃用</item>
			<item name="path" type="path" array="false" required="false">回调的请求地址</item>
			<item name="description" type="richtext" array="false" required="false">对于回调的详细介绍</item>
			<item name="response" type="request" array="true" required="false">定义可能的返回信息</item>
			<item name="request" type="request" array="true" required="true">定义可用的请求信息</item>
			<item name="header" type="param" array="true" required="false">传递的报头内容</item>
		</type>
		<type name="string">
			<usage>普通的字符串类型</usage>
		</type>
		<type name="date">
			<usage>采用 <a href="https://tools.ietf.org/html/rfc3339">RFC3339</a> 格式表示的时间，比如：<samp>2019-12-16T00:35:48+08:00</samp></usage>
		</type>
		<type name="version">
			<usage>版本号，格式遵守 <a href="https://semver.org/lang/zh-CN/">semver</a> 规则</usage>
		</type>
		<type name="bool">
			<usage>布尔值类型，取值为 <var>true</var> 或是 <var>false</var></usage>
		</type>
		<type name="number">
			<usage>普通的数值类型</usage>
		</type>
	</spec>
	<commands>
		<command name="build">生成文档内容</command>
		<command name="detect">根据目录下的内容生成配置文件</command>
		<command name="help">显示帮助信息</command>
		<command name="lang">显示所有支持的语言</command>
		<command name="locale">显示所有支持的本地化内容</command>
		<command name="lsp">启动 language server protocol 服务</command>
		<command name="mock">启用 mock 服务</command>
		<command name="static">启用静态文件服务</command>
		<command name="test">测试语法的正确性</command>
		<command name="version">显示版本信息</command>
	</commands>
	<config>
		<item name="version" type="string" array="false" required="true">此配置文件的所使用的文档版本</item>
		<item name="inputs" type="object" array="true" required="true">指定输入的数据，同一项目只能解析一种语言。</item>
		<item name="inputs.lang" type="string" array="false" required="true">源文件类型。具体支持的类型可通过 -l 参数进行查找。</item>
		<item name="inputs.dir" type="string" array="false" required="true">需要解析的源文件所在目录</item>
		<item name="inputs.exts" type="string" array="true" required="false">只从这些扩展名的文件中查找文档</item>
		<item name="inputs.recursive" type="bool" array="false" required="false">是否解析子目录下的源文件</item>
		<item name="inputs.encoding" type="string" array="false" required="false">编码，默认为 <var>utf-8</var>，值可以是 <a href="https://www.iana.org/assignments/character-sets/character-sets.xhtml">character-sets</a> 中的内容。</item>
		<item name="output" type="object" array="false" required="true">控制输出行为</item>
		<item name="output.type" type="string" array="false" required="false">输出的类型，目前可以 <var>apidoc+xml</var>、<var>openapi+json</var> 和 <var>openapi+yaml</var>。</item>
		<item name="output.path" type="string" array="false" required="true">指定输出的文件名，包含路径信息。</item>
		<item name="output.tags" type="string" array="true" required="false">只输出与这些标签相关联的文档，默认为全部。</item>
		<item name="output.style" type="string" array="false" required="false">为 XML 文件指定的 XSL 文件</item>
	</config>
</locale>
`),
	},
	{
		Name:        "locale.cmn-Hant.xml",
		ContentType: "application/xml; charset=utf-8",
		Content: []byte(`<?xml version="1.0" encoding="UTF-8"?>

<!-- 该文件由工具自动生成，请勿手动修改！ -->

<locale>
	<spec>
		<type name="apidoc">
			<usage>用於描述整個文檔的相關內容，只能出現壹次。</usage>
			<item name="@apidoc" type="string" array="false" required="false">文檔的版本要號</item>
			<item name="@lang" type="string" array="false" required="false">文檔內容的本地化 ID，比如 <var>zh-Hans</var>、<var>en-US</var> 等。</item>
			<item name="@logo" type="string" array="false" required="false">文檔的圖標，僅可使用 SVG 格式圖標。</item>
			<item name="@created" type="date" array="false" required="false">文檔的創建時間</item>
			<item name="@version" type="version" array="false" required="false">文檔的版本號</item>
			<item name="title" type="string" array="false" required="true">文檔的標題</item>
			<item name="description" type="richtext" array="false" required="false">文檔的整體描述內容</item>
			<item name="contact" type="contact" array="false" required="false">文檔作者的聯系方式</item>
			<item name="license" type="link" array="false" required="false">文檔的版權信息</item>
			<item name="tag" type="tag" array="true" required="false">文檔中定義的所有標簽</item>
			<item name="server" type="server" array="true" required="false">API 基地址列表，每個 API 最少應該有壹個 server。</item>
			<item name="api" type="api" array="true" required="false">文檔中的 API 文檔</item>
			<item name="header" type="param" array="true" required="false">文檔中所有 API 都包含的公共報頭</item>
			<item name="response" type="request" array="true" required="false">文檔中所有 API 文檔都需要支持的返回內容</item>
			<item name="mimetype" type="string" array="true" required="true">文檔所支持的 mimetype</item>
		</type>
		<type name="richtext">
			<usage>富文本內容</usage>
			<item name="@type" type="string" array="false" required="true">指定富文本內容的格式，目前支持 <var>html</var> 和 <var>markdown</var>。</item>
			<item name="." type="string" array="false" required="true">富文本的實際內容</item>
		</type>
		<type name="contact">
			<usage>用於描述聯系方式</usage>
			<item name="@name" type="string" array="false" required="true">聯系人的名稱</item>
			<item name="url" type="string" array="false" required="false">聯系人的 URL</item>
			<item name="email" type="string" array="false" required="false">聯系人的電子郵件</item>
		</type>
		<type name="link">
			<usage>用於描述鏈接信息，壹般轉換為 HTML 的 <code>a</code> 標簽。</usage>
			<item name="@text" type="string" array="false" required="true">鏈接的字面文字</item>
			<item name="@url" type="string" array="false" required="true">鏈接指向的文本</item>
		</type>
		<type name="tag">
			<usage>用於對各個 API 進行分類</usage>
			<item name="@name" type="string" array="false" required="true">標簽的唯壹 ID</item>
			<item name="@title" type="string" array="false" required="true">標簽的字面名稱</item>
			<item name="@deprecated" type="version" array="false" required="false">該標簽在大於該版本時被棄用</item>
		</type>
		<type name="server">
			<usage>用於指定各個 API 的服務器地址</usage>
			<item name="@name" type="string" array="false" required="true">服務唯壹 ID</item>
			<item name="@url" type="string" array="false" required="true">服務的基地址，與該服務關聯的 API，訪問地址都是相對於此地址的。</item>
			<item name="@deprecated" type="version" array="false" required="false">服務在大於該版本時被棄用</item>
			<item name="@summary" type="string" array="false" required="false">服務的摘要信息</item>
			<item name="description" type="richtext" array="false" required="false">服務的詳細描述</item>
		</type>
		<type name="api">
			<usage>用於定義單個 API 接口的具體內容</usage>
			<item name="@version" type="version" array="false" required="false">表示此接口在該版本中添加</item>
			<item name="@method" type="string" array="false" required="true">當前接口所支持的請求方法</item>
			<item name="@id" type="string" array="false" required="false">接口的唯壹 ID</item>
			<item name="@summary" type="string" array="false" required="false">簡要介紹</item>
			<item name="@deprecated" type="version" array="false" required="false">在此版本之後將會被棄用</item>
			<item name="path" type="path" array="false" required="true">定義路徑信息</item>
			<item name="description" type="richtext" array="false" required="false">該接口的詳細介紹，為 HTML 內容。</item>
			<item name="request" type="request" array="true" required="false">定義可用的請求信息</item>
			<item name="response" type="request" array="true" required="false">定義可能的返回信息</item>
			<item name="callback" type="callback" array="false" required="false">定義回調接口內容</item>
			<item name="header" type="param" array="true" required="false">傳遞的報頭內容，如果是某個 mimetype 專用的，可以放在 request 元素中。</item>
			<item name="tag" type="string" array="true" required="false">關聯的標簽</item>
			<item name="server" type="string" array="true" required="false">關聯的服務</item>
		</type>
		<type name="path">
			<usage>用於定義請求時與路徑相關的內容</usage>
			<item name="@path" type="string" array="false" required="true">接口地址</item>
			<item name="param" type="param" array="true" required="false">地址中的參數</item>
			<item name="query" type="param" array="true" required="false">地址中的查詢參數</item>
		</type>
		<type name="param">
			<usage>參數類型，基本上可以作為 request 的子集使用。</usage>
			<item name="@xml-attr" type="bool" array="false" required="false">是否作為父元素的屬性，僅作用於 XML 元素。是否作為父元素的屬性，僅用於 XML 的請求。</item>
			<item name="@xml-extract" type="bool" array="false" required="false">將當前元素的內容作為父元素的內容，要求父元素必須為 <var>object</var>。</item>
			<item name="@xml-cdata" type="bool" array="false" required="false">當前內容為 CDATA，与 xml-attr 互斥。</item>
			<item name="@xml-ns" type="string" array="false" required="false">XML 標簽的命名空間</item>
			<item name="@xml-ns-prefix" type="string" array="false" required="false">XML 標簽的命名空間名稱前綴</item>
			<item name="@xml-wrapped" type="string" array="false" required="false">如果當前元素的 <code>@array</code> 為 <var>true</var>，是否將其包含在 wrapped 指定的標簽中。</item>
			<item name="@name" type="string" array="false" required="true">值的名稱</item>
			<item name="@type" type="string" array="false" required="true">值的類型，可以是 <var>string</var>、<var>number</var>、<var>bool</var> 和 <var>object</var></item>
			<item name="@deprecated" type="version" array="false" required="false">表示在大於等於該版本號時不再啟作用</item>
			<item name="@default" type="string" array="false" required="false">默認值</item>
			<item name="@optional" type="bool" array="false" required="false">是否為可選的參數</item>
			<item name="@array" type="bool" array="false" required="false">是否為數組</item>
			<item name="@summary" type="string" array="false" required="false">簡要介紹</item>
			<item name="@array-style" type="bool" array="false" required="false">以數組的方式展示數據</item>
			<item name="param" type="param" array="true" required="false">子類型，比如對象的子元素。</item>
			<item name="enum" type="enum" array="true" required="false">當前參數可用的枚舉值</item>
			<item name="description" type="richtext" array="false" required="false">詳細介紹，為 HTML 內容。</item>
		</type>
		<type name="enum">
			<usage>定義枚舉類型的數所的枚舉值</usage>
			<item name="@deprecated" type="version" array="false" required="false">該屬性棄用的版本號</item>
			<item name="@value" type="string" array="false" required="true">枚舉值</item>
			<item name="@summary" type="string" array="false" required="false">枚舉值的說明</item>
			<item name="description" type="richtext" array="false" required="false">枚舉值的詳細說明</item>
		</type>
		<type name="request">
			<usage>定義了請求和返回的相關內容</usage>
			<item name="@xml-attr" type="bool" array="false" required="false">是否作為父元素的屬性，僅作用於 XML 元素。是否作為父元素的屬性，僅用於 XML 的請求。</item>
			<item name="@xml-extract" type="bool" array="false" required="false">將當前元素的內容作為父元素的內容，要求父元素必須為 <var>object</var>。</item>
			<item name="@xml-cdata" type="bool" array="false" required="false">當前內容為 CDATA，与 xml-attr 互斥。</item>
			<item name="@xml-ns" type="string" array="false" required="false">XML 標簽的命名空間</item>
			<item name="@xml-ns-prefix" type="string" array="false" required="false">XML 標簽的命名空間名稱前綴</item>
			<item name="@xml-wrapped" type="string" array="false" required="false">如果當前元素的 <code>@array</code> 為 <var>true</var>，是否將其包含在 wrapped 指定的標簽中。</item>
			<item name="@name" type="string" array="false" required="false">當 mimetype 為 <var>application/xml</var> 時，此值表示 XML 的頂層元素名稱，否則無用。</item>
			<item name="@type" type="string" array="false" required="false">值的類型，可以是 <var>string</var>、<var>number</var>、<var>bool</var>、<var>object</var> 和空值；空值表示不輸出任何內容。</item>
			<item name="@deprecated" type="version" array="false" required="false">表示在大於等於該版本號時不再啟作用</item>
			<item name="@array" type="bool" array="false" required="false">是否為數組</item>
			<item name="@summary" type="string" array="false" required="false">簡要介紹</item>
			<item name="@status" type="number" array="false" required="false">狀態碼。在 request 中，該值不可用，否則為必填項。</item>
			<item name="@mimetype" type="string" array="false" required="false">媒體類型，比如 <var>application/json</var> 等。</item>
			<item name="enum" type="enum" array="true" required="false">當前參數可用的枚舉值</item>
			<item name="param" type="param" array="true" required="false">子類型，比如對象的子元素。</item>
			<item name="example" type="example" array="true" required="false">示例代碼</item>
			<item name="header" type="param" array="true" required="false">傳遞的報頭內容</item>
			<item name="description" type="richtext" array="false" required="false">詳細介紹，為 HTML 內容。</item>
		</type>
		<type name="example">
			<usage>示例代碼</usage>
			<item name="@mimetype" type="string" array="false" required="true">特定於類型的示例代碼</item>
			<item name="@summary" type="string" array="false" required="false">示例代碼的概要信息</item>
			<item name="." type="string" array="false" required="true">示例代碼的內容，需要使用 CDATA 包含代碼。</item>
		</type>
		<type name="callback">
			<usage>定義接口的回調內容</usage>
			<item name="@method" type="string" array="false" required="true">回調的請求方法</item>
			<item name="@summary" type="string" array="false" required="false">簡要介紹</item>
			<item name="@deprecated" type="version" array="false" required="false">在此版本之後將會被棄用</item>
			<item name="path" type="path" array="false" required="false">回調的請求地址</item>
			<item name="description" type="richtext" array="false" required="false">對於回調的詳細介紹</item>
			<item name="response" type="request" array="true" required="false">定義可能的返回信息</item>
			<item name="request" type="request" array="true" required="true">定義可用的請求信息</item>
			<item name="header" type="param" array="true" required="false">傳遞的報頭內容</item>
		</type>
		<type name="string">
			<usage>普通的字符串類型</usage>
		</type>
		<type name="date">
			<usage>采用 <a href="https://tools.ietf.org/html/rfc3339">RFC3339</a> 格式表示的時間，比如：<samp>2019-12-16T00:35:48+08:00</samp></usage>
		</type>
		<type name="version">
			<usage>版本號，格式遵守 <a href="https://semver.org/lang/zh-TW/">semver</a> 規則</usage>
		</type>
		<type name="bool">
			<usage>布爾值類型，取值為 <var>true</var> 或是 <var>false</var></usage>
		</type>
		<type name="number">
			<usage>普通的數值類型</usage>
		</type>
	</spec>
	<commands>
		<command name="build">生成文檔內容</command>
		<command name="detect">根據目錄下的內容生成配置文件</command>
		<command name="help">顯示幫助信息</command>
		<command name="lang">顯示所有支持的語言</command>
		<command name="locale">顯示所有支持的本地化內容</command>
		<command name="lsp">啟動 language server protocol 服務</command>
		<command name="mock">啟用 mock 服務</command>
		<command name="static">啟用靜態文件服務</command>
		<command name="test">測試語法的正確性</command>
		<command name="version">顯示版本信息</command>
	</commands>
	<config>
		<item name="version" type="string" array="false" required="true">此配置文件的所使用的文档版本</item>
		<item name="inputs" type="object" array="true" required="true">指定輸入的數據，同壹項目只能解析壹種語言。</item>
		<item name="inputs.lang" type="string" array="false" required="true">源文件類型。具體支持的類型可通過 -l 參數進行查找。</item>
		<item name="inputs.dir" type="string" array="false" required="true">需要解析的源文件所在目錄</item>
		<item name="inputs.exts" type="string" array="true" required="false">只從這些擴展名的文件中查找文檔</item>
		<item name="inputs.recursive" type="bool" array="false" required="false">是否解析子目錄下的源文件</item>
		<item name="inputs.encoding" type="string" array="false" required="false">編碼，默認為 <var>utf-8</var>，值可以是 <a href="https://www.iana.org/assignments/character-sets/character-sets.xhtml">character-sets</a> 中的內容。</item>
		<item name="output" type="object" array="false" required="true">控制輸出行為</item>
		<item name="output.type" type="string" array="false" required="false">輸出的類型，目前可以 <var>apidoc+xml</var>、<var>openapi+json</var> 和 <var>openapi+yaml</var>。</item>
		<item name="output.path" type="string" array="false" required="true">指定輸出的文件名，包含路徑信息。</item>
		<item name="output.tags" type="string" array="true" required="false">只輸出與這些標簽相關聯的文檔，默認為全部。</item>
		<item name="output.style" type="string" array="false" required="false">為 XML 文件指定的 XSL 文件</item>
	</config>
</locale>
`),
	},
	{
		Name:        "site.xml",
		ContentType: "application/xml; charset=utf-8",
		Content: []byte(`<?xml version="1.0" encoding="UTF-8"?>

<!-- 该文件由工具自动生成，请勿手动修改！ -->

<site>
	<name>apidoc</name>
	<version>6.1.0</version>
	<repo>https://github.com/caixw/apidoc</repo>
	<url>https://apidoc.tools</url>
	<languages>
		<language id="c#">C#</language>
		<language id="c++">C/C++</language>
		<language id="d">D</language>
		<language id="erlang">Erlang</language>
		<language id="go">Go</language>
		<language id="groovy">Groovy</language>
		<language id="java">Java</language>
		<language id="javascript">JavaScript</language>
		<language id="kotlin">Kotlin</language>
		<language id="pascal">Pascal/Delphi</language>
		<language id="perl">Perl</language>
		<language id="php">PHP</language>
		<language id="python">Python</language>
		<language id="ruby">Ruby</language>
		<language id="rust">Rust</language>
		<language id="scala">Scala</language>
		<language id="swift">Swift</language>
	</languages>
	<locales>
		<locale id="cmn-Hans" href="index.xml" title="简体中文" doc="locale.cmn-Hans.xml"></locale>
		<locale id="cmn-Hant" href="index.cmn-Hant.xml" title="繁體中文" doc="locale.cmn-Hant.xml"></locale>
	</locales>
</site>
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

<xsl:template match="processing-instruction('xml-stylesheet')[1]">
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

table tr {
    cursor: pointer;
}

table th, table td {
    font-weight: normal;
    text-align: left;
    border-bottom: 1px solid transparent;
}

table tr:hover th,
table tr:hover td {
    opacity: .7;
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
    cursor: pointer;
}

header .menu ul li:hover label {
    opacity: .7;
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
    opacity: .8;
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
}

.goto-top {
    border: solid var(--color);
    border-width: 0 5px 5px 0;
    display: block;
    padding: 5px;
    transform: rotate(-135deg);
    position: fixed;
    bottom: 3rem;
    right: calc(((100% - var(--max-width)) / 2) - 3rem);
}

/*
 * 用于计算 goto-top 按钮的位置始终保持在内容主体的右侧边上。
 * 2048px 为 --max-width 的值，但是 CSS 并不支持直接在 @media
 * 使用 var() 函数。
 */
@media (max-width: calc(2048px + 6rem)) {
    .goto-top {
        right: 3rem;
    }
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
    initGotoTop();
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
            <p><xsl:copy-of select="$locale-generator" /></p>
        </div>
        <a href="#" class="goto-top" title="{$locale-goto-top}" aria-label="{$locale-goto-top}" />
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
        <img alt="logo" src="{$icon}" />
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
<xsl:variable name="id" select="concat(server, @method, translate(path/@path, $id-from, $id-to))" />

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
                <xsl:with-param name="headers" select="header | /apidoc/header" />
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
<div>
<xsl:choose>
    <xsl:when test="$response/description">
        <xsl:attribute name="data-type">
            <xsl:value-of select="$response/description/@type" />
        </xsl:attribute>
        <pre><xsl:copy-of select="$response/description/node()" /></pre>
    </xsl:when>
    <xsl:otherwise><xsl:value-of select="$response/@summary" /></xsl:otherwise>
</xsl:choose>
</div>
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

<xsl:template match="processing-instruction('xml-stylesheet')[1]">
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
<!-- NOTE: id 一律用小写，之后需要用到比较字符串！ -->
<l:locales>
    <locale id="cmn-hans">简体中文</locale>
    <locale id="cmn-hant">繁體中文</locale>
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
        <xsl:with-param name="lang" select="'cmn-hans'" />
        <xsl:with-param name="text" select="document('')/xsl:stylesheet/l:locales/locale[@id='cmn-hans']" />
    </xsl:call-template>

    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'cmn-hant'" />
        <xsl:with-param name="text" select="document('')/xsl:stylesheet/l:locales/locale[@id='cmn-hant']" />
    </xsl:call-template>
</xsl:variable>

<!-- server -->
<xsl:variable name="locale-server">
    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'cmn-hans'" />
        <xsl:with-param name="text" select="'服务'" />
    </xsl:call-template>

    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'cmn-hant'" />
        <xsl:with-param name="text" select="'服務'" />
    </xsl:call-template>
</xsl:variable>

<!-- tag -->
<xsl:variable name="locale-tag">
    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'cmn-hans'" />
        <xsl:with-param name="text" select="'标签'" />
    </xsl:call-template>

    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'cmn-hant'" />
        <xsl:with-param name="text" select="'標簽'" />
    </xsl:call-template>
</xsl:variable>

<!-- expand -->
<xsl:variable name="locale-expand">
    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'cmn-hans'" />
        <xsl:with-param name="text" select="'展开'" />
    </xsl:call-template>

    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'cmn-hant'" />
        <xsl:with-param name="text" select="'展開'" />
    </xsl:call-template>
</xsl:variable>

<!-- method -->
<xsl:variable name="locale-method">
    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'cmn-hans'" />
        <xsl:with-param name="text" select="'请求方法'" />
    </xsl:call-template>

    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'cmn-hant'" />
        <xsl:with-param name="text" select="'請求方法'" />
    </xsl:call-template>
</xsl:variable>

<!-- request -->
<xsl:variable name="locale-request">
    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'cmn-hans'" />
        <xsl:with-param name="text" select="'请求'" />
    </xsl:call-template>

    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'cmn-hant'" />
        <xsl:with-param name="text" select="'請求'" />
    </xsl:call-template>
</xsl:variable>

<!-- response -->
<xsl:variable name="locale-response">
    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'cmn-hans'" />
        <xsl:with-param name="text" select="'返回'" />
    </xsl:call-template>

    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'cmn-hant'" />
        <xsl:with-param name="text" select="'返回'" />
    </xsl:call-template>
</xsl:variable>

<!-- callback -->
<xsl:variable name="locale-callback">
    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'cmn-hans'" />
        <xsl:with-param name="text" select="'回调'" />
    </xsl:call-template>

    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'cmn-hant'" />
        <xsl:with-param name="text" select="'回調'" />
    </xsl:call-template>
</xsl:variable>

<!-- path param -->
<xsl:variable name="locale-path-param">
    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'cmn-hans'" />
        <xsl:with-param name="text" select="'路径参数'" />
    </xsl:call-template>

    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'cmn-hant'" />
        <xsl:with-param name="text" select="'路徑參數'" />
    </xsl:call-template>
</xsl:variable>

<!-- query -->
<xsl:variable name="locale-query">
    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'cmn-hans'" />
        <xsl:with-param name="text" select="'查询参数'" />
    </xsl:call-template>

    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'cmn-hant'" />
        <xsl:with-param name="text" select="'查詢參數'" />
    </xsl:call-template>
</xsl:variable>

<!-- header -->
<xsl:variable name="locale-header">
    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'cmn-hans'" />
        <xsl:with-param name="text" select="'报头'" />
    </xsl:call-template>

    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'cmn-hant'" />
        <xsl:with-param name="text" select="'報頭'" />
    </xsl:call-template>
</xsl:variable>

<!-- body -->
<xsl:variable name="locale-body">
    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'cmn-hans'" />
        <xsl:with-param name="text" select="'报文'" />
    </xsl:call-template>

    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'cmn-hant'" />
        <xsl:with-param name="text" select="'報文'" />
    </xsl:call-template>
</xsl:variable>

<!-- example -->
<xsl:variable name="locale-example">
    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'cmn-hans'" />
        <xsl:with-param name="text" select="'示例代码'" />
    </xsl:call-template>

    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'cmn-hant'" />
        <xsl:with-param name="text" select="'示例代碼'" />
    </xsl:call-template>
</xsl:variable>

<!-- var -->
<xsl:variable name="locale-var">
    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'cmn-hans'" />
        <xsl:with-param name="text" select="'变量'" />
    </xsl:call-template>

    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'cmn-hant'" />
        <xsl:with-param name="text" select="'變量'" />
    </xsl:call-template>
</xsl:variable>

<!-- type -->
<xsl:variable name="locale-type">
    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'cmn-hans'" />
        <xsl:with-param name="text" select="'类型'" />
    </xsl:call-template>

    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'cmn-hant'" />
        <xsl:with-param name="text" select="'類型'" />
    </xsl:call-template>
</xsl:variable>

<!-- value -->
<xsl:variable name="locale-value">
    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'cmn-hans'" />
        <xsl:with-param name="text" select="'值'" />
    </xsl:call-template>

    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'cmn-hant'" />
        <xsl:with-param name="text" select="'值'" />
    </xsl:call-template>
</xsl:variable>

<!-- goto-top -->
<xsl:variable name="locale-goto-top">
    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'cmn-hans'" />
        <xsl:with-param name="text" select="'返回顶部'" />
    </xsl:call-template>

    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'cmn-hant'" />
        <xsl:with-param name="text" select="'返回頂部'" />
    </xsl:call-template>
</xsl:variable>

<!-- description -->
<xsl:variable name="locale-description">
    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'cmn-hans'" />
        <xsl:with-param name="text" select="'描述'" />
    </xsl:call-template>

    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'cmn-hant'" />
        <xsl:with-param name="text" select="'描述'" />
    </xsl:call-template>
</xsl:variable>

<!-- enum -->
<xsl:variable name="locale-enum">
    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'cmn-hans'" />
        <xsl:with-param name="text" select="'枚举'" />
    </xsl:call-template>

    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'cmn-hant'" />
        <xsl:with-param name="text" select="'枚舉'" />
    </xsl:call-template>
</xsl:variable>

<!-- generator -->
<xsl:variable name="locale-generator">
    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'cmn-hans'" />
        <xsl:with-param name="text">
            <xsl:if test="apidoc/license">
            文档版权为 <a href="{apidoc/license/@url}"><xsl:value-of select="apidoc/license/@text" /></a>。
            </xsl:if>
            包含了 <xsl:value-of select="count(apidoc/api)" /> 个接口声明，由 <a href="https://apidoc.tools">apidoc</a> 生成于 <time><xsl:value-of select="apidoc/@created" /></time>。
        </xsl:with-param>
    </xsl:call-template>

    <xsl:call-template name="build-locale">
        <xsl:with-param name="lang" select="'cmn-hant'" />
        <xsl:with-param name="text">
            <xsl:if test="apidoc/license">
            文檔版權為 <a href="{apidoc/license/@url}"><xsl:value-of select="apidoc/license/@text" /></a>。
            </xsl:if>
            包含了 <xsl:value-of select="count(apidoc/api)" /> 個接口聲明，由 <a href="https://apidoc.tools">apidoc</a> 生成於 <time><xsl:value-of select="apidoc/@created" /></time>。
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
