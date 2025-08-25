# Security Policy

The Websoft9 project team takes security issues very seriously. We are committed to ensuring our software is secure and reliable, and we appreciate security researchers and users helping us maintain the security of our project.

## Supported Versions

We provide security update support for the following versions:

| Version | Support Status |
| --- | --- |
| 1.x.x | :white_check_mark: Full support |
| 0.9.x | :warning: Limited support (critical security issues only) |
| < 0.9 | :x: No longer supported |

**Explanation:**

- **Full support**: Regular security updates and patches
- **Limited support**: Patches only for critical and high-severity security vulnerabilities
- **No longer supported**: No security updates provided, upgrade to supported version recommended

## Reporting Security Vulnerabilities

### How to Report

If you discover a security vulnerability, please **do not** report it through public GitHub Issues. Instead, please contact us privately through the following methods:

#### Preferred Method: GitHub Security Advisories

1. Visit our [GitHub Security Advisories](https://github.com/websoft9/websoft9/security/advisories) page
2. Click the "Report a vulnerability" button
3. Fill out the detailed vulnerability report form

#### Alternative Method: Email Report

Send email to: **<security@websoft9.com>**

Email subject format: `[SECURITY] Brief description of vulnerability`

### Report Content

Please include the following information in your security report:

#### Required Information

- **Vulnerability Type**: Such as SQL injection, XSS, privilege escalation, etc.
- **Affected Component**: Specific module or functionality affected
- **Vulnerability Description**: Detailed description of the security issue
- **Reproduction Steps**: Clear steps explaining how to trigger the vulnerability
- **Impact Assessment**: Potential security risks and scope of impact
- **Environment Information**: Operating system, version, configuration, etc.

#### Optional Information

- **Proof of Concept**: Code or screenshots demonstrating the vulnerability (if safe)
- **Suggested Fix**: Potential fix solutions you think might work
- **References**: Related CVEs, documentation, or research materials

#### Report Template

```markdown
## Vulnerability Overview
Brief description of the discovered security vulnerability

## Vulnerability Details
- **Vulnerability Type**: [e.g., SQL injection, XSS, privilege escalation, etc.]
- **Severity**: [Low/Medium/High/Critical]
- **Affected Component**: [Specific module or file]
- **Affected Versions**: [Range of affected versions]

## Reproduction Steps
1. First step
2. Second step
3. ...

## Impact Assessment
Describe the potential security impact of the vulnerability

## Environment Information
- Operating System:
- Websoft9 Version:
- Other relevant configurations:

## Additional Information
Other information helpful for understanding and fixing the vulnerability
```

### Response Time Commitment

We commit to responding to security reports within the following timeframes:

| Severity | Initial Response | Status Updates | Fix Target |
|----------|------------------|----------------|------------|
| Critical | 24 hours | Every 48 hours | 7 days |
| High | 48 hours | Weekly | 30 days |
| Medium | 5 business days | Bi-weekly | 90 days |
| Low | 10 business days | Monthly | Next major version |

## Security Vulnerability Severity Classification

We use the CVSS 3.1 standard to classify security vulnerabilities:

### Critical - CVSS 9.0-10.0

- Remote code execution vulnerabilities
- Complete system privilege escalation
- Large-scale data breaches
- Vulnerabilities affecting core authentication mechanisms

**Examples:**

- Unauthenticated remote code execution
- SQL injection leading to complete database access
- Authentication bypass vulnerabilities

### High - CVSS 7.0-8.9

- Privilege escalation vulnerabilities
- Sensitive data disclosure
- Denial of service attacks
- Security bypass of important functionality

**Examples:**

- Local privilege escalation
- Cross-site scripting (XSS) attacks
- Sensitive file reading

### Medium - CVSS 4.0-6.9

- Information disclosure
- Lower impact permission issues
- Configuration-related security issues

**Examples:**

- Information disclosure vulnerabilities
- CSRF attacks
- Use of weak encryption algorithms

### Low - CVSS 0.1-3.9

- Minor information disclosure
- Vulnerabilities requiring special conditions
- Best practice related issues

**Examples:**

- Version information disclosure
- Insecure default configurations
- Missing security headers

## Security Update Process

### Handling Process After Vulnerability Confirmation

1. **Vulnerability Verification** (1-3 days)
   - Technical team verifies the authenticity and impact scope of the vulnerability
   - Assess the severity and priority of the vulnerability

2. **Impact Analysis** (1-2 days)
   - Analyze affected versions and components
   - Evaluate the complexity and risk of the fix

3. **Fix Development** (Based on severity)
   - Develop security patches
   - Conduct internal testing and verification

4. **Security Advisory Preparation**
   - Prepare security advisories and release notes
   - Coordinate release timing

5. **Patch Release**
   - Release security updates
   - Publish security advisories
   - Notify users to upgrade

### Release Channels

Security updates will be released through the following channels:

- **GitHub Releases**: Primary release channel
- **GitHub Security Advisories**: Security announcements
- **Project Website**: Security notifications
- **Mailing List**: Subscriber notifications
- **Social Media**: Important security update notifications

## Security Best Practices

### User Security Recommendations

#### Deployment Security

- **Use Latest Version**: Always use the latest stable version
- **Regular Updates**: Apply security patches and updates promptly
- **Security Configuration**: Follow security configuration guidelines
- **Network Isolation**: Deploy in trusted network environments
- **Access Control**: Implement principle of least privilege

#### Configuration Security

```yaml
# Security configuration example
security:
  # Enable HTTPS
  tls:
    enabled: true
    cert_file: "/path/to/cert.pem"
    key_file: "/path/to/key.pem"

  # JWT configuration
  jwt:
    secret: "${JWT_SECRET}"  # Use environment variables
    expiration: "24h"

  # Password policy
  password:
    min_length: 8
    require_uppercase: true
    require_lowercase: true
    require_numbers: true
    require_special: true

  # Session security
  session:
    timeout: "30m"
    secure_cookie: true
    http_only: true
    same_site: "strict"
```

#### Monitoring and Auditing

- **Enable Audit Logs**: Record all security-related operations
- **Monitor Abnormal Activities**: Set up security monitoring and alerts
- **Regular Security Scans**: Use security scanning tools to check for vulnerabilities
- **Backup Strategy**: Regularly backup important data

### Development Security Guidelines

#### Secure Coding Practices

**Input Validation**

```go
// Correct input validation example
func (h *UserHandler) CreateUser(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
        return
    }

    // Validate and sanitize input
    if err := h.validator.Validate(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // HTML escape to prevent XSS
    req.Username = html.EscapeString(strings.TrimSpace(req.Username))
    req.Email = strings.ToLower(strings.TrimSpace(req.Email))

    // Handle business logic...
}
```

**SQL Injection Protection**

```go
// Use parameterized queries
func (r *userRepository) GetByUsername(ctx context.Context, username string) (*User, error) {
    var user User
    err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
    if err != nil {
        return nil, errors.Wrap(err, "failed to get user by username")
    }
    return &user, nil
}
```

**Sensitive Data Handling**

```go
// Password encryption storage
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", errors.Wrap(err, "failed to hash password")
    }
    return string(bytes), nil
}

// Use environment variables for sensitive configuration
type Config struct {
    DBPassword string `env:"DB_PASSWORD,required"`
    JWTSecret  string `env:"JWT_SECRET,required"`
    APIKey     string `env:"API_KEY,required"`
}
```

#### Security Testing

**Security Testing Checklist**

- [ ] Input validation testing
- [ ] Authentication and authorization testing
- [ ] SQL injection testing
- [ ] XSS attack testing
- [ ] CSRF attack testing
- [ ] Sensitive data leakage testing
- [ ] Privilege escalation testing
- [ ] Session management testing

**Automated Security Scanning**

```yaml
# GitHub Actions security scanning
name: Security Scan

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main, develop ]
  schedule:
    - cron: '0 2 * * *'

jobs:
  security-scan:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3

    - name: Run Gosec Security Scanner
      uses: securecodewarrior/github-action-gosec@master
      with:
        args: '-fmt sarif -out gosec-results.sarif ./...'

    - name: Run Trivy vulnerability scanner
      uses: aquasecurity/trivy-action@master
      with:
        scan-type: 'fs'
        scan-ref: '.'
        format: 'sarif'
        output: 'trivy-results.sarif'

    - name: Upload scan results to GitHub Security tab
      uses: github/codeql-action/upload-sarif@v2
      with:
        sarif_file: 'gosec-results.sarif'
```

## Security Contact Information

### Security Team

- **Security Lead**: [Name] <security-lead@websoft9.com>
- **Technical Lead**: [Name] <tech-lead@websoft9.com>
- **Security Email**: <security@websoft9.com>

### PGP Public Key

For encrypted communication, please use our PGP public key:

```text
-----BEGIN PGP PUBLIC KEY BLOCK-----
[PGP public key content]
-----END PGP PUBLIC KEY BLOCK-----
```

**Key Fingerprint**: `XXXX XXXX XXXX XXXX XXXX XXXX XXXX XXXX XXXX XXXX`

## Acknowledgments

### Security Researcher Acknowledgments

We thank the following security researchers for their contributions to the security of the Websoft9 project:

- [Researcher Name] - Discovered and reported [Vulnerability Type]
- [Researcher Name] - Discovered and reported [Vulnerability Type]

*If you wish to remain anonymous in this list, please indicate so when reporting.*

### Responsible Disclosure

We commit to:

- **Timely Response**: Respond to security reports within committed timeframes
- **Transparent Communication**: Regularly update on fix progress
- **Public Acknowledgment**: Thank reporters in security advisories (unless anonymity requested)
- **No Legal Action**: Not pursue legal action against good-faith security research

We expect security researchers to:

- **Report Privately**: Not publicly disclose unpatched vulnerabilities
- **Allow Time**: Provide reasonable time for fixes
- **Avoid Damage**: Not exploit vulnerabilities to cause harm
- **Follow Laws**: Conduct research within legal boundaries

## Security Resources

### Related Documentation

- [Development Standards](./webox/docs/开发规范.md) - Contains detailed secure coding standards
- [Contributing Guide](./CONTRIBUTING.md) - Code contribution and review process

### Security Tools

- **Code Scanning**: Gosec, CodeQL
- **Dependency Scanning**: Trivy, Snyk
- **Container Scanning**: Docker Scout, Clair
- **Penetration Testing**: OWASP ZAP, Burp Suite

### Security Standards

We follow these security standards and frameworks:

- **OWASP Top 10** - Web application security risks
- **NIST Cybersecurity Framework** - Cybersecurity framework
- **ISO 27001** - Information security management system
- **CWE/SANS Top 25** - Most dangerous software errors

For any security-related questions, please contact: <security@websoft9.com>
