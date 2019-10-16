<?xml version="1.0" encoding="UTF-8"?>

<xsl:stylesheet version="1.0" xmlns:xsl="http://www.w3.org/1999/XSL/Transform">

<xsl:output
    method="html"
    encoding="utf-8"
    indent="yes"
    version="5.0"
    doctype-system="about:legacy-compat" />

<xsl:template match="/">
    <html>
        <head>
            <title><xsl:value-of select="docs/title" /></title>
            <meta charset="UTF-8" />
            <meta name="viewport" content="width=device-width,initial-scale=1,maximum-scale=1"/>
            <meta name="generator" content="https://apidoc.tools" />
            <link rel="icon" type="image/png" href="./icon.png" />
            <link rel="canonical" href="{document('config.xml')/config/url}" />
            <link rel="stylesheet" type="text/css" href="./index.css" />
            <link rel="license" href="{/docs/liense/@url}" />
            <script src="./index.js"></script>
        </head>
        <body>
            <header>
                <h1>
                    <img src="./icon.svg" />
                    <xsl:value-of select="document('config.xml')/config/name" />
                    <span class="version">&#160;(<xsl:value-of select="document('config.xml')/config/version" />)</span>
                </h1>

                <a href="{document('config.xml')/config/repo}">Github</a>
                <xsl:for-each select="docs/doc[not(@parent)]">
                    <xsl:sort select="position()" order="descending" data-type="number" />
                    <a href="#{@id}"><xsl:value-of select="@title" /></a>
                </xsl:for-each>
            </header>

            <main>
                <xsl:for-each select="docs/doc[not(@parent)]">
                    <xsl:call-template name="article">
                        <xsl:with-param name="doc" select="." />
                    </xsl:call-template>
                </xsl:for-each>
            </main>

            <footer>
                <xsl:copy-of select="docs/footer/node()" />
            </footer>
        </body>
    </html>
</xsl:template>

<xsl:template name="article">
    <xsl:param name="doc" />

    <article>
        <h2><xsl:value-of select="$doc/@title" /></h2>
        <xsl:copy-of select="$doc/node()" />

        <xsl:for-each select="/docs/doc[@parent=$doc/@id]">
            <xsl:call-template name="article">
                <xsl:with-param name="doc" select="." />
            </xsl:call-template>
        </xsl:for-each>

        <xsl:for-each select="/docs/types[@parent=$doc/@id]/type">
            <table>
                <thead>
                    <tr>
                        <th><xsl:value-of select="/docs/type-locale/header/name" /></th>
                        <th><xsl:value-of select="/docs/type-locale/header/type" /></th>
                        <th><xsl:value-of select="/docs/type-locale/header/optional" /></th>
                        <th><xsl:value-of select="/docs/type-locale/header/description" /></th>
                    </tr>
                </thead>
                <tbody>
                    <xsl:for-each select="item">
                        <tr>
                            <th><xsl:value-of select="@name" /></th>
                            <td><xsl:value-of select="@type" /></td>
                            <td><xsl:value-of select="@optional" /></td>
                            <td><xsl:copy-of select="." /></td>
                        </tr>
                    </xsl:for-each>

                    <xsl:for-each select="group">
                        <tr data-group="true">
                        <xsl:choose> <!-- 如果 group 也有类型信息，则按照 type 的方式输出 -->
                            <xsl:when test="@type">
                                <th><xsl:value-of select="@name" /></th>
                                <td><xsl:value-of select="@type" /></td>
                                <td><xsl:value-of select="@optional" /></td>
                                <td><xsl:copy-of select="description" /></td>
                            </xsl:when>
                            <xsl:otherwise>
                                <th colspan="4"><xsl:value-of select="@name" /></th>
                            </xsl:otherwise>
                        </xsl:choose>
                        </tr>

                        <xsl:for-each select="item">
                        <tr>
                            <th><xsl:value-of select="@name" /></th>
                            <td><xsl:value-of select="@type" /></td>
                            <td><xsl:value-of select="@optional" /></td>
                            <td><xsl:copy-of select="." /></td>
                        </tr>
                        </xsl:for-each>
                    </xsl:for-each>
                </tbody>
            </table>
        </xsl:for-each>
    </article>
</xsl:template>

</xsl:stylesheet>
