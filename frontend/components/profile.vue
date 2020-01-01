<template>
  <el-dialog
    class="profile-dialog"
    title="账户管理"
    :visible.sync="visible"
    :before-close="handleClose"
    destroy-on-close
  >
    <el-container>
      <el-aside width="200px">
        <el-menu
          class="aside-menu"
          :default-active="activedMenu"
          @select="handleSelectMenu"
        >
          <el-menu-item index="info">个人信息</el-menu-item>
          <el-menu-item index="password">账号密码</el-menu-item>
          <el-menu-item index="bind" disabled="">绑定第三方登录</el-menu-item>
          <el-menu-item index="secrecy" disabled="">隐私设置</el-menu-item>
          <el-menu-item index="push" disabled="">邮件推送</el-menu-item>
        </el-menu>
      </el-aside>
      <el-main>
        <el-container v-show="activedMenu === 'info'">
          <el-form
            :model="profileForm"
            status-icon
            :rules="profileFormRules"
            :ref="profileFormName"
            label-width="80px"
            style="width: 100%"
          >
            <el-form-item label="ID">
              {{ user ? user.id : '' }}
            </el-form-item>
            <el-form-item label="用户名">
              <div>{{ user ? user.username : '' }}</div>
            </el-form-item>
            <el-form-item label="头像">
              <el-upload
                class="avatar-uploader"
                action="https://jsonplaceholder.typicode.com/posts/"
                :show-file-list="false"
                :on-success="handleAvatarSuccess"
                :before-upload="beforeAvatarUpload"
              >
                <img v-if="user.avatar" :src="user.avatar" class="avatar" />
                <i v-else class="el-icon-plus avatar-uploader-icon" />
              </el-upload>
            </el-form-item>
            <el-form-item prop="nickname" label="昵称">
              <el-input v-model="profileForm.nickname" auto-complete="off">
                <i class="el-icon-edit el-input__icon" slot="prefix" />
              </el-input>
            </el-form-item>
            <el-form-item prop="sex" label="性别">
              <el-select v-model="profileForm.gender">
                <el-option
                  v-for="v in genders"
                  :key="v.value"
                  :label="v.label"
                  :value="v.value"
                >
                </el-option>
              </el-select>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="submitUpdateProfile()"
                >确定</el-button
              >
            </el-form-item>
          </el-form>
        </el-container>
        <el-container v-show="activedMenu === 'password'">
          <el-form
            :model="passwordForm"
            status-icon
            :rules="passwordFormRules"
            :ref="passwordFormName"
            label-width="80px"
            style="width: 100%"
          >
            <el-form-item prop="old" label="旧密码" required>
              <el-input
                type="password"
                v-model="passwordForm.old"
                auto-complete="off"
                placeholder="请输入旧密码"
              >
                <i
                  class="i-icon i-icon-22 el-icon-edit el-input__icon"
                  slot="prefix"
                />
              </el-input>
            </el-form-item>
            <el-form-item prop="new" label="新密码" required>
              <el-input
                type="password"
                v-model="passwordForm.new"
                auto-complete="off"
                placeholder="请输入新密码"
              >
                <i
                  class="i-icon i-icon-22 el-icon-edit el-input__icon"
                  slot="prefix"
                />
              </el-input>
            </el-form-item>
            <el-form-item prop="confirm" label="确认密码" required>
              <el-input
                type="password"
                v-model="passwordForm.confirm"
                auto-complete="off"
                placeholder="请再次输入新密码"
              >
                <i
                  class="i-icon i-icon-22 el-icon-edit el-input__icon"
                  slot="prefix"
                />
              </el-input>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="submitChangePassword()"
                >确定</el-button
              >
            </el-form-item>
          </el-form>
        </el-container>
        <el-container v-show="activedMenu === 'bind'"
          >绑定第三方登录</el-container
        >
        <el-container v-show="activedMenu === 'secrecy'">隐私设置</el-container>
        <el-container v-show="activedMenu === 'push'">邮件推送</el-container>
      </el-main>
    </el-container>
    <div slot="footer" class="profile-dialog__footer">
      <el-button @click="handleClose">关闭</el-button>
    </div>
  </el-dialog>
</template>

<script>
import { mapGetters, mapActions } from 'vuex'

export default {
  props: {
    visible: {
      type: Boolean,
      default() {
        return false
      }
    }
  },
  computed: {
    ...mapGetters({
      user: 'user'
    })
  },
  data() {
    return {
      activedMenu: 'info',
      passwordForm: {
        old: '',
        new: '',
        confirm: ''
      },
      passwordFormName: 'password',
      passwordFormRules: {
        old: [{ required: true, message: '请输入旧密码' }],
        new: [
          { required: true, message: '请输入新密码' },
          {
            validator: (rule, value, callback) => {
              if (value === '') {
                callback(new Error('请再次输入密码'))
              } else if (value === this.passwordForm.old) {
                callback(new Error('新密码不能和旧密码相同!'))
              } else if (value !== this.passwordForm.confirm) {
                callback(new Error('两次输入密码不一致!'))
              } else {
                callback()
              }
            }
          }
        ],
        confirm: [
          { required: true, message: '请重新输入新密码' },
          {
            validator: (rule, value, callback) => {
              if (value === '') {
                callback(new Error('请重新输入新密码'))
              } else if (value !== this.passwordForm.new) {
                callback(new Error('两次输入密码不一致!'))
              } else {
                const form = this.$refs[this.passwordFormName]
                form.validateField('new')
                callback()
              }
            }
          }
        ]
      },

      profileForm: {
        nickname: '',
        gender: '',
        avatar: ''
      },
      profileFormName: 'profile',
      profileFormRules: {},
      genders: [
        {
          label: '未知',
          value: 0
        },
        {
          label: '男',
          value: 1
        },
        {
          label: '女',
          value: 2
        }
      ]
    }
  },
  watch: {
    visible(visible) {
      // reset
      for (let key in this.passwordForm) {
        this.passwordForm[key] = ''
      }

      this.profileForm.nickname = this.user?.nickname
      this.profileForm.gender = this.user?.gender
      this.profileForm.avatar = this.user?.avatar
    }
  },
  methods: {
    ...mapActions({
      getProfile: 'getProfile'
    }),
    handleClose() {
      this.$emit('update:visible', false)
    },
    handleSelectMenu(index, indexPath) {
      this.activedMenu = index
    },
    submitChangePassword() {
      const form = this.$refs[this.passwordFormName]
      form.validate(valid => {
        if (valid) {
          const { old: old_password, new: new_password } = this.passwordForm
          const body = {
            old_password,
            new_password
          }

          this.$axios
            .$put('/user/password', body)
            .then(() => {
              this.$success('修改成功...')
              form.resetFields()
            })
            .catch(err => {
              this.$error(err.message)
            })
        }
      })
    },
    submitUpdateProfile() {
      const form = this.$refs[this.profileFormName]
      form.validate(valid => {
        if (valid) {
          const { nickname, gender, avatar } = this.profileForm
          const body = {
            nickname,
            gender,
            avatar
          }

          this.$axios
            .$put('/user/profile', body)
            .then(() => {
              this.$success('修改成功...')

              // 更新资料
              this.getProfile().then(() => {
                form.clearValidate()
              })
            })
            .catch(err => {
              this.$error(err.message)
            })
        }
      })
    },
    handleAvatarSuccess(res, file) {
      this.imageUrl = URL.createObjectURL(file.raw)
    },
    beforeAvatarUpload(file) {
      const isJPG = file.type === 'image/jpeg'
      const isLt2M = file.size / 1024 / 1024 < 2

      if (!isJPG) {
        this.$message.error('上传头像图片只能是 JPG 格式!')
      }
      if (!isLt2M) {
        this.$message.error('上传头像图片大小不能超过 2MB!')
      }
      return isJPG && isLt2M
    }
  }
}
</script>

<style lang="less">
.profile-dialog {
  .el-dialog__body {
    padding: 0;
  }

  .aside-menu {
    height: 100%;
  }
}
</style>

<style lang="less" scoped>
.avatar-uploader .el-upload {
  border: 1px dashed #d9d9d9;
  border-radius: 6px;
  cursor: pointer;
  position: relative;
  overflow: hidden;
}
.avatar-uploader .el-upload:hover {
  border-color: #409eff;
}
.avatar-uploader-icon {
  font-size: 28px;
  color: #8c939d;
  width: 178px;
  height: 178px;
  line-height: 178px;
  text-align: center;
  border: 1px dashed #e5e5e5;
}
.avatar {
  width: 178px;
  height: 178px;
  display: block;
}
</style>
