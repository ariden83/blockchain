"use strict";

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
                errorTodo: "",
                errorSend: "",
                terms: false,
                step: 1,
                pincode: "",
                pincode_confirm: "",
                pincodeError: false,
                pincodeSuccess: false,
                paraphrase: "",
                understood: false,
                saveseed: false,
                errorStep2: "",
                errorStep3: "",
                seedOver: false,
                seed: "",
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
                return this.pincodeError || this.pincodeSuccess;
            },
            buttonSubmitDisabled() {
                return this.step < 5;
            },
            buttonSubmitStep2Disabled() {
                return this.step == 6;
            },
            buttonSubmitStep3Disabled() {
                return this.step == 8;
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
                if (this.terms && this.step > 3) {
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
                if (this.terms) {
                    this.step = 5
                } else {
                    this.step = 4
                }
            },
            submit() {
                this.errorSend = '';
                if (this.step === 5 && this.pincodeLength == 6 && this.terms) {
                    this.callAPILogin(this.pincode)
                } else if (!this.terms) {
                    this.errorSend = 'Vous devez accepter les conditions générales';
                    setTimeout(() => this.resetErrorMessageForAPI(), 5000);
                }
            },
            async encrypt(code) {
                return await cryptGcm(this.paraphrase, code);
            },
            async decrypt(seed) {
                return await DecryptGcm(this.paraphrase, seed);
            },
            callAPILogin(code) {
                this.errorSend = '... wait ...';
                let t = this
                this.step = 3
                let cipher;

                return t.encrypt(code)
                    .then(c => {
                        cipher = c;
                        return grecaptcha.execute('6LfmdSAeAAAAAPf5oNQ1UV0wf6QhnH9dQFDSop7V', {action: 'submit'});
                    })
                    .then(token => axios.post('/api/inscription', {
                        password: cipher,
                        recaptcha: token,
                    }))
                    .then(function (response) {
                        if (response.data && response.data.status === 'ok') {
                            t.step = 6;
                            t.errorSend = '';
                            return t.decrypt(response.data.seed)
                            .then(seed => {
                                t.seed = seed;
                            });
                        } else {
                            t.errorSend = 'Error! Could not reach the API. ' + error;
                            console.log(error);
                            setTimeout(() => t.resetErrorMessageForAPI(), 3500);
                        }
                    })
                    .catch(function (error) {
                        console.log(error)
                        t.errorSend = 'Error! Could not reach the API. ' + error;
                        t.step = 1;
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
            validunderstood() {
                if (this.understood) {
                    this.step = 7;
                } else {
                    this.step = 6;
                }
            },
            submitStep2() {
                if (this.step == 7) {
                    this.step = 8;
                    this.cutSeed(this.seed)
                } else {
                    this.errorStep2 = "You must agree to the terms that you can no longer access your crypto wallet if you lose your recovery phrase";
                }
            },
            copierpress(){
                let x = document.createElement("INPUT");
                x.setAttribute("type", "text");
                x.setAttribute("id", "myseed");
                document.getElementById("copy").appendChild(x);
                document.getElementById("myseed").value = this.seed;
                let copyText = document.getElementById("myseed");
                copyText.select();
                document.execCommand("copy");
                document.getElementById("myseed").remove();
                this.errorStep2 = "Your Secret Recovery Phrase is copied to the clipboard";
            },
            cutSeed(str) {
                const words = str.split(' ');
                for (const w of words) {
                    let c = document.createElement("CANVAS");
                    let ctx = c.getContext("2d");
                    ctx.font = "30px Calibri";
                    ctx.setTransform((Math.random() / 10) + 0.9,    //scalex
                        0.1 - (Math.random() / 5),      //skewx
                        0.1 - (Math.random() / 5),      //skewy
                        (Math.random() / 10) + 0.9,     //scaley
                        (Math.random() * 3) + 3,      //transx
                        0);                           //transy
                    ctx.fillText(w, 10, 50);

                    let div = document.createElement("DIV");
                    div.className = "col-xs-3"
                    div.appendChild(c);
                    document.getElementById("my-seed").appendChild(div);
                }
            },
            validsaveseed() {
                if (this.saveseed) {
                    this.step = 9;
                } else {
                    this.step = 8;
                }
            },
            mouseOver(){
                this.seedOver = true;
            },
            mouseOut(){
                this.seedOver = false;
            },
            submitStep3() {
                var t = this;
                if (this.step == 9) {
                   return axios.post('/api/inscription/validate', {})
                    .then(function (response) {
                        if (response.data && response.data.status === 'ok') {
                            window.location.replace("/authorize");
                        } else {
                            t.errorSend = 'Error! Could not reach the API. ' + error;
                            console.log(error);
                            setTimeout(() => t.resetErrorMessageForAPI(), 3500);
                        }
                    })
                    .catch(function (error) {
                        console.log(error)
                        t.errorSend = 'Error! Could not reach the API. ' + error;
                        setTimeout(() => t.resetErrorMessageForAPI(), 3500);
                    });
                } else {
                    this.errorStep2 = "You must confirm that you have saved your Secret Recovery Phrase in a safe place";
                }
            },
        }
    });

    app.config.productionTip = false;
    app.config.devtools = false;
    app.mount(".todoapp");
});
