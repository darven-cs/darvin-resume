# Phase 2 Plan 02-01 SUMMARY: Monaco Editor 集成

## 执行状态: 已完成

## 变更文件

### 新增文件

1. **`/home/darven/桌面/dev_app/darvin-resume/frontend/src/components/MonacoEditor.vue`**
   - Vue 3 组件封装 Monaco Editor
   - 支持 v-model 双向绑定 (modelValue + update:modelValue emit)
   - 支持 change emit
   - Markdown 语法高亮 (language: 'markdown')
   - VS Code 兼容配置: 14px 字体, 1.6行高, undo/redo, multi-cursor, find/replace
   - glyphMargin: true (行级交互需要)
   - minimap 禁用
   - automaticLayout: true
   - defineExpose 暴露 getEditor() 和 getMonaco()
   - onUnmounted 正确 dispose 编辑器实例
   - watch modelValue 变化同步到编辑器

### 修改文件

2. **`/home/darven/桌面/dev_app/darvin-resume/frontend/package.json`**
   - 新增依赖:
     - `@monaco-editor/loader@^1.4.0`
     - `monaco-editor@^0.45.0`
   - 新增 devDependencies:
     - `vite-plugin-monaco-editor@^1.1.0`

3. **`/home/darven/桌面/dev_app/darvin-resume/frontend/vite.config.ts`**
   - 导入 vite-plugin-monaco-editor
   - 配置 languageWorkers: ['editorWorkerService']
   - 配置 optimizeDeps.include: ['monaco-editor']
   - 注意: markdown 不是有效的 worker 类型，插件只支持 editorWorkerService, css, html, json, typescript

4. **`/home/darven/桌面/dev_app/darvin-resume/frontend/src/vite-env.d.ts`**
   - 新增 Monaco Editor 全局类型声明

5. **`/home/darven/桌面/dev_app/darvin-resume/frontend/src/views/HomeView.vue`** (预先存在的问题修复)
   - 修复 wailsjs 导入路径从 `'../../wailsjs/go/main/App'` 改为 `'../wailsjs/wailsjs/go/main/App'`

## 验证结果

- `npm run build` 成功，无 TypeScript 错误
- 构建产物:
  - dist/index.html (1.32 KiB)
  - dist/assets/EditorView.9825bd1b.js (0.34 KiB)
  - dist/assets/EditorView.ca16cbef.css (0.16 KiB)
  - dist/assets/index.f81c3031.css (1.53 KiB)
  - dist/assets/index.0d249789.js (88.66 KiB)

## 技术说明

### Monaco Editor 加载方式
使用 `@monaco-editor/loader` 动态加载 Monaco Editor，配置 CDN 源:
```typescript
loader.config({
  paths: {
    vs: 'https://cdn.jsdelivr.net/npm/monaco-editor@0.45.0/min/vs',
  },
})
```

### vite-plugin-monaco-editor 配置
- 使用 `editorWorkerService` worker 支持编辑器基础功能
- markdown 语言不需要专用 worker，Monaco 内置支持
- 配置 `optimizeDeps.include` 预打包 monaco-editor 依赖

## 待办事项 (后续 Phase)
- EditorView.vue 需要使用 MonacoEditor.vue 组件替换占位符
- 可能需要调整 Monaco 配置以匹配简历编辑场景
