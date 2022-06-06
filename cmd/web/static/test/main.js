const identifier = 'my-vue-app'

const {loadModule} = window['vue3-sfc-loader'];

const load = (() => {
    const options = {
        moduleCache: {
            vue: Vue,
            vuex: Vuex
        },

        async getFile(url) {
            const response = await fetch(url);
            if (!response.ok)
                throw Object.assign(new Error(res.statusText + ' ' + url), {response});
            return {
                getContentData: asBinary => asBinary ? response.arrayBuffer() : response.text(),
            }
        },

        addStyle(textContent) {
            const style = Object.assign(document.createElement('style'), {textContent});
            const ref = document.head.getElementsByTagName('style')[0] || null;
            document.head.insertBefore(style, ref);
        },
    };
    return path => Vue.defineAsyncComponent(() => loadModule(path, options));
})();


const setCache = (object, keySuffix) => {
    const key = `${identifier}${keySuffix ? `-${keySuffix}` : ''}`;
    const value = JSON.stringify(object);
    return window.localStorage.setItem(key, value);
};

const getCache = (keySuffix) => {
    return JSON.parse(window.localStorage.getItem(`${identifier}${keySuffix ? `-${keySuffix}` : ''}`));
};

const CachedState = load('/static/test/currencies/CachedState.vue');
const CurrencySection = load('/static/test/currencies/Index5.vue');

const store = Vuex.createStore({
    state() {
        const cached = getCache();
        if (cached) return cached;
        else return {
            count: 0
        }
    },
    mutations: {
        increment(state) {
            state.count++
        },
        cache: setCache,
    },
    actions: {
        add(context) {
            context.commit('increment');
            context.commit('cache');
        }
    }
});

const app = Vue.createApp({
    components: {
        CachedState,
        CurrencySection,
    },
    template: `<cached-state></cached-state><CurrencySection></CurrencySection>`,
    store
});

window[identifier] = {
    ...window[identifier],
    app,
    store
}

window[identifier].app.mount('#app');

