// 由 make.go 生成，请勿修改！

package html

var data = []*static{
	{
		name:        "apidoc.xsl",
		contentType: "application/xslt+xml",
		data: []byte(`<?xml version="1.0" encoding="UTF-8"?>

<xsl:stylesheet version="1.0" xmlns:xsl="http://www.w3.org/1999/XSL/Transform">
<xsl:output
    method="html"
    encoding="utf-8"
    indent="yes"
    version="5.0"
    doctype-system="about:legacy-compat" />

<!-- 替换字符串中特定的字符 -->
<xsl:template name="replace">
    <xsl:param name="text" />
    <xsl:param name="old" />
    <xsl:param name="new" />
    <xsl:choose>
        <xsl:when test="contains($text, $old)">
            <xsl:value-of select="substring-before($text, $old)" />
            <xsl:value-of select="$new" />
            <xsl:call-template name="replace">
                <xsl:with-param name="text" select="substring-after($text, $old)" />
                <xsl:with-param name="old" select="$old" />
                <xsl:with-param name="new" select="$new" />
            </xsl:call-template>
        </xsl:when>
        <xsl:otherwise>
            <xsl:value-of select="$text" />
        </xsl:otherwise>
    </xsl:choose>
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
            <xsl:value-of select="'弃用于 '" />
            <xsl:value-of select="$deprecated" />
        </xsl:attribute>
    </xsl:if>
</xsl:template>

<!-- 根据 method 和 path 生成唯一的 ID -->
<xsl:template name="get-api-id">
    <xsl:param name="method" />
    <xsl:param name="path" />

    <xsl:variable name="p1">
        <xsl:call-template name="replace">
            <xsl:with-param name="text" select="$path" />
            <xsl:with-param name="old" select="'/'" />
            <xsl:with-param name="new" select="'_'" />
        </xsl:call-template>
    </xsl:variable>
    <xsl:variable name="p2">
        <xsl:call-template name="replace">
            <xsl:with-param name="text" select="$p1" />
            <xsl:with-param name="old" select="'{'" />
            <xsl:with-param name="new" select="'-'" />
        </xsl:call-template>
    </xsl:variable>

    <xsl:value-of select="$method" />
    <xsl:call-template name="replace">
        <xsl:with-param name="text" select="$p2" />
        <xsl:with-param name="old" select="'}'" />
        <xsl:with-param name="new" select="'-'" />
    </xsl:call-template>
</xsl:template>

<!-- path param, path query, header 等的界面 -->
<xsl:template name="param">
    <xsl:param name="title" />
    <xsl:param name="param" />
    <xsl:param name="example" /> <!-- 示例代码 -->

    <div class="param">
        <h4 class="title">&#x27a4;&#160;<xsl:value-of select="$title" /></h4>
        <table>
            <thead>
                <tr>
                    <th>变量</th>
                    <th>类型</th>
                    <th title="是否为必填，以及默认值。">值</th>
                    <th>描述</th>
                </tr>
            </thead>
            <tbody>
                <xsl:call-template name="param-list">
                <xsl:with-param name="param" select="$param" />
                </xsl:call-template>
            </tbody>
        </table>
    </div>
</xsl:template>

<!-- 列顺序必须要与 param 中的相同 -->
<xsl:template name="param-list">
    <xsl:param name="param" />
    <xsl:param name="parent" /> <!-- 上一级的名称，嵌套对象时可用 -->

    <xsl:for-each select="$param">
    <tr>
        <th>
            <xsl:if test="$parent">
                <xsl:value-of select="$parent" />
                <xsl:value-of select="'.'" />
            </xsl:if>
            <xsl:value-of select="@name" />
        </th>

        <td>
            <xsl:call-template name="deprecated">
                <xsl:with-param name="deprecated" select="@deprecated" />
            </xsl:call-template>

            <xsl:value-of select="@type" />
            <xsl:if test="@array = 'true'">
                <xsl:value-of select="'[]'" />
            </xsl:if>
        </td>

        <td>
            <xsl:choose>
                <xsl:when test="@required = 'true'"><xsl:value-of select="'R'" /></xsl:when>
                <xsl:otherwise><xsl:value-of select="'O'" /></xsl:otherwise>
            </xsl:choose>

            <xsl:if test="@default">
                <xsl:value-of select="' '" />
                <xsl:value-of select="@default" />
            </xsl:if>
        </td>

        <td>
            <xsl:choose>
                <xsl:when test="description"><xsl:value-of select="description" /></xsl:when>
                <xsl:otherwise><xsl:value-of select="@summary" /></xsl:otherwise>
            </xsl:choose>
            <xsl:if test="./enum">
                <p>可以使用以下枚举值：</p>
                <ul>
                <xsl:for-each select="./enum">
                    <li>
                    <xsl:call-template name="deprecated">
                        <xsl:with-param name="deprecated" select="@deprecated" />
                    </xsl:call-template>

                    <xsl:value-of select="@value" />:<xsl:value-of select="." />
                    </li>
                </xsl:for-each>
                </ul>
            </xsl:if>
        </td>
    </tr>

    <xsl:if test="./param">
        <xsl:call-template name="param-list">
            <xsl:with-param name="param" select="./param" />
            <xsl:with-param name="parent" select="./@name" />
        </xsl:call-template>
    </xsl:if>
    </xsl:for-each>
</xsl:template>

<!-- api/request 的介面元素 -->
<xsl:template name="request">
<xsl:param name="request" />
<xsl:param name="path" />
<div class="request">
    <xsl:if test="$path/param">
        <xsl:call-template name="param">
            <xsl:with-param name="title" select="'路径参数'" />
            <xsl:with-param name="param" select="$path/param" />
        </xsl:call-template>
    </xsl:if>

    <xsl:if test="$path/query">
        <xsl:call-template name="param">
            <xsl:with-param name="title" select="'查询参数'" />
            <xsl:with-param name="param" select="$path/query" />
        </xsl:call-template>
    </xsl:if>
    
    <xsl:if test="$request/header">
        <xsl:call-template name="param">
            <xsl:with-param name="title" select="'请求报头'" />
            <xsl:with-param name="param" select="$request/header" />
        </xsl:call-template>
    </xsl:if>

    <xsl:call-template name="param">
        <xsl:with-param name="title" select="'请求报文'" />
        <xsl:with-param name="param" select="$request" />
    </xsl:call-template>
</div>
</xsl:template>

<!-- api/response 的界面 -->
<xsl:template name="response">
    <xsl:param name="response" />

    <h5 class="status"><xsl:value-of select="$response/@status" /></h5>

    <xsl:if test="$response/header">
        <xsl:call-template name="param">
            <xsl:with-param name="title" select="'返回报头'" />
            <xsl:with-param name="param" select="$response/header" />
        </xsl:call-template>
    </xsl:if>

    <xsl:call-template name="param">
        <xsl:with-param name="title" select="'返回报文'" />
        <xsl:with-param name="param" select="$response" />
    </xsl:call-template>
</xsl:template>

<!-- api 界面元素 -->
<xsl:template match="/apidoc/api">
    <details class="api">
    <xsl:attribute name="data-method">
        <xsl:value-of select="@method" />
    </xsl:attribute>

    <xsl:attribute name="data-tags">
        <xsl:for-each select="tag">
            <xsl:value-of select="." />
            <xsl:value-of select="','" />
        </xsl:for-each>
    </xsl:attribute>

        <summary>
            <a class="link"> <!-- 链接符号 -->
            <xsl:attribute name="href">
                <xsl:value-of select="'#'" />
                <xsl:call-template name="get-api-id">
                    <xsl:with-param name="path" select="path/@path" />
                    <xsl:with-param name="method" select="@method" />
                </xsl:call-template>
            </xsl:attribute>
            &#128279;
            </a>

            <span class="action"><xsl:value-of select="@method" /></span>
            <span>
                <xsl:call-template name="deprecated">
                    <xsl:with-param name="deprecated" select="@deprecated" />
                </xsl:call-template>

                <xsl:value-of select="path/@path" />
            </span>

            <span class="summary">
            <xsl:value-of select="@summary" />
            </span>
        </summary>
        <div class="description">
            <xsl:if test="./description">
            <xsl:value-of select="./description" />
            </xsl:if>
        </div>

        <div class="body">
            <div class="requests">
                <h4 class="title">请求</h4>
                <xsl:for-each select="request">
                    <xsl:call-template name="request">
                        <xsl:with-param name="request" select="." />
                        <xsl:with-param name="path" select="../path" />
                    </xsl:call-template>
                </xsl:for-each>
            </div>
            <div class="responses">
                <h4 class="title">返回</h4>
                <xsl:for-each select="response">
                    <xsl:call-template name="response">
                        <xsl:with-param name="response" select="." />
                    </xsl:call-template>
                </xsl:for-each>
            </div>
        </div>

        <xsl:if test="./callback">
        <div class="callback">
            <xsl:attribute name="data-method">
                <xsl:value-of select="./callback/@method" />
            </xsl:attribute>
            <h3>回调</h3>
            <div class="description">
                <xsl:value-of select="./callback/@summary" />
                <xsl:if test="./callback/description">
                    <br />
                    <xsl:value-of select="./callback/description" />
                </xsl:if>
            </div>

            <div class="body">
                <div class="requests">
                    <h4 class="title">请求</h4>
                    <xsl:for-each select="./callback/request">
                        <xsl:call-template name="request">
                            <xsl:with-param name="request" select="." />
                            <xsl:with-param name="path" select="../path" />
                        </xsl:call-template>
                    </xsl:for-each>
                </div>

                <xsl:if test="./callback/response">
                    <div class="responses">
                        <h4 class="title">返回</h4>
                        <xsl:for-each select="./callback/response">
                            <xsl:call-template name="response">
                                <xsl:with-param name="response" select="." />
                            </xsl:call-template>
                        </xsl:for-each>
                    </div>
                </xsl:if>
            </div> <!-- end .body -->
        </div> <!-- end .callback -->
        </xsl:if>
    </details>
</xsl:template>

<!-- header 界面元素 -->
<xsl:template name="header">
    <header>
        <h1>
            <img src="./icon.svg" />
            <xsl:value-of select="/apidoc/title" />
            <span class="version">(<xsl:value-of select="/apidoc/@version" />)</span>
        </h1>

        <div class="menu tags-selector">
            <h2>标签</h2>
            <ul>
                <xsl:for-each select="apidoc/tag">
                <li>
                    <xsl:attribute name="data-tag">
                        <xsl:value-of select="@name" />
                    </xsl:attribute>
                    <label><input type="checkbox" checked="checked" /><xsl:value-of select="@title" /></label>
                </li>
                </xsl:for-each>
            </ul>
        </div>

        <div class="menu methods-selector">
            <h2>请求方法</h2>
            <ul>
                <!-- 浏览器好像都不支持 xpath 2.0，所以无法使用 distinct-values -->
                <!-- xsl:for-each select="distinct-values(/apidoc/api/@method)" -->
                <xsl:for-each select="/apidoc/api/@method[not(../preceding-sibling::api/@method = .)]">
                <li>
                    <xsl:attribute name="data-method">
                        <xsl:value-of select="." />
                    </xsl:attribute>
                    <label><input type="checkbox" checked="true" /><xsl:value-of select="." /></label>
                </li>
                </xsl:for-each>
            </ul>
        </div>
    </header>
</xsl:template>

<xsl:template match="/">
    <html>
        <head>
            <title><xsl:value-of select="apidoc/title" /></title>
            <meta charset="UTF-8" />
            <meta name="viewport" content="width=device-width,initial-scale=1,maximum-scale=1" />
            <meta name="generator" content="https://apidoc.tools" />
            <link rel="stylesheet" type="text/css" href="./apidoc.css" />
            <link rel="icon" type="image/png" href="./icon.png" />
            <link rel="license">
                <xsl:attribute name="href"><xsl:value-of select="apidoc/license/@url"/></xsl:attribute>
            </link>
            <script src="./apidoc.js"></script>
        </head>
        <body>
            <xsl:call-template name="header" />

            <main>
                <div class="content">
                    <xsl:value-of select="/apidoc/content" />
                </div>

                <xsl:for-each select="apidoc/api">
                <xsl:sort select="path/@path"/>
                    <xsl:apply-templates select="." />
                </xsl:for-each>
            </main>

            <footer>
            <p>文档版权为
            <a>
                <xsl:attribute name="href"><xsl:value-of select="apidoc/license/@url" /></xsl:attribute>
                <xsl:value-of select="apidoc/license" />
            </a>。
            由 <a href="https://github.com/caixw/apidoc">apidoc</a> 生成。
            </p>
            </footer>
        </body>
    </html>
</xsl:template>
</xsl:stylesheet>
`),
	},
	{
		name:        "apidoc.css",
		contentType: "text/css",
		data: []byte(`@charset "utf-8";

:root {
    --max-width: 100%;
    --header-height: 54px;
    --border-color: #e0e0e0;
    --padding: 1rem;
    --article-padding: calc(var(--padding) / 2);

    --background: white;
    --delete-color: red;

    /* method */
    --method-get-color: green;
    --method-options-color: green;
    --method-post-color: darkorange;
    --method-put-color: darkorange;
    --method-patch-color: darkorange;
    --method-delete-color: red;
}

html {
    height: 100%;
}

body {
    padding: 0;
    margin: 0;
    height: 100%;
    background: var(--background);
    text-align: center;
}

table {
    width: 100%;
}

table th, table td {
    font-weight: normal;
    text-align: left;
    border-bottom: 1px solid var(--border-color);
}

table caption {
    text-align: left;
}

ul, ol {
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
}

.del {
    text-decoration: line-through;
    text-decoration-color: var(--delete-color);
}

/*************************** header ***********************/

header {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    display: block;
    z-index: 1000;
    padding: var(--padding);
    padding-bottom: var(--article-padding);
    height: var(--header-height);
    box-sizing: border-box;

    background: var(--background);
    border-bottom: 1px solid var(--border-color);

    margin: 0 auto;
    max-width: var(--max-width);
    text-align: left;
}

header h1, header h2 {
    margin: 0;
    display: inline-block;
}

header h1 .version {
    font-size: 1rem;
}

header h1 img {
    height: 1.5rem;
    margin-right: .5rem;
}

header .menu {
    float: right;
    margin-right: var(--padding);
    position: relative;
    margin-top: 9px;
}
header .menu:first-of-type {
    margin-right: 0;
}

header .menu h2 {
    font-size: 1rem;
    line-height: 1;
}

header .menu ul {
    position: absolute;
    min-width: 4rem;
    right: 0;
    display: none;
    list-style: none;
    background: var(--background);
    border-bottom: 1px solid var(--border-color);
    border-left: 1px solid var(--border-color);
    border-right: 1px solid var(--border-color);
    padding: var(--article-padding);
}

header .menu ul li {
    margin-top: var(--article-padding);
}

header .menu:hover ul {
    display: block;
}

/*************************** main ***********************/

main {
    padding: 0rem var(--padding);
    top: calc(var(--header-height) + var(--padding));
    left: 0;
    right: 0;
    position: relative;

    margin: 0 auto;
    max-width: var(--max-width);
    text-align: left;
}

main .content {
    padding: var(--article-padding);
}

main .api {
    margin-bottom: var(--article-padding);
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
}

main .api summary .link {
    margin-right: 10px;
    text-decoration: none;
}

main .api .description {
    padding:var(--article-padding);
    margin: 0;
    border-bottom: 1px solid var(--border-color);
}

main .api summary .summary {
    float: right;
    font-weight: 400;
    opacity: .5;
}

main .api[data-method=GET],
main .callback[data-method=GET] h3 {
    border: 1px solid var(--method-get-color);
}
main .api[data-method=GET] summary {
    border-bottom: 1px solid var(--method-get-color);
}

main .api[data-method=POST],
main .callback[data-method=POST] h3 {
    border: 1px solid var(--method-post-color);
}
main .api[data-method=POST] summary {
    border-bottom: 1px solid var(--method-post-color);
}

main .api[data-method=PUT],
main .callback[data-method=PUT] h3 {
    border: 1px solid var(--method-put-color);
}
main .api[data-method=PUT] summary {
    border-bottom: 1px solid var(--method-put-color);
}

main .api[data-method=PATCH],
main .callback[data-method=PATCH] h3 {
    border: 1px solid var(--method-patch-color);
}
main .api[data-method=PATCH] summary {
    border-bottom: 1px solid var(--method-patch-color);
}

main .api[data-method=DELETE],
main .callback[data-method=DELETE] h3 {
    border: 1px solid var(--method-delete-color);
}
main .api[data-method=DELETE] summary {
    border-bottom: 1px solid var(--method-delete-color);
}

main .api[data-method=OPTIONS],
main .callback[data-method=OPTIONS] h3 {
    border: 1px solid var(--method-options-color);
}
main .api[data-method=OPTIONS] summary {
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
}

main .api .body .requests,
main .api .body .responses {
    flex: 1 1 50%;
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

main .api .body .responses .status {
    margin: calc(var(--padding) + var(--article-padding)) 0 var(--article-padding);
    border-bottom: 1px solid var(--border-color);
}


/*************************** footer ***********************/

footer {
    margin-top: 4rem;
    padding: var(--padding);

    margin: 4rem auto;
    max-width: var(--max-width);
    text-align: left;
}
`),
	},
	{
		name:        "apidoc.js",
		contentType: "application/javascript",
		data: []byte(`'use strict';

window.onload = function () {
    this.registerMethodFilter();
    this.registerTagFilter();
}

function registerMethodFilter() {
    const list = document.querySelectorAll('.methods-selector li input');
    list.forEach((val) => {
        val.addEventListener('change', (event) => {
            const chk = event.target.checked;
            const method = event.target.parentNode.parentNode.getAttribute('data-method');

            const apis = this.document.querySelectorAll('.api');
            apis.forEach((api) => {
                if (api.getAttribute('data-method') != method) {
                    return;
                }

                api.style.display = chk ? 'block' : 'none';
            });
        });
    });
}

function registerTagFilter() {
    const list = document.querySelectorAll('.tags-selector li input');
    list.forEach((val) => {
        val.addEventListener('change', (event) => {
            const chk = event.target.checked;
            const tag = event.target.parentNode.parentNode.getAttribute('data-tag');

            const apis = this.document.querySelectorAll('.api');
            apis.forEach((api) => {
                if (!api.getAttribute('data-tags').includes(tag+',')) {
                    return;
                }

                api.style.display = chk ? 'block' : 'none';
            });
        });
    });
}
`),
	},
	{
		name:        "icon.svg",
		contentType: "image/svg+xml",
		data: []byte(`<?xml version="1.0" encoding="utf-8" ?>
<svg viewBox="0 0 512 512" xmlns="http://www.w3.org/2000/svg" stroke="red" fill="grey">
    <circle cx="256" cy="256" r="252" fill-opacity="0" stroke-width="8" />
    <circle cx="585" cy="395" r="350" fill-opacity="0" stroke-width="8" />
    <circle cx="-70" cy="395" r="350" fill-opacity="0" stroke-width="8" />
</svg>
`),
	},
}
