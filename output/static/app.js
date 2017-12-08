"use strict";

// 代码缩进的空格数量。
let indentSize = 4

// 对应 vars.JSONDataDirName 的值
let dataDirName = 'data'

$(()=>{
    initTemplate()
})

function initTemplate() {
    Handlebars.registerPartial('examples', $('#examples').html())
    Handlebars.registerPartial('params', $('#params').html())
    Handlebars.registerPartial('headers', $('#headers').html())
    Handlebars.registerPartial('response', $('#response').html())

    Handlebars.registerHelper('dateFormat', formatDate)
    Handlebars.registerHelper('elapsedFormat', formatElapsed)

    let pageTpl = Handlebars.compile($('#page').html())
    let apiTpl = Handlebars.compile($('#api').html())

    fetch('./'+dataDirName+'/page.json').then((resp)=>{
        return resp.json();
    }).then((json)=>{
        $('#app').html(pageTpl(json))
        document.title = json.title + ' | ' + json.appName

        loadApis(json)
    })

    // 加载 api 模板内容，json 为对应的数据
    function loadApis(json) {
        let menu = $('aside .menu')
        menu.find('li.content').on('click', (event)=>{
            $('main').html(json.content)
        })

        menu.find('li.api').on('click', (event)=>{
            let path = $(event.target).attr('data-path')
            fetch(path).then((resp)=>{
                return resp.json()
            }).then((json)=>{
                $('main').html(apiTpl(json))

                indentCode()
                prettifyParams()
                highlightCode()
            }).catch((reason)=>{
                console.error(reason)
            })
        })
    } // end loadApis
}

// 美化带有子元素的参数显示
function prettifyParams() {
    $('.request .params tbody th,.response .params tbody th').each(function(index, elem){
        let text = $(elem).text()
        text = text.replace(/(.*\.{1})/,'<span class="parent">$1</span>')
        $(elem).html(text)
    })
}

// 代码高亮，依赖于是否能访问网络。
function highlightCode() {
    if (typeof(Prism) != 'undefined') {
        Prism.plugins.autoloader.languages_path='https://cdn.bootcss.com/prism/1.5.1/components/'
        Prism.highlightAll(false)
    }
}

// 调整缩进
function indentCode() {
    $('pre code').each((index, elem)=>{
        let code = $(elem).text()
        $(elem).text(alignCode(code))
    })
}

// 对齐代码。
function alignCode(code) {
    return code.replace(/^\s*/gm, (word)=>{
        word = word.replace('\t', repeatSpace(indentSize))

        // 按 indentSize 的倍数取得缩进的量
        let len = Math.ceil((word.length-2)/indentSize)*indentSize
        return repeatSpace(len)
    })
}

function repeatSpace(len) {
    var code = []
    while(code.length < len) {
        code.push(' ')
    }

    return code.join('')
}

function formatDate(unix) {
    let date = new Date(unix*1000)

    let str = []
    str.push(date.getFullYear(), '-')
    str.push(date.getMonth(), '-')
    str.push(date.getDate(), ' ')
    str.push(date.getHours(), ':')
    str.push(date.getMinutes(), ':')
    str.push(date.getSeconds())
    return str.join('')
}

function formatElapsed(number) {
    return (number / 100000000).toFixed(4) + '秒'
}
