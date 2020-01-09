<template>
  <div style="position: relative">
    <div class="ctrl-bar">
      <div
        class="progress"
        :style="`width: ${progress * 100}%;`"
        v-if="records.length"
      >
        &nbsp;
        <!--        <span class="text">{{ playingTime }} / {{ endTime }}</span>-->
      </div>
    </div>
    <div style="height: 100%;" ref="container"></div>
  </div>
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
  if (rowStr === undefined) return undefined
  const date = rowStr.substr(1, 23) // 1~23 为日期时间戳字符串，带毫秒
  const content = new TextEncoder().encode(atob(rowStr.substr(26))) // 前面 25个字符为时间戳，第26个字符为空格，其余后面为正文内容

  return {
    date: new Date(date),
    content
  }
}

const STATE = {
  init: 0,
  ready: 0.9,
  playing: 1,
  pause: 2,
  stop: 3,
  finish: 4
}

export default {
  data() {
    return {
      terminal: null, // 终端的实例
      // 当前状态
      // init: 组件刚加载完
      // ready: 初始化完毕
      // playing: 正在播放
      // pause: 暂停中
      // stop: 已停止
      // finish: 已播放完毕
      state: STATE.init,
      step: 0, // 当前执行的步骤
      totalStep: this.records?.length || 0 // 总步数
    }
  },
  props: {
    // 重方的记录列表
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
    },
    // 是否自动播放
    autoplay: {
      type: Boolean,
      default() {
        return true
      }
    }
  },
  watch: {
    records(records) {
      this.totalStep = records?.length || 0
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
    },
    // 当前播放进度, 1 为 100%
    progress() {
      if (!this.records.length || !this.step) {
        return 0
      }

      const percent = this.step / this.totalStep

      if (percent > 1) {
        console.warn(`进度条大于 100% 当前步骤 ${this.step}/${this.totalStep}`)
      }

      return Math.min(percent, 1)
    }
  },
  methods: {
    // 初始化
    async init() {
      if (this.state !== STATE.init) {
        return
      }

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

      const term = (this.terminal = new Terminal({
        cols, // 每列 9 个像素宽
        rows // 每行 17个 像素
      }))

      term.loadAddon(new WebLinksAddon())
      term.loadAddon(new FitAddon())
      term.loadAddon(new SearchAddon())

      const container = this.$refs.container
      term.open(container)

      this.state = STATE.ready

      term.focus()
    },
    // 开始播放, 从第 n 步开始, 最小为 0
    async play(startStep = 0) {
      await this.init()

      const term = this.terminal

      if (!term) {
        return
      }

      this.state = STATE.ready

      this.state = STATE.playing

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

      for (let i = startStep; i < records.length; i++) {
        if (this.state !== STATE.playing) {
          return
        }
        this.step = i + 1
        const record = records[i]
        await new Promise(resolve => {
          term.write(record.content, () => {
            resolve()
          })
        })
        await sleep(record.duration)
      }

      this.state = STATE.finish
    },
    // 停止播放
    async stop() {
      const term = this.terminal

      this.state = STATE.stop
      this.step = 0

      if (!term) {
        return
      }

      term.reset()

      this.state = STATE.stop
    },
    // 暂停播放
    async pause() {
      this.pause = STATE.pause
    },
    dispose() {
      this.conn && this.conn.close()
      this.terminal && this.terminal.dispose()
    }
  },
  async mounted() {
    if (this.autoplay) {
      await this.play()
    }
  },
  destroyed() {
    this.stop()
    this.dispose()
  }
}
</script>

<style lang="less" scoped>
.ctrl-bar {
  position: absolute;
  bottom: -5px;
  left: 0;
  width: 100%;
  height: 5px;
  z-index: 99999;
  color: #fff;
}

.progress {
  height: 100%;
  background-color: #306a27;
  display: inline-block;

  .text {
    position: absolute;
    bottom: 0;
    left: 0;
  }
}
</style>
