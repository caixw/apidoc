"use strict";

$(document).ready(function(){
    $('pre').each(function(index, elem){
        // 去多余的前导空格，由data-indent属性指定缩进多少个tab
        var indent = $(this).attr('data-indent');
        var pattern = '';
        for(var i=0; i<indent; i++){
            pattern += '    ';
        }
        var text = $(elem).html().replace(new RegExp(pattern, 'g'), '');

        text = text.replace(/ /g, '&#160;');
        text = text.replace(/(@\w+)/g, '<span class="keyword">$1</span>');
        $(elem).html(text);
    });
});
