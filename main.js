"use strict";

$(document).ready(function(){
    // 按4的倍数去掉pre中多余的空格
    $('pre').each(function(index, elem){
        var text = $(elem).html().replace(/    /g,'');
        $(elem).html(text);
    });
});
