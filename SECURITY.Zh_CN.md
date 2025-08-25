# Security Policy

## 安全政策

Websoft9 项目团队非常重视安全问题。我们致力于确保我们的软件安全可靠，并感谢安全研究人员和用户帮助我们保持项目的安全性。

## 支持的版本

我们为以下版本提供安全更新支持：

| 版本 | 支持状态 |
| --- | --- |
| 1.x.x | :white_check_mark: 完全支持 |
| 0.9.x | :warning: 有限支持（仅关键安全问题） |
| < 0.9 | :x: 不再支持 |

**说明：**

- **完全支持**：定期安全更新和补丁
- **有限支持**：仅针对关键和高危安全漏洞提供补丁
- **不再支持**：不再提供安全更新，建议升级到支持的版本

## 报告安全漏洞

### 如何报告

如果您发现了安全漏洞，请**不要**通过公开的 GitHub Issues 报告。相反，请通过以下方式私下联系我们：

#### 首选方式：GitHub Security Advisories

1. 访问我们的 [GitHub Security Advisories](https://github.com/websoft9/websoft9/security/advisories) 页面
2. 点击 "Report a vulnerability" 按钮
3. 填写详细的漏洞报告表单

#### 备选方式：邮件报告

发送邮件至：**<security@websoft9.com>**

邮件主题格式：`[SECURITY] 简要描述漏洞`

### 报告内容

请在您的安全报告中包含以下信息：

#### 必需信息

- **漏洞类型**：如 SQL 注入、XSS、权限提升等
- **影响组件**：受影响的具体模块或功能
- **漏洞描述**：详细描述安全问题
- **复现步骤**：清晰的步骤说明如何触发漏洞
- **影响评估**：潜在的安全风险和影响范围
- **环境信息**：操作系统、版本、配置等

#### 可选信息

- **概念验证**：演示漏洞的代码或截图（如果安全）
- **建议修复**：您认为可能的修复方案
- **参考资料**：相关的 CVE、文档或研究资料

#### 报告模板

```markdown
## 漏洞概述
简要描述发现的安全漏洞

## 漏洞详情
- **漏洞类型**: [如：SQL注入、XSS、权限提升等]
- **严重程度**: [低/中/高/严重]
- **影响组件**: [具体的模块或文件]
- **影响版本**: [受影响的版本范围]

## 复现步骤
1. 第一步
2. 第二步
3. ...

## 影响评估
描述漏洞可能造成的安全影响

## 环境信息
- 操作系统：
- Websoft9 版本：
- 其他相关配置：

## 附加信息
其他有助于理解和修复漏洞的信息
```

### 响应时间承诺

我们承诺在以下时间内响应安全报告：

| 严重程度 | 初始响应 | 状态更新 | 修复目标 |
|----------|----------|----------|----------|
| 严重 (Critical) | 24 小时 | 每 48 小时 | 7 天 |
| 高 (High) | 48 小时 | 每周 | 30 天 |
| 中 (Medium) | 5 个工作日 | 每两周 | 90 天 |
| 低 (Low) | 10 个工作日 | 每月 | 下一个主要版本 |

## 安全漏洞严重程度分级

我们使用 CVSS 3.1 标准对安全漏洞进行分级：

### 严重 (Critical) - CVSS 9.0-10.0

- 远程代码执行漏洞
- 完全系统权限提升
- 大规模数据泄露
- 影响核心认证机制的漏洞

**示例**：

- 未经认证的远程代码执行
- SQL 注入导致完整数据库访问
- 认证绕过漏洞

### 高 (High) - CVSS 7.0-8.9

- 权限提升漏洞
- 敏感数据泄露
- 拒绝服务攻击
- 重要功能的安全绕过

**示例**：

- 本地权限提升
- 跨站脚本攻击 (XSS)
- 敏感文件读取

### 中 (Medium) - CVSS 4.0-6.9

- 信息泄露
- 较低影响的权限问题
- 配置相关的安全问题

**示例**：

- 信息泄露漏洞
- CSRF 攻击
- 弱加密算法使用

### 低 (Low) - CVSS 0.1-3.9

- 轻微的信息泄露
- 需要特殊条件的漏洞
- 最佳实践相关问题

**示例**：

- 版本信息泄露
- 不安全的默认配置
- 缺少安全头

## 安全更新流程

### 漏洞确认后的处理流程

1. **漏洞验证** (1-3 天)
   - 技术团队验证漏洞的真实性和影响范围
   - 评估漏洞的严重程度和优先级

2. **影响分析** (1-2 天)
   - 分析受影响的版本和组件
   - 评估修复的复杂度和风险

3. **修复开发** (根据严重程度)
   - 开发安全补丁
   - 进行内部测试和验证

4. **安全公告准备**
   - 准备安全公告和发布说明
   - 协调发布时间

5. **补丁发布**
   - 发布安全更新
   - 发布安全公告
   - 通知用户升级

### 发布渠道

安全更新将通过以下渠道发布：

- **GitHub Releases**：主要发布渠道
- **GitHub Security Advisories**：安全公告
- **项目官网**：安全通知
- **邮件列表**：订阅用户通知
- **社交媒体**：重要安全更新通知

## 安全最佳实践

### 用户安全建议

#### 部署安全

- **使用最新版本**：始终使用最新的稳定版本
- **定期更新**：及时应用安全补丁和更新
- **安全配置**：遵循安全配置指南
- **网络隔离**：在受信任的网络环境中部署
- **访问控制**：实施最小权限原则

#### 配置安全

```yaml
# 安全配置示例
security:
  # 启用 HTTPS
  tls:
    enabled: true
    cert_file: "/path/to/cert.pem"
    key_file: "/path/to/key.pem"
  
  # JWT 配置
  jwt:
    secret: "${JWT_SECRET}"  # 使用环境变量
    expiration: "24h"
    
  # 密码策略
  password:
    min_length: 8
    require_uppercase: true
    require_lowercase: true
    require_numbers: true
    require_special: true
    
  # 会话安全
  session:
    timeout: "30m"
    secure_cookie: true
    http_only: true
    same_site: "strict"
```

#### 监控和审计

- **启用审计日志**：记录所有安全相关操作
- **监控异常活动**：设置安全监控和告警
- **定期安全扫描**：使用安全扫描工具检查漏洞
- **备份策略**：定期备份重要数据

### 开发安全指南

#### 安全编码实践

**输入验证**

```go
// 正确的输入验证示例
func (h *UserHandler) CreateUser(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
        return
    }
    
    // 验证和清理输入
    if err := h.validator.Validate(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    // HTML 转义防止 XSS
    req.Username = html.EscapeString(strings.TrimSpace(req.Username))
    req.Email = strings.ToLower(strings.TrimSpace(req.Email))
    
    // 处理业务逻辑...
}
```

**SQL 注入防护**

```go
// 使用参数化查询
func (r *userRepository) GetByUsername(ctx context.Context, username string) (*User, error) {
    var user User
    err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
    if err != nil {
        return nil, errors.Wrap(err, "failed to get user by username")
    }
    return &user, nil
}
```

**敏感数据处理**

```go
// 密码加密存储
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", errors.Wrap(err, "failed to hash password")
    }
    return string(bytes), nil
}

// 敏感配置使用环境变量
type Config struct {
    DBPassword string `env:"DB_PASSWORD,required"`
    JWTSecret  string `env:"JWT_SECRET,required"`
    APIKey     string `env:"API_KEY,required"`
}
```

#### 安全测试

**安全测试检查清单**

- [ ] 输入验证测试
- [ ] 认证和授权测试
- [ ] SQL 注入测试
- [ ] XSS 攻击测试
- [ ] CSRF 攻击测试
- [ ] 敏感数据泄露测试
- [ ] 权限提升测试
- [ ] 会话管理测试

**自动化安全扫描**

```yaml
# GitHub Actions 安全扫描
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

## 安全联系信息

### 安全团队

- **安全负责人**：[姓名] <security-lead@websoft9.com>
- **技术负责人**：[姓名] <tech-lead@websoft9.com>
- **安全邮箱**：<security@websoft9.com>

### PGP 公钥

如需加密通信，请使用我们的 PGP 公钥：

```text
-----BEGIN PGP PUBLIC KEY BLOCK-----
[PGP 公钥内容]
-----END PGP PUBLIC KEY BLOCK-----
```

**密钥指纹**：`XXXX XXXX XXXX XXXX XXXX XXXX XXXX XXXX XXXX XXXX`

## 致谢

### 安全研究人员致谢

我们感谢以下安全研究人员对 Websoft9 项目安全性的贡献：

- [研究人员姓名] - 发现并报告了 [漏洞类型]
- [研究人员姓名] - 发现并报告了 [漏洞类型]

*如果您希望在此列表中保持匿名，请在报告时说明。*

### 负责任的披露

我们承诺：

- **及时响应**：在承诺时间内响应安全报告
- **透明沟通**：定期更新修复进展
- **公开致谢**：在安全公告中感谢报告者（除非要求匿名）
- **不追究责任**：对善意的安全研究不采取法律行动

我们期望安全研究人员：

- **私下报告**：不公开披露未修复的漏洞
- **给予时间**：允许合理的修复时间
- **避免破坏**：不利用漏洞造成损害
- **遵守法律**：在法律允许的范围内进行研究

## 安全资源

### 相关文档

- [开发规范](./webox/docs/开发规范.md) - 包含详细的安全编码规范
- [贡献指南](./CONTRIBUTOR.md) - 代码贡献和审查流程

### 安全工具

- **代码扫描**：Gosec, CodeQL
- **依赖扫描**：Trivy, Snyk
- **容器扫描**：Docker Scout, Clair
- **渗透测试**：OWASP ZAP, Burp Suite

### 安全标准

我们遵循以下安全标准和框架：

- **OWASP Top 10** - Web 应用安全风险
- **NIST Cybersecurity Framework** - 网络安全框架
- **ISO 27001** - 信息安全管理体系
- **CWE/SANS Top 25** - 最危险的软件错误

如有任何安全相关问题，请联系：<security@websoft9.com>
