<template>
  <el-container style="height: 100vh;">
    <el-aside width="200px" class="aside">
      <el-menu :default-active="$route.path" router>
        <!-- <el-menu-item index="/">
          <i class="el-icon-s-home" />
          首页
        </el-menu-item> -->
        <el-menu-item index="/host">
          <i class="el-icon-s-promotion" />
          服务器管理
        </el-menu-item>
        <el-submenu index="/team">
          <template slot="title">
            <i class="el-icon-location" />
            <span>团队管理</span>
          </template>
          <el-menu-item index="/team"
            >{{ currentWorkspace === undefined ? '我的团队' : '团队信息' }}
          </el-menu-item>
          <el-menu-item index="/team/invite"
            >{{ currentWorkspace === undefined ? '邀请我的' : '团队邀请' }}
          </el-menu-item>
        </el-submenu>
      </el-menu>
    </el-aside>

    <el-container>
      <el-header class="header text-right">
        <div class="workspace">
          <span>
            <i class="el-icon-menu" />
            工作空间 |</span
          >
          <el-select
            v-model="workspace"
            placeholder="请选择工作区"
            @change="onChangeWorkspace"
          >
            <el-option
              v-for="item in workspaces"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            >
            </el-option>
          </el-select>
        </div>

        <el-dropdown @command="handleCommand" trigger="hover">
          <span class="username">
            <i class="el-icon-user" />
            欢迎您，{{ user ? user.nickname : '' }}
          </span>
          <el-dropdown-menu slot="dropdown">
            <el-dropdown-item command="profile">个人资料</el-dropdown-item>
            <el-dropdown-item command="logout">登出</el-dropdown-item>
          </el-dropdown-menu>
        </el-dropdown>
      </el-header>

      <el-main>
        <nuxt />
      </el-main>
    </el-container>

    <profile :visible.sync="profileDialogVisible" />

    <console />
  </el-container>
</template>

<script>
import { mapGetters, mapActions } from 'vuex'
import ClientCookie from 'js-cookie'
import Console from '../components/console'
import Profile from '../components/profile'

export default {
  components: {
    Console,
    Profile
  },
  data() {
    return {
      profileDialogVisible: false,
      workspace: this.currentWorkspace
    }
  },
  computed: {
    ...mapGetters({
      user: 'user',
      hosts: 'console/hosts',
      workspaces: 'workspace/workspaces',
      currentWorkspace: 'workspace/current'
    })
  },
  watch: {
    currentWorkspace(val) {
      this.workspace = val
    }
  },
  methods: {
    ...mapActions({
      getWorkspaces: 'workspace/getWorkspaces',
      switchWorkspace: 'workspace/switchWorkspace'
    }),
    async fetchWorkspaces() {
      await this.getWorkspaces()
      this.workspace = this.currentWorkspace
    },
    handleCommand(command) {
      switch (command) {
        case 'logout':
          this.logout()
          break
        case 'profile':
          this.profile()
          break
        default:
        //
      }
    },
    logout() {
      this.$router.push({ name: 'login' })
      this.$store.commit('SET_USER', null)
      this.onChangeWorkspace(undefined)

      ClientCookie.remove('Authorization')
      ClientCookie.remove('workspace')
    },
    profile() {
      this.profileDialogVisible = true
    },
    onChangeWorkspace(workspaceID) {
      this.switchWorkspace(workspaceID)
    }
  },
  mounted() {
    this.fetchWorkspaces()
  }
}
</script>

<style lang="less">
@bg-color: #2e323b;
@bg-color-darken1: darken(@bg-color, 5%);
@bg-color-darken2: darken(@bg-color, 10%);

html,
body,
#__nuxt,
#__program,
#__layout {
  padding: 0;
  margin: 0;
  height: 100%;
  overflow: hidden;
}

.aside {
  width: 200px;
  color: #333;
  background-color: @bg-color;

  .el-menu {
    border-right: solid 1px @bg-color;
  }

  .el-menu-item:hover {
    background-color: @bg-color-darken1;
  }

  .el-menu-item,
  .el-submenu {
    background-color: @bg-color;
    color: #fff;

    &.is-active {
      background-color: @bg-color-darken2;
    }

    .el-submenu__title {
      &.active,
      &:hover {
        background-color: @bg-color-darken1;
      }

      &,
      i {
        color: #fff !important;
      }
    }
  }
}
</style>

<style scoped lang="less">
@bg-color: #2e323b;
.header {
  position: relative;
  background-color: @bg-color;
  color: #fff;
  line-height: 60px;

  .username {
    cursor: pointer;
    color: #fff;

    i {
      margin-right: 5px;
    }
  }

  .workspace {
    position: absolute;
    left: 50%;
    top: 50%;
    min-width: 250px;
    transform: translateX(-50%) translateY(-50%);
    margin-left: -100px;
  }
}
</style>
