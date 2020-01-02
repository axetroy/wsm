<template>
  <div class="main">
    <el-card shadow="never">
      <div slot="header">
        <h4>{{ type === 'create' ? '添加' : '修改' }}服务器</h4>
      </div>
      <el-form
        label-width="160px"
        status-icon
        :model="form"
        :ref="formName"
        :rules="formRules"
      >
        <el-form-item label="服务器名称" required prop="name">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="服务器地址" required prop="host">
          <el-input v-model="form.host" />
        </el-form-item>
        <el-form-item label="服务器端口" required prop="port">
          <el-input type="number" v-model="form.port" />
        </el-form-item>
        <el-form-item label="用户名" required prop="username">
          <el-input v-model="form.username" />
        </el-form-item>
        <el-form-item
          label="连接方式"
          :required="type === 'create'"
          prop="passport"
        >
          <el-input
            v-if="form.connect_type === 'password'"
            placeholder="请输入连接密码"
            v-model="form.passport"
            type="password"
          >
            <el-select
              v-model="form.connect_type"
              slot="prepend"
              placeholder="请选择连接方式"
              style="width: 90px"
            >
              <el-option
                v-for="v of connectTypes"
                :key="v.value"
                :label="v.label"
                :value="v.value"
              />
            </el-select>
          </el-input>
          <template v-else>
            <el-select
              v-model="form.connect_type"
              placeholder="请选择连接方式"
              style="width: 90px"
            >
              <el-option
                v-for="v of connectTypes"
                :key="v.value"
                :label="v.label"
                :value="v.value"
              />
            </el-select>

            <el-input
              type="textarea"
              placeholder="请输入连接密钥"
              :autosize="{ minRows: 5, maxRows: 10 }"
              v-model="form.passport"
              style="margin-top: 20px;"
            />
          </template>
        </el-form-item>
        <el-form-item label="描述" prop="remark">
          <el-input
            type="textarea"
            :autosize="{ minRows: 5 }"
            v-model="form.remark"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="onSubmit">提交</el-button>
          <el-button @click="$router.go(-1)">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'

export default {
  async asyncData({ $axios, query, store }) {
    const currentWorkspace = store.getters['workspace/current']
    const type = query.id !== undefined ? 'update' : 'create'

    const formRules = {
      name: [{ required: true, message: '请输入服务器名称' }],
      host: [{ required: true, message: '请输入服务器地址，例如 192.168.0.0' }],
      port: [{ required: true, message: '请输入服务器的 SSH 端口，例如 22' }],
      username: [
        { required: true, message: '请输入连接服务器的用户名，例如 root' }
      ],
      connect_type: [{ required: true, message: '请选择连接服务器的方式' }],
      passport: [{ required: true, message: '请输入密码/私钥' }]
    }

    let form = {
      name: '',
      host: '',
      port: 22,
      username: '',
      connect_type: 'password',
      passport: '',
      remark: ''
    }

    if (type === 'update') {
      const { data } = await $axios.$get(
        currentWorkspace
          ? `/team/_/${currentWorkspace}/host/_/${query.id}`
          : `/host/_/${query.id}`
      )
      form = data
      if (!form) {
        return redirect('/404')
      }
    }

    return {
      type,
      formName: 'form',
      formRules,
      form
    }
  },
  data() {
    return {
      connectTypes: [
        {
          label: '密码',
          value: 'password'
        },
        {
          label: '密钥',
          value: 'private_key'
        }
      ]
    }
  },
  computed: {
    ...mapGetters({
      currentWorkspace: 'workspace/current'
    })
  },
  methods: {
    onSubmit() {
      this.$refs[this.formName].validate(valid => {
        if (valid) {
          switch (this.type) {
            case 'create':
              return this.create()
            case 'update':
              return this.update()
          }
        }
      })
    },
    async update() {
      const form = this.form
      form.port = +form.port
      const currentWorkspace = this.currentWorkspace
      const id = this.$route.query.id
      try {
        await this.$axios.$put(
          currentWorkspace
            ? `/team/_/${currentWorkspace}/host/_/${id}`
            : `/host/_/${id}`,
          form
        )
        this.$success('更新成功.')
        this.$router.back()
      } catch (err) {
        this.$error(err.message)
      }
    },
    async create() {
      const form = this.form
      form.port = +form.port
      const currentWorkspace = this.currentWorkspace
      try {
        await this.$axios.$post(
          currentWorkspace ? `/team/_/${currentWorkspace}/host` : '/host',
          form
        )
        this.$success('创建成功.')
        this.$router.back()
      } catch (err) {
        console.dir(err)
        this.$error(err.message)
      }
    }
  }
}
</script>
