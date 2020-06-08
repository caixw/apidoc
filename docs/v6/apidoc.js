'use strict';

window.onload = function () {
    registerFilter('method');
    registerFilter('server');
    registerFilter('tag');
    registerExpand();
    registerLanguageFilter();
    prettyDescription();
    initGotoTop();
};

const uncategorizedID = '_____';

function registerFilter(type) {
    const menu = document.querySelector('.' + type + '-selector');
    if (menu === null) { // 可能为空，表示不存在该过滤项
        return;
    }

    menu.style.display = 'block'; // 有 JS 的情况下，展示过滤菜单

    menu.querySelectorAll('ul>li').forEach((e) => {
        if (e.getAttribute(`data-${type}`) === '') {
            e.setAttribute(`data-${type}`, uncategorizedID);
        }
    });

    const apis = document.querySelectorAll('.api');
    apis.forEach((api) => { // 复制 data-${type} 至 data-hidden-${type}
        let attr = api.getAttribute(`data-${type}`);
        if (attr === '') {
            attr = uncategorizedID;
            api.setAttribute(`data-${type}`, attr);
        }
        api.setAttribute(`data-hidden-${type}`, attr);
    });

    menu.querySelectorAll('li input').forEach((val) => {
        val.addEventListener('change', (event) => {
            if (event.target === null) {
                return;
            }

            const chk = event.target.checked;
            const tag = event.target.parentNode.parentNode.getAttribute('data-' + type);
            apis.forEach((api) => {
                const attr = api.getAttribute(`data-${type}`);
                if (!attr.split(',').includes(tag)) {
                    return;
                }

                const hAttr = api.getAttribute(`data-hidden-${type}`).split(',');
                const index = hAttr.indexOf(tag);
                if (chk) {
                    if (-1 === index) {
                        hAttr.push(tag);
                        api.setAttribute(`data-hidden-${type}`, hAttr.join(','));
                    }
                } else {
                    if (-1 < index) {
                        hAttr.splice(index, 1);
                        api.setAttribute(`data-hidden-${type}`, hAttr.join(','));
                    }
                }

                const hidden = api.getAttribute('data-hidden-tag') === '' ||
                    api.getAttribute('data-hidden-server') === '' ||
                    api.getAttribute('data-hidden-method') === '';
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
        document.querySelectorAll('details.api').forEach((elem) => {
            elem.open = chk;
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
