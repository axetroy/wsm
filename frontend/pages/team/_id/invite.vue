<template>
  <div class="main">
    <el-card shadow="never">
      <div slot="header">
        <h4>邀请成员加入</h4>
      </div>
      <el-form
        label-width="160px"
        status-icon
        :model="form"
        :ref="formName"
        :rules="formRules"
      >
        <el-form-item label="邀请成员">
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
export default {
  async asyncData({ $axios, query }) {},
  data() {
    return {
      teamID: this.$route.params.id,
      formName: 'form',
      formRules: {},
      form: {
        username: '',
        role: 'member'
      },
      selectedMembers: [],
      roles: [
        {
          label: '拥有者',
          value: 'owner'
        },
        {
          label: '管理员',
          value: 'administrator'
        },
        {
          label: '成员',
          value: 'member'
        },
        {
          label: '访客',
          value: 'visitor'
        }
      ]
    }
  },
  methods: {
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
    async onSubmit() {
      if (!this.selectedMembers.length) {
        return
      }

      const members = this.selectedMembers

      try {
        await this.$axios.$post(`/team/_/${this.teamID}/member/invite`, {
          members: members.map(v => {
            return {
              id: v.id,
              role: 'member'
            }
          })
        })
        this.$success('邀请成功')
        this.$router.back()
      } catch (err) {
        this.$error(err.message)
      }
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
