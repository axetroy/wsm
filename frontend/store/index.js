import { parse as cookieParse } from 'cookie'

export const state = () => ({})

export const getters = {}

export const mutations = {}

export const actions = {
  // 初始化
  async nuxtServerInit(store, { req }) {
    const TOKEN_KEY = 'Authorization'
    const WORKSPACE_KEY = 'workspace'

    const cookieMap = cookieParse(req?.headers?.cookie || '')

    const token = cookieMap[TOKEN_KEY]
    const workspace = cookieMap[WORKSPACE_KEY]

    if (token) {
      store.dispatch('user/setToken', token)
      await store.dispatch('user/getProfile')
    }

    if (workspace) {
      store.dispatch('workspace/switchWorkspace', workspace)
      await Promise.all([
        store.dispatch('workspace/getCurrentTeamMemberProfile'),
        store.dispatch('workspace/getWorkspaces')
      ])
    }
  }
}
