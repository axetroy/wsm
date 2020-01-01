import cookie from 'js-cookie'

export const state = () => ({
  workspaces: [], // 可用的工作区列表，没有翻页，所有的都在这里
  current: undefined, // 当前工作区 ID
  profile: undefined, // 我在该工作区中的成员身份
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
  ]
})

export const getters = {
  workspaces(state) {
    return state.workspaces
  },
  current(state) {
    return state.current
  },
  profile(state) {
    return state.profile
  },
  roles(state) {
    return state.roles
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
  async updateWorkspaces(store, { $axios }) {
    const { data: workspaces } = await $axios.$get('/team')
    store.commit('UPDATE_WORKSPACES', [defaultWorkspace].concat(workspaces))
  },
  async getProfile(store, { $axios }) {
    const { profile } = await $axios.$get(
      `/team/_/${store.getters.current}/profile`
    )
    store.commit('UPDATE_PROFILE', profile)

    return profile
  },
  switchWorkspace(store, workspaceID) {
    store.commit('SWITCH_WORKSPACE', workspaceID)
  }
}
