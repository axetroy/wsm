<template>
  <el-dialog title="终端回放" :visible.sync="isShow" width="1255px">
    <replayer
      v-if="connection"
      :rows="20"
      :cols="135"
      class="terminal"
      ref="player"
      :records="connection.records"
    />
  </el-dialog>
</template>

<script>
import { mapActions } from 'vuex'
import Replayer from '~/components/replayer'

export default {
  components: {
    Replayer
  },
  props: {
    visible: {
      type: Boolean
    },
    connection: Object
  },
  data() {
    return {
      isShow: this.visible || false // 是否显示弹窗
    }
  },
  watch: {
    visible(visible) {
      this.isShow = visible
    },
    isShow(val) {
      if (!val) {
        this.$emit('update:connection', undefined)
      }
      this.$emit('update:visible', val)
    }
  },
  computed: {},
  methods: {
    ...mapActions({
      getWorkspaces: 'workspace/getWorkspaces'
    })
  }
}
</script>

<style lang="less" scope>
.tag {
  &:not(:last-child) {
    margin-right: 20px;
  }
}
</style>
