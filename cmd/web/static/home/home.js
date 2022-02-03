"use strict";

// Full spec-compliant TodoMVC with localStorage persistence
// and hash-based routing in ~120 effective lines of JavaScript.
document.addEventListener("DOMContentLoaded", () => {
    // app Vue instance
    const app = Vue.createApp({
        name: "metamask",
        delimiters: ['${', '}'],
        mounted() {
            let t = this;
            if (typeof window.ethereum !== 'undefined') {
                t.metamaskSuccess = true;
            }
            window.ethereum.on('accountsChanged', (newAccounts) => {
                t.account = newAccounts;
                t.updateButton();
            });
        },
        data() {
            return {
                errorTodo: "",
                account: "",
                metamaskSuccess: false,
            };
        },
        watch: {
        },
        methods: {
            clickButton() {
                let t = this;
                this.getEtherAccount()
                    .then(accounts => {
                        if (accounts && accounts.length > 0) {
                            t.account = accounts[0];
                            t.updateButton();
                            t.callAPILogin();
                            return
                        }
                        t.errorTodo = "Account metamask not found";
                        setTimeout(() => t.resetErrorMessageForAPI(), 5000);
                        this.metamaskSuccess = false;
                    })
                    .catch(err => {
                        console.log(err)
                        t.errorTodo = "The authorization request to Metamask is made. Validate authentication on metamask popup that appeared";
                        setTimeout(() => t.resetErrorMessageForAPI(), 5000);
                        return
                    });
            },
            updateButton() {
                this.metamaskSuccess = false;
            },
            async getEtherAccount() {
                return await ethereum.request({ method: 'eth_requestAccounts' });
            },
            resetErrorMessageForAPI() {
                this.errorTodo = '';
            },
            callAPILogin() {
                this.errorTodo = '... wait ...';
                let t = this;
                return axios.post('/api/metamask', {
                    account: t.account,
                })
                .then(function (response) {
                    if (response.data && response.data.status === 'ok') {
                        window.location.replace("/wallet");
                    } else {
                        t.errorSend = 'Error! Could not reach the API. ' + error;
                        setTimeout(() => t.resetErrorMessageForAPI(), 3500);
                    }
                })
                .catch(function (error) {
                    t.errorSend = 'Error! Could not reach the API. ' + error;
                    setTimeout(() => t.resetErrorMessageForAPI(), 3500);
                });
            },

        }
    });

    app.config.productionTip = false;
    app.config.devtools = false;
    app.mount(".todoapp");
});
