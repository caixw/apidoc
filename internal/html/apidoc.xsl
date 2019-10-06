<?xml version="1.0" encoding="UTF-8"?>

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
    <xsl:param name="example" /><!-- 示例代码，即在 parent 为空时，才有可能有此值 -->

    <h5><xsl:value-of select="$title" /></h5>
    <table>
        <thead>
            <tr>
                <th>变量</th>
                <th>类型</th>
                <th>必须</th>
                <th>默认值</th>
                <th>描述</th>
            </tr>
        </thead>
        <tbody>
            <xsl:call-template name="param-list">
            <xsl:with-param name="param" select="$param" />
            </xsl:call-template>
        </tbody>
    </table>
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
        <td><xsl:value-of select="@type" /></td>
        <td>
            <xsl:if test="@required = 'true'">
                <xsl:value-of select="'&#10003;'" />
            </xsl:if>
        </td>
        <td><xsl:value-of select="@default" /></td>
        <td>
            <xsl:value-of select="@summary" />
            <xsl:if test="./enum">
                <p>可以使用以下枚举值：</p>
                <ul>
                <xsl:for-each select="./enum">
                    <li>
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

<xsl:template match="/apidoc/api/request">
<div class="request">
    <h5 class="mimetype"><xsl:value-of select="@mimetype" /></h5>
    <xsl:if test="../path/param">
        <xsl:call-template name="param">
            <xsl:with-param name="title" select="'路径参数'" />
            <xsl:with-param name="param" select="../path/param" />
        </xsl:call-template>
    </xsl:if>
    <xsl:if test="../path/query">
        <xsl:call-template name="param">
            <xsl:with-param name="title" select="'查询参数'" />
            <xsl:with-param name="param" select="../path/query" />
        </xsl:call-template>
    </xsl:if>
    <xsl:if test="./header">
        <xsl:call-template name="param">
            <xsl:with-param name="title" select="'请求报头'" />
            <xsl:with-param name="param" select="header" />
        </xsl:call-template>
    </xsl:if>

    <xsl:call-template name="param">
        <xsl:with-param name="param" select="." />
    </xsl:call-template>
</div>
</xsl:template>

<xsl:template match="/apidoc/api/response">
    <xsl:call-template name="param">
        <xsl:with-param name="param" select="." />
    </xsl:call-template>
</xsl:template>

<!-- api 元素界面 -->
<xsl:template match="/apidoc/api">
    <article class="api">
    <xsl:attribute name="data-method">
        <xsl:value-of select="@method" />
    </xsl:attribute>
    <xsl:attribute name="id">
        <xsl:call-template name="get-api-id">
            <xsl:with-param name="path" select="path/@path" />
            <xsl:with-param name="method" select="@method" />
        </xsl:call-template>
    </xsl:attribute>

        <h3>
            <span class="action">
            <xsl:value-of select="@method" />
            </span>
            <xsl:value-of select="path/@path" />
        </h3>
        <div class="summary">
            <xsl:value-of select="@summary" />

            <xsl:if test="./description">
            <br />
            <xsl:value-of select="./description" />
            </xsl:if>
        </div>

        <div class="body">
            <div class="requests">
                <h4>请求</h4>
                <xsl:for-each select="request">
                    <xsl:apply-templates select="." />
                </xsl:for-each>
            </div>
            <div class="responses">
                <h4>返回</h4>
                <xsl:for-each select="response">
                    <xsl:apply-templates select="." />
                </xsl:for-each>
            </div>
        </div>
    </article>
</xsl:template>

<xsl:template match="/">
    <html>
        <head>
            <title><xsl:value-of select="apidoc/title" /></title>
            <meta charset="UTF-8" />
            <meta name="viewport" content="width=device-width,initial-scale=1,maximum-scale=1"/>
            <link rel="stylesheet" type="text/css" href="./apidoc.css" />
            <link rel="icon" type="image/png" href="./icon.png" />
            <link rel="license">
                <xsl:attribute name="href"><xsl:value-of select="apidoc/license/@url"/></xsl:attribute>
            </link>
            <script src="./apidoc.js"></script>
        </head>
        <body>
            <header>
                <h1>
                    <span class="app-name"><img src="./icon.svg" />apidoc</span>
                </h1>
            </header>

            <aside>
                <h2>Servers</h2>
                <ul class="servers-selector">
                    <xsl:for-each select="apidoc/server">
                    <li>
                        <xsl:attribute name="data-server"><!-- chrome 和 safari 必须要在其它元素之前 -->
                            <xsl:value-of select="@name" />
                        </xsl:attribute>
                        <label><input type="checkbox" /><xsl:value-of select="@name" /></label>
                    </li>
                    </xsl:for-each>
                </ul>

                <h2>Tags</h2>
                <ul class="tags-selector">
                    <xsl:for-each select="apidoc/tag">
                    <li>
                        <xsl:attribute name="data-server">
                            <xsl:value-of select="@name" />
                        </xsl:attribute>
                        <label><input type="checkbox" /><xsl:value-of select="@name" /></label>
                    </li>
                    </xsl:for-each>
                </ul>

                <h2>Methods</h2>
                <ul class="methods-selector">
                    <!-- 浏览器好像都不支持 xpath 2.0，所以无法使用 distinct-values -->
                    <!-- xsl:for-each select="distinct-values(/apidoc/api/@method)" -->
                    <xsl:for-each select="/apidoc/api/@method[not(../preceding-sibling::api/@method = .)]">
                    <li>
                        <xsl:attribute name="data-method">
                            <xsl:value-of select="." />
                        </xsl:attribute>
                        <label><input type="checkbox" /><xsl:value-of select="." /></label>
                    </li>
                    </xsl:for-each>
                </ul>
            </aside>

            <main>
                <div class="content">
                    <xsl:value-of select="/apidoc/content" />
                </div>

                <xsl:for-each select="/apidoc/server">
                <article class="server">
                    <h3>
                        <span class="action">
                        <xsl:value-of select="@name"/>
                        </span>
                        <xsl:value-of select="@url"/>
                    </h3>
                    <div class="summary">
                        <xsl:value-of select="."/>
                    </div>
                </article>
                </xsl:for-each>

                <xsl:for-each select="apidoc/api">
                <xsl:sort select="path/@path"/>
                    <xsl:apply-templates select="." />
                </xsl:for-each>
            </main>

            <footer>
            <p>文档由 <a href="https://github.com/caixw/apidoc">apidoc</a> 生成！</p>
            </footer>
        </body>
    </html>
</xsl:template>
</xsl:stylesheet>
