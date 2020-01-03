<template>
  <div id="team-page" class="main">
    <template v-if="currentWorkspace">
      <el-row :gutter="20">
        <el-col :span="8">
          <el-row>
            <el-col :span="24">
              <el-card>
                <div slot="header">
                  <div>
                    <i class="el-icon-s-home" />
                    团队信息
                  </div>
                </div>
                <el-row>
                  <el-col class="mb20" :span="12">
                    <label>
                      <i class="el-icon-view" />
                      团队名称:
                    </label>
                    <div>{{ stat.name }}</div>
                  </el-col>
                  <el-col class="mb20" :span="12">
                    <label>
                      <i class="el-icon-s-custom" />
                      成员数量:
                    </label>
                    <div>{{ stat.member_num }}</div>
                  </el-col>
                  <el-col :span="12">
                    <label>
                      <i class="el-icon-copy-document" />
                      团队服务器:
                    </label>
                    <div>{{ stat.host_num }}</div>
                  </el-col>
                  <el-col :span="12">
                    <label>
                      <i class="el-icon-time" />
                      创建日期:
                    </label>
                    <div>{{ stat.created_at | dateformat }}</div>
                  </el-col>
                </el-row>
              </el-card>
            </el-col>

            <el-col :span="24" style="margin-top: 30px">
              <el-card>
                <div slot="header">
                  <div>
                    <i class="el-icon-user" />
                    我的信息
                  </div>
                </div>
                <el-row>
                  <el-col :span="12">
                    <label>
                      <i class="el-icon-user" />
                      我的身份:</label
                    >
                    <div>
                      <span
                        v-for="v in roles"
                        :key="v.value"
                        v-if="v.value === memberProfile.role"
                      >
                        {{ v.label }}
                      </span>
                    </div>
                  </el-col>
                  <el-col :span="12">
                    <label>
                      <i class="el-icon-time" />
                      加入时间:
                    </label>
                    <div>{{ memberProfile.created_at | dateformat }}</div>
                  </el-col>
                </el-row>
              </el-card>
            </el-col>
          </el-row>
        </el-col>

        <el-col :span="16">
          <el-card
            v-for="v of members"
            :key="v.id"
            style="margin-bottom: 20px;"
          >
            <div slot="header">
              <i class="el-icon-user" /> {{ v.username }} ({{ v.nickname }})
            </div>

            <div class="meta-info">
              <div>
                于 {{ v.created_at | dateformat }} 加入团队, 成为
                <span
                  v-for="item in roles"
                  :key="item.value"
                  v-if="item.value === v.role"
                  >{{ item.label }}
                </span>
              </div>
            </div>
            <div class="action-block">
              <el-popconfirm
                v-if="
                  memberProfile &&
                    memberProfile.id !== v.id &&
                    memberProfile.role === 'owner'
                "
                title="你确定要踢出团队吗"
                v-on="{ onConfirm: () => kickMemberFromTeam(v) }"
              >
                <el-button type="primary" size="small" slot="reference">
                  踢出团队
                </el-button>
              </el-popconfirm>
              <el-dropdown
                split-button
                type="primary"
                size="small"
                v-if="
                  memberProfile &&
                    memberProfile.id !== v.id &&
                    getRoles(v.role).length
                "
                @command="changeRole"
              >
                变更身份
                <el-dropdown-menu slot="dropdown">
                  <el-dropdown-item
                    v-for="r in getRoles(v.role)"
                    :key="r.value"
                    :command="{ user: v, role: r.value }"
                  >
                    {{ r.label }}
                  </el-dropdown-item>
                </el-dropdown-menu>
              </el-dropdown>
              <el-popconfirm
                v-if="
                  memberProfile &&
                    memberProfile.id !== v.id &&
                    memberProfile.role === 'owner'
                "
                title="你确定要把团队转让给他吗？这个会把他升级为拥有者，而您变成团队成员"
                v-on="{ onConfirm: () => transferTeam(v) }"
              >
                <el-button type="primary" size="small" slot="reference">
                  转让给他
                </el-button>
              </el-popconfirm>
            </div>
          </el-card>

          <div class="pagination">
            <el-pagination
              @current-change="changeMemberPage"
              background
              layout="prev, pager, next"
              :page-size="meta.limit"
              :total="meta.total"
              hide-on-single-page
            >
            </el-pagination>
          </div>
        </el-col>
      </el-row>
    </template>
    <template v-else>
      <el-card>
        <div slot="header">
          <div>所属团队</div>
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
            hide-on-single-page
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
        store.dispatch('workspace/getCurrentTeamMemberProfile')
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
        this.getCurrentTeamMemberProfile()
      } else {
        this.changeTeamPage(0)
      }
    }
  },
  methods: {
    ...mapActions({
      getCurrentTeamMemberProfile: 'workspace/getCurrentTeamMemberProfile'
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
    getRoles(role) {
      if (!this.memberProfile) {
        return []
      }

      switch (role) {
        case 'owner':
          return []
        case 'administrator':
          if (this.memberProfile.role !== 'owner') {
            return []
          }
          break
        case 'member':
          if (
            ['owner', 'administrator'].includes(this.memberProfile.role) ===
            false
          ) {
            return []
          }
          break
      }

      const roles = this.roles.filter(v => !!v.value)

      switch (this.memberProfile.role) {
        case 'owner':
          return roles.filter(v =>
            ['administrator', 'member', 'visitor'].includes(v.value)
          )
        case 'administrator':
          return roles.filter(v => ['member', 'visitor'].includes(v.value))
        case 'member':
          return roles.filter(v => ['visitor'].includes(v.value))
        default:
          return []
      }
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
    async transferTeam(user) {
      const currentWorkspace = this.currentWorkspace
      try {
        await this.$axios.$put(
          '/team/_/' + currentWorkspace + '/transfer/' + user.id
        )
        this.$success('转让成功')
        this.changeMemberPage(0)
        this.getCurrentTeamMemberProfile()
      } catch (err) {
        this.$error(`转让失败: ${err.message}`)
      }
    },
    // 从团队中踢出成员
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
    async inviteTeamMember(team) {},
    // 变更身份
    async changeRole({ user, role }) {
      const currentWorkspace = this.currentWorkspace
      try {
        await this.$axios.$put(`/team/_/${currentWorkspace}/role/${user.id}`, {
          role
        })
        this.$success('变更成功')
        this.changeMemberPage(0)
      } catch (err) {
        this.$error(`变更失败: ${err.message}`)
      }
    }
  }
}
</script>

<style lang="less">
#team-page {
  .mb20 {
    margin-bottom: 20px;
    border-bottom: 1px solid #ebeef5;
    padding-bottom: 20px;
  }

  .el-card__body {
    label {
      display: inline-block;
      margin-bottom: 5px;
    }
  }

  .meta-info {
    padding-bottom: 20px;
  }

  .action-block {
    padding-top: 20px;
    border-top: 1px solid #ebeef5;
  }
}
</style>
