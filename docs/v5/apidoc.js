'use strict';

window.onload = function () {
    registerFilter('method');
    registerFilter('server');
    registerFilter('tag');
    registerExpand();
    registerLanguageFilter();

    initExample();

    prettyDescription();
};

function registerFilter(type) {
    const menu = document.querySelector('.' + type + '-selector');
    if (menu === null) { // 可能为空，表示不存在该过滤项
        return;
    }

    menu.style.display = 'block'; // 有 JS 的情况下，展示过滤菜单

    menu.querySelectorAll('li input').forEach((val) => {
        val.addEventListener('change', (event) => {
            if (event.target === null) {
                return;
            }

            const chk = event.target.checked;
            const tag = event.target.parentNode.parentNode.getAttribute('data-' + type);
            document.querySelectorAll('.api').forEach((api) => {
                if (!api.getAttribute('data-' + type).includes(tag + ',')) {
                    return;
                }

                api.setAttribute("data-hidden-" + type, chk ? "" : "true");

                const hidden = api.getAttribute('data-hidden-tag') === 'true' ||
                    api.getAttribute('data-hidden-server') === 'true' ||
                    api.getAttribute('data-hidden-method') === 'true';
                api.style.display = hidden ? 'none' : 'block';
            });
        }); // end addEventListener('change')
    }); // end forEach('li input')
}

function registerLanguageFilter() {
    const menu = document.querySelector('.languages-selector');

    menu.style.display = 'block';

    menu.querySelectorAll('li input').forEach((val) => {
        val.addEventListener('change', (event) => {
            if (event.target === null || !event.target.checked) {
                return;
            }

            const lang = event.target.parentNode.parentNode.getAttribute('lang');
            document.querySelectorAll('[data-locale]').forEach((elem) => {
                elem.className = elem.getAttribute('lang') === lang ? '' : 'hidden';
            });
        }); // end addEventListener('change')
    }); // end forEach('li input')
}

function registerExpand() {
    const expand = document.querySelector('.expand-selector');
    if (expand === null) {
        return;
    }

    expand.style.display = 'block';

    expand.querySelector('input').addEventListener('change', (event) => {
        const chk = event.target.checked;
        document.querySelectorAll('details').forEach((elem) => {
            elem.open = chk;
        });
    });
}

function initExample() {
    document.querySelectorAll('.toggle-example').forEach((btn) => {
        btn.addEventListener('click', (event) => {
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

// 美化描述内容
//
// 即将 html 内容转换成真的 HTML 格式，而 markdown 则依然是 pre 显示。
function prettyDescription() {
    document.querySelectorAll('[data-type]').forEach((elem) => {
        const type = elem.getAttribute('data-type');
        if (type !== 'html') {
            return;
        }

        elem.innerHTML = elem.getElementsByTagName('pre')[0].innerText;
    });
}
