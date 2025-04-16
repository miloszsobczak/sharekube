# Security Policy

## Supported Versions

We currently provide security updates for the following versions of ShareKube:

| Version | Supported          |
| ------- | ------------------ |
| 1.0.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

We take the security of ShareKube seriously. If you believe you've found a security vulnerability, please follow these steps:

1. **Do not disclose the vulnerability publicly**
2. **Email the maintainer directly** at milosz.sobczak@gmail.com
3. **Include details in your report**:
   - Type of issue
   - Full paths of source file(s) related to the issue
   - Location of the affected source code
   - Any special configuration required to reproduce the issue
   - Step-by-step instructions to reproduce the issue
   - Proof-of-concept or exploit code (if possible)
   - Impact of the issue, including how an attacker might exploit it

### What to expect

- We will acknowledge receipt of your vulnerability report within 3 business days
- We will provide a more detailed response within 7 days, indicating next steps
- We will keep you informed of our progress
- We will treat your report with strict confidentiality, and not pass on your personal details to third parties without your permission

## Public Disclosure Timing

We believe in responsible disclosure. We ask that you do not share information about the vulnerability with others until:

- We've had sufficient time to address the vulnerability
- A patch has been made available to users
- A coordinated disclosure date has been agreed upon

## Security-Related Configuration

ShareKube includes security-conscious defaults, but you may wish to review the following to enhance your security posture:

- Review RBAC settings when deploying in a shared cluster
- Follow the principle of least privilege when setting up service accounts
- Consider network policies to restrict traffic between namespaces where ShareKube operates

## Related Documents

- [README.md](README.md): Project overview and main documentation
- [CONTRIBUTING.md](CONTRIBUTING.md): Guidelines for contributing to ShareKube
- [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md): Our community guidelines 