"use strict";


// Full spec-compliant TodoMVC with localStorage persistence
// and hash-based routing in ~120 effective lines of JavaScript.
document.addEventListener("DOMContentLoaded", () => {
    // app Vue instance
    const app = Vue.createApp({
        data: function() {
            return {
                random_pct: 0,
                increasing_pct: 0,
                decreasing_pct: 100
            }
        },
        mounted: function() {
            setInterval(() => {
                if (this.is_paused)
                    return

                this.random_pct = Math.ceil(Math.random() * 100)
                this.increasing_pct = Math.min(this.increasing_pct + 5, 100)
                this.decreasing_pct = Math.max(this.decreasing_pct - 5, 0)
            }, 2000)
        }
    });

    app.config.productionTip = false;
    app.config.devtools = false;
    app.mount(".todoapp");
});
