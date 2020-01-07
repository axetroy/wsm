<template>
  <div style="height: 100%;" ref="container"></div>
</template>

<script>
import { mapGetters } from 'vuex'
import 'xterm/css/xterm.css'

export default {
  data() {
    return {
      connecting: false, // 是否正在连接
      conn: null, // ws 连接的实例
      terminal: null // 终端的实例
    }
  },
  props: {
    host: {
      type: Object
    },
    rows: {
      type: Number,
      default() {
        return 0
      }
    },
    cols: {
      type: Number,
      default() {
        return 0
      }
    }
  },
  computed: {
    ...mapGetters({
      API_HOST: 'API_HOST'
    }),
    _cols() {
      if (this.cols <= 0) {
        return parseInt(window.innerWidth / 9 + '')
      }
      return this.cols || 80
    },
    _rows() {
      // auto rows
      if (this.rows <= 0) {
        return parseInt(window.innerHeight / 17 + '')
      }

      return this.row || 35
    }
  },
  methods: {
    async connect(host) {
      const [
        { Terminal },
        { WebLinksAddon },
        { AttachAddon },
        { FitAddon },
        { SearchAddon }
      ] = await Promise.all([
        import('xterm'),
        import('xterm-addon-web-links'),
        import('xterm-addon-attach'),
        import('xterm-addon-fit'),
        import('xterm-addon-search')
      ])

      const token = this.$axios.defaults.headers.common['Authorization']
      const apiHost = this.API_HOST.replace(/^http/, 'ws')

      const rows = this._rows
      const cols = this._cols

      // Open the websocket connection to the backend
      const socketUrl = `${apiHost}/v1/shell/connect/${host.id}?Authorization=${token}&cols=${cols}&rows=${rows}`

      const socket = new WebSocket(socketUrl)

      this.conn = socket

      const term = new Terminal({
        cols, // 每列 9 个像素宽
        rows // 每行 17个 像素
      })

      term.loadAddon(new WebLinksAddon())
      term.loadAddon(new AttachAddon(socket))
      term.loadAddon(new FitAddon())
      term.loadAddon(new SearchAddon())

      // term.writeln('----------CONNECTING TO SERVER----------')

      const container = this.$refs.container
      term.open(container)

      term.focus()

      // Attach the socket to the terminal
      socket.onopen = ev => {
        // term.writeln('----------CONNECT SUCCESS----------')
      }

      socket.onerror = ev => {
        term.writeln('----------ERROR----------')
      }

      socket.onclose = ev => {
        term.writeln('----------CLOSE----------')
      }
    },
    dispose() {
      this.conn && this.conn.close()
      this.terminal && this.terminal.dispose()
      this.connecting = false
    }
  },
  mounted() {
    this.connect(this.host)
  },
  beforeDestroy() {
    this.dispose()
  }
}
</script>
