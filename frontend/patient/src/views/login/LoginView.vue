<template>
  <div class="login-container">
    <div class="login-card">
      <div class="login-header">
        <HeartOutlined class="login-logo" />
        <h1>患者服务中心</h1>
        <p>HIS-Go Hospital Information System</p>
      </div>
      <a-form :model="formState" :rules="rules" @finish="handleLogin" layout="vertical" size="large">
        <a-form-item name="username">
          <a-input v-model:value="formState.username" placeholder="用户名">
            <template #prefix><UserOutlined /></template>
          </a-input>
        </a-form-item>
        <a-form-item name="password">
          <a-input-password v-model:value="formState.password" placeholder="密码">
            <template #prefix><LockOutlined /></template>
          </a-input-password>
        </a-form-item>
        <a-form-item>
          <a-button type="primary" html-type="submit" :loading="loading" block>登 录</a-button>
        </a-form-item>
      </a-form>
      <div class="login-tips">
        <p>演示账号: demo-patient / demo123</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import { UserOutlined, LockOutlined, HeartOutlined } from '@ant-design/icons-vue'
import { useAuthStore } from '@/store/auth'

const router = useRouter()
const authStore = useAuthStore()
const loading = ref(false)

const formState = reactive({ username: '', password: '' })
const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
}

async function handleLogin() {
  loading.value = true
  try {
    await authStore.login(formState)
    message.success('登录成功')
    router.push('/dashboard')
  } catch (err: any) {
    message.error(err.message || '登录失败')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-container {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #1890ff 0%, #36cfc9 100%);
  padding: 24px;
}

.login-card {
  width: 100%;
  max-width: 380px;
  padding: 40px 32px;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.15);
}

.login-header {
  text-align: center;
  margin-bottom: 32px;
}

.login-logo {
  font-size: 48px;
  color: #1890ff;
  margin-bottom: 12px;
}

.login-header h1 {
  font-size: 22px;
  color: #1a1a2e;
  margin-bottom: 4px;
}

.login-header p {
  color: #999;
  font-size: 12px;
}

.login-tips {
  margin-top: 20px;
  text-align: center;
  color: #ccc;
  font-size: 11px;
}
</style>
