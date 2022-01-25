// Full spec-compliant TodoMVC with localStorage persistence
// and hash-based routing in ~120 effective lines of JavaScript.
console.log("******************************* v0.0.8")

document.addEventListener("DOMContentLoaded", () => {
    // app Vue instance
    const app = Vue.createApp({
        name: "VuePincode",
        delimiters: ['${', '}'],
        data() {
            return {
                pincode: "",
                pincodeError: false,
                pincodeSuccess: false
            };
        },
        computed: {
            pincodeLength() {
                return this.pincode.length;
            },
            buttonDisabled() {
                return this.pincodeError || this.pincodeSuccess;
            }
        },
        watch: {
            pincode() {
                if (this.pincodeLength === 6) {
                    this.$emit("pincode", this.pincode);
                }
            }
        },
        destroyed() {
            this.resetPincode();
        },
        methods: {
            submit() {
                console.log(this.pincode);
            },
            clickPinButton(pressedNumber) {
                if (this.pincodeLength < 6) {
                    this.pincode = this.pincode + pressedNumber;
                }
            },
            resetPincode() {
                this.pincode = "";
                this.pincodeError = false;
                this.pincodeSuccess = false;
            },
            triggerMiss() {
                this.pincodeError = true;
                setTimeout(() => this.resetPincode(), 800);
            },
            triggerSuccess() {
                this.pincodeSuccess = true;
                setTimeout(() => this.resetPincode(), 2500);
            }
        }
    });

    app.config.productionTip = false;
    app.config.devtools = false;
    app.mount(".todoapp");
});
