/// <reference types="vite/client" />

declare module '*.vue' {
    import type {DefineComponent} from 'vue'
    const component: DefineComponent<{}, {}, any>
    export default component
}

// Monaco Editor 全局类型声明
declare module 'monaco-editor' {
    export * from 'monaco-editor/esm/vs/editor/editor.api'
}
