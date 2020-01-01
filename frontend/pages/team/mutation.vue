<template>
  <div class="main">
    <el-card shadow="never">
      <div slot="header">
        <h4>{{ type === 'create' ? '创建' : '修改' }}团队</h4>
      </div>
      <el-form
        label-width="160px"
        status-icon
        :model="form"
        :ref="formName"
        :rules="formRules"
      >
        <el-form-item label="团队名称" required prop="name">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item v-if="type === 'create'" label="团队成员">
          <el-autocomplete
            v-model="form.username"
            :fetch-suggestions="querySearchAsync"
            placeholder="请输入用户名"
            @select="handleSelect"
          >
            <el-select
              v-model="form.role"
              slot="prepend"
              placeholder="请选择"
              style="width: 90px"
            >
              <el-option label="管理员" value="administrator" />
              <el-option label="成员" value="member" />
              <el-option label="访客" value="visitor" />
            </el-select>
          </el-autocomplete>
          <div>
            <el-tag
              class="tag"
              v-for="member in selectedMembers"
              :key="member.id"
              closable
              type="success"
              @close="removeMember(member)"
            >
              {{ member.username }} :
              <span
                v-for="v of roles"
                :key="v.value"
                v-if="member.role === v.value"
                >{{ v.label }}</span
              >
            </el-tag>
          </div>
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
import { mapActions, mapGetters } from 'vuex'

export default {
  async asyncData({ $axios, query }) {
    const type = query.id !== undefined ? 'update' : 'create'

    const formRules = {
      name: [{ required: true, message: '请输入团队名称' }]
    }

    let form = {
      name: '',
      role: 'member'
    }

    if (type === 'update') {
      const { data } = await $axios.$get('/team/_/' + query.id)
      form = data
      if (!form) {
        return redirect('/404')
      }
    }

    return {
      type,
      formName: 'form',
      formRules,
      form,
      selectedMembers: []
    }
  },
  computed: {
    ...mapGetters({
      _roles: 'workspace/roles'
    }),
    roles() {
      return this._roles.filter(v => v.value !== undefined)
    }
  },
  methods: {
    ...mapActions({
      getWorkspaces: 'workspace/getWorkspaces'
    }),
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
      try {
        await this.$axios.$put('/team/_/' + this.$route.query.id, form)
        this.$success('更新成功.')
        this.$router.back()
      } catch (err) {
        this.$error(err.message)
      }
    },
    async create() {
      const form = this.form
      form.port = +form.port
      try {
        await this.$axios.$post('/team', {
          ...form,
          members: this.selectedMembers.map(v => v.id)
        })
        this.$success('创建成功.')
        this.getWorkspaces()
        this.$router.back()
      } catch (err) {
        console.dir(err)
        this.$error(err.message)
      }
    },
    handleSelect(user) {
      const index = this.selectedMembers.findIndex(v => v.id === user.id)

      if (index >= 0) {
        this.selectedMembers.splice(index, 1)
      }

      this.selectedMembers.push({ ...user, role: this.form.role })
      this.form.username = ''
    },
    querySearchAsync(queryString, cb) {
      if (queryString === '') {
        cb([])
        return
      }

      this.$axios
        .$get('/user/search', {
          params: {
            account: queryString
          }
        })
        .then(({ data: users }) => {
          const list = users.map(v => {
            return {
              ...v,
              value: v.username
            }
          })
          cb(list)
        })
    },
    removeMember(member) {
      const index = this.selectedMembers.findIndex(v => v.id === member.id)
      this.selectedMembers.splice(index, 1)
    }
  }
}
</script>

<style lang="less" scope>
.tag {
  &:not(:last-child) {
    margin-right: 20px;
  }
}
</style>
