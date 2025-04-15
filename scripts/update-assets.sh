#!/bin/bash
set -e

# Colors for better output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}Updating shared assets across all applications...${NC}"

# Source directory for shared assets
ASSETS_DIR="$(pwd)/packages/shared-assets"

if [ ! -d "$ASSETS_DIR" ]; then
  echo -e "${RED}Error: Shared assets directory not found at ${ASSETS_DIR}${NC}"
  exit 1
fi

# Update Docusaurus assets
DOCS_STATIC_DIR="$(pwd)/apps/docs/static/img"
if [ -d "$DOCS_STATIC_DIR" ]; then
  echo -e "${GREEN}Updating Docusaurus assets...${NC}"
  cp -v "$ASSETS_DIR/logo/logo.svg" "$DOCS_STATIC_DIR/"
  cp -v "$ASSETS_DIR/logo/logo_square.svg" "$DOCS_STATIC_DIR/"
  cp -v "$ASSETS_DIR/logo/favicon.ico" "$DOCS_STATIC_DIR/"
  echo -e "${GREEN}✅ Docusaurus assets updated.${NC}"
else
  echo -e "${YELLOW}Warning: Docusaurus static directory not found at ${DOCS_STATIC_DIR}${NC}"
fi

# Add other apps here as needed
# Example:
# UI_ASSETS_DIR="$(pwd)/apps/ui/public"
# if [ -d "$UI_ASSETS_DIR" ]; then
#   echo -e "${GREEN}Updating UI assets...${NC}"
#   cp -v "$ASSETS_DIR/logo/logo.svg" "$UI_ASSETS_DIR/"
#   cp -v "$ASSETS_DIR/logo/logo_square.svg" "$UI_ASSETS_DIR/"
#   cp -v "$ASSETS_DIR/logo/favicon.ico" "$UI_ASSETS_DIR/favicon.ico"
#   echo -e "${GREEN}✅ UI assets updated.${NC}"
# fi

echo -e "${GREEN}✅ All assets have been updated successfully!${NC}"
echo -e "${YELLOW}Remember to commit these changes if needed.${NC}" 