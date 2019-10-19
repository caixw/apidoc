'use strict';

window.onload = function() {
    initGotoTop();
};

function initGotoTop() {
    const top = document.querySelector('.goto-top');

    // 在最顶部时，隐藏按钮
    window.addEventListener('scroll', (e) => {
        const body = document.querySelector('html');
        if (body.scrollTop > 50) {
            top.style.display = 'block';
        } else {
            top.style.display = 'none';
        }
    });
}
