<template>
  <div id="host-page">
    <el-row :gutter="20">
      <el-col :span="12">
        <div style="margin-bottom: 20px;" v-if="editable">
          <nuxt-link to="/host/mutation">
            <el-button type="primary" size="small">添加服务器 </el-button>
          </nuxt-link>
        </div>

        <el-pagination
          class="text-center"
          style="margin-bottom: 20px;"
          @current-change="changePage"
          background
          layout="prev, pager, next"
          :page-size="meta.limit"
          :total="meta.total"
          hide-on-single-page
        />

        <el-scrollbar v-if="data.length" style="height: 580px;" :native="false">
          <el-card v-for="v of data" :key="v.id" style="margin-bottom: 20px;">
            <div slot="header" class="host-header">
              <i class="el-icon-s-promotion" />
              {{ v.name }}
            </div>

            <div class="meta-info">
              {{ v.remark || '赶紧给它备注一下吧' }}
            </div>
            <div class="action-block">
              <el-dropdown
                split-button
                type="primary"
                size="small"
                @click="connect(v)"
              >
                连接
                <el-dropdown-menu slot="dropdown">
                  <el-dropdown-item @click.native="test(v)"
                    >测试连接
                  </el-dropdown-item>
                  <el-dropdown-item>
                    <a
                      style="text-decoration: none;color: inherit;"
                      target="_blank"
                      :href="
                        '/tty/' +
                          v.id +
                          (currentWorkspace
                            ? '?team_id=' + currentWorkspace
                            : '')
                      "
                    >
                      新窗口连接
                    </a>
                  </el-dropdown-item>
                </el-dropdown-menu>
              </el-dropdown>

              <el-popconfirm
                v-if="editable"
                title="你确定要删除这个服务器吗? 该操作不可恢复"
                v-on="{ onConfirm: () => deleteHost(v) }"
              >
                <el-button type="danger" size="small" slot="reference">
                  删除
                </el-button>
              </el-popconfirm>
              <nuxt-link :to="'/host/mutation?id=' + v.id" v-if="editable">
                <el-button type="warning" size="small">修改</el-button>
              </nuxt-link>
            </div>
          </el-card>
        </el-scrollbar>
        <div v-else>
          暂无服务器，请管理员先添加服务器
        </div>
      </el-col>

      <el-col :span="12">
        <el-tabs v-model="hostTab">
          <el-tab-pane label="连接记录" name="connection">
            <el-pagination
              class="text-center"
              style="margin-bottom: 20px;"
              @current-change="changeConnectionPage"
              background
              layout="prev, pager, next"
              :page-size="connectionsMeta.limit"
              :total="connectionsMeta.total"
              hide-on-single-page
            />
            <el-scrollbar
              v-if="connections.length"
              style="height: 580px;"
              :native="false"
            >
              <el-card
                v-for="v of connections"
                :key="v.id"
                style="margin-bottom: 20px;"
              >
                <div slot="header">
                  <i class="el-icon-user" />
                  {{ v.user.username }} ({{ v.user.nickname }})
                </div>

                <div class="meta-info"></div>
                于 {{ v.created_at | dateformat }} 连接服务器 {{ v.host.name }}
                <div class="action-block">
                  <a :href="`/replay/${v.id}`" target="_blank">
                    <el-button type="primary" size="small">
                      <i class="el-icon-video-play" />
                      重放记录
                    </el-button>
                  </a>
                </div>
              </el-card>
            </el-scrollbar>
            <div v-else>
              暂无连接记录
            </div>
          </el-tab-pane>
        </el-tabs>
      </el-col>
    </el-row>
  </div>
</template>

<script>
import { mapGetters, mapActions } from 'vuex'

export default {
  async asyncData({ $axios, store }) {
    const currentWorkspace = store.getters['workspace/current']

    if (currentWorkspace) {
      const [
        { meta, data },
        { data: memberProfile },
        { data: connections, meta: connectionsMeta }
      ] = await Promise.all([
        $axios.$get(`/team/_/${currentWorkspace}/host`),
        $axios.$get(`/team/_/${currentWorkspace}/profile`),
        $axios.$get(`/team/_/${currentWorkspace}/connection`)
      ])

      return {
        data,
        meta,
        memberProfile,
        connections,
        connectionsMeta
      }
    } else {
      const [
        { meta, data },
        { data: connections, meta: connectionsMeta }
      ] = await Promise.all([
        $axios.$get('/host'),
        $axios.$get(`/host/connection`)
      ])

      return {
        data,
        meta,
        memberProfile: {},
        connections,
        connectionsMeta
      }
    }
  },
  data() {
    return {
      hostTab: 'connection',
      data: [],
      meta: {},
      connections: [],
      connectionsMeta: {}
    }
  },
  computed: {
    ...mapGetters({
      currentWorkspace: 'workspace/current'
    }),
    editable() {
      return this.currentWorkspace
        ? this.memberProfile.role === 'owner' ||
            this.memberProfile.role === 'administrator'
        : true
    }
  },
  watch: {
    currentWorkspace(currentWorkspace) {
      if (currentWorkspace) {
        this.changePage(0)
        this.$axios
          .$get(`/team/_/${currentWorkspace}/profile`)
          .then(({ data: memberProfile }) => {
            this.memberProfile = memberProfile
          })
      } else {
        this.changePage(0)
      }

      this.changeConnectionPage(0)
    }
  },
  methods: {
    ...mapActions({
      appendHost: 'console/appendHost',
      activeHost: 'console/activeHost',
      showConsole: 'console/show'
    }),
    async changePage(page) {
      const { meta, data } = await this.$axios.$get(
        this.currentWorkspace
          ? `/team/_/${this.currentWorkspace}/host`
          : '/host',
        {
          params: {
            page
          }
        }
      )

      this.data = data
      this.meta = meta
    },
    async changeConnectionPage(page) {
      const { meta, data } = await this.$axios.$get(
        this.currentWorkspace
          ? `/team/_/${this.currentWorkspace}/connection`
          : '/host/connection',
        {
          params: {
            page
          }
        }
      )

      this.connections = data
      this.connectionsMeta = meta
    },
    async deleteHost(host) {
      const currentWorkspace = this.currentWorkspace
      try {
        await this.$axios.$delete(
          currentWorkspace
            ? `/team/_/${currentWorkspace}/host/_/${host.id}`
            : `/host/_/${host.id}`
        )

        // get list
        this.$success('删除成功')
        const { meta, data } = await this.$axios.$get(
          currentWorkspace ? `/team/_/${currentWorkspace}/host` : '/host'
        )
        this.data = data
        this.meta = meta
      } catch (err) {
        this.$error(`删除失败: ${err.message}`)
      }
    },
    // 连接服务器，打开终端
    async connect(host) {
      this.appendHost(host) // 添加这个服务器
      this.activeHost(host) // 激活当前服务器
      this.showConsole()
    },
    // 测试服务器是否可用
    async test(host) {
      try {
        const { data } = await this.$axios.$get(`/shell/test/${host.id}`)
        if (data === true) {
          this.$success('服务器可用')
          return true
        } else {
          this.$warning('服务器不可用')
          return false
        }
      } catch (err) {
        this.$error(err.message)
        return false
      }
    }
  }
}
</script>

<style lang="less">
#host-page {
  .mb20 {
    margin-bottom: 20px;
    border-bottom: 1px solid #ebeef5;
    // padding-bottom: 20px;
  }

  .el-card__header {
    &.active {
      background: #419dfb;
      color: #fff;
    }
  }

  .el-card__body {
    label {
      display: inline-block;
      margin-bottom: 5px;
    }
  }

  .meta-info {
    color: #adadad;
  }

  .action-block {
    padding-top: 20px;
    /*border-top: 1px solid #ebeef5;*/
  }
}
</style>
