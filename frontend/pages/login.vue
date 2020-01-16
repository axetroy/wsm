<!-- Copyright 2019-2020 Axetroy. All rights reserved. Apache License 2.0. -->
<template>
  <div class="login-container">
    <div id="login-form">
      <h3>登陆</h3>
      <el-form
        :model="loginForm"
        status-icon
        :rules="loginFormRules"
        :ref="formName"
        label-width="0"
      >
        <el-form-item prop="account">
          <el-input
            type="text"
            v-model="loginForm.account"
            auto-complete="off"
            placeholder="用户名"
            @keyup.enter.native="submitForm()"
            clearable
          >
            <i class="el-icon-user" slot="prefix" />
          </el-input>
        </el-form-item>
        <el-form-item prop="password">
          <el-input
            type="password"
            v-model="loginForm.password"
            auto-complete="off"
            placeholder="密码"
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
            >{{ loading ? '登录中...' : '登 录' }}</el-button
          >
        </el-form-item>

        <div>
          <el-button
            type="text"
            size="small"
            @click="$router.push({ name: 'register' })"
          >
            注册帐号
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
import { mapActions } from 'vuex'

export default {
  layout: 'entry',
  data() {
    return {
      formName: 'form',
      loginForm: {
        account: '',
        password: ''
      },
      loginFormRules: {
        account: [{ required: true, message: '请输入用户名' }],
        password: [{ required: true, message: '请输入密码' }]
      },
      loading: false
    }
  },
  methods: {
    ...mapActions({
      login: 'user/login'
    }),
    submitForm() {
      const form = this.loginForm
      this.$refs[this.formName].validate(valid => {
        if (valid) {
          const body = { ...form }
          this.loading = true

          this.login(body)
            .then(() => {
              this.loading = false
              this.$success('登陆成功...')
              this.$router.push({ path: '/' })
            })
            .catch(err => {
              this.loading = false
              this.$error(err.message)
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
    padding: 20px 0;
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
