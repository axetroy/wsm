<!-- Copyright 2019-2020 Axetroy. All rights reserved. Apache License 2.0. -->
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
                      <i class="el-icon-user" />
                      团长用户名:
                    </label>
                    <div>
                      {{ team.owner.username }} ({{ team.owner.nickname }})
                    </div>
                  </el-col>
                  <el-col class="mb20" :span="12">
                    <label>
                      <i class="el-icon-user" />
                      团长用户名:
                    </label>
                    <div>
                      {{ team.owner.username }} ({{ team.owner.nickname }})
                    </div>
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

                <div style="margin-top: 20px;">
                  <el-popconfirm
                    v-if="memberProfile && memberProfile.role === 'owner'"
                    title="解散团队会删除团队的所有相关信息，该操作不可恢复，你确定要继续吗?"
                    v-on="{ onConfirm: () => deleteTeam(team) }"
                  >
                    <el-button type="danger" size="small" slot="reference">
                      解散团队
                    </el-button>
                  </el-popconfirm>
                  <el-button
                    type="primary"
                    size="small"
                    @click="showInviteDialog = true"
                  >
                    邀请成员
                  </el-button>

                  <invite-dialog
                    :visible.sync="showInviteDialog"
                    :teamid="team ? team.id : ''"
                    @invite-success="changeInvitePage(0)"
                  />
                </div>
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
                        v-if="memberProfile && memberProfile.role === v.value"
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
                    <div>
                      {{
                        (memberProfile ? memberProfile.created_at : 0)
                          | dateformat
                      }}
                    </div>
                  </el-col>
                </el-row>

                <div style="margin-top: 20px;">
                  <el-popconfirm
                    title="你确定要退出团队吗"
                    v-on="{ onConfirm: () => quitTeam(team) }"
                  >
                    <el-button type="danger" size="small" slot="reference">
                      退出团队
                    </el-button>
                  </el-popconfirm>
                </div>
              </el-card>
            </el-col>
          </el-row>
        </el-col>

        <el-col :span="16">
          <el-tabs v-model="teamTab">
            <el-tab-pane label="成员列表" name="member">
              <el-pagination
                class="text-center"
                style="margin-bottom: 20px;"
                @current-change="changeMemberPage"
                background
                layout="prev, pager, next"
                :page-size="membersMeta.limit"
                :total="membersMeta.total"
                hide-on-single-page
              />

              <el-scrollbar style="height: 580px;" :native="false">
                <el-card
                  v-for="v of members"
                  :key="v.id"
                  style="margin-bottom: 20px;"
                >
                  <div slot="header">
                    <i class="el-icon-user" /> {{ v.username }} ({{
                      v.nickname
                    }})
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
                      <el-button type="danger" size="small" slot="reference">
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
                          getUpdatedRoles(v.role).length
                      "
                      @command="changeRole"
                    >
                      变更身份
                      <el-dropdown-menu slot="dropdown">
                        <el-dropdown-item
                          v-for="r in getUpdatedRoles(v.role)"
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
                      <el-button type="danger" size="small" slot="reference">
                        转让给他
                      </el-button>
                    </el-popconfirm>
                  </div>
                </el-card>
              </el-scrollbar>
            </el-tab-pane>
            <el-tab-pane label="团队邀请" name="invite">
              <el-pagination
                class="text-center"
                style="margin-bottom: 20px;"
                @current-change="changeInvitePage"
                background
                layout="prev, pager, next"
                :page-size="invitesMeta.limit"
                :total="invitesMeta.total"
                hide-on-single-page
              />

              <el-scrollbar style="height: 580px;" :native="false">
                <el-card
                  v-for="v of invites"
                  :key="v.id"
                  style="margin-bottom: 20px;"
                >
                  <div slot="header">
                    <i class="el-icon-user" /> {{ v.user.username }} ({{
                      v.user.nickname
                    }})
                  </div>

                  <div class="meta-info">
                    <div>
                      {{ v.invitor.username }} ({{ v.invitor.nickname }}) 于
                      {{ v.created_at | dateformat }} 邀请加入团队, 成为
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
                          (memberProfile.role === 'owner' ||
                            memberProfile.role === 'administrator') &&
                          v.state === 'init'
                      "
                      title="你确定要撤销这个邀请吗? 该操作不可恢复"
                      v-on="{ onConfirm: () => cancelInvite(v) }"
                    >
                      <el-button type="primary" size="small" slot="reference">
                        撤销邀请
                      </el-button>
                    </el-popconfirm>
                    <!--                    <template v-if="v.state === 'init'">-->
                    <!--                      <el-button-->
                    <!--                        type="primary"-->
                    <!--                        size="small"-->
                    <!--                        slot="reference"-->
                    <!--                        @click="handleInvite(v, true)"-->
                    <!--                      >-->
                    <!--                        接受-->
                    <!--                      </el-button>-->
                    <!--                      <el-button-->
                    <!--                        size="small"-->
                    <!--                        slot="reference"-->
                    <!--                        @click="handleInvite(v, false)"-->
                    <!--                      >-->
                    <!--                        拒绝-->
                    <!--                      </el-button>-->
                    <!--                    </template>-->
                    <template v-if="v.state !== 'init'">
                      <span
                        v-for="s of inviteState"
                        :key="s.value"
                        v-if="s.value === v.state"
                      >
                        {{ s.label }}
                      </span>
                    </template>
                  </div>
                </el-card>
              </el-scrollbar>
            </el-tab-pane>
          </el-tabs>
        </el-col>
      </el-row>
    </template>

    <template v-else>
      <el-tabs v-model="myTab">
        <el-tab-pane label="我所属的团队" name="team">
          <el-button
            class="mb20"
            size="small"
            type="primary"
            @click="showTeamMutationDialog = true"
            >创建团队</el-button
          >

          <team-mutation-dialog :visible.sync="showTeamMutationDialog" />

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
                <el-button
                  type="primary"
                  size="small"
                  @click="switchWorkspace(scope.row.id)"
                >
                  进入工作区
                </el-button>
              </template>
            </el-table-column>
          </el-table>

          <el-pagination
            class="pagination"
            @current-change="changeTeamPage"
            background
            layout="prev, pager, next"
            :page-size="teamsMeta.limit"
            :total="teamsMeta.total"
            hide-on-single-page
          />
        </el-tab-pane>
        <el-tab-pane label="邀请我的" name="invite">
          <el-pagination
            class="text-center"
            style="margin-bottom: 20px;"
            @current-change="changeMyInvitePage"
            background
            layout="prev, pager, next"
            :page-size="myInvitesMeta.limit"
            :total="myInvitesMeta.total"
            hide-on-single-page
          />

          <el-scrollbar style="height: 580px;" :native="false">
            <el-card
              v-for="v of myInvites"
              :key="v.id"
              style="margin-bottom: 20px;"
            >
              <div slot="header">
                <i class="el-icon-user" /> {{ v.team.name }}
              </div>

              <div class="meta-info">
                <div>
                  {{ v.invitor.username }} ({{ v.invitor.nickname }}) 于
                  {{ v.created_at | dateformat }} 邀请加入团队 【{{
                    v.team.name
                  }}】, 成为
                  <span
                    v-for="item in roles"
                    :key="item.value"
                    v-if="item.value === v.role"
                    >{{ item.label }}
                  </span>
                </div>
              </div>
              <div class="action-block">
                <template v-if="v.state === 'init'">
                  <el-button
                    type="primary"
                    size="small"
                    slot="reference"
                    @click="handleInvite(v, true)"
                  >
                    接受
                  </el-button>
                  <el-button
                    size="small"
                    slot="reference"
                    @click="handleInvite(v, false)"
                  >
                    拒绝
                  </el-button>
                </template>
                <template v-if="v.state !== 'init'">
                  <span
                    v-for="s of inviteState"
                    :key="s.value"
                    v-if="s.value === v.state"
                  >
                    {{ s.label }}
                  </span>
                </template>
              </div>
            </el-card>
          </el-scrollbar>
        </el-tab-pane>
      </el-tabs>
    </template>
  </div>
</template>

<script>
import { mapGetters, mapActions } from 'vuex'
import invite_dialog from '../../components/invite_dialog'
import team_mutation_dialog from '../../components/team_mutation_dialog'

export default {
  components: {
    'invite-dialog': invite_dialog,
    'team-mutation-dialog': team_mutation_dialog
  },
  async asyncData({ store }) {
    const currentWorkspace = store.getters['workspace/current']

    if (currentWorkspace) {
      const [
        _,
        stat,
        { data: members, meta },
        { data: invites, meta: invitesMeta }
      ] = await Promise.all([
        store.dispatch('workspace/getCurrentTeamMemberProfile'),
        store.dispatch('workspace/getCurrentTeamStat'),
        store.dispatch('workspace/getCurrentTeamMembers'),
        store.dispatch('workspace/getCurrentTeamInvites')
      ])

      return {
        stat,
        members,
        membersMeta: meta,
        invites,
        invitesMeta
      }
    } else {
      const [
        { meta: teamsMeta, data: teams },
        { meta: myInvitesMeta, data: myInvites }
      ] = await Promise.all([
        store.dispatch('workspace/getTeams'),
        store.dispatch('user/getInvites')
      ])

      return {
        teams,
        teamsMeta,
        myInvites,
        myInvitesMeta
      }
    }
  },
  data() {
    return {
      stat: {}, // 团队统计
      teamTab: 'member',
      myTab: 'team',
      teams: [], // 我所在的团队
      teamsMeta: {},
      members: [], // 团队成员
      membersMeta: {},
      invites: [], // 团队邀请
      invitesMeta: {},
      myInvites: [], // 邀请我的列表
      myInvitesMeta: {},
      showInviteDialog: false, // 是否显示团队邀请弹窗
      showTeamMutationDialog: false // 是否显示团队邀请弹窗
    }
  },
  computed: {
    ...mapGetters({
      roles: 'workspace/roles',
      currentWorkspace: 'workspace/current',
      memberProfile: 'workspace/profile',
      team: 'workspace/team',
      inviteState: 'workspace/inviteState'
    })
  },
  watch: {
    currentWorkspace(val) {
      if (val) {
        this.changeMemberPage(0)
        this.changeInvitePage(0)
        this.getCurrentTeamStat().then(stat => {
          this.stat = stat
        })
        this.getCurrentTeamMemberProfile()
      } else {
        this.changeTeamPage(0)
        this.changeMyInvitePage(0)
      }
    }
  },
  methods: {
    ...mapActions({
      getCurrentTeamMemberProfile: 'workspace/getCurrentTeamMemberProfile',
      getTeams: 'workspace/getTeams',
      getCurrentTeamInvites: 'workspace/getCurrentTeamInvites',
      getCurrentTeamMembers: 'workspace/getCurrentTeamMembers',
      getCurrentTeamStat: 'workspace/getCurrentTeamStat',
      getMyInvites: 'user/getInvites',
      switchWorkspace: 'workspace/switchWorkspace'
    }),
    async changeTeamPage(page) {
      const { meta, data: teams } = await this.getTeams({ page })
      this.teams = teams
      this.teamsMeta = meta
    },
    async changeMemberPage(page) {
      const { meta, data: members } = await this.getCurrentTeamMembers({ page })

      this.members = members
      this.membersMeta = meta
    },
    async changeInvitePage(page) {
      const { meta, data: invites } = await this.getCurrentTeamInvites({ page })

      this.invites = invites
      this.invitesMeta = meta
    },
    async changeMyInvitePage(page) {
      const { meta, data: invites } = await this.getMyInvites({ page })

      this.myInvites = invites
      this.myInvitesMeta = meta
    },
    // 获取可修改的角色列表
    getUpdatedRoles(role) {
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
        this.switchWorkspace(undefined)
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
        this.switchWorkspace(undefined)
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
    },
    // 撤销邀请
    async cancelInvite(invite) {
      const currentWorkspace = this.currentWorkspace
      try {
        await this.$axios.$delete(
          `/team/_/${currentWorkspace}/member/invite/_/${invite.id}`
        )
        this.$success('撤销成功')
        this.changeInvitePage(0)
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
        this.changeMyInvitePage(0)

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

<style lang="less">
#team-page {
  .mb20 {
    margin-bottom: 20px;
    border-bottom: 1px solid #ebeef5;
    // padding-bottom: 20px;
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
