import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'path'
import Inspector from 'unplugin-vue-dev-locator/vite'
import traeBadgePlugin from 'vite-plugin-trae-solo-badge'

// https://vite.dev/config/
export default defineConfig(({ mode }) => {
  const isDev = mode === 'development'
  return {
    build: {
      sourcemap: 'hidden',
    },
    plugins: [
      vue(),
      // 仅开发环境启用定位器
      ...(isDev ? [Inspector()] : []),
      // badge 插件自身已通过 prodOnly 控制，仅在生产渲染
      traeBadgePlugin({
        variant: 'dark',
        position: 'bottom-right',
        prodOnly: true,
        clickable: true,
        clickUrl: 'https://www.trae.ai/solo?showJoin=1',
        autoTheme: true,
        autoThemeTarget: '#app',
      }),
    ],
    resolve: {
      alias: {
        '@': path.resolve(__dirname, './src'), // ✅ 定义 @ = src
      },
    },
    server: {
      host: '127.0.0.1',
      port: 3333,
      proxy: {
        '/api': {
          target: 'http://127.0.0.1:8888',
          changeOrigin: true,
          secure: false,
        },
      },
    },
  }
})
