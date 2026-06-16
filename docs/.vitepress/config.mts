import { defineConfig } from 'vitepress'

export default defineConfig({
  title: "binmate",
  description: "A cross-platform binary manager",
  srcDir: './src',
  vite: {
    server: {
      port: 31560
    }
  },
  themeConfig: {
    nav: [
      { text: 'Home', link: '/' },
      { text: 'Examples', link: '/markdown-examples' }
    ],

    sidebar: [
      {
        text: 'Examples',
        items: [
          { text: 'Markdown Examples', link: '/markdown-examples' },
          { text: 'Runtime API Examples', link: '/api-examples' }
        ]
      }
    ],

    socialLinks: [
      { icon: 'github', link: 'https://github.com/vuejs/vitepress' }
    ]
  }
})
