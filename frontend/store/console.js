export const state = () => ({
  hosts: [], // 当前正在连接的服务器
  currentHost: null, // 当前正在连接的服务器
  consoleShow: false // 是否显示终端
})

export const getters = {
  /**
   * @param {{ hosts: any; }} state
   */
  hosts(state) {
    return state.hosts
  },
  /**
   * @param {{ currentHost: any; }} state
   */
  currentHost(state) {
    return state.currentHost
  },
  /**
   * @param {{ consoleShow: any; }} state
   */
  consoleShow(state) {
    return state.consoleShow
  }
}

export const mutations = {
  /**
   * @param {{ hosts: any[]; }} state
   * @param {{ id: string; }} host
   */
  APPEND_HOST(state, host) {
    const index = state.hosts.findIndex(v => v.id === host.id)
    if (index >= 0) {
      return
    }
    state.hosts = state.hosts.concat(host)
  },
  /**
   * @param {{ hosts: any[]; currentHost: string|null; consoleShow:boolean }} state
   * @param {string} hostID
   */
  REMOVE_HOST(state, hostID) {
    const hosts = state.hosts
    const index = hosts.findIndex(v => v.id === hostID)

    hosts.splice(index, 1)

    // reset currentHost
    if (!hosts.length) {
      state.currentHost = null
      state.consoleShow = false
    } else {
      const prevHost = state.hosts[index - 1]
      const nextHost = state.hosts[index + 1]
      if (nextHost) {
        state.currentHost = nextHost.id
      } else {
        state.currentHost = prevHost.id
      }
    }

    state.hosts = hosts
  },
  /**
   * @param {{ currentHost: any; }} state
   * @param {string} currentHostID
   */
  SET_CURRENT_HOST(state, currentHostID) {
    state.currentHost = currentHostID
  },
  /**
   * @param {{ consoleShow: boolean; }} state
   * @param {Boolean} isShow
   */
  SET_CONSOLE_SHOW(state, isShow) {
    state.consoleShow = !!isShow
  },
  /**
   * @param {{ consoleShow: boolean; }} state
   */
  TOGGLE_CONSOLE(state) {
    state.consoleShow = !state.consoleShow
  }
}

export const actions = {}
