"use strict";

// https://github.com/tvjsx/trading-vue-demo
// Full spec-compliant TodoMVC with localStorage persistence
// and hash-based routing in ~120 effective lines of JavaScript.
document.addEventListener("DOMContentLoaded", () => {
    // app Vue instance
    const app = Vue.createApp({
        name: "VuePincode",
        delimiters: ['${', '}'],
        data: {
            data: new TradingVueJs.DataCube(Data),
            width: window.innerWidth,
            height: window.innerHeight,
            night: true
        },
        mounted() {
            window.addEventListener('resize', this.onResize)
            window.DataCube = this.data
        },
        methods: {
            onResize(event) {
                this.width = window.innerWidth
                this.height = window.innerHeight
            }
        },
        computed: {
            colors() {
                return this.night ? {} : {
                    colorBack: '#fff',
                    colorGrid: '#eee',
                    colorText: '#333'
                }
            },
        },
        beforeDestroy() {
            window.removeEventListener('resize', this.onResize)
        },
    });

    app.config.productionTip = false;
    app.config.devtools = false;
    app.mount(".todoapp");
});