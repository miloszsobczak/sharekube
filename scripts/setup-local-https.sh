#!/bin/bash
set -e

# Colors for better output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}Setting up local HTTPS certificates for ShareKube development...${NC}"

# Check if mkcert is installed
if ! command -v mkcert &> /dev/null; then
    echo -e "${RED}mkcert is not installed. Please install it first:${NC}"
    echo -e "${YELLOW}macOS: brew install mkcert${NC}"
    echo -e "${YELLOW}Linux: https://github.com/FiloSottile/mkcert#linux${NC}"
    echo -e "${YELLOW}Windows: https://github.com/FiloSottile/mkcert#windows${NC}"
    exit 1
fi

# Create certificates directory if it doesn't exist
CERT_DIR="$(pwd)/apps/docs/certs"
mkdir -p "$CERT_DIR"

# Install mkcert root CA if not already done
echo -e "${GREEN}Installing mkcert root CA (you may be prompted for your password)...${NC}"
mkcert -install

# Generate certificates for local domains
echo -e "${GREEN}Generating certificates for local.sharekube.dev and docs.local.sharekube.dev...${NC}"
cd "$CERT_DIR"
mkcert local.sharekube.dev docs.local.sharekube.dev "*.local.sharekube.dev"

# Create key.pem and cert.pem symlinks (conventional names used by many Node.js servers)
echo -e "${GREEN}Creating key.pem and cert.pem symlinks...${NC}"
ln -sf "local.sharekube.dev+2-key.pem" key.pem
ln -sf "local.sharekube.dev+2.pem" cert.pem

# Remind about /etc/hosts
echo -e "${YELLOW}Don't forget to add the following lines to your /etc/hosts file:${NC}"
echo -e "${YELLOW}127.0.0.1 local.sharekube.dev${NC}"
echo -e "${YELLOW}127.0.0.1 docs.local.sharekube.dev${NC}"

echo -e "${GREEN}✅ Certificates generated successfully in ${CERT_DIR}${NC}"
echo -e "${GREEN}To start the dev server with HTTPS, run:${NC}"
echo -e "${YELLOW}cd apps/docs && yarn dev${NC}"
echo -e "${GREEN}The site will be available at:${NC}"
echo -e "${YELLOW}https://local.sharekube.dev:3000${NC} - Main site"
echo -e "${YELLOW}https://local.sharekube.dev:3000/docs${NC} - Documentation" 