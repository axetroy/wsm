// Copyright 2019-2020 Axetroy. All rights reserved. Apache License 2.0.
import cookie from 'js-cookie'

export const state = () => ({
  token: '', // 用户的 token
  profile: undefined, // 我的用户资料
  invites: [] // 邀请我加入团队的列表
})

export const getters = {
  token(state) {
    return state.token
  },
  profile(state) {
    return state.profile
  },
  invites(state) {
    return state.invites
  }
}

export const mutations = {
  SET_TOKEN(state, token) {
    state.token = token
  },
  SET_USER(state, profile) {
    state.profile = profile || undefined
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
  async setToken(store, token) {
    store.commit('SET_TOKEN', token)
  },
  // 登录
  async login(store, body) {
    const { data: profile } = await this.$axios.$post('/auth/signin', body)

    const { token } = profile
    cookie.set('Authorization', token)
    store.commit('SET_TOKEN', token)
    store.commit('SET_USER', profile)
    return profile
  },
  // 登出
  logout(store) {
    if (process.client) {
      cookie.remove('Authorization')
      cookie.remove('workspace')
    } else {
      const res = this?.app?.context?.res
      res.setHeader('Set-Cookie', [
        `Authorization=; path=/; expires=Thu, 01 Jan 1970 00:00:00 GMT`,
        `workspace=; path=/; expires=Thu, 01 Jan 1970 00:00:00 GMT`
      ])
    }
    store.commit('SET_TOKEN', '') // 删除 token
    store.commit('SET_USER', undefined) // 清空用户资料
    store.dispatch('workspace/switchWorkspace', undefined, { root: true }) // 清空工作区
  }
}
