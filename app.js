'use strict';

$(document).ready(function(){
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
