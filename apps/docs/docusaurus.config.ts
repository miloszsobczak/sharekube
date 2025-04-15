import {themes as prismThemes} from 'prism-react-renderer';
import type {Config} from '@docusaurus/types';
import type * as Preset from '@docusaurus/preset-classic';

// This runs in Node.js - Don't use client-side code here (browser APIs, JSX...)

const config: Config = {
  title: 'ShareKube',
  tagline: 'Create temporary preview environments in Kubernetes',
  favicon: 'img/favicon.ico',

  // Set the production url of your site here
  url: 'http://docs.sharekube.dev',
  // Set the /<baseUrl>/ pathname under which your site is served
  baseUrl: '/',

  // GitHub pages deployment config.
  // If you aren't using GitHub pages, you don't need these.
  organizationName: 'miloszsobczak', // Usually your GitHub org/user name.
  projectName: 'sharekube', // Usually your repo name.

  onBrokenLinks: 'warn',
  onBrokenMarkdownLinks: 'warn',

  // Enable Mermaid diagrams
  markdown: {
    mermaid: true,
  },

  // Add Mermaid theme
  themes: ['@docusaurus/theme-mermaid'],

  // Even if you don't use internationalization, you can use this field to set
  // useful metadata like html lang. For example, if your site is Chinese, you
  // may want to replace "en" with "zh-Hans".
  i18n: {
    defaultLocale: 'en',
    locales: ['en'],
  },

  presets: [
    [
      'classic',
      {
        docs: {
          sidebarPath: './sidebars.ts',
          // No edit links
          routeBasePath: '/', // Set docs as the root
        },
        blog: false,
        theme: {
          customCss: './src/css/custom.css',
        },
      } satisfies Preset.Options,
    ],
  ],

  themeConfig: {
    // Mermaid specific configuration
    mermaid: {
      theme: { light: 'default', dark: 'dark' },
      options: {
        maxTextSize: 100000,
        securityLevel: 'loose',
        flowchart: {
          curve: 'basis',
          useMaxWidth: false,
        },
        fontSize: 16,
      },
    },
    // Replace with your project's social card
    image: 'img/sharekube-social-card.jpg',
    navbar: {
      title: 'ShareKube',
      logo: {
        alt: 'ShareKube Logo',
        src: 'img/logo_square.svg',
      },
      items: [
        {
          href: 'https://github.com/miloszsobczak/sharekube',
          label: 'GitHub',
          position: 'right',
        },
      ],
    },
    footer: {
      style: 'light',
      logo: {
        alt: 'ShareKube Logo',
        src: 'img/logo.svg',
        href: '/',
        width: 200,
      },
      copyright: `Copyright © ${new Date().getFullYear()} ShareKube Project. Built with Docusaurus.`,
    },
    prism: {
      theme: prismThemes.github,
      darkTheme: prismThemes.dracula,
      additionalLanguages: ['bash', 'yaml'],
    },
  } satisfies Preset.ThemeConfig,
};

export default config;
