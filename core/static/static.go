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
    background-color:#fafafa;
    border:1px solid #eee;
    padding:3em 2.3em 0.5em 2.3em;
    left:0;
    right:0;
    top:0;
    position:fixed;

}

header .title{
    font-size:2em;
}

header .filter{
    margin-top:0.8em;
}

header label{
    margin-left:1em;
    vertical-align: bottom;
}

/*=============== main ================*/

.main{
    padding:0em 2em;
    margin-top:8em;
}

.main section{
    padding:1em;
    margin:1em 0em;
    border:1px solid #eee;
}

.main .get:hover{
    border:1px solid rgba(0,255,0,0.5);
}

.main .delete:hover{
    border:1px solid rgba(255,0,0,0.5);
}

.main .put:hover,.main .patch:hover{
    border:1px solid rgba(193,174,49,0.5);
}

.main .post:hover{
    border:1px solid rgba(240,114,11,0.5);
}

.main h3{
    margin-top:0em;
}

.main h4{
    margin-bottom:0em;
    border-bottom:1px solid #eee;
    padding-bottom:0.2em;
}

.main .version{
    margin-right:1em;
}

.main .methods{
    margin-right:2em;
    text-transform:uppercase;
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
        <script src="http://code.jquery.com/jquery-2.1.4.min.js"></script>
    </head>
    <body>
        <header>
            <span class="title">APIDOC</span> {{.Version}}
            <div class="fr filter">
                <label><input type="checkbox" checked="checked" value="get">GET</label>
                <label><input type="checkbox" checked="checked" value="post">POST</label>
                <label><input type="checkbox" checked="checked" value="put">PUT</label>
                <label><input type="checkbox" checked="checked" value="patch">PATCH</label>
                <label><input type="checkbox" checked="checked" value="delete">DELETE</label>
            </div>
        </header>

        <div class="main">
            {{range $key, $docs:=.Docs}}
            <div>
                <h2>{{$key}}</h2>
                {{range $docs}}
                    {{if .}}
                    <section class="{{.Methods}}">
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
        <script>
            $(document).ready(function(){
                $('header .filter input').on('change', function(){
                    var val = $(this).attr('value');
                    $('.main section.'+val).each(function(index, elem){
                        $(elem).slideToggle();
                    });
                })
            });
        </script>
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
        <h4>Query</h4>
        {{template "param" .Queries}}
    {{end}}

    {{if .Params}}
        <h4>Param</h4>
        {{template "param" .Params}}
    {{end}}

    {{if .Request}}
    <div>
        <h4>Request   {{.Request.Type}}</h4>
        <div>
            {{range $k,$v:=.Request.Headers}}
            <p>{{$k}}:{{$v}}</p>
            {{end}}
            {{range .Request.Params}}
            <p>{{.Name}}:{{.Type}},{{.Description}}</p>
            {{end}}
            {{range .Request.Examples}}
            <pre class="code" data-type="{{.Type}}">{{.Code}}</pre>
            {{end}}
        </div>
    </div>
    {{end}}

    {{if .Status}}
    <div>
        <h4>Response</h4>
        {{range .Status}}
        <div>
            <p>status:{{.Code}},  {{.Summary}},  {{.Type}}</p>
            {{range $k,$v:=.Headers}}
            <p>{{$k}}:{{$v}}</p>
            {{end}}
            {{range .Params}}
            <p>{{.Name}}:{{.Type}},{{.Description}}</p>
            {{end}}
            {{range .Examples}}
            <pre class="code" data-type="{{.Type}}">{{.Code}}</pre>
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
