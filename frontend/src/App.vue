<template>
  <div class="app-layout">
    <aside class="sidebar">
      <div class="sidebar-header">
        <h2>Open-Resume</h2>
      </div>
      <nav class="sidebar-nav">
        <RouterLink to="/">简历列表</RouterLink>
      </nav>
    </aside>
    <main class="main-content">
      <router-view />
    </main>
  </div>
</template>

<script setup lang="ts">
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html, body, #app {
  height: 100%;
  width: 100%;
  overflow: hidden;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, sans-serif;
  background-color: #1e1e1e;
  color: #e0e0e0;
}

.app-layout {
  display: flex;
  height: 100%;
  width: 100%;
}

.sidebar {
  width: 240px;
  min-width: 240px;
  background-color: #252526;
  border-right: 1px solid #3c3c3c;
  display: flex;
  flex-direction: column;
}

.sidebar-header {
  padding: 1rem;
  border-bottom: 1px solid #3c3c3c;
}

.sidebar-header h2 {
  font-size: 1.25rem;
  font-weight: 600;
  color: #fff;
}

.sidebar-nav {
  padding: 1rem 0;
}

.sidebar-nav a {
  display: block;
  padding: 0.75rem 1rem;
  color: #e0e0e0;
  text-decoration: none;
  transition: background-color 0.2s;
}

.sidebar-nav a:hover,
.sidebar-nav a.router-link-active {
  background-color: #37373d;
}

.main-content {
  flex: 1;
  overflow: auto;
  background-color: #1e1e1e;
}
</style>
