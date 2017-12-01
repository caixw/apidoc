"use strict";

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
