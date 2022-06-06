<template>
  <main class="container my-4">
    <currencies-top v-if="!slug" />
    <currencies-detail v-else :slug="slug" />
  </main>
</template>

<script>
import { defineAsyncComponent } from 'vue'
export default {
  name: 'Currencies',
  components: {
    CurrenciesTop: defineAsyncComponent(() =>
      import(/* webpackChunkName: "currencies-top" */ '/static/test/currencies/Top1.vue')
    ),
    CurrenciesDetail: defineAsyncComponent(() =>
      import(/* webpackChunkName: "currencies-detail" */ '/static/test/currencies/Detail.vue')
    )
  },
  props: {
    slug: {
      type: String,
      default: null
    }
  }
}
</script>
