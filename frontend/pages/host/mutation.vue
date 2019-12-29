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
          label="密码"
          :required="type === 'create'"
          prop="password"
        >
          <el-input
            type="password"
            autocomplete="new-password"
            v-model="form.password"
          />
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

    const formRules = {}

    let form = {
      name: '',
      host: '',
      port: 22,
      username: '',
      password: '',
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
      console.log('当前工作区', currentWorkspace)
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
