<template>
  <div class="main">
    <template v-if="currentWorkspace">
      <el-card shadow="never">
        <div slot="header">
          <h4>已发出的团队邀请</h4>
        </div>

        <el-table :data="data" border style="width: 100%">
          <el-table-column prop="id" label="ID" width="180" />
          <el-table-column label="用户">
            <el-table-column prop="user.id" label="ID" width="180" />
            <el-table-column prop="user.username" label="用户名" />
          </el-table-column>
          <el-table-column label="角色">
            <template slot-scope="scope">
              <span
                v-for="v of roles"
                :key="v.value"
                v-if="v.value === scope.row.role"
              >
                {{ v.label }}
              </span>
            </template>
          </el-table-column>
          <el-table-column label="状态">
            <template slot-scope="scope">
              <span
                v-for="v of state"
                :key="v.value"
                v-if="v.value === scope.row.state"
              >
                {{ v.label }}
              </span>
            </template>
          </el-table-column>
          <el-table-column prop="remark" label="备注"> </el-table-column>
          <el-table-column label="邀请时间">
            <template slot-scope="scope">
              {{ scope.row.created_at | dateformat }}
            </template>
          </el-table-column>
          <el-table-column label="操作">
            <template slot-scope="scope">
              <el-popconfirm
                v-if="
                  (memberProfile.role === 'owner' ||
                    memberProfile.role === 'administrator') &&
                    scope.row.state === 'init'
                "
                title="你确定要撤销这个邀请吗? 该操作不可恢复"
                v-on="{ onConfirm: () => cancelInvite(scope.row) }"
              >
                <el-button type="text" size="small" slot="reference">
                  撤销邀请
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
    </template>
    <template v-else>
      <el-card shadow="never">
        <div slot="header">
          <h4>我收到的团队邀请</h4>
        </div>

        <el-table :data="data" border style="width: 100%">
          <el-table-column prop="id" label="ID" width="180" />
          <el-table-column label="团队">
            <el-table-column prop="team.id" label="ID" width="180" />
            <el-table-column prop="team.name" label="团队名" />
            <el-table-column label="拥有者">
              <el-table-column prop="team.owner.id" label="ID" width="180" />
              <el-table-column prop="team.owner.username" label="用户名" />
            </el-table-column>
          </el-table-column>
          <el-table-column label="角色">
            <template slot-scope="scope">
              <span
                v-for="v of roles"
                :key="v.value"
                v-if="v.value === scope.row.role"
              >
                {{ v.label }}
              </span>
            </template>
          </el-table-column>
          <el-table-column label="状态">
            <template slot-scope="scope">
              <span
                v-for="v of state"
                :key="v.value"
                v-if="v.value === scope.row.state"
              >
                {{ v.label }}
              </span>
            </template>
          </el-table-column>
          <el-table-column prop="remark" label="备注"> </el-table-column>
          <el-table-column label="邀请时间">
            <template slot-scope="scope">
              {{ scope.row.created_at | dateformat }}
            </template>
          </el-table-column>
          <el-table-column label="操作">
            <template slot-scope="scope">
              <template v-if="scope.row.state === 'init'">
                <el-button
                  type="primary"
                  size="small"
                  slot="reference"
                  @click="handleInvite(scope.row, true)"
                >
                  接受
                </el-button>
                <el-button
                  size="small"
                  slot="reference"
                  @click="handleInvite(scope.row, false)"
                >
                  拒绝
                </el-button>
              </template>
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
    </template>
  </div>
</template>

<script>
import { mapGetters, mapActions } from 'vuex'

export default {
  async asyncData({ $axios, store }) {
    const currentWorkspace = store.getters['workspace/current']

    if (currentWorkspace) {
      const [{ meta, data }, { data: memberProfile }] = await Promise.all([
        $axios.$get(`/team/_/${currentWorkspace}/member/invite`),
        $axios.$get(`/team/_/${currentWorkspace}/profile`)
      ])

      return {
        data,
        meta,
        memberProfile
      }
    } else {
      const [{ meta, data }] = await Promise.all([$axios.$get('/team/invite')])

      return {
        data,
        meta,
        memberProfile: {}
      }
    }
  },
  data() {
    return {
      state: [
        {
          label: '全部',
          value: undefined
        },
        {
          label: '未处理',
          value: 'init'
        },
        {
          label: '已接受',
          value: 'accept'
        },
        {
          label: '已拒绝',
          value: 'refuse'
        },
        {
          label: '已撤销',
          value: 'cancel'
        },
        {
          label: '已弃用',
          value: 'deprecated'
        }
      ]
    }
  },
  computed: {
    ...mapGetters({
      currentWorkspace: 'workspace/current',
      roles: 'workspace/roles'
    })
  },
  watch: {
    currentWorkspace() {
      this.changePage(0)
    }
  },
  methods: {
    ...mapActions({
      getWorkspaces: 'workspace/getWorkspaces'
    }),
    async changePage(page) {
      const { meta, data } = await this.$axios.$get(
        this.currentWorkspace
          ? `/team/_/${this.currentWorkspace}/member/invite` // 团队的邀请列表
          : '/team/invite', // 我的收邀列表
        {
          params: {
            page
          }
        }
      )

      this.data = data
      this.meta = meta
    },
    // 撤销邀请
    async cancelInvite(invite) {
      const currentWorkspace = this.currentWorkspace
      try {
        await this.$axios.$delete(
          `/team/_/${currentWorkspace}/member/invite/_/${invite.id}`
        )
        this.$success('撤销成功')
        this.changePage(0)
      } catch (err) {
        this.$error(`撤销失败: ${err.message}`)
      }
    },
    // 处理邀请，是接受还是拒绝
    async handleInvite(invite, isAccept) {
      try {
        await this.$axios.$put(
          `/team/_/${invite.team.id}/member/invite/_/${invite.id}`,
          {
            state: isAccept ? 'accept' : 'refuse'
          }
        )
        this.$success('成功')
        this.changePage(0)

        if (isAccept === true) {
          this.getWorkspaces()
        }
      } catch (err) {
        this.$error(`失败: ${err.message}`)
      }
    }
  }
}
</script>
