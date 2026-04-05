<template>
  <div class="home-view">
    <h1>Open-Resume</h1>
    <p class="subtitle">Markdown原生、隐私优先的本地化简历工具</p>
    <button class="create-btn">新建简历</button>
  </div>
</template>

<script setup lang="ts">
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
</style>
