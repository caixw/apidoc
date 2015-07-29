// 该文件由make.go自动生成，请勿手动修改！

package static

var files=map[string][]byte{
"./style.css":[]byte(`@charset "utf-8";

/*============== reset =================*/
body{
    margin:0em;
}

a{
    text-decoration:none;
    color:#3b8bba;
}

.fl{
    float:left;
}

.fr{
    float:right;
}

/*=============== header ================*/

header{
    color:#777;
    background-color:white;
    border-bottom:1px solid #eee;
    padding:3em 1em 0.3em 1em;
    margin:0em 2em;
}


/*=============== main ================*/

.main{
    padding:1em 3em;
}

.main .version{
    margin-right:1em;
}

.main .methods{
    margin-right:2em;
}

.main .request-type{
    margin-left:2em;
}

.main table caption{
    text-align:left;
}
.main table thead th{
    background-color:#fafafa;
    border:1px solid #eee;
    text-align:left;
    padding:0.2em;
}

.main table td{
    border:1px solid #eee;
    padding:0.2em;
}


/*=============== footer ================*/
footer{
    border-top:1px solid #eee;
    background-color:#fafafa;
    color:#777;
    padding:1em;
    margin-top:2em;
}

footer p{
    margin:0.5em;
    padding:0px;
}
`),}
var Templates=map[string]string{
"./main.html":`{{define "main"}}<!doctype html>
<html>
    <head>
        <meta charset="utf-8" />
        <meta name="viewport" content="width=device-width, initial-scale=1">
        <meta http-equiv="X-UA-Compatible" content="IE=edge">
        <meta name="generator" content="https://github.com/caixw/apidoc">
        <title>apidoc</title>
        <link rel="stylesheet" href="./style.css" />
    </head>
    <body>
        <header>APIDOC</header>

        <div class="main">
            {{range $key, $docs:=.Docs}}
            <div>
                <h2>{{$key}}</h2>
                {{range $docs}}
                    {{if .}}
                    <section>
                        {{template "doc" .}}
                    </section>
                    {{end}}
                {{end}}
            </div>
            {{end}}
        </div>

        <footer>
            <p>内容由<a href="https://github.com/caixw/apidoc">apidoc</a>于<time id="date">{{.Date}}</time>编译完成。</p>
        </footer>
    </body>
</html>
{{end}}
`,"./doc.html":`{{define "doc"}}
    <h3>
        <span class="methods">{{.Methods}}</span>
        <span class="url">{{.URL}}</span>
        <span class="fr">{{.Summary}}</span>
    </h3>

    {{if .Description}}
    <p class="description">{{.Description}}</p>
    {{end}}

    {{if .Queries}}
        <h5>Query</h5>
        {{template "param" .Queries}}
    {{end}}

    {{if .Params}}
        <h5>Param</h5>
        {{template "param" .Params}}
    {{end}}

    {{if .Request}}
        <div>
            <h5>Request</h5>
            <div>
                <p>数据类型:{{.Request.Type}}</p>
                {{range $k,$v:=.Request.Headers}}
                <p>{{$k}}:{{$v}}</p>
                {{end}}
                {{range .Request.Params}}
                <p>{{.Name}}:{{.Type}},{{.Description}}</p>
                {{end}}
                {{range .Request.Examples}}
                <div class="code" data-type="{{.Type}}">{{.Code}}</div>
                {{end}}
            </div>
        </div>
    {{end}}

    {{if .Status}}
        <div>
            <h5>Response</h5>
            {{range .Status}}
                <div>
                <p>status:{{.Code}},{{.Summary}}</p>
                <p>数据类型:{{.Type}}</p>
                {{range $k,$v:=.Headers}}
                <p>{{$k}}:{{$v}}</p>
                {{end}}
                {{range .Params}}
                <p>{{.Name}}:{{.Type}},{{.Description}}</p>
                {{end}}
                {{range .Examples}}
                <div class="code" data-type="{{.Type}}">{{.Code}}</div>
                {{end}}
            </div>
            {{end}}
        </div>
    {{end}}

{{end}}

{{define "param"}}
    <table>
        <thead>
            <tr>
                <th>Name</th>
                <th>Type</th>
                <th>Optional</th>
                <th>Description</th>
            </tr>
        </thead>
        <tbody>
        {{range .}}
            <tr>
                <td>{{.Name}}</td>
                <td>{{.Type}}</td>
                <td>{{.Optional}}</td>
                <td>{{.Description}}</td>
            </tr>
        {{end}}
        </tbody>
    </table>
{{end}}

`,}
