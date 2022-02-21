"use strict";

document.addEventListener("DOMContentLoaded", () => {
    const app = Vue.createApp({
        data() {
            return {
                title: "Gerador QRCode",
                text: "",
                qrcode: "",
            }
        },
        methods: {
            newQRCode() {
                let qrcode = new QRious({size: 300})
                qrcode.value = "test test test";
                this.qrcode = qrcode.toDataURL();
            },
        },
    });

    app.config.productionTip = false;
    app.config.devtools = false;
    app.mount(".todoapp");
});