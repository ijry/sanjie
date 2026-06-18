import { defineConfig } from 'vite'
import uni from '@dcloudio/vite-plugin-uni'
import UniUpRoot from 'uview-plus/libs/root/index.js'

function rootIndexFallback() {
  return {
    name: 'sanjie-root-index-fallback',
    configureServer(server) {
      server.middlewares.use((req, _res, next) => {
        if (req.url === '/') {
          req.url = '/index.html'
        }
        next()
      })
    }
  }
}

export default defineConfig({
  plugins: [
    rootIndexFallback(),
    UniUpRoot({
      rootFileName: 'App.up',
      autoCreateRootFile: true
    }),
    uni()
  ],
  css: {
    preprocessorOptions: {
      scss: {
        additionalData: [
          '@import "uview-plus/libs/css/mixin.scss";',
          '@import "uview-plus/libs/css/theme-vars.scss";',
          '@import "uview-plus/libs/css/theme-legacy-bridge.scss";',
          '@import "@/uni.scss";'
        ].join('\n'),
        silenceDeprecations: ['legacy-js-api', 'import']
      }
    }
  }
})
