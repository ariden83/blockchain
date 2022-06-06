"use strict";

import binance from './helpers/binance.js'
import chart from './helpers/chart.js'
import store from './store/index.js'

// https://github.com/tvjsx/trading-vue-demo
// Full spec-compliant TodoMVC with localStorage persistence
// and hash-based routing in ~120 effective lines of JavaScript.
document.addEventListener("DOMContentLoaded", () => {

    // import chart from '/static/test/helpers/chart.js'

    // app Vue instance
    const app = Vue.createApp({
        name: 'Chart',
        delimiters: ['${', '}'],
        props: {
            symbol: {
                type: String,
                default: null
            },
            priceScale: {
                type: Number,
                default: null
            },
            priceUnit: {
                type: String,
                default: null
            },
            tick: {
                type: [String, Number],
                default: null
            }
        },
        computed: {
            interval () {
                return this.$store.state.graphOptions.interval
            }
        },
        watch: {
            tick (value) {
                if (!this.fn) {
                    return false
                }
                this.fn.update(value)
            }
        },
        methods: {
            onSuccess (response) {
                this.loading = false
                const ticks = response.data.map(([t, o, h, l, c]) => [t, c])
                const [low, high] = ticks.reduce((acc, val) => {
                    val = parseFloat(val[1])
                    acc[0] = (acc[0] === undefined || val < acc[0]) ? val : acc[0]
                    acc[1] = (acc[1] === undefined || val > acc[1]) ? val : acc[1]
                    return acc
                }, [])
                this.$store.commit('UPDATE_GRAPH_STATS', {
                    low: low,
                    high: high
                })
                this.$nextTick(() => {
                    this.fn = chart(ticks, {
                        priceScale: this.priceScale,
                        priceUnit: this.priceUnit,
                        width: document.getElementById('price-chart').offsetWidth
                    })
                })
            },
            request () {
                this.loading = true
                binance.klines(this.symbol, this.interval).then(this.onSuccess).catch(console.log)
            },
            onClick (event) {
                const selection = event.target.innerHTML
                if ((this.interval === selection) || (this.intervals.indexOf(selection) === -1)) {
                    return false
                }
                this.$store.commit('UPDATE_GRAPH_OPTIONS', {
                    interval: event.target.innerHTML
                })
                this.request()
            }
        },
        mounted: function () {
            this.request(this.symbol, this.interval)
        },
        data: () => ({
            loading: false,
            fn: null,
            intervals: [
                '1m',
                '5m',
                '15m',
                '1h',
                '4h',
                '6h',
                '1d',
                '1w'
            ]
        })
    });

    app.config.productionTip = false;
    app.config.devtools = false;
    app.mount(".todoapp").use(store);
});