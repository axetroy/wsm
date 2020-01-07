import { parse as cookieParse } from 'cookie'

export const state = () => ({
  API_HOST: 'http://localhost:80'
})

export const getters = {
  API_HOST(state) {
    return state.API_HOST
  }
}

export const mutations = {
  SET_API_HOST(state, host) {
    state.API_HOST = host
  }
}

export const actions = {
  // 初始化
  async nuxtServerInit(store, { req }) {
    const TOKEN_KEY = 'Authorization'
    const WORKSPACE_KEY = 'workspace'

    const cookieMap = cookieParse(req?.headers?.cookie || '')

    const token = cookieMap[TOKEN_KEY]
    const workspace = cookieMap[WORKSPACE_KEY]

    const apiHost = process.env.API_HOST

    if (apiHost) {
      store.commit('SET_API_HOST', apiHost)
    }

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
