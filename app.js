'use strict';

$(document).ready(function(){
    /* sticky */
    if (!navigator.userAgent.match(/firefox/i)){
        var header = $('header');
        var top = header.offset().top;
        $(document).on('scroll', function(e){
            window.scrollY > top ? header.addClass('sticky') : header.removeClass('sticky');
        });
    }

    // 根据与页面顶部的距离，控制是否显示 top 按钮。
    $(window).on('scroll', function(){
        var button = $('#top');
        if($(document).scrollTop() > 30){
            button.fadeIn();
        }else{
            button.fadeOut();
        }
    }).trigger('scroll'); // end $(window).onscroll

    // 滚动到顶部
    $('.goto-top').on('click', function(){
        var times = 20;
        var height = $('#top').offset().top;
        var offset = height / times;

        var tick = window.setInterval(function(){
            height -= offset;
            window.scrollTo(0, height);

            if (height <= 0){
                window.clearInterval(tick);
            }
        }, 10);

        return false;
    });

    /* 高亮和格式化代码块 */
    $('pre').each(function(index, elem){
        // 去多余的前导空格，由 data-indent 属性指定缩进多少个 tab
        var indent = $(this).attr('data-indent');
        var pattern = '';
        for(var i=0; i<indent; i++){
            pattern += '    ';
        }
        var text = $(elem).html().replace(new RegExp(pattern, 'g'), '');
        // 替换空格必须在高亮之前，高亮代码会添加 html 属性，用到空格
        text = text.replace(/ /g, '&#160;');

        // 高亮代码
        var lang = $(elem).attr('data-language');
        if (lang == 'json'){
            text = text.replace(/(\"[^"']*\")(:)(.+)(,?)/g, '<span class="keyword">$1</span>$2$3$4').
                replace(/(:[&#160;]*)\"([^"']*)\"(,?)/g, '$1<span class="string">"$2"</span>$3').
                replace(/(:[&#160;]*)(true|false)(,?)/g, '$1<span class="bool">$2</span>$3')

        }else{ // 未指定
            text = text.replace(/(@\w+)/g, '<span class="keyword">$1</span>');
        }

        $(elem).html(text);
    });
});
