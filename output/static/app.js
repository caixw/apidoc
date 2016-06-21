"use strict";

$(document).ready(function(){
    // 隐藏不当前页面用不到的过滤器
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
});
