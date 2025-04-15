# ShareKube Shared Assets

This package contains shared assets used across ShareKube projects.

## Usage

### Installation

```bash
# From within your app
yarn add @sharekube/shared-assets
```

### Using the Logo in JavaScript/TypeScript

```js
// When using webpack or similar bundlers that support importing assets
import logoSvg from '@sharekube/shared-assets/logo/logo.svg';

// Use in your component
<img src={logoSvg} alt="ShareKube Logo" />
```

### Using the Logo in CSS/SCSS

```css
.logo {
  background-image: url('~@sharekube/shared-assets/logo/logo.svg');
}
```

### Using in HTML

```html
<img src="/node_modules/@sharekube/shared-assets/logo/logo.svg" alt="ShareKube Logo" />
```

### Using in Docusaurus

For Docusaurus projects, you may need to copy the assets to your static directory as part of your build process:

```js
// In your build script or webpack config
const fs = require('fs-extra');
const path = require('path');

// Copy logo to static directory
fs.copySync(
  path.resolve(__dirname, 'node_modules/@sharekube/shared-assets/logo'),
  path.resolve(__dirname, 'static/img')
);
```

## Available Assets

### Logo

- `logo/logo.svg` - Main logo in SVG format
- `logo/favicon.ico` - Favicon for web browsers

## License

This package is licensed under the Apache 2.0 License. 