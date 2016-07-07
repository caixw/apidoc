"use strict";

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

    // 隐藏当前页面用不到的过滤器
    $('header .filter input').each(function(index, elem){
        var val = $(elem).val();
        if ($('.main .api .method.'+val).length == 0) {
            $(elem).parent().hide();
        }
    });

    // 按请求方法过滤
    $('header .filter input').on('change', function(){
        var val = $(this).attr('value');
        $('.main .api .method.'+val).parents('.api').each(function(index, elem){
            $(elem).slideToggle();
        });
    })

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
