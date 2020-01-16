<!-- Copyright 2019-2020 Axetroy. All rights reserved. Apache License 2.0. -->
<template>
  <div class="login-container">
    <div id="login-form">
      <h3>注册帐号</h3>

      <el-form
        :model="loginForm"
        status-icon
        :rules="loginFormRules"
        :ref="formName"
        label-width="80px"
        label-position="top"
      >
        <el-form-item prop="username" label="用户名">
          <el-input
            type="text"
            v-model="loginForm.username"
            auto-complete="off"
            placeholder="用户名"
            @keyup.enter.native="submitForm()"
            clearable
          >
            <i class="el-icon-user" slot="prefix" />
          </el-input>
        </el-form-item>
        <el-form-item prop="password" label="密码">
          <el-input
            type="password"
            v-model="loginForm.password"
            auto-complete="off"
            placeholder="密码"
            clearable
          >
            <i class="el-icon-lock" slot="prefix" />
          </el-input>
        </el-form-item>
        <el-form-item prop="password_confirm" label="确认密码">
          <el-input
            type="password"
            v-model="loginForm.password_confirm"
            auto-complete="off"
            placeholder="确认密码"
            @keyup.enter.native="submitForm()"
            clearable
          >
            <i class="el-icon-lock" slot="prefix" />
          </el-input>
        </el-form-item>
        <el-form-item>
          <el-button
            :loading="loading"
            type="primary"
            style="width:100%"
            @click="submitForm()"
            >{{ loading ? '注册中...' : '注 册' }}</el-button
          >
        </el-form-item>

        <div>
          <el-button
            type="text"
            size="small"
            @click="$router.push({ name: 'login' })"
          >
            登 录
          </el-button>
          <el-button
            class="pull-right"
            type="text"
            size="small"
            @click="$router.push({ name: 'forgot' })"
          >
            忘记密码
          </el-button>
        </div>
      </el-form>
    </div>
  </div>
</template>

<script>
export default {
  layout: 'entry',
  data() {
    return {
      formName: 'form',
      loginForm: {
        username: '',
        password: '',
        password_confirm: ''
      },
      loginFormRules: {
        username: [{ required: true, message: '请输入用户名' }],
        password: [{ required: true, message: '请输入密码' }],
        password_confirm: [
          { required: true, message: '请输入确认密码' },
          {
            validator: (rule, value, callback) => {
              if (value === '') {
                callback(new Error('请再次输入密码'))
              } else if (value !== this.loginForm.password) {
                callback(new Error('两次输入密码不一致!'))
              } else {
                callback()
              }
            }
          }
        ]
      },
      loading: false
    }
  },
  methods: {
    submitForm() {
      const form = this.loginForm
      this.$refs[this.formName].validate(valid => {
        if (valid) {
          const body = { ...form }
          this.loading = true

          this.$axios
            .$post('/auth/signup', body)
            .then(({ data: profile }) => {
              this.loading = false
              this.$success('注册成功...')
              this.$router.push({ name: 'login' })
            })
            .catch(err => {
              this.loading = false
              this.$error(err instanceof Error ? err.message : String(err))
            })
        }
      })
    }
  }
}
</script>

<style lang="less" scoped>
.login-container {
  position: absolute;
  width: 100%;
  height: 100%;
  margin: 0;
  padding: 0;
  top: 0;
  left: 0;
  background-repeat: no-repeat;
  background-clip: initial;
  background-position: center;
  background-color: #2e323b;
}

#login-form {
  position: relative;
  top: 50%;
  left: 0;
  width: 420px;
  margin: 0 auto;
  transform: translateY(-50%);
  background-color: #fff;

  h3 {
    padding-top: 20px;
    margin: 0;
    text-align: center;
  }

  .el-form {
    background-color: #fff;
    padding: 20px;
  }
}

label {
  color: #fff;
}

.no-selected {
  user-select: none;
}

.align-middle {
  vertical-align: middle;
  cursor: pointer;
  color: #748597;
}
</style>
