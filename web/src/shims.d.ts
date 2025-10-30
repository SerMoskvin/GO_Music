declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent<{}, {}, any>
  export default component
}

// Для алиасов
declare module '@/*' {
  const value: any
  export default value
}