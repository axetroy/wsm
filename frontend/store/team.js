// Copyright 2019-2020 Axetroy. All rights reserved. Apache License 2.0.
export const state = () => ({})

export const getters = {}

export const mutations = {}

export const actions = {
  async fetchTeamById(store, teamId) {
    const { data } = await this.$axios.$get(`/team/_/${teamId}`)

    return data
  },
  // 获取团队服务器的ID
  async fetchTeamHostById(store, { teamId, hostId }) {
    const { data } = await this.$axios.$get(
      `/team/_/${teamId}/host/_/${hostId}`
    )

    return data
  }
}
