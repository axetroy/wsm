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
        <el-menu :default-active="activedMenu" @select="handleSelectMenu">
          <el-menu-item index="info">个人信息</el-menu-item>
          <el-menu-item index="password">账号密码</el-menu-item>
          <el-menu-item index="bind" disabled="">绑定第三方登录</el-menu-item>
          <el-menu-item index="secrecy" disabled="">隐私设置</el-menu-item>
          <el-menu-item index="push" disabled="">邮件推送</el-menu-item>
        </el-menu>
      </el-aside>
      <el-main>
        <el-container v-show="activedMenu === 'info'">个人信息</el-container>
        <el-container v-show="activedMenu === 'password'"
          >账号密码</el-container
        >
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
export default {
  props: {
    visible: {
      type: Boolean,
      default() {
        return false
      }
    }
  },
  data() {
    return {
      activedMenu: 'info'
    }
  },
  methods: {
    handleClose() {
      this.$emit('update:visible', false)
    },
    handleSelectMenu(index, indexPath) {
      // console.log(index, indexPath)
      this.activedMenu = index
    }
  }
}
</script>

<style lang="less">
.profile-dialog {
  .el-dialog__body {
    padding: 0;
  }
}
</style>
