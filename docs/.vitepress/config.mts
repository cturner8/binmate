import { defineConfig } from "vitepress";

// https://vitepress.dev/reference/site-config
export default defineConfig({
  title: "binmate",
  description:
    "A cross-platform binary manager for installing and switching between versions of command-line tools from GitHub releases.",
  srcDir: "./src",
  cleanUrls: true,
  lastUpdated: true,

  head: [
    ["meta", { name: "theme-color", content: "#F97316" }],
    ["meta", { name: "og:type", content: "website" }],
    ["meta", { name: "og:title", content: "binmate" }],
    [
      "meta",
      {
        name: "og:description",
        content: "A cross-platform binary manager for command-line tools.",
      },
    ],
  ],

  vite: {
    server: {
      port: 31560,
      host: "127.0.0.1",
    },
  },

  themeConfig: {
    // https://vitepress.dev/reference/default-theme-config
    nav: [
      { text: "Home", link: "/" },
      { text: "Guide", link: "/guide/introduction" },
      { text: "Reference", link: "/reference/cli" },
      {
        text: "Releases",
        link: "https://github.com/cturner8/binmate/releases",
      },
    ],

    sidebar: {
      "/guide/": [
        {
          text: "Getting Started",
          items: [
            { text: "Introduction", link: "/guide/introduction" },
            { text: "Installation", link: "/guide/installation" },
            { text: "Usage", link: "/guide/usage" },
          ],
        },
        {
          text: "Using binmate",
          items: [
            { text: "The Interface", link: "/guide/interface" },
            { text: "Configuration", link: "/guide/configuration" },
          ],
        },
        {
          text: "Contributing",
          items: [{ text: "Releasing", link: "/contributing/releasing" }],
        },
      ],
      "/reference/": [
        {
          text: "Reference",
          items: [
            { text: "CLI Commands", link: "/reference/cli" },
            { text: "Configuration", link: "/reference/configuration" },
            { text: "Database", link: "/reference/database" },
            { text: "Architecture", link: "/reference/architecture" },
          ],
        },
      ],
      "/contributing/": [
        {
          text: "Contributing",
          items: [{ text: "Releasing", link: "/contributing/releasing" }],
        },
      ],
    },

    socialLinks: [
      { icon: "github", link: "https://github.com/cturner8/binmate" },
    ],

    editLink: {
      pattern: "https://github.com/cturner8/binmate/edit/dev/docs/src/:path",
      text: "Edit this page on GitHub",
    },

    search: {
      provider: "local",
    },

    footer: {
      message: "Released under the MIT License.",
      copyright: "Copyright © 2026 cturner8",
    },
  },
});
