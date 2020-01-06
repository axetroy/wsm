export const state = () => ({})

export const getters = {}

export const mutations = {}

export const actions = {
  // 获取我的服务器ID
  async fetchHostById(store, hostId) {
    const { data } = await this.$axios.$get(`/host/_/${hostId}`)

    return data
  }
}
