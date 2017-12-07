// 该文件由 make.go 自动生成，请勿手动修改！

package static

var assets=map[string][]byte{
"./style.css":[]byte(`@charset "utf-8";

/*============== reset =================*/
body {
    margin: 0
}

a {
    text-decoration: none;
    color: #3b8bba;
}

/*=============== aside ================*/

aside {
    background: rgb(189,189,189);

    position: fixed;
    top: 0;
    left: 0;
    bottom: 0;
    width: 350px;
}

aside header {
    padding: 1rem;
    position: sticky;
    position: -webkit-sticky;
}

aside menu {
}

aside footer {
    padding: 1rem;
    position: fixed;
    bottom: 0;
    left: 0;
}

/*=============== main ================*/

main {
    margin-left: 350px;
    padding: 1rem;
}

/*=============== .api ================*/

main .api{
    padding:1rem;
    margin:1rem 0rem;
    border:1px solid #eee;
}

.api h3{
    cursor:pointer;
    margin:0rem;
    display:flex;
    align-items:center;
}

.api h4{
    font-size:1.1rem;
    margin-bottom:.2rem;
    border-bottom:1px solid #eee;
    padding-bottom:.2rem;
}

.api h5{
    margin:.8rem 0rem .2rem 0rem;
    font-size:1rem;
}

.api h3 .method{
    width:5rem;
    font-weight:bold;
    text-transform:uppercase;
}

.api h3 .get{
    color:green;
}

.api h3 .options{
    color:green;
}

.api h3 .delete{
    color:red;
}

.api h3 .put,.api h3 .patch{
    color:rgb(193,174,49);
}

.api h3 .post{
    color:rgb(240,114,11);
}

.api h3 .url{
    margin-right:2rem;
}

.api h4 .success{
    color:green;
    margin-right:1rem;
}

.api h4 .error{
    color:red;
    margin-right:1rem;
}

.api table{
    text-align:left;
    border-collapse:collapse;
    border:1px solid #ddd;
}

.api table thead tr{
    background:#eee;
}

.api table tr{
    border-bottom:1px solid #ddd;
    line-height:1.5rem;
}

.api table tbody th .parent{
    color:#ccc;
}

.api table th, .api table td{
    padding:.3rem 1rem;
}

`),"./app.js":[]byte(`"use strict";

// 代码缩进的空格数量。
var indentSize = 4;

window.onload = function() {
    initTemplate()
}

function initTemplate() {
    var source = document.querySelector('#page').innerHTML
    let pageTpl = Handlebars.compile(source)

    source = document.querySelector('#examples').innerHTML
    Handlebars.registerPartial('examples', source)

    source = document.querySelector('#params').innerHTML
    Handlebars.registerPartial('params', source)

    source = document.querySelector('#headers').innerHTML
    Handlebars.registerPartial('headers', source)

    source = document.querySelector('#response').innerHTML
    Handlebars.registerPartial('response', source)

    source = document.querySelector('#api').innerHTML
    let apiTpl = Handlebars.compile(source)


    fetch('./data/page.json').then((resp)=>{
        return resp.json();
    }).then((json)=>{
        document.querySelector('#app').innerHTML = pageTpl(json)
        document.title = json.title + ' | ' + json.appName

        loadApis(json)
    })

    // 加载 api 模板内容，json 为对应的数据
    function loadApis(json) {
        document.querySelector('.menu>li.content').addEventListener('click', (event)=>{
            document.querySelector('main').innerHTML = json.content
        })

        document.querySelectorAll('.menu>li.api').forEach((elem, index, list)=>{
            elem.addEventListener('click', (event)=>{
                let path = event.target.getAttribute('data-path')
                fetch(path).then((resp)=>{
                    return resp.json()
                }).then((json)=>{
                    document.querySelector('main').innerHTML = apiTpl(json)

                    indentCode()
                    prettifyParams()
                    highlightCode()
                }).catch((reason)=>{
                    console.error(reason)
                })
            })
        })
    } // end loadApis
}

// 美化带有子元素的参数显示
function prettifyParams() {
    $('.request .params tbody th,.response .params tbody th').each(function(index, elem){
        var text = $(elem).text();
        text = text.replace(/(.*\.{1})/,'<span class="parent">$1</span>');
        $(elem).html(text);
    });
}

// 代码高亮，依赖于是否能访问网络。
function highlightCode() {
    if (typeof(Prism) != 'undefined') {
        Prism.plugins.autoloader.languages_path='https://cdn.bootcss.com/prism/1.5.1/components/';
        Prism.highlightAll(false);
    }
}

// 调整缩进
function indentCode() {
    $('pre code').each((index, elem)=>{
        let code = $(elem).text();
        $(elem).text(alignCode(code));
    });
}

// 对齐代码。
function alignCode(code) {
    return code.replace(/^\s*/gm, (word)=>{
        word = word.replace('\t', repeatSpace(indentSize));

        // 按 indentSize 的倍数取得缩进的量
        let len = Math.ceil((word.length-2)/indentSize)*indentSize;
        return repeatSpace(len);
    });
}

function repeatSpace(len) {
    var code = [];
    while(code.length < len) {
        code.push(' ');
    }

    return code.join('');
}
`),"./index.html":[]byte(`<!DOCTYPE html>
<html lang="zh-cmn-Hans">
    <head>
        <meta charset="utf-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <meta http-equiv="X-UA-Compatible" content="IE=edge" />
        <title>APIDOC</title>
        <link rel="stylesheet" href="./style.css" />
        <link href="https://cdn.bootcss.com/prism/1.5.1/themes/prism.min.css" rel="stylesheet" />

        <script src="https://cdn.bootcss.com/jquery/3.2.1/jquery.min.js"></script>
        <script src="https://cdn.bootcss.com/handlebars.js/4.0.11/handlebars.min.js"></script>
        <script src="https://cdn.bootcss.com/prism/1.5.1/prism.min.js" data-manual></script>
        <script src="https://cdn.bootcss.com/prism/1.5.1/plugins/autoloader/prism-autoloader.min.js"></script>
    </head>
    <body>
        <div id="app"></div>

        <script id="page" type="text/x-handlebars-template">
            <aside>
                <header><p>{{title}}</p></header>

                <ul class="menu">
                    <li class="menu-item content" data-path="content">home</li>
                    {{#each groups}}
                    <li class="menu-item api" data-path="{{this}}">{{@key}}</li>
                    {{/each}}
                </ul>

                <footer>
                    <p>内容由 <a href="{{appURL}}">{{appName}}</a> 编译于 <time>{{date}}</time>，用时{{elapsed}}。</p>

                    {{#if licenseName}}
                    <p>内容采用<a href="{{licenseURL}}">{{licenseName}}</a>进行许可。</p>
                    {{/if}}
                </footer>
            </aside>

            <main id="main">{{{content}}}</main>
        </script>


        <script id="examples" type="text/x-handlebars-template">
            {{#each examples}}
            <pre><code class="language-{{type}}">{{code}}
            </code></pre>
            {{/each}}
        </script>

        <script id="params" type="text/x-handlebars-template">
            <table class="params">
                <thead>
                    <tr><th>名称</th><th>类型</th><th>描述</th></tr>
                </thead>
                <tbody>
                {{#each params}}
                <tr>
                    <th>{{name}}</th>
                    <td>{{type}}</td>
                    <td>{{summary}}</td>
                </tr>
                {{/each}}
                </tbody>
            </table>
        </script>

        <script id="headers" type="text/x-handlebars-template">
            <table>
                <thead>
                    <tr><th>名称</th><th>描述</th></tr>
                </thead>
                <tbody>
                {{#each headers}}
                <tr>
                    <th>{{@key}}</th>
                    <td>{{this}}</td>
                </tr>
                {{/each}}
                </tbody>
            </table>
        </script>

        <script id="response" type="text/x-handlebars-template">
        {{#if response.headers}}
            <h5>请求头</h5>
            {{> headers headers=response.headers}}
        {{/if}}

        {{#if response.params}}
            <h5>参数:</h5>
            {{> params params=response.params}}
        {{/if}}

        {{#if response.examples}}
            <h5>示例:</h5>
            {{> examples examples=response.examples}}
        {{/if}}
        </script>


        <script id="api" type="text/x-handlebars-template">
            <h2>{{name}}</h2>

            {{#each apis}}
            <section class="api">
                <h3>
                    <span class="method {{method}}">{{method}}</span>
                    <span class="url">{{url}}</span>
                    <span class="summary">{{summary}}</span>
                </h3>

                <div class="content">
                    {{#if description}}
                    <p class="description">{{description}}</p>
                    {{/if}}

                    {{#if queries}}
                        <h5>查询参数</h5>
                        {{> params params=queries}}
                    {{/if}}

                    {{#if params}}
                        <h5>参数</h5>
                        {{> params params=params}}
                    {{/if}}

                    {{#if request}}
                    <div class="request">
                        <h4>请求{{#if request.type}}:&#160;{{request.type}}{{/if}}</h4>
                        <div>
                            {{#if request.headers}}
                                <h5>报头:</h5>
                                {{> headers headers=request.headers}}
                            {{/if}}

                            {{#if request.params}}
                                <h5>参数:</h5>
                                {{> params params=request.params}}
                            {{/if}}

                            {{#if request.examples}}
                                <h5>示例:</h5>
                                {{> examples examples=request.examples}}
                            {{/if}}
                        </div>
                    </div>
                    {{/if}}

                    {{#if success}}
                    <div class="response success">
                        <h4><span class="success">SUCCESS:</span>{{success.code}},&#160;{{success.summary}}</h4>
                        {{> response response=success}}
                    </div>
                    {{/if}}

                    {{#if error}}
                    <div class="response error">
                        <h4><span class="error">ERROR:</span>{{error.code}},&#160;{{error.summary}}</h4>
                        {{> response response=error}}
                    </div>
                    {{/if}}
                </div>
            </section>
            {{/each}}
        </script>

        <script src="./app.js"></script>
    </body>
</html>
`),}
