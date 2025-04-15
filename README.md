# ShareKube

ShareKube is a Kubernetes extension that allows users to create temporary preview environments by copying explicitly defined resources from one namespace to another within the same cluster.

## Project Structure

This is a monorepo containing all ShareKube projects:

- `apps/` - Applications
  - `docs/` - Documentation site (Docusaurus)
- `packages/` - Shared packages
  - `shared-assets/` - Shared assets like logos and icons

## Getting Started

### Prerequisites

- Node.js 16+
- Yarn
- [mkcert](https://github.com/FiloSottile/mkcert) (for local HTTPS development)

### Installation

```bash
# Clone the repository
git clone https://github.com/miloszsobczak/sharekube.git
cd sharekube

# Install dependencies
yarn

# Set up local HTTPS certificates (recommended)
./scripts/setup-local-https.sh
```

### Development

#### Documentation

The documentation site is hosted at docs.sharekube.dev and is a pure documentation site (no separate landing page).

To run the documentation site locally:

```bash
# Add docs.local.sharekube.dev to your hosts file
# See hosts.txt for the required entry

# Start the documentation site
cd apps/docs
yarn dev

# The site will be available at https://docs.local.sharekube.dev:3000
```

#### Local HTTPS Development

For security reasons and to prevent ERR_SSL_PROTOCOL_ERROR issues, the development server uses HTTPS. 

See [LOCAL_HTTPS_SETUP.md](LOCAL_HTTPS_SETUP.md) for detailed instructions on setting up HTTPS for local development.

## License

This project is licensed under the Apache 2.0 License - see the LICENSE file for details. 