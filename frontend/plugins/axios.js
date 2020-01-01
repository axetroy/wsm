const CookieParse = require('cookie').parse
const ClientCookie = require('js-cookie')

const TOKEN_KEY = 'Authorization'

export default function({ $axios, redirect, req, res }) {
  $axios.defaults.baseHost = '0.0.0.0:9000'
  $axios.defaults.baseURL = `http://${$axios.defaults.baseHost}/v1/`
  $axios.defaults.headers.post['Content-Type'] = 'application/json'
  $axios.defaults.headers.put['Content-Type'] = 'application/json'

  $axios.defaults.headers.common = {
    get [TOKEN_KEY]() {
      const tokenRaw =
        // @ts-ignore
        process.client
          ? ClientCookie.get(TOKEN_KEY)
          : CookieParse(req.headers.cookie || '')[TOKEN_KEY]

      return tokenRaw || ''
    }
  }

  $axios.onRequest(config => {
    console.log('Making request to ' + config.url)
  })

  // Add a response interceptor
  $axios.interceptors.response.use(
    function(response) {
      const data = response.data
      if (!data) {
        return Promise.reject(new Error('unknown error'))
      }

      const { status, message } = data

      if (status != 1) {
        switch (status) {
          // TOKEN 无效
          case 999999:
            // @ts-ignore
            if (process.client) {
              ClientCookie.remove(TOKEN_KEY)
            } else {
              // TODO: 删除 Cookie
              res.setHeader(
                'Set-Cookie',
                `${TOKEN_KEY}=deleted; path=/; expires=Thu, 01 Jan 1970 00:00:00 GMT`
              )
            }
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
