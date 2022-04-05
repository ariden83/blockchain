"use strict";

// https://github.com/pokhrelashok/vue-step-progress-indicator/blob/main/src/vue-step-progress-indicator.vue
// Full spec-compliant TodoMVC with localStorage persistence
// and hash-based routing in ~120 effective lines of JavaScript.
document.addEventListener("DOMContentLoaded", () => {
    // app Vue steps
    const app = Vue.createApp({
        name: "ProgressBar",
        delimiters: ['${', '}'],
        props: {
            reactivityType: {
                type: String,
                required: false,
                default: "single-step",
                validator: (propValue) => {
                    const types = ["all", "backward", "forward", "single-step"];
                    return types.includes(propValue);
                },
            },
            showLabel: {
                type: Boolean,
                required: false,
                default: true,
            },
            showBridge: {
                type: Boolean,
                required: false,
                default: false,
            },
            showBridgeOnSmallDevices: {
                type: Boolean,
                required: false,
                default: true,
            },
            colors: {
                type: Object,
                required: false,
                default: function () {
                    return {};
                },
            },
            styles: {
                type: Object,
                required: false,
                default: function () {
                    return {};
                },
            },
        },
        data() {
            return {
                steps: [
                    'Start',
                    'Minage',
                    'Création',
                    'Validation',
                    'Done',
                ],
                ws: "",
                delai: 15000,
                errorToken : "",
                inactiveColor: false,
                activeStep: 0,
                maxRetry: 5,
                currentRetry: 0,
                isReactive: true,
                currentStep: this.activeStep,
                styleData: {
                    progress__wrapper: {
                        flexWrap: "wrap",
                        display: "flex",
                        justifyContent: "flex-start",
                        margin: "1rem 0",
                    },
                    progress__block: {
                        display: "flex",
                        alignItems: "center",
                    },
                    progress__bridge: {
                        backgroundColor: "grey",
                        height: "2px",
                        flexGrow: "1",
                        width: "20px",
                    },
                    progress__bubble: {
                        margin: "0",
                        padding: "0",
                        lineHeight: "30px",
                        display: "flex",
                        justifyContent: "center",
                        alignItems: "center",
                        height: "30px",
                        width: "30px",
                        borderRadius: "100%",
                        backgroundColor: "transparent",
                        border: "2px grey solid",
                        textAlign: "center",
                    },
                    progress__label: {
                        margin: "0 0.8rem",
                    },
                },
                colorData: {
                    progress__bubble: {
                        active: {
                            color: "#fff",
                            backgroundColor: "#e74c3c",
                            borderColor: "#e74c3c",
                        },
                        inactive: {
                            color: "#fff",
                            backgroundColor: "#34495e",
                            borderColor: "#34495e",
                        },
                        completed: {
                            color: "#fff",
                            borderColor: "#27ae60",
                            backgroundColor: "#27ae60",
                        },
                    },
                    progress__label: {
                        active: {
                            color: "#e74c3c",
                        },
                        inactive: {
                            color: "#34495e",
                        },
                        completed: {
                            color: "#27ae60",
                        },
                    },
                    progress__bridge: {
                        active: {
                            backgroundColor: "#e74c3c",
                        },
                        inactive: {
                            backgroundColor: "#34495e",
                        },
                        completed: {
                            backgroundColor: "#27ae60",
                        },
                    },
                },
            }
        },
        methods: {
            callStartCreateToken: function () {
                this.callPageChange(0);
            },
            callPageChange: function (step) {
                if (!this.isReactive) return;
                this.currentStep = step;
                // this.$emit("onStepChanged", step);
                this.isReactive = false;
                this.callAPICreateToken();
                // if (step == this.steps.length - 1) this.$emit("onEnterFinalStep", step);
            },
            isActive: function (index) {
                return index === this.currentStep;
            },
            getColors: function (className, index) {
                let styles = {};
                if (index < this.currentStep) {
                    styles["color"] = this.colorData[className]["completed"]["color"];
                    styles["backgroundColor"] = this.inactiveColor
                        ? this.inactiveColor
                        : this.colorData[className]["completed"]["backgroundColor"];
                    styles["borderColor"] = this.colorData[className]["completed"]["borderColor"];
                } else if (index == this.currentStep) {
                    styles["color"] = this.colorData[className]["active"]["color"];
                    styles["backgroundColor"] = this.colorData[className]["active"]["backgroundColor"];
                    styles["borderColor"] = this.colorData[className]["active"]["borderColor"];
                } else {
                    styles["color"] = this.colorData[className]["inactive"]["color"];
                    styles["backgroundColor"] = this.colorData[className]["inactive"]["backgroundColor"];
                    styles["borderColor"] = this.colorData[className]["inactive"]["borderColor"];
                }
                return styles;
            },
            overwriteStyle: function (style1, style2) {
                for (const property in style1) {
                    if (Object.hasOwnProperty.call(style1, property)) {
                        const element = style1[property];
                        for (const value in element) {
                            if (Object.hasOwnProperty.call(element, value)) {
                                style2[property][value] = element[value];
                            }
                        }
                    }
                }
                return style2;
            },
            retryCallApi() {
                this.currentRetry ++;
                if ( this.currentRetry >= this.maxRetry) {
                    return
                }
                this.callPageChange(0);
            },
            websocket() {
                let t = this;
                if (this.ws) {
                    return false;
                }
                let loc = window.location, new_uri;
                if (loc.protocol === "https:") {
                    new_uri = "wss:";
                } else {
                    new_uri = "ws:";
                }
                new_uri += "//" + loc.host;
                new_uri += '/p/block/state';
                console.log(new_uri);
                this.ws = new WebSocket(new_uri);
                this.ws.onopen = function(evt) {
                    console.log("OPEN");
                    t.ws.send("test");
                }
                this.ws.onclose = function(evt) {
                    console.log("CLOSE");
                    t.ws = null;
                }
                this.ws.onmessage = function(evt) {
                    console.log("RESPONSE: " + evt.data);
                }
                this.ws.onerror = function(evt) {
                    console.log("ERROR: " + evt.data);
                }
                return false;
            },
            callAPICreateToken() {
                let t = this;
                t.currentStep = 1;
                t.websocket();
            },
        },
        watch: {
            activeStep: function (newVal) {
                if (this.activeStep < this.steps.length) this.currentStep = newVal;
            },
        },
        mounted() {
            this.styleData = this.overwriteStyle(this.styles, this.styleData);
            this.colorData = this.overwriteStyle(this.colors, this.colorData);
        },
    });

    app.config.productionTip = false;
    app.config.devtools = false;
    app.mount(".todoapp");
});
