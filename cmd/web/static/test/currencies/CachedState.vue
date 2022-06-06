<template>
  <p>The button increments the count</p>
  <strong>Count: {{count}}</strong>
  <button @click="add()">yes</button>
</template>

<script>
const identifier = 'my-vue-app';

const store = window[identifier].store;

export default {
  name: "CachedState",
  computed: {
    count() {
      return store.state.count;
    }
  },
  methods: {
    add: () => store.dispatch('add'),
  },
  store,
}
</script>

<style scoped>
p {
  font-size: 16px;
}
</style>
