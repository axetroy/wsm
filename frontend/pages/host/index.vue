<template>
  <div class="main">
    <el-card shadow="never">
      <div slot="header">
        <h4>可操作服务器列表</h4>
        <nuxt-link
          v-if="
            currentWorkspace
              ? memberProfile.role === 'owner' ||
                memberProfile.role === 'administrator'
              : true
          "
          to="/host/mutation"
          ><el-button type="primary" size="small" round
            >新增</el-button
          ></nuxt-link
        >
      </div>

      <el-table :data="data" border style="width: 100%">
        <el-table-column prop="name" label="名称" width="180">
        </el-table-column>
        <el-table-column prop="host" label="服务器" width="180">
        </el-table-column>
        <el-table-column prop="port" label="端口" width="180">
        </el-table-column>
        <el-table-column prop="username" label="用户名" width="180">
        </el-table-column>
        <el-table-column prop="remark" label="备注"> </el-table-column>
        <el-table-column label="创建时间">
          <template slot-scope="scope">
            {{ scope.row.created_at | dateformat }}
          </template>
        </el-table-column>
        <el-table-column label="操作">
          <template slot-scope="scope">
            <el-button type="text" size="small" @click="test(scope.row)">
              测试
            </el-button>
            <el-button type="text" size="small" @click="connect(scope.row)">
              连接
            </el-button>
            <nuxt-link :to="'/host/mutation?id=' + scope.row.id">
              <el-button type="text" size="small">编辑</el-button>
            </nuxt-link>
            <el-popconfirm
              title="你确定要删除这个服务器吗? 该操作不可恢复"
              v-on="{ onConfirm: () => deleteHost(scope.row) }"
            >
              <el-button type="text" size="small" slot="reference">
                删除
              </el-button>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          @current-change="changePage"
          background
          layout="prev, pager, next"
          :page-size="meta.limit"
          :total="meta.total"
        >
        </el-pagination>
      </div>
    </el-card>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'

export default {
  async asyncData({ $axios, store }) {
    const currentWorkspace = store.getters['workspace/current']

    if (currentWorkspace) {
      const [{ meta, data }, { data: memberProfile }] = await Promise.all([
        $axios.$get(`/team/_/${currentWorkspace}/host`),
        $axios.$get(`/team/_/${currentWorkspace}/profile`)
      ])

      return {
        data,
        meta,
        memberProfile
      }
    } else {
      const { meta, data } = await $axios.$get('/host')

      return {
        data,
        meta,
        memberProfile: {}
      }
    }
  },
  data() {
    return {}
  },
  computed: {
    ...mapGetters({
      currentWorkspace: 'workspace/current'
    })
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
    }
  },
  methods: {
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
      try {
        const { data } = await this.$axios.$get(`/shell/test/${host.id}`)

        if (data !== true) {
          this.$warning('服务器不可用')
          return
        }
      } catch (err) {
        this.$error(err.message)
        return false
      }

      this.$store.commit('console/APPEND_HOST', host) // 添加这个服务器
      this.$store.commit('console/SET_CURRENT_HOST', host.id) // 设置当前的 ID
      this.$store.commit('console/SET_CONSOLE_SHOW', true) // 显示 console
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