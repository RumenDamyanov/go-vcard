# Security Policy

## Supported Versions

We currently support the following versions with security updates:

| Version | Supported          |
| ------- | ------------------ |
| 1.x.x   | :white_check_mark: |

Build with **Go 1.26.10+** or **1.27.3+**. Go 1.25 and older are end of life and do not receive the October 2026 security releases. Those fixes are in the toolchain (`crypto/x509`, `net/http`, `html/template`), not in this module. The `go` line asks the `go` command for at least Go 1.26.10 when `GOTOOLCHAIN=auto`, which is the default. `GOTOOLCHAIN=local` fails the build until that toolchain is installed.

- [CVE-2026-111943](https://pkg.go.dev/vuln/GO-2026-6809) (`crypto/x509`): untrustworthy eMudhra intermediates are rejected. Fixed in Go 1.26.10 and 1.27.3 ([announcement](https://groups.google.com/g/golang-announce/c/66pojueaQek/m/tmWX2Pq9AgAJ)).
- Go 1.26.9 and 1.27.2 also fixed HTTP/2 server crashes, trailer memory use, and `html/template` context tracking ([announcement](https://groups.google.com/g/golang-announce/c/U2fTuyDJznI)).
- Gin and Fiber decompress Brotli with `github.com/andybalholm/brotli`. Versions before v1.2.7 panic on crafted input ([CVE-2026-111497](https://www.wiz.io/blog/cve-2026-111497-brotli-rce-node-js)). Adapter modules that depend on those frameworks require v1.2.7.

## Reporting a Vulnerability

If you discover a security vulnerability in go-vcard, please report it responsibly:

### Private Reporting (Recommended)

1. **Email**: Send details to [contact@rumenx.com](mailto:contact@rumenx.com)
2. **Subject**: Include "SECURITY" and "go-vcard" in the subject line
3. **Details**: Provide:
   - Description of the vulnerability
   - Steps to reproduce
   - Potential impact
   - Suggested fix (if any)

### GitHub Security Advisories

You can also use [GitHub's security advisory feature](https://github.com/rumendamyanov/go-vcard/security/advisories) to report vulnerabilities privately.

## Response Timeline

- **Initial Response**: Within 24 hours
- **Assessment**: Within 72 hours
- **Fix & Release**: As soon as possible (typically within 7 days for critical issues)

## Security Best Practices

When using go-vcard:

1. **Validate Input**: Always validate and sanitize user input before creating vCards
2. **Update Regularly**: Keep the library updated to the latest version
3. **Secure Storage**: If storing vCard data, ensure proper encryption and access controls
4. **Content Filtering**: Be cautious with user-generated content in vCard fields

## Acknowledgments

We appreciate security researchers and users who responsibly disclose vulnerabilities. Contributors will be acknowledged (with their permission) in:

- Security advisories
- Release notes
- This security policy

Thank you for helping keep go-vcard secure!
