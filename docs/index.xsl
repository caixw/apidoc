<?xml version="1.0" encoding="UTF-8"?>

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
                <a href="#" class="goto-top" />
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
                <img src="./icon.svg" />
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
<xsl:variable name="header-locale" select="/docs/type-locale/header" />

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
            <input type="checkbox" checked="true" disabled="true" />
        </xsl:when>
        <xsl:otherwise>
            <input type="checkbox" disabled="true" />
        </xsl:otherwise>
    </xsl:choose>
</xsl:template>

</xsl:stylesheet>
