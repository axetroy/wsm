const TOKEN_KEY = 'Authorization'

export default function({ store, $axios, redirect, req, res }) {
  $axios.defaults.baseURL = process.client
    ? `/v1`
    : `http://0.0.0.0:${process.env.PORT}/v1`

  $axios.defaults.headers.common = {
    'Content-Type': 'application/json',
    get [TOKEN_KEY]() {
      return store.getters['user/token']
    }
  }

  $axios.onRequest(config => {
    console.log(`Making request to /v1${config.url}`)
  })

  // Add a response interceptor
  $axios.interceptors.response.use(
    function(response) {
      const data = response.data
      if (!data) {
        return Promise.reject(new Error('unknown error'))
      }

      const { status, message } = data

      if (status !== 1) {
        switch (status) {
          // TOKEN 无效
          case 999999:
            store.dispatch('user/logout')
            // @ts-ignore
            // if (process.client) {
            //   cookie.remove(TOKEN_KEY)
            // } else {
            //   res.setHeader(
            //     'Set-Cookie',
            //     `${TOKEN_KEY}=deleted; path=/; expires=Thu, 01 Jan 1970 00:00:00 GMT`
            //   )
            // }
            redirect('/login')
            break
          default:
            break
        }
        return Promise.reject(new Error(message || 'unknown error'))
      }

      // Any status code that lie within the range of 2xx cause this function to trigger
      // Do something with response data
      return response
    },
    function(error) {
      // Any status codes that falls outside the range of 2xx cause this function to trigger
      // Do something with response error
      return Promise.reject(error)
    }
  )

  $axios.onError(error => {
    const code = parseInt(error.response && error.response.status)
    if (code === 400) {
      redirect('/400')
    }
  })
}
