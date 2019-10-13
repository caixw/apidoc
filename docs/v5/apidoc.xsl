<?xml version="1.0" encoding="UTF-8"?>

<xsl:stylesheet version="1.0" xmlns:xsl="http://www.w3.org/1999/XSL/Transform">

<xsl:import href="./locales.xsl" />

<xsl:output
    method="html"
    encoding="utf-8"
    indent="yes"
    version="5.0"
    doctype-system="about:legacy-compat" />

<xsl:template match="/">
    <html>
    <xsl:if test="@lang">
        <xsl:attribute name="lang"><xsl:value-of select="@lang" /></xsl:attribute>
    </xsl:if>

        <head>
            <title><xsl:value-of select="apidoc/title" /></title>
            <meta charset="UTF-8" />
            <meta name="viewport" content="width=device-width,initial-scale=1,maximum-scale=1" />
            <meta name="generator" content="https://apidoc.tools" />
            <link rel="stylesheet" type="text/css" href="./apidoc.css" />
            <link rel="icon" type="image/png" href="../icon.png" />
            <link rel="license" href="{apidoc/license/@url}" />
            <script src="./apidoc.js"></script>
        </head>
        <body>
            <xsl:call-template name="header" />

            <main>
                <div class="content"><xsl:value-of select="/apidoc/content" /></div>

                <xsl:for-each select="apidoc/api">
                <xsl:sort select="path/@path"/>
                    <xsl:apply-templates select="." />
                </xsl:for-each>
            </main>

            <footer>
            <p><xsl:copy-of select="$locale-footer" /></p>
            </footer>
        </body>
    </html>
</xsl:template>

<!-- header 界面元素 -->
<xsl:template name="header">
    <header>
        <h1>
            <img src="../icon.svg" />
            <xsl:value-of select="/apidoc/title" />
            <span class="version">(<xsl:value-of select="/apidoc/@version" />)</span>
        </h1>

        <div class="menu tags-selector">
            <h2><xsl:copy-of select="$locale-tag" /></h2>
            <ul>
                <xsl:for-each select="apidoc/tag">
                <li data-tag="{@name}">
                    <label><input type="checkbox" checked="checked" /><xsl:value-of select="@title" /></label>
                </li>
                </xsl:for-each>
            </ul>
        </div>

        <div class="menu methods-selector">
            <h2><xsl:copy-of select="$locale-method" /></h2>
            <ul>
                <!-- 浏览器好像都不支持 xpath 2.0，所以无法使用 distinct-values -->
                <!-- xsl:for-each select="distinct-values(/apidoc/api/@method)" -->
                <xsl:for-each select="/apidoc/api/@method[not(../preceding-sibling::api/@method = .)]">
                <li data-method="{.}">
                    <label><input type="checkbox" checked="true" /><xsl:value-of select="." /></label>
                </li>
                </xsl:for-each>
            </ul>
        </div>

        <div class="menu languages-selector">
            <h2><xsl:copy-of select="$locale-language" /></h2>
            <ul><xsl:call-template name="languages" /></ul>
        </div>
    </header>
</xsl:template>

<!-- api 界面元素 -->
<xsl:template match="/apidoc/api">
    <xsl:variable name="id">
        <xsl:call-template name="get-api-id">
            <xsl:with-param name="path" select="path/@path" />
            <xsl:with-param name="method" select="@method" />
        </xsl:call-template>
    </xsl:variable>

    <details id="{$id}" class="api" data-method="{@method}">
    <xsl:attribute name="data-tags">
        <xsl:for-each select="tag"><xsl:value-of select="concat(., ',')" /></xsl:for-each>
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
                <h4 class="title"><xsl:copy-of select="$locale-request" /></h4>
                <xsl:for-each select="request">
                    <xsl:call-template name="request">
                        <xsl:with-param name="request" select="." />
                        <xsl:with-param name="path" select="../path" />
                    </xsl:call-template>
                </xsl:for-each>
            </div>
            <div class="responses">
                <h4 class="title"><xsl:copy-of select="$locale-response" /></h4>
                <xsl:for-each select="response">
                    <xsl:call-template name="response">
                        <xsl:with-param name="response" select="." />
                    </xsl:call-template>
                </xsl:for-each>
            </div>
        </div>

        <xsl:if test="./callback">
        <div class="callback" data-method="{./callback/@method}">
            <h3><xsl:copy-of select="$locale-callback" /></h3>
            <div class="description">
                <xsl:value-of select="./callback/@summary" />
                <xsl:if test="./callback/description">
                    <br />
                    <xsl:value-of select="./callback/description" />
                </xsl:if>
            </div>

            <div class="body">
                <div class="requests">
                    <h4 class="title"><xsl:copy-of select="$locale-request" /></h4>
                    <xsl:for-each select="./callback/request">
                        <xsl:call-template name="request">
                            <xsl:with-param name="request" select="." />
                            <xsl:with-param name="path" select="../path" />
                        </xsl:call-template>
                    </xsl:for-each>
                </div>

                <xsl:if test="./callback/response">
                    <div class="responses">
                        <h4 class="title"><xsl:copy-of select="$locale-response" /></h4>
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

<!-- api/request 的介面元素 -->
<xsl:template name="request">
<xsl:param name="request" />
<xsl:param name="path" />
<div class="request">
    <xsl:if test="$path/param">
        <xsl:call-template name="param">
            <xsl:with-param name="title">
                <xsl:copy-of select="$locale-path-param" />
            </xsl:with-param>
            <xsl:with-param name="param" select="$path/param" />
        </xsl:call-template>
    </xsl:if>

    <xsl:if test="$path/query">
        <xsl:call-template name="param">
            <xsl:with-param name="title">
                <xsl:copy-of select="$locale-query" />
            </xsl:with-param>
            <xsl:with-param name="param" select="$path/query" />
        </xsl:call-template>
    </xsl:if>
    
    <xsl:if test="$request/header">
        <xsl:call-template name="param">
            <xsl:with-param name="title">
                <xsl:copy-of select="$locale-header" />
            </xsl:with-param>
            <xsl:with-param name="param" select="$request/header" />
        </xsl:call-template>
    </xsl:if>

    <xsl:call-template name="param">
        <xsl:with-param name="title">
            <xsl:copy-of select="$locale-body" />
        </xsl:with-param>
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
            <xsl:with-param name="title">
                <xsl:copy-of select="$locale-header" />
            </xsl:with-param>
            <xsl:with-param name="param" select="$response/header" />
        </xsl:call-template>
    </xsl:if>

    <xsl:call-template name="param">
        <xsl:with-param name="title">
            <xsl:copy-of select="$locale-body" />
        </xsl:with-param>
        <xsl:with-param name="param" select="$response" />
    </xsl:call-template>
</xsl:template>

<!-- path param, path query, header 等的界面 -->
<xsl:template name="param">
    <xsl:param name="title" />
    <xsl:param name="param" />
    <xsl:param name="example" /> <!-- 示例代码 -->

    <div class="param">
        <h4 class="title">&#x27a4;&#160;<xsl:copy-of select="$title" /></h4>
        <table>
            <thead>
                <tr>
                    <th><xsl:copy-of select="$locale-var" /></th>
                    <th><xsl:copy-of select="$locale-type" /></th>
                    <th><xsl:copy-of select="$locale-value" /></th>
                    <th><xsl:copy-of select="$locale-description" /></th>
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
    <xsl:call-template name="deprecated">
        <xsl:with-param name="deprecated" select="@deprecated" />
    </xsl:call-template>

        <th>
            <xsl:if test="$parent"><xsl:value-of select="concat($parent, '.')" /></xsl:if>
            <xsl:value-of select="@name" />
        </th>

        <td>
            <xsl:value-of select="@type" />
            <xsl:if test="@array = 'true'"><xsl:value-of select="'[]'" /></xsl:if>
        </td>

        <td>
            <xsl:choose>
                <xsl:when test="@optional = 'true'"><xsl:value-of select="'O'" /></xsl:when>
                <xsl:otherwise><xsl:value-of select="'R'" /></xsl:otherwise>
            </xsl:choose>

            <xsl:if test="@default">
                <xsl:value-of select="concat(' ', @default)" />
            </xsl:if>
        </td>

        <td>
            <xsl:choose>
                <xsl:when test="description"><xsl:value-of select="description" /></xsl:when>
                <xsl:otherwise><xsl:value-of select="@summary" /></xsl:otherwise>
            </xsl:choose>
            <xsl:if test="./enum">
                <p><xsl:copy-of select="$locale-enum" /></p>
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
            <xsl:value-of select="concat('弃用于 ', $deprecated)" />
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

<!-- 替换字符串中特定的字符 -->
<xsl:template name="replace">
    <xsl:param name="text" />
    <xsl:param name="old" />
    <xsl:param name="new" />
    <xsl:choose>
        <xsl:when test="contains($text, $old)">
            <xsl:value-of select="concat(substring-before($text, $old), $new)" />
            <xsl:call-template name="replace">
                <xsl:with-param name="text" select="substring-after($text, $old)" />
                <xsl:with-param name="old" select="$old" />
                <xsl:with-param name="new" select="$new" />
            </xsl:call-template>
        </xsl:when>
        <xsl:otherwise><xsl:value-of select="$text" /></xsl:otherwise>
    </xsl:choose>
</xsl:template>

</xsl:stylesheet>
