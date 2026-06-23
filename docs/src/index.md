---
# https://vitepress.dev/reference/default-theme-home-page
layout: home

hero:
  name: "binmate"
  text: "A cross-platform binary manager"
  tagline: Install, manage, and switch between versions of command-line tools from GitHub releases - straight from your terminal.
  actions:
    - theme: brand
      text: Get Started
      link: /guide/introduction
    - theme: alt
      text: Installation
      link: /guide/installation
    - theme: alt
      text: View on GitHub
      link: https://github.com/cturner8/binmate

features:
  - icon: 📦
    title: Interactive TUI
    details: Browse, install, and manage your binaries through a friendly Terminal User Interface built with Bubble Tea.
  - icon: ⌨️
    title: CLI Commands
    details: Automate binary management with a full set of command-line commands for scripting and CI pipelines.
  - icon: 🔀
    title: Version Management
    details: Install multiple versions side by side and switch the active version with a single command.
  - icon: 🐙
    title: GitHub Integration
    details: Automatically fetch and install releases directly from GitHub repositories, including pre-releases.
  - icon: 🗄️
    title: Database Tracking
    details: A local SQLite database keeps track of every installation, version, and download.
  - icon: 🔒
    title: Checksum Verification
    details: Downloads are verified against published checksums to ensure the integrity of every binary.
---
