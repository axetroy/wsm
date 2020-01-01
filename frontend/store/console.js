export const state = () => ({
  hosts: [], // 当前正在连接的服务器
  currentHost: null, // 当前正在连接的服务器
  isShow: false // 是否显示终端
})

export const getters = {
  hosts(state) {
    return state.hosts
  },
  currentHost(state) {
    return state.currentHost
  },
  isShow(state) {
    return state.isShow
  }
}

export const mutations = {
  APPEND_HOST(state, host) {
    const index = state.hosts.findIndex(v => v.id === host.id)
    if (index >= 0) {
      return
    }
    state.hosts = state.hosts.concat(host)
  },
  REMOVE_HOST(state, hostID) {
    const hosts = state.hosts
    const index = hosts.findIndex(v => v.id === hostID)

    hosts.splice(index, 1)

    if (!hosts.length) {
      state.currentHost = null
      state.isShow = false
      return
    }

    // reset currentHost
    const prevHost = state.hosts[index - 1]
    const nextHost = state.hosts[index + 1]
    if (nextHost) {
      state.currentHost = nextHost.id
    } else if (prevHost) {
      state.currentHost = prevHost.id
    } else {
      state.currentHost = null
    }

    state.hosts = hosts
  },
  SET_CURRENT_HOST(state, host) {
    state.currentHost = host.id
  },
  SET_CONSOLE_SHOW(state, isShow) {
    state.isShow = !!isShow
  }
}

export const actions = {
  toggle({ commit, state }) {
    commit('SET_CONSOLE_SHOW', !state.isShow)
  },
  show({ commit }) {
    commit('SET_CONSOLE_SHOW', true)
  },
  hide({ commit }) {
    commit('SET_CONSOLE_SHOW', false)
  },
  appendHost({ commit }, host) {
    commit('APPEND_HOST', host)
  },
  removeHost({ commit }, hostID) {
    commit('REMOVE_HOST', hostID)
  },
  activeHost({ commit }, host) {
    commit('SET_CURRENT_HOST', host)
  }
}
