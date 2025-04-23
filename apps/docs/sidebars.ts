import type {SidebarsConfig} from '@docusaurus/plugin-content-docs';

/**
 * ShareKube documentation sidebar configuration
 */
const sidebars: SidebarsConfig = {
  docsSidebar: [
    {
      type: 'doc',
      id: 'home',
      label: 'Home',
    },
    {
      type: 'doc',
      id: 'overview',
      label: 'Overview',
    },
    {
      type: 'doc',
      id: 'getting-started',
      label: 'Getting Started',
    },
    {
      type: 'doc',
      id: 'api-reference',
      label: 'CRD API Reference',
    },
    {
      type: 'doc',
      id: 'dynamic-permissions',
      label: 'Dynamic Permissions',
    },
    {
      type: 'doc',
      id: 'future-roadmap',
      label: 'Future Roadmap',
    },
  ],
};

export default sidebars;
