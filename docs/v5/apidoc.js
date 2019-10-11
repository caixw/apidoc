'use strict';

// 深色模式的颜色定义
const darkColor = {
    '--color': 'white',
    '--background': 'black',
    '--border-color': '#e0e0e0',
    '--delete-color': 'red',

    /* method */
    '--method-get-color': 'green',
    '--method-options-color': 'green',
    '--method-post-color': 'darkorange',
    '--method-put-color': 'darkorange',
    '--method-patch-color': 'darkorange',
    '--method-delete-color': 'red',
}

// 浅色模式的颜色定义，由代码从 css 中获取
const lightColor = {}

window.onload = function () {
    this.registerMethodFilter();
    this.registerTagFilter();
    this.initColors()
}

function registerMethodFilter() {
    const list = document.querySelectorAll('.methods-selector li input');
    list.forEach((val) => {
        val.addEventListener('change', (event) => {
            const chk = event.target.checked;
            const method = event.target.parentNode.parentNode.getAttribute('data-method');

            const apis = this.document.querySelectorAll('.api');
            apis.forEach((api) => {
                if (api.getAttribute('data-method') != method) {
                    return;
                }

                api.style.display = chk ? 'block' : 'none';
            });
        });
    });
}

function registerTagFilter() {
    const list = document.querySelectorAll('.tags-selector li input');
    list.forEach((val) => {
        val.addEventListener('change', (event) => {
            const chk = event.target.checked;
            const tag = event.target.parentNode.parentNode.getAttribute('data-tag');

            const apis = this.document.querySelectorAll('.api');
            apis.forEach((api) => {
                if (!api.getAttribute('data-tags').includes(tag + ',')) {
                    return;
                }

                api.style.display = chk ? 'block' : 'none';
            });
        });
    });
}

function initColors() {
    // 备份原来的颜色至 lightColor
    for (const key in darkColor) {
        lightColor[key] = document.documentElement.style.getPropertyValue(key);
    }

    const elem = document.querySelector('.colors-selector input');
    elem.addEventListener('change', (event) => {
        if (event.target.checked) {
            setColors(darkColor);
        } else {
            setColors(lightColor);
        }
    });
}

// 设置颜色值，可用的颜色值列表可参考 apidoc.css 中的 root 元素内容。
function setColors(colors) {
    const root = document.documentElement;
    for (const key in colors) {
        root.style.setProperty(key, colors[key]);
    }
}
