---
layout: home

hero:
  name: "gitinspect"
  text: "Git repos, AI-ready"
  tagline: Turn any Git repository into a structured, token-efficient snapshot optimized for LLMs and AI agents.
  actions:
    - theme: brand
      text: Get Started
      link: /guide/getting-started
    - theme: alt
      text: MCP Integration
      link: /guide/mcp
    - theme: alt
      text: GitHub
      link: https://github.com/your-username/gitinspect

features:
  - title: Token Budget Management
    details: Estimate tokens as len/4. Sort files by priority. Stop when budget exceeded. Never waste context window space.
    icon: 📊
  - title: AI Agent Ready
    details: Built-in MCP server over stdio. Works with Claude, Cursor, Windsurf, VS Code Copilot, and any MCP-compatible agent.
    icon: 🤖
  - title: Any Git Host
    details: GitHub, GitLab, Bitbucket, self-hosted — works with any Git repository without requiring a full clone.
    icon: 🌐
  - title: Smart File Priority
    details: README first, then main files, then manifests, then build files. The most important code is always in context.
    icon: 🎯
  - title: Dependency Extraction
    details: Automatically parse go.mod, package.json, Cargo.toml, requirements.txt, and Gemfile for dependency awareness.
    icon: 📦
  - title: Multiple Output Formats
    details: JSON (default), plain text, or YAML. Perfect for piping into LLMs, saving to files, or programmatic consumption.
    icon: 📄
---
