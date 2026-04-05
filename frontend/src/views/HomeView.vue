<template>
  <div class="home-view">
    <h1>Open-Resume</h1>
    <p class="subtitle">Markdown原生、隐私优先的本地化简历工具</p>
    <button class="create-btn" @click="createResume">新建简历</button>
    <p v-if="message" class="message">{{ message }}</p>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { CreateResume } from '../wailsjs/wailsjs/go/main/App'

const router = useRouter()
const message = ref('')

async function createResume() {
  try {
    const resume = await CreateResume('我的简历')
    message.value = `创建成功！ID: ${resume.id}`
    // 跳转到编辑器
    setTimeout(() => {
      router.push(`/editor/${resume.id}`)
    }, 500)
  } catch (err) {
    message.value = `创建失败: ${err}`
    console.error(err)
  }
}
</script>

<style scoped>
.home-view {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 100%;
  padding: 2rem;
}

h1 {
  font-size: 2.5rem;
  font-weight: 700;
  color: #fff;
  margin: 0 0 0.5rem 0;
}

.subtitle {
  font-size: 1rem;
  color: #8b949e;
  margin: 0 0 2rem 0;
}

.create-btn {
  padding: 0.75rem 2rem;
  font-size: 1rem;
  font-weight: 500;
  color: #fff;
  background-color: #238636;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.create-btn:hover {
  background-color: #2ea043;
}

.message {
  margin-top: 1rem;
  padding: 0.5rem 1rem;
  border-radius: 6px;
  font-size: 0.875rem;
}

.message:empty {
  display: none;
}
</style>
