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
    background: rgb(119,119,119);

    position: fixed;
    top: 0;
    left: 0;
    bottom: 0;
    width: 350px;
}

aside header {
    padding: 1rem;
}

aside footer {
    padding: 1rem;
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

.api .content{
    display:none;
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

$(document).ready(function(){
    // 调整缩进
    $('pre code').each(function(index, elem){
        var code = $(elem).text();
        $(elem).text(alignCode(code));
    });

    // 美化带有子元素的参数显示
    $('.request .params tbody th,.response .params tbody th').each(function(index, elem){
        var text = $(elem).text();
        text = text.replace(/(.*\.{1})/,'<span class="parent">$1</span>');
        $(elem).html(text);
    });

    // 按分组跳转页面
    $('#groups').on('change', function(){
        window.location.href = $(this).find('option:selected').val();
    });

    $('.api h3').on('click', function(){
        $(this).siblings('.content').slideToggle();
    });


    /* sticky */
    if (!navigator.userAgent.match(/firefox/i)){
        var header = $('header');
        var top = header.offset().top;
        $(document).on('scroll', function(e){
            window.scrollY > top ? header.addClass('sticky') : header.removeClass('sticky');
        });
    }


    // 代码高亮，依赖于是否能访问网络。
    if (typeof(Prism) != 'undefined') {
        Prism.plugins.autoloader.languages_path='https://cdn.bootcss.com/prism/1.5.1/components/';
        Prism.highlightAll(false);
    }
});

// 对齐代码。
function alignCode(code) {
    return code.replace(/^\s*/gm, function(word) {
        word = word.replace('\t', repeatSpace(indentSize));

        // 按 indentSize 的倍数取得缩进的量
        var len = Math.ceil((word.length-2)/indentSize)*indentSize;
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
`),}
var Templates=map[string]string{
"./index.html":`<!DOCTYPE html>
<html lang="zh-cmn-Hans">
    <head>
        <meta charset="utf-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1" />
        <meta http-equiv="X-UA-Compatible" content="IE=edge" />
        <meta name="generator" content="{{.AppOfficialURL}}" />
        <title></title>
        <link rel="stylesheet" href="./style.css" />
        <link href="https://cdn.bootcss.com/prism/1.5.1/themes/prism.min.css" rel="stylesheet" />

        <script src="./jquery-3.0.0.min.js"></script>
        <script src="https://cdn.bootcss.com/prism/1.5.1/prism.min.js" data-manual></script>
        <script src="https://cdn.bootcss.com/prism/1.5.1/plugins/autoloader/prism-autoloader.min.js"></script>
        <script src="./app.js"></script>
    </head>
    <body>
        <aside>
            <header>
                <h1>apidoc</h1>
            </header>

            <div>
                <ul>
                    <li>123</li>
                    <li>123</li>
                    <li>123</li>
                </ul>
            </div>

            <footer>
                <p>
                    内容由 <a href="{{.AppOfficialURL}}">{{.AppName}}</a> 编译于 <time>{{.Date|dateFormat}}</time>，
                    用时{{.Elapsed}}。
                </p>

                {{if .LicenseName}}
                <p>
                    内容采用
                    {{ if .LicenseURL}}<a href="{{.LicenseURL}}">{{end}}
                        {{- .LicenseName -}}
                    {{ if .LicenseURL}}</a>{{end}}
                    进行许可。
                </p>
                {{end}}
            </footer>
        </aside>

        <main>123</main>

        <script>
            import page from './data/page.json'
        </script>
    </body>
</html>
`,"./api.html":`{{- define "api" -}}
<section class="api">
    <h3>
        <span class="method {{.Method}}">{{.Method}}</span>
        <span class="url">{{.URL}}</span>
        <span class="summary">{{.Summary}}</span>
    </h3>

    <div class="content">
        {{if .Description}}
        <p class="description">{{.Description}}</p>
        {{end}}

        {{if .Queries}}
            <h5>查询参数</h5>
            {{template "params" .Queries}}
        {{end}}

        {{if .Params}}
            <h5>参数</h5>
            {{template "params" .Params}}
        {{end}}

        {{if .Request}}
        <div class="request">
            <h4>请求{{if .Request.Type}}:&#160;{{.Request.Type}}{{end}}</h4>
            <div>
                {{if .Request.Headers}}
                    <h5>报头:</h5>
                    {{template "headers" .Request.Headers}}
                {{end}}

                {{if .Request.Params}}
                    <h5>参数:</h5>
                    {{template "params" .Request.Params}}
                {{end}}

                {{if .Request.Examples}}
                    <h5>示例:</h5>
                    {{template "examples" .Request.Examples}}
                {{end}}
            </div>
        </div>
        {{end}}

        {{if .Success}}
        <div class="response success">
            <h4><span class="success">SUCCESS:</span>{{.Success.Code}},&#160;{{.Success.Summary}}</h4>
            {{template "response" .Success}}
        </div>
        {{end}}

        {{if .Error}}
        <div class="response error">
            <h4><span class="error">ERROR:</span>{{.Error.Code}},&#160;{{.Error.Summary}}</h4>
            {{template "response" .Error}}
        </div>
        {{end}}
    </div>
</section>
{{- end -}}



{{- define "examples" -}}
{{range .}}
<pre><code class="language-{{.Type|lower}}">{{- .Code -}}
</code></pre>
{{end}}
{{- end -}}



{{/* @apiParam 和 @apiQuery */}}
{{- define "params" -}}
<table class="params">
    <thead>
        <tr><th>名称</th><th>类型</th><th>描述</th></tr>
    </thead>
    <tbody>
    {{range . -}}
    <tr>
        <th>{{.Name}}</th>
        <td>{{.Type}}</td>
        <td>{{.Summary}}</td>
    </tr>
    {{- end}}
    </tbody>
</table>
{{- end -}}


{{/* @apiHeader */}}
{{- define "headers" -}}
<table>
    <thead>
        <tr><th>名称</th><th>描述</th></tr>
    </thead>
    <tbody>
    {{range $k, $v := . -}}
    <tr>
        <th>{{$k}}</th>
        <td>{{$v}}</td>
    </tr>
    {{- end}}
    </tbody>
</table>
{{- end -}}



{{/* @apiSuccess 和 @apiError */}}
{{- define "response" -}}
        {{if .Headers -}}
            <h5>请求头</h5>
            {{template "headers" .Headers}}
        {{- end}}

        {{- if .Params -}}
            <h5>参数:</h5>
            {{template "params" .Params}}
        {{- end}}

        {{- if .Examples -}}
            <h5>示例:</h5>
            {{template "examples" .Examples}}
        {{- end}}
{{- end -}}
`,}
