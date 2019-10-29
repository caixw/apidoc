'use strict';

window.onload = function () {
    registerMethodFilter();
    registerTagFilter();
    registerLanguageFilter();
    initExample()
};

function registerMethodFilter() {
    const menu = document.querySelector('.methods-selector');

    menu.style.display = 'block';

    const list = menu.querySelectorAll('li input');
    list.forEach((val) => {
        val.addEventListener('change', (event) => {
            if (event.target === null) {
                return;
            }

            const chk = event.target.checked;
            const method = event.target.parentNode.parentNode.getAttribute('data-method');

            const apis = document.querySelectorAll('.api');
            apis.forEach((api) => {
                if (api.getAttribute('data-method') !== method) {
                    return;
                }

                api.setAttribute("data-hidden-method", chk ? "" : "true");
                toggleAPIVisible(api);
            });
        });
    });
}

function registerTagFilter() {
    const menu = document.querySelector('.tags-selector');

    menu.style.display = 'block';

    const list = menu.querySelectorAll('li input');
    list.forEach((val) => {
        val.addEventListener('change', (event) => {
            if (event.target === null) {
                return;
            }

            const chk = event.target.checked;
            const tag = event.target.parentNode.parentNode.getAttribute('data-tag');

            const apis = document.querySelectorAll('.api');
            apis.forEach((api) => {
                if (!api.getAttribute('data-tags').includes(tag + ',')) {
                    return;
                }

                api.setAttribute("data-hidden-tag", chk ? "" : "true");
                toggleAPIVisible(api);
            });
        });
    });
}

function toggleAPIVisible(api) {
    const hidden = api.getAttribute('data-hidden-tag') === 'true' ||
        api.getAttribute('data-hidden-server') === 'true' ||
        api.getAttribute('data-hidden-method') === 'true';

    api.style.display = hidden ? 'none' : 'block';
}

function registerLanguageFilter() {
    const menu = document.querySelector('.languages-selector');

    menu.style.display = 'block';

    const list = menu.querySelectorAll('li input');
    list.forEach((val) => {
        val.addEventListener('change', (event) => {
            if (event.target === null || !event.target.checked) {
                return;
            }

            const lang = event.target.parentNode.parentNode.getAttribute('lang');
            const elems = document.querySelectorAll('[data-locale]');

            elems.forEach((elem) => {
                if (elem.getAttribute('lang') === lang) {
                    elem.className = '';
                } else {
                    elem.className = 'hidden';
                }
            });
        });
    });
}

// 初始化示例代码的初始化功能
function initExample() {
    const buttons = document.querySelectorAll('.toggle-example');

    buttons.forEach((btn)=>{
        btn.addEventListener('click', (event)=> {
            if (event.target === null) {
                return;
            }

            const parent = event.target.parentNode.parentNode.parentNode;
            const table = parent.querySelector('table');
            const pre = parent.querySelector('pre');

            if (table.getAttribute('data-visible') === 'true') {
                table.setAttribute('data-visible', 'false');
                table.style.display = 'none';
            } else {
                table.setAttribute('data-visible', 'true');
                table.style.display = 'table';
            }

            if (pre.getAttribute('data-visible') === 'true') {
                pre.setAttribute('data-visible', 'false');
                pre.style.display = 'none';
            } else {
                pre.setAttribute('data-visible', 'true');
                pre.style.display = 'block';
            }
        });
    });
}
