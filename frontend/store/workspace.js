import cookie from 'js-cookie'

export const state = () => ({
  workspaces: [], // 可用的工作区列表
  current: undefined // 当前工作区 ID
})

export const getters = {
  workspaces(state) {
    return state.workspaces
  },
  current(state) {
    return state.current
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
  }
}

export const actions = {}
