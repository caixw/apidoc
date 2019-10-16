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
                <a href="#{docs/docs/@id}"><xsl:value-of select="docs/docs/@title" /></a>
                <a href="#{docs/usage/@id}"><xsl:value-of select="docs/usage/@title" /></a>
                <a href="#{docs/about/@id}"><xsl:value-of select="docs/about/@title" /></a>
            </header>

            <main>
                <article id="{docs/about/@id}">
                    <h2><xsl:value-of select="docs/about/@title" /></h2>
                    <xsl:copy-of select="docs/about" />
                </article>

                <article id="{docs/usage/@id}">
                    <h2><xsl:value-of select="docs/usage/@title" /></h2>
                    <xsl:copy-of select="docs/usage" />
                </article>

                <article id="{docs/docs/@id}">
                    <h2><xsl:value-of select="docs/docs/@title" /></h2>
                    <xsl:copy-of select="docs/docs" />

                    <xsl:for-each select="docs/types[@parent='docs']/type">
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
                                    <tr><th colspan="4"><xsl:value-of select="@name" /></th></tr>
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
            </main>

            <footer>
                <xsl:copy-of select="docs/footer/node()" />
            </footer>
        </body>
    </html>
</xsl:template>
</xsl:stylesheet>
