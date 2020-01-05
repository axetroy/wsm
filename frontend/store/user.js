import cookie from 'js-cookie'

export const state = () => ({
  profile: undefined, // 我的用户资料
  invites: [] // 邀请我加入团队的列表
})

export const getters = {
  profile(state) {
    return state.profile
  },
  invites(state) {
    return state.invites
  }
}

export const mutations = {
  SET_USER(state, user) {
    state.user = user || null
  }
}

export const actions = {
  // 获取邀请我的列表
  async getInvites(store, params) {
    const { data, meta } = await this.$axios.$get('/team/invite', { params })
    return { data, meta }
  },
  // 更新用户的资料
  async getProfile(store) {
    const { data: profile } = await this.$axios.$get('/user/profile')
    store.commit('SET_USER', profile)
    return profile
  },
  // 登录
  async login(store, body) {
    const { data: profile } = await this.$axios.$post('/auth/signin', body)

    const { token } = profile
    cookie.set('Authorization', token)
    store.commit('SET_USER', profile)
    return profile
  }
}
