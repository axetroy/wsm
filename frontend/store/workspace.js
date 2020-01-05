import cookie from 'js-cookie'

export const state = () => ({
  workspaces: [], // 可用的工作区列表，没有翻页，所有的都在这里
  current: undefined, // 当前工作区 ID
  profile: undefined, // 我在该工作区中的成员身份
  // 成员角色
  roles: [
    {
      label: '全部',
      value: undefined
    },
    {
      label: '拥有者',
      value: 'owner'
    },
    {
      label: '管理员',
      value: 'administrator'
    },
    {
      label: '成员',
      value: 'member'
    },
    {
      label: '访客',
      value: 'visitor'
    }
  ],
  // 邀请记录的状态
  inviteState: [
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
})

export const getters = {
  workspaces(state) {
    return state.workspaces
  },
  current(state) {
    return state.current
  },
  team(state) {
    return state.workspaces.find(v => v.id === state.current)
  },
  profile(state) {
    return state.profile
  },
  roles(state) {
    return state.roles
  },
  inviteState(state) {
    return state.inviteState
  }
}

export const mutations = {
  UPDATE_WORKSPACES(state, workspaces) {
    state.workspaces = workspaces
  },
  SWITCH_WORKSPACE(state, workspaceID) {
    state.current = workspaceID
    // 设置 cookie，有效期为 24 小时
    const expAt = new Date(new Date().getTime() + 1000 * 3600 * 24)
    if (workspaceID) {
      cookie.set('workspace', workspaceID, expAt)
    } else {
      cookie.remove('workspace')
    }
  },
  UPDATE_PROFILE(state, profile) {
    state.profile = profile
  }
}

const defaultWorkspace = {
  id: undefined,
  name: '当前帐号'
}

export const actions = {
  // 获取所有工作区
  async getWorkspaces(store) {
    const { data: workspaces } = await this.$axios.$get('/team/all')
    store.commit('UPDATE_WORKSPACES', [defaultWorkspace].concat(workspaces))
  },
  // 获取团队列表
  async getTeams(store, params) {
    const { data, meta } = await this.$axios.$get('/team', { params })
    return { data, meta }
  },
  // 获取当前工作区下的我成员资料
  async getCurrentTeamMemberProfile(store) {
    const { data: profile } = await this.$axios.$get(
      `/team/_/${store.getters.current}/profile`
    )
    store.commit('UPDATE_PROFILE', profile)

    return profile
  },
  // 切换工作区
  switchWorkspace(store, workspaceID) {
    store.commit('SWITCH_WORKSPACE', workspaceID)
  },
  // 获取当前工作区的邀请记录
  async getCurrentTeamInvites({ state }, params = {}) {
    const { data, meta } = await this.$axios.$get(
      `/team/_/${state.current}/member/invite`,
      {
        params
      }
    )

    return { data, meta }
  },
  // 获取当前工作区的统计信息
  async getCurrentTeamStat({ state }) {
    const { data } = await this.$axios.$get(`/team/_/${state.current}/stat`)

    return data
  },
  // 获取当前工作区的成员列表
  async getCurrentTeamMembers({ state }, params = {}) {
    const {
      data,
      meta
    } = await this.$axios.$get(`/team/_/${state.current}/member`, { params })

    return { data, meta }
  }
}
