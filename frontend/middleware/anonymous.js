// Copyright 2019-2020 Axetroy. All rights reserved. Apache License 2.0.

// 要求用户匿名
export default function ({ store, redirect }) {
  if (store.getters['user/profile']) {
    return redirect('/')
  }
}
