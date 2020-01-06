<template>
  <terminal :rows="0" class="terminal" ref="terminal" :host="host" />
</template>

<script>
import { mapActions } from 'vuex'
import Terminal from '~/components/terminal'

export default {
  layout: 'empty',
  components: {
    Terminal
  },
  head() {
    return {
      title: `WSM: ${this.host.name} - ${
        this.team ? this.team.name : this.profile.username
      }`
    }
  },
  async asyncData({ store, params, query }) {
    const hostId = params.id
    const teamId = query.team_id

    const profile = store.getters['user/profile']

    if (teamId) {
      const [team, host] = await Promise.all([
        store.dispatch('team/fetchTeamById', teamId),
        store.dispatch('team/fetchTeamHostById', {
          teamId,
          hostId
        })
      ])

      return {
        profile,
        host,
        team
      }
    } else {
      const host = await store.dispatch('host/fetchHostById', hostId)

      return {
        profile,
        host,
        team: undefined
      }
    }
  },
  methods: {
    ...mapActions({
      fetchTeamHostById: 'team/fetchTeamHostById'
    })
  }
}
</script>

<style>
#__layout {
  background-color: #000;
}
</style>
