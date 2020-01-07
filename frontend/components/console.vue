<template>
  <div
    v-if="hosts && hosts.length"
    id="console-panel"
    :class="isShow ? 'active' : ''"
  >
    <div id="console-header">
      <div id="console-title">
        <template v-if="isShow">
          <span class="tab-icon">
            <i class="icon-fn icon-close" @click="closeAll()" />
            <i
              class="icon-fn"
              :class="{ 'icon-min': isShow }"
              @click="hide()"
            />
            <i class="icon-fn icon-max" @click="show()" />
          </span>
          <span @click="toggle" class="title"
            ><i class="el-icon-download" />终端控制台</span
          >
        </template>
        <template v-else>
          <span class="tab-icon">
            <i class="icon-fn icon-close" @click="closeAll()" />
            <i
              class="icon-fn"
              :class="{ 'icon-min': isShow }"
              @click="hide()"
            />
            <i class="icon-fn icon-max" @click="show()" />
          </span>
          <span @click="toggle" class="title"
            ><i class="el-icon-download" />终端控制台</span
          >
        </template>
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
            :rows="35"
            class="terminal"
            :ref="'terminal-' + v.id"
            :host="v"
          />
        </el-tab-pane>
      </el-tabs>
    </div>
  </div>
</template>

<script>
import { mapGetters, mapActions } from 'vuex'
import Terminal from './terminal'

export default {
  components: {
    Terminal
  },
  data() {
    const hosts = this.$store.getters['workspace/hosts'] || []

    return {
      active: true,
      activeName: hosts.length ? hosts[0].id : ''
    }
  },
  computed: {
    ...mapGetters({
      isShow: 'console/isShow',
      currentHostId: 'console/currentHostId',
      hosts: 'console/hosts'
    })
  },
  watch: {
    hosts(val) {
      if (val && val.length) {
        this.show()
      }
    },
    currentHostId(val) {
      this.activeName = val
      this.show()
    }
  },
  methods: {
    ...mapActions({
      toggle: 'console/toggle',
      show: 'console/show',
      hide: 'console/hide',
      removeHost: 'console/removeHost'
    }),
    handleTabsEdit(hostId, action) {
      switch (action) {
        case 'remove':
          this.close(hostId)
      }
    },
    close(hostId) {
      const terminalInstance = this.$refs['terminal-' + hostId][0]

      if (terminalInstance) {
        terminalInstance.dispose()
      }

      this.removeHost(hostId)
    },
    // 关闭终端，断开所有连接
    closeAll() {
      const hosts = this.hosts.map(v => v)
      for (const host of hosts) {
        this.close(host.id)
      }
      this.hide()
    }
  }
}
</script>

<style lang="less">
@bg-color: #2e323b;
@bg-color-darken1: darken(@bg-color, 10%);

#console-panel {
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

#console-panel {
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

    .tab-icon {
      position: absolute;
      left: 5px;

      .icon-fn {
        width: 12px;
        height: 12px;
        border-radius: 50%;
        display: inline-block;
        cursor: pointer;
        background-color: grey;

        &.icon-close {
          background-color: red;
        }

        &.icon-min {
          background-color: yellow;
        }

        &.icon-max {
          background-color: green;
        }
      }
    }

    .title {
      cursor: pointer;
    }
  }
}

.terminal {
  height: @console-height;
  background-color: #000;
}
</style>
