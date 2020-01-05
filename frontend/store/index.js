import cookie from 'js-cookie'
import { parse as CookieParse } from 'cookie'

export const state = () => ({})

export const getters = {}

export const mutations = {}

export const actions = {
  // 初始化
  async nuxtServerInit(store, context) {
    const { req } = context
    const TOKEN_KEY = 'Authorization'
    const WORKSPACE_KEY = 'workspace'

    const cookieMap = CookieParse(req.headers['cookie'] || '')

    const token = cookieMap[TOKEN_KEY]
    const workspace = cookieMap[WORKSPACE_KEY]

    if (token) {
      await store.dispatch('user/getProfile', context)
    }

    if (workspace) {
      store.dispatch('workspace/switchWorkspace', workspace)
      await store.dispatch('workspace/getCurrentTeamMemberProfile')
      await store.dispatch('workspace/getWorkspaces')
    }
  }
}
