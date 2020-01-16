// Copyright 2019-2020 Axetroy. All rights reserved. Apache License 2.0.
import Vue from 'vue'
import format from 'date-fns/format'

Vue.filter('dateformat', (value, layout = 'yyyy-MM-dd HH:mm:ss') => {
  if (!value) return ''
  const type = typeof value
  if (type === 'number' || type === 'string') {
    value = new Date(value)
  }
  return format(value, layout)
})

export default function(context, inject) {
  context.$dateformat = format
  inject('dateformat', context.$dateformat)
}
