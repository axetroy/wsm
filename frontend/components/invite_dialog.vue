<!-- Copyright 2019-2020 Axetroy. All rights reserved. Apache License 2.0. -->
<template>
  <el-dialog title="邀请成员加入" :visible.sync="isShow">
    <el-form
      label-width="100px"
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
        <el-button type="primary" @click="onSubmit()">提交</el-button>
        <el-button @click="cancel()">取消</el-button>
      </el-form-item>
    </el-form>
  </el-dialog>
</template>

<script>
import { mapGetters } from 'vuex'

export default {
  props: {
    visible: {
      type: Boolean
    },
    teamid: {
      type: String
    }
  },
  data() {
    return {
      formName: 'form',
      formRules: {},
      form: {
        username: '',
        role: 'member'
      },
      selectedMembers: [],
      isShow: this.visible || false // 是否显示弹窗
    }
  },
  watch: {
    visible(val) {
      this.isShow = val
    },
    isShow(val) {
      this.$emit('update:visible', val)
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
    },
    async onSubmit() {
      if (!this.selectedMembers.length) {
        return
      }

      const members = this.selectedMembers.map(v => {
        return {
          id: v.id,
          role: 'member'
        }
      })

      try {
        await this.$axios.$post(`/team/_/${this.teamid}/member/invite`, {
          members
        })
        this.$success('邀请成功')
        this.$emit('invite-success', members)
        this.cancel()
      } catch (err) {
        this.$error(err.message)
      }
    },
    cancel() {
      this.isShow = false
      this.selectedMembers = []
      this.form.username = ''
      this.$refs[this.formName].resetFields()
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
