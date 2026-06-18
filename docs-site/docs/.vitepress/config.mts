import { defineConfig } from 'vitepress'

export default defineConfig({
  title: 'Sanjie',
  description: '玩笑性质的三界移动办公系统',
  base: '/sanjie/',
  cleanUrls: true,
  themeConfig: {
    logo: '/logo.svg',
    siteTitle: 'Sanjie',
    nav: [
      { text: '指南', link: '/guide/features' },
      { text: '演示', link: '/demo/' },
      { text: '接口', link: '/api/' },
      { text: '部署', link: '/deploy/' },
      { text: '开发', link: '/dev/roadmap' },
      { text: '贡献', link: '/dev/contributing' }
    ],
    sidebar: [
      {
        text: '开始',
        items: [
          { text: '项目首页', link: '/' },
          { text: '快速开始', link: '/README' },
          { text: '演示', link: '/demo/' }
        ]
      },
      {
        text: '指南',
        items: [
          { text: '功能介绍', link: '/guide/features' },
          { text: '架构说明', link: '/guide/architecture' }
        ]
      },
      {
        text: '参考',
        items: [
          { text: '接口文档', link: '/api/' },
          { text: '部署说明', link: '/deploy/' },
          { text: '开发路线', link: '/dev/roadmap' },
          { text: '贡献指南', link: '/dev/contributing' }
        ]
      }
    ],
    socialLinks: [
      { icon: 'github', link: 'https://github.com/ijry/sanjie' }
    ],
    search: {
      provider: 'local'
    },
    footer: {
      message: 'Released under the MIT License.',
      copyright: 'Copyright © 2026 ijry'
    }
  }
})
