<template>
  <div class="main">
    <template v-if="currentWorkspace">
      <el-row :gutter="20">
        <el-col :span="6">
          <el-card>
            团队名称:
            <div>{{ stat.name }}</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card
            >团队成员:
            <div>{{ stat.member_num }}</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card
            >团队服务器:
            <div>{{ stat.host_num }}</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card>
            创建日期:
            <div>{{ stat.created_at | dateformat }}</div>
          </el-card>
        </el-col>
      </el-row>

      <el-card style="margin-top: 30px">
        <div slot="header">
          <h4>团队成员</h4>
          <nuxt-link
            v-if="
              memberProfile &&
                (memberProfile.role === 'owner' ||
                  memberProfile.role === 'administrator')
            "
            :to="`/team/${currentWorkspace}/invite`"
          >
            <el-button type="primary" size="small" round>邀请成员 </el-button>
          </nuxt-link>
        </div>

        <el-table :data="members" border style="width: 100%">
          <el-table-column prop="id" label="ID" />
          <el-table-column prop="username" label="用户名" />
          <el-table-column label="角色">
            <template slot-scope="scope">
              <span
                v-for="v in roles"
                :key="v.value"
                v-if="v.value === scope.row.role"
              >
                {{ v.label }}
              </span>
            </template>
          </el-table-column>
          <el-table-column label="加入时间">
            <template slot-scope="scope">
              {{ scope.row.created_at | dateformat }}
            </template>
          </el-table-column>
          <el-table-column label="操作">
            <template slot-scope="scope">
              <el-popconfirm
                v-if="
                  memberProfile &&
                    memberProfile.id !== scope.row.id &&
                    memberProfile.role === 'owner'
                "
                title="你确定要踢出团队吗"
                v-on="{ onConfirm: () => kickMemberFromTeam(scope.row) }"
              >
                <el-button type="text" size="small" slot="reference">
                  踢出团队
                </el-button>
              </el-popconfirm>
            </template>
          </el-table-column>
        </el-table>

        <div class="pagination">
          <el-pagination
            @current-change="changeTeamPage"
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
      <el-card>
        <div slot="header">
          <h4>所属团队</h4>
          <nuxt-link to="/team/mutation">
            <el-button type="primary" size="small" round>创建团队 </el-button>
          </nuxt-link>
        </div>

        <el-table :data="teams" border style="width: 100%">
          <el-table-column prop="id" label="ID" width="180" />
          <el-table-column prop="name" label="名称" width="180" />
          <el-table-column label="拥有者">
            <el-table-column prop="owner.id" label="ID" width="180" />
            <el-table-column prop="owner.username" label="用户名" />
          </el-table-column>
          <el-table-column label="角色">
            <template slot-scope="scope">
              <span
                v-for="v in roles"
                :key="v.value"
                v-if="v.value === scope.row.role"
              >
                {{ v.label }}
              </span>
            </template>
          </el-table-column>
          <el-table-column label="加入时间">
            <template slot-scope="scope">
              {{ scope.row.created_at | dateformat }}
            </template>
          </el-table-column>
          <el-table-column label="操作">
            <template slot-scope="scope">
              <el-popconfirm
                v-if="scope.row.role === 'owner'"
                title="解散团队会删除团队的所有相关信息，该操作不可恢复，你确定要继续吗?"
                v-on="{ onConfirm: () => deleteTeam(scope.row) }"
              >
                <el-button type="text" size="small" slot="reference">
                  解散团队
                </el-button>
              </el-popconfirm>
              <el-button
                v-if="scope.row.role === 'owner'"
                type="text"
                size="small"
                @click="inviteTeamMember(scope.row)"
              >
                转让团队
              </el-button>
              <el-popconfirm
                title="你确定要退出团队吗"
                v-on="{ onConfirm: () => quitTeam(scope.row) }"
              >
                <el-button type="text" size="small" slot="reference">
                  退出团队
                </el-button>
              </el-popconfirm>
              <nuxt-link
                v-if="
                  scope.row.role === 'owner' ||
                    scope.row.role === 'administrator'
                "
                :to="`/team/${scope.row.id}/invite`"
              >
                <el-button type="text" size="small">邀请成员</el-button>
              </nuxt-link>
              <nuxt-link
                v-if="
                  scope.row.role === 'owner' ||
                    scope.row.role === 'administrator'
                "
                :to="'/team/mutation?id=' + scope.row.id"
              >
                <el-button type="text" size="small">编辑</el-button>
              </nuxt-link>
            </template>
          </el-table-column>
        </el-table>

        <div class="pagination">
          <el-pagination
            @current-change="changeTeamPage"
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
  async asyncData(context) {
    const { $axios, store } = context
    const currentWorkspace = store.getters['workspace/current']

    if (currentWorkspace) {
      const [{ data: stat }, { data: members, meta }] = await Promise.all([
        $axios.$get(`/team/_/${currentWorkspace}/stat`),
        $axios.$get(`/team/_/${currentWorkspace}/member`),
        store.dispatch('workspace/getProfile', context)
      ])

      return {
        stat,
        members,
        teams: [],
        meta: meta
      }
    } else {
      const { meta, data: teams } = await $axios.$get('/team')

      return {
        stat: {},
        members: [],
        teams,
        meta
      }
    }
  },
  computed: {
    ...mapGetters({
      roles: 'workspace/roles',
      currentWorkspace: 'workspace/current',
      memberProfile: 'workspace/profile'
    })
  },
  watch: {
    currentWorkspace(val) {
      if (val) {
        this.changeMemberPage(0)
        this.$axios.$get(`/team/_/${val}/stat`).then(({ data: stat }) => {
          this.stat = stat
        })
        this.getTeamMemberProfile(this)
      } else {
        this.changeTeamPage(0)
      }
    }
  },
  methods: {
    ...mapActions({
      getTeamMemberProfile: 'workspace/getProfile'
    }),
    async changeTeamPage(page) {
      const { meta, data: teams } = await this.$axios.$get('/team', {
        params: {
          page
        }
      })

      this.teams = teams
      this.meta = meta
    },
    async changeMemberPage(page) {
      const currentWorkspace = this.currentWorkspace
      const { meta, data: members } = await this.$axios.$get(
        `/team/_/${currentWorkspace}/member`,
        {
          params: {
            page
          }
        }
      )

      this.members = members
      this.meta = meta
    },
    // 解散团队
    async deleteTeam(team) {
      try {
        await this.$axios.$delete('/team/_/' + team.id)

        // get list
        this.$success('解散成功')
        this.changeTeamPage(0)
      } catch (err) {
        this.$error(`解散失败: ${err.message}`)
      }
    },
    // 退出团队
    async quitTeam(team) {
      try {
        await this.$axios.$delete('/team/_/' + team.id + '/quit')

        // get list
        this.$success('退出成功')
        this.changeTeamPage(0)
      } catch (err) {
        this.$error(`退出失败: ${err.message}`)
      }
    },
    // 转让团队
    async transferTeam(team) {
      try {
        await this.$axios.$delete('/team/_/' + team.id + '/quit')

        // get list
        this.$success('退出成功')
        this.changeTeamPage(0)
      } catch (err) {
        this.$error(`退出失败: ${err.message}`)
      }
    },
    async kickMemberFromTeam(member) {
      const currentWorkspace = this.currentWorkspace
      try {
        await this.$axios.$delete(
          `/team/_/${currentWorkspace}/member/_/${member.id}`
        )
        this.$success('踢出成功')
        this.changeMemberPage(0)
      } catch (err) {
        this.$error(`踢出失败: ${err.message}`)
      }
    },
    // 邀请成员
    async inviteTeamMember(team) {}
  }
}
</script>
