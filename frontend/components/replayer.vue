<template>
  <div style="height: 100%;" ref="container"></div>
</template>

<script>
import { mapGetters } from 'vuex'
import 'xterm/css/xterm.css'

function sleep(ms) {
  return new Promise(resolve => {
    setTimeout(() => {
      resolve()
    }, ms)
  })
}

function parseRow(rowStr) {
  const date = rowStr.substr(1, 23) // 1~23 为日期时间戳字符串，带毫秒
  const content = new TextEncoder().encode(atob(rowStr.substr(26))) // 前面 25个字符为时间戳，第26个字符为空格，其余后面为正文内容

  return {
    date: new Date(date),
    content
  }
}

export default {
  data() {
    return {
      terminal: null // 终端的实例
    }
  },
  props: {
    // 重防的记录 ID
    records: {
      type: Array
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
    async connect() {
      const [
        { Terminal },
        { WebLinksAddon },
        { FitAddon },
        { SearchAddon }
      ] = await Promise.all([
        import('xterm'),
        import('xterm-addon-web-links'),
        import('xterm-addon-fit'),
        import('xterm-addon-search')
      ])

      const rows = this._rows
      const cols = this._cols

      const term = new Terminal({
        cols, // 每列 9 个像素宽
        rows // 每行 17个 像素
      })

      term.loadAddon(new WebLinksAddon())
      term.loadAddon(new FitAddon())
      term.loadAddon(new SearchAddon())

      const container = this.$refs.container
      term.open(container)

      term.focus()

      const records = this.records.map((v, index) => {
        let duration = 0

        const { date, content } = parseRow(v)

        const nextRow = this.records[index + 1]

        if (nextRow) {
          const { date: nextDate } = parseRow(nextRow)

          duration = nextDate.getTime() - date.getTime()
        }

        return {
          date, // 日期时间
          content, // 输出内容
          duration // 开始时间, 以毫秒为单位
        }
      })

      for (let i = 0; i < records.length; i++) {
        const record = records[i]
        await new Promise(resolve => {
          term.write(record.content, () => {
            resolve()
          })
        })
        await sleep(record.duration)
      }
    },
    dispose() {
      this.conn && this.conn.close()
      this.terminal && this.terminal.dispose()
    }
  },
  mounted() {
    this.connect()
  },
  beforeDestroy() {
    this.dispose()
  }
}
</script>
