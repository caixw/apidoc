'use strict';

window.onload = function () {
    this.registerMethodFilter();
    this.registerTagFilter();
};

function registerMethodFilter() {
    const list = document.querySelectorAll('.methods-selector li input');
    list.forEach((val) => {
        val.addEventListener('change', (event) => {
            const chk = event.target.checked;
            const method = event.target.parentNode.parentNode.getAttribute('data-method');

            const apis = this.document.querySelectorAll('.api');
            apis.forEach((api) => {
                if (api.getAttribute('data-method') !== method) {
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
