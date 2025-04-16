# ShareKube GitHub Workflows

This document describes GitHub Actions workflows for automating various tasks in the ShareKube project.

## Workflows

### Documentation Deployment (`docs-deploy.yml`)

This workflow automatically builds and deploys the Docusaurus documentation site to GitHub Pages.

#### Trigger
The workflow is triggered when:
- Changes are pushed to the `main` branch that affect files in the `apps/docs` directory
- Manually triggered through the GitHub Actions UI (workflow_dispatch event)

#### What it does
1. Builds the Docusaurus site using Node.js and Yarn
2. Deploys the built site to GitHub Pages

#### Requirements
For this workflow to function properly:
1. GitHub Pages must be enabled for the repository
2. The permissions for GitHub Actions must be set to "Read and write permissions" in the repository settings

## Setting up GitHub Pages with Custom Domain

To complete the setup:

1. Go to your repository on GitHub
2. Navigate to Settings > Pages
3. Under "Source", select "GitHub Actions"
4. Under "Custom domain", enter `docs.sharekube.dev` and click "Save"
5. GitHub will verify ownership of the domain

### DNS Configuration for Custom Domain

Add the following DNS records to your domain registrar for the `sharekube.dev` domain:

For an apex domain:
```
Type    Name             Value
A       docs             185.199.108.153
A       docs             185.199.109.153
A       docs             185.199.110.153
A       docs             185.199.111.153
AAAA    docs             2606:50c0:8000::153
AAAA    docs             2606:50c0:8001::153
AAAA    docs             2606:50c0:8002::153
AAAA    docs             2606:50c0:8003::153
```

Or, alternatively, you can use a CNAME record:
```
Type    Name             Value
CNAME   docs             miloszsobczak.github.io.
```

After setting up the DNS records and GitHub configuration, your documentation will be available at: `https://docs.sharekube.dev/` 