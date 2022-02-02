"use strict";

// Full spec-compliant TodoMVC with localStorage persistence
// and hash-based routing in ~120 effective lines of JavaScript.
document.addEventListener("DOMContentLoaded", () => {
    // app Vue instance
    const app = Vue.createApp({
        name: "metamask",
        delimiters: ['${', '}'],
        mounted() {
            if (typeof window.ethereum !== 'undefined') {
                this.metamaskSuccess = true;
            }
        },
        data() {
            return {
                errorTodo: "",
                metamaskSuccess: false,
            };
        },
        watch: {
        },
        methods: {
            clickButton() {
                ethereum.request({ method: 'eth_requestAccounts' });
            },
        }
    });

    app.config.productionTip = false;
    app.config.devtools = false;
    app.mount(".todoapp");
});
