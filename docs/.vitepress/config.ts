import { defineConfig } from 'vitepress'

export default defineConfig({
  title: 'gitinspect',
  description: 'Turn any Git repo into an AI-friendly snapshot',
  lang: 'en-US',
  base: '/gitinspect/',
  cleanUrls: true,
  head: [
    ['link', { rel: 'icon', type: 'image/svg+xml', href: '/logo.svg' }],
    ['link', { rel: 'icon', type: 'image/png', sizes: '16x16', href: '/favicon-16x16.png' }],
    ['link', { rel: 'icon', type: 'image/png', sizes: '32x32', href: '/favicon-32x32.png' }],
    ['link', { rel: 'icon', type: 'image/png', sizes: '48x48', href: '/favicon-48x48.png' }],
    ['link', { rel: 'icon', type: 'image/x-icon', href: '/favicon.ico' }],
    ['link', { rel: 'apple-touch-icon', sizes: '180x180', href: '/apple-touch-icon.png' }],
    ['link', { rel: 'manifest', href: '/site.webmanifest' }],
  ],
  themeConfig: {
    logo: '/logo.svg',
    nav: [
      { text: 'Guide', link: '/guide/getting-started' },
      { text: 'Examples', link: '/examples/local-repo' },
      { text: 'MCP', link: '/guide/mcp' },
      {
        text: 'v0.1.0',
        items: [
          { text: 'Changelog', link: '/guide/changelog' },
          { text: 'GitHub', link: 'https://github.com/richie-rich90454/gitinspect' },
        ],
      },
    ],
    sidebar: [
      {
        text: 'Introduction',
        items: [
          { text: 'What is gitinspect?', link: '/' },
          { text: 'Getting Started', link: '/guide/getting-started' },
        ],
      },
      {
        text: 'Guide',
        items: [
          { text: 'Installation', link: '/guide/installation' },
          { text: 'CLI Reference', link: '/guide/cli' },
          { text: 'Output Formats', link: '/guide/output-formats' },
          { text: 'File Priority', link: '/guide/priority' },
          { text: 'Token Budget', link: '/guide/token-budget' },
          { text: 'Dependency Extraction', link: '/guide/dependencies' },
          { text: 'Caching', link: '/guide/caching' },
        ],
      },
      {
        text: 'AI Agents (MCP)',
        items: [
          { text: 'MCP Server', link: '/guide/mcp' },
          { text: 'Agent Config', link: '/guide/agent-config' },
        ],
      },
      {
        text: 'Examples',
        items: [
          { text: 'Local Repository', link: '/examples/local-repo' },
          { text: 'Remote Repository', link: '/examples/remote-repo' },
          { text: 'HTTP Server', link: '/examples/http-server' },
          { text: 'AI Agent Usage', link: '/examples/ai-agent' },
        ],
      },
      {
        text: 'Reference',
        items: [
          { text: 'Feature Overview', link: '/guide/comparison' },
        ],
      },
    ],
    search: {
      provider: 'local',
    },
    socialLinks: [
      { icon: 'github', link: 'https://github.com/richie-rich90454/gitinspect' },
    ],
    footer: {
      message: 'Released under the Apache-2.0 License.',
      copyright: 'Copyright 2026-present gitinspect contributors',
    },
  },
})
