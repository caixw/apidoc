<?xml version="1.0" encoding="UTF-8"?>

<xsl:stylesheet version="1.0" xmlns:xsl="http://www.w3.org/1999/XSL/Transform">

<xsl:output
    method="html"
    encoding="utf-8"
    indent="yes"
    version="5.0"
    doctype-system="about:legacy-compat" />

<xsl:variable name="config">
    <xsl:value-of select="document('config.xml')/config/name" />
</xsl:variable>

<xsl:template match="/">
    <html>
        <head>
            <title><xsl:value-of select="docs/title" /></title>
            <meta charset="UTF-8" />
            <meta name="viewport" content="width=device-width,initial-scale=1,maximum-scale=1"/>
            <meta name="generator" content="https://apidoc.tools" />
            <link rel="icon" type="image/png" href="./icon.png" />
            <link rel="stylesheet" type="text/css" href="./index.css" />
            <link rel="license" href="{/docs/liense/@url}" />
            <script src="./index.js"></script>
        </head>
        <body>
            <header>
                <h1>
                    <img src="./icon.svg" />
                    <xsl:value-of select="document('config.xml')/config/item[@key='name']" />
                    <span class="version">(<xsl:value-of select="document('config.xml')/config/item[@key='version']" />)</span>
                </h1>
            </header>

            <main>
            <xsl:for-each select="docs/article">
                <article>
                    <h2><xsl:value-of select="@title" /></h2>
                    <xsl:copy-of select="." />
                </article>
            </xsl:for-each>
            </main>

            <footer>
                <xsl:copy-of select="docs/footer/node()" />
            </footer>
        </body>
    </html>
</xsl:template>
</xsl:stylesheet>
