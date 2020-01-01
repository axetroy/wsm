const CookieParse = require('cookie').parse

export const state = () => ({
  user: null
})

export const getters = {
  user(state) {
    return state.user
  }
}

export const mutations = {
  SET_USER(state, user) {
    state.user = user || null
  }
}

export const actions = {
  // 初始化
  async nuxtServerInit(store, context) {
    const { req } = context
    const TOKEN_KEY = 'Authorization'
    const WORKSPACE_KEY = 'workspace'

    const cookieMap = CookieParse(req.headers['cookie'] || '')

    const token = cookieMap[TOKEN_KEY]
    const workspace = cookieMap[WORKSPACE_KEY]

    if (workspace) {
      store.commit('workspace/SWITCH_WORKSPACE', workspace)
    }

    if (token) {
      await store.dispatch('updateProfile', context)
    }
  },
  // 更新用户的资料
  async updateProfile(store, { $axios }) {
    const { data: profile } = await $axios.$get('/user/profile')
    store.commit('SET_USER', profile)
  }
}
