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
                        <label><input type="checkbox" /><xsl:value-of select="@name" /></label>
                        <xsl:attribute name="data-server">
                            <xsl:value-of select="@name" />
                        </xsl:attribute>
                    </li>
                    </xsl:for-each>
                </ul>

                <h2>Tags</h2>
                <ul class="tags-selector">
                    <xsl:for-each select="apidoc/tag">
                    <li>
                        <label><input type="checkbox" /><xsl:value-of select="@name" /></label>
                        <xsl:attribute name="data-server">
                            <xsl:value-of select="@name" />
                        </xsl:attribute>
                    </li>
                    </xsl:for-each>
                </ul>

                <h2>Methods</h2>
                <ul class="methods-selector">
                    <!-- 浏览器好像都不支持 xpath 2.0，所以无法使用 distinct-values -->
                    <!-- xsl:for-each select="distinct-values(/apidoc/api/@method)" -->
                    <xsl:for-each select="/apidoc/api/@method[not(../preceding-sibling::api/@method = .)]">
                    <li>
                        <label><input type="checkbox" /><xsl:value-of select="." /></label>
                        <xsl:attribute name="data-method">
                            <xsl:value-of select="." />
                        </xsl:attribute>
                    </li>
                    </xsl:for-each>
                </ul>
                <!-- methods -->
            </aside>

            <main>
                <div class="content">
                    <xsl:value-of select="/apidoc/content" />
                </div>

                <div class="servers">
                    <xsl:for-each select="/apidoc/server">
                    <article class="server">
                        <h3>
                            <xsl:value-of select="@url"/>
                            <xsl:value-of select="@name"/>
                        </h3>
                        <div class="summary">
                            <xsl:value-of select="."/>
                        </div>
                    </article>
                    </xsl:for-each>
                </div>

                <xsl:for-each select="apidoc/api">
                <article>
                <xsl:attribute name="data-method">
                    <xsl:value-of select="@method" />
                </xsl:attribute>

                    <h3>
                        <span class="method">
                        <xsl:value-of select="@method" />
                        </span>
                        <xsl:value-of select="path/@path" />
                    </h3>
                    <div class="summary">
                        <xsl:value-of select="@summary" />
                    </div>

                    <div>
                        <div class="request">
                            <!-- req -->
                        </div>
                        <div class="response">
                            <!-- response -->
                        </div>
                    </div>
                </article>
                </xsl:for-each>
            </main>

            <footer>
            <p>文档由 <a href="https://github.com/caixw/apidoc">apidoc</a> 生成！</p>
            </footer>
        </body>
    </html>
</xsl:template>
</xsl:stylesheet>
