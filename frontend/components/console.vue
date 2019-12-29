<template>
  <div id="console-pannel" :class="consoleShow ? 'active' : ''">
    <div id="console-header">
      <div @click="toggle" id="console-title">
        <template v-if="consoleShow">
          <i class="el-icon-download" />终端控制台
        </template>
        <template v-else> <i class="el-icon-upload2" />终端控制台 </template>
      </div>
      <el-tabs
        id="console-tab"
        v-model="activeName"
        type="border-card"
        editable
        @edit="handleTabsEdit"
      >
        <el-tab-pane
          v-for="v in hosts"
          :label="v.name"
          :name="v.id"
          :key="v.id"
          @click="toggle"
        >
          <terminal
            class="terminal"
            :ref="'terminal-' + v.id"
            :host="v"
          ></terminal>
        </el-tab-pane>
      </el-tabs>
    </div>
  </div>
</template>

<script>
import 'xterm/css/xterm.css'
import Terminal from './terminal'
import { mapGetters } from 'vuex'

export default {
  components: {
    Terminal
  },
  data() {
    const hosts = this.$props.hosts || []

    return {
      active: true,
      activeName: hosts.length ? hosts[0].id : ''
    }
  },
  props: {
    hosts: {
      type: Array,
      default() {
        return []
      }
    },
    currentHost: {
      type: String,
      default() {
        return ''
      }
    }
  },
  computed: {
    ...mapGetters({
      consoleShow: 'console/consoleShow'
    })
  },
  watch: {
    hosts(val) {
      if (val && val.length) {
        this.$store.commit('console/SET_CONSOLE_SHOW', true)
      }
    },
    currentHost(val) {
      this.activeName = val
      this.$store.commit('console/SET_CONSOLE_SHOW', true)
      return
    }
  },
  methods: {
    toggle() {
      this.$store.commit('console/SET_CONSOLE_SHOW', !this.consoleShow)
    },
    handleTabsEdit(hostId, action) {
      switch (action) {
        case 'remove':
          const terminalInstance = this.$refs['terminal-' + hostId][0]

          if (terminalInstance) {
            terminalInstance.dispose()
            this.$store.commit('console/REMOVE_HOST', hostId)
            // const newHosts = this.props.slice(0)
          }
      }
    }
  }
}
</script>

<style lang="less">
@bg-color: #2e323b;
@bg-color-darken1: darken(@bg-color, 10%);

#console-pannel {
  .el-tabs__content {
    padding: 0;
    width: 100%;
  }
  .el-tabs--border-card {
    border: 0;
    .el-tabs__header {
      background-color: @bg-color;
      border-bottom: 1px solid @bg-color;
    }

    .el-tabs__item {
      background-color: #c5c5c5;
      color: #8b8b8b;
      border-bottom: 0;
      border-top: 0;
      &.is-active {
        background-color: white;
        color: #000;
      }
    }
  }
  .el-tabs__new-tab {
    display: none;
  }
}
</style>

<style lang="less" scoped>
@bg-color: #2e323b;
@bg-color-darken1: darken(@bg-color, 10%);
@title-height: 30px;
@tab-height: 40px;
@console-height: 595px;
@height: @title-height + @tab-height + @console-height;

#console-pannel {
  position: absolute;
  bottom: -@height + @title-height;
  left: 0;
  width: 100%;
  height: @height;
  transition: all 0.3s ease-in-out;
  z-index: 9999999;

  &.active {
    bottom: 0;
  }
}

#console-header {
  #console-tab {
    height: @tab-height;
  }

  #console-title {
    height: @title-height;
    background: @bg-color-darken1;
    color: #fff;
    line-height: @title-height;
    text-align: center;
  }
}

.terminal {
  height: @console-height;
  background-color: #000;
}
</style>
