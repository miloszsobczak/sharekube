# Setting Up HTTPS for Local Development

This guide will help you set up HTTPS for local development of the ShareKube documentation site.

## Prerequisites

- [mkcert](https://github.com/FiloSottile/mkcert) - A simple tool for making locally-trusted development certificates

### Installing mkcert

**macOS (using Homebrew):**
```bash
brew install mkcert
brew install nss # if you use Firefox
```

**Windows (using Chocolatey):**
```bash
choco install mkcert
```

**Linux:**
Follow the instructions at https://github.com/FiloSottile/mkcert#linux

## Setup Process

### 1. Automated Setup

We've provided a script to automate the setup process:

```bash
# From the project root
./scripts/setup-local-https.sh
```

This script will:
1. Check if mkcert is installed
2. Create a certificate directory
3. Install the mkcert root CA (if not already done)
4. Generate certificates for docs.local.sharekube.dev
5. Create the necessary symlinks

### 2. Manual Setup

If you prefer to do it manually, follow these steps:

1. Install the mkcert CA:
   ```bash
   mkcert -install
   ```

2. Create a certificate directory:
   ```bash
   mkdir -p apps/docs/certs
   cd apps/docs/certs
   ```

3. Generate certificates:
   ```bash
   mkcert docs.local.sharekube.dev "*.local.sharekube.dev"
   ```

4. Create symlinks (optional, but makes configuration easier):
   ```bash
   ln -sf "docs.local.sharekube.dev+1-key.pem" key.pem
   ln -sf "docs.local.sharekube.dev+1.pem" cert.pem
   ```

### 3. Configure /etc/hosts

Add the following line to your `/etc/hosts` file:

```
127.0.0.1 docs.local.sharekube.dev
```

## Starting the Development Server

```bash
cd apps/docs
yarn dev
```

This will start the server with HTTPS enabled. You should now be able to access the site at:

https://docs.local.sharekube.dev:3000

## Troubleshooting

### ERR_SSL_PROTOCOL_ERROR

If you're seeing an SSL protocol error:

1. Ensure mkcert is properly installed and the CA has been installed (`mkcert -install`)
2. Check that the certificate files are in the correct location (apps/docs/certs/)
3. Verify your /etc/hosts file has the correct entry
4. Try restarting your browser or clearing SSL state
5. Check that the server is running with the correct HTTPS configuration

### Certificate Issues in Chrome

If Chrome doesn't recognize the certificate:
1. Type `chrome://flags/#allow-insecure-localhost` in your address bar
2. Enable the "Allow invalid certificates for resources loaded from localhost" option
3. Restart Chrome

## Why We Use HTTPS in Development

Using HTTPS in development:
1. More closely mimics production environments
2. Tests security features that require HTTPS
3. Avoids mixed content warnings when integrating external services
4. Allows using features that require secure contexts (like Service Workers) 