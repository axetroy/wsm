<template>
  <replayer
    :rows="0"
    class="terminal"
    ref="terminal"
    :records="record.records"
  />
</template>

<script>
import Replayer from '~/components/replayer'

export default {
  layout: 'empty',
  components: {
    Replayer
  },
  async asyncData({ store, params, query, $axios }) {
    const hostConnectionRecordId = params.id

    // const profile = store.getters['user/profile']

    const { data } = await $axios.$get(
      '/host/connection/' + hostConnectionRecordId
    )

    return {
      record: data
    }
  }
}
</script>

<style>
#__layout {
  background-color: #000;
}
</style>
