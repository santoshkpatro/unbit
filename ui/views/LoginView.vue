<script setup>
import { reactive } from 'vue'
import { useRouter } from 'vue-router'
import { loginAPI } from '@/api/auth'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()

const router = useRouter()

const form = reactive({
  email: '',
  password: '',
})

async function onFinish() {
  const user = await loginAPI({
    email: form.email,
    password: form.password,
  })

  authStore.setLoggedInUser(user)
  router.push({ name: 'issue-list' })
}
</script>

<template>
  <div class="login-wrap">
    <a-form layout="vertical" class="login-form" :model="form" @finish="onFinish" autocomplete="on">
      <a-typography-title :level="3" class="title">Sign in</a-typography-title>

      <a-form-item name="email" :rules="[{ required: true, message: 'Please enter your email' }]">
        <a-input
          v-model:value="form.email"
          type="email"
          size="large"
          placeholder="Email"
          :bordered="false"
          autofocus
        />
      </a-form-item>

      <a-form-item
        name="password"
        :rules="[{ required: true, message: 'Please enter your password' }]"
      >
        <a-input-password
          v-model:value="form.password"
          size="large"
          placeholder="Password"
          :bordered="false"
        />
      </a-form-item>

      <a-button type="primary" html-type="submit" size="large" block> Log in </a-button>
    </a-form>
  </div>
</template>

<style scoped>
.login-wrap {
  min-height: 100dvh;
  display: flex;
  align-items: center;
  justify-content: center;
}

.login-form {
  width: min(92vw, 360px);
}

.title {
  margin: 0 0 8px;
  text-align: center;
}

:deep(.ant-form-item) {
  margin-bottom: 16px;
}

:deep(.ant-input-affix-wrapper),
:deep(.ant-input) {
  box-shadow: none !important;
  background: transparent;
}

:deep(.ant-btn:focus-visible) {
  outline: none;
  box-shadow: none;
}
</style>
