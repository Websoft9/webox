# 审计日志功能 - AI Coding实施计划

## 1. 功能分析

```yaml
feature_name: "audit-logs"
business_context: "审计日志是平台提供的对敏感业务操作的日志记录，以便于用户进行业务审计、业务事件追踪等管理需求。记录平台用户发生的所有业务操作，按照功能模块划分日志类别，支持查询筛选、导出和统计分析功能。"
user_scenarios: 
  - "管理员查看平台用户的操作记录，进行安全审计和合规检查"
  - "通过日志类别、时间范围、用户等条件筛选查询特定操作记录"
  - "导出审计日志明细用于合规报告和事件追踪分析"
  - "查看操作统计信息，分析平台使用情况和安全风险"
functional_scope:
  includes: 
    - "审计日志的自动记录（所有业务操作）"
    - "多条件查询筛选（日志类别、操作类型、用户、时间范围等）"
    - "审计日志详情查看"
    - "日志导出功能（CSV、Excel、JSON格式）"
    - "操作统计分析和可视化展示"
    - "敏感信息脱敏处理"
  excludes: 
    - "审计日志的删除操作（仅支持查看和导出）"
    - "实时日志流式传输功能"
    - "日志数据的修改功能"
dependencies: 
  - "用户管理模块（记录操作用户信息）"
  - "权限管理模块（控制日志访问权限）"
  - "现有业务模块（集成审计日志记录）"
```

## 2. 设计引用

```yaml
design_sources:
  requirements: "5.7 审计日志 - 产品需求说明书V1.1.md"
  technical_spec: "2.4.7 审计日志 - 详细设计说明书V1.1.md"
  database_schema: "审计日志表（audit_logs）- 数据库设计说明书V1.1.md"
  api_endpoints: "4.4.7 审计日志 - API接口设计说明书V1.1.md"

architecture_alignment:
  patterns: ["分层架构", "手动依赖注入", "RESTful API", "基于接口的设计"]
  frameworks: ["Gin HTTP", "GORM ORM", "Zap日志", "JWT认证"]
  rationale: "遵循现有项目的分层架构模式，通过接口定义实现松耦合，使用GORM进行数据访问，集成现有的认证和日志系统"
  
project_structure:
  interface_layer: "internal/interface/ - 为每层定义契约"
  implementation_layer: "internal/repository/, internal/service/, internal/controller/ - 实现接口"
  dependency_flow: "Controller -> Service Interface -> Repository Interface"
  injection_point: "main.go - 通过构造函数进行手动依赖注入"
```

## 3. 数据层设计

```yaml
models:
  AuditLog:
    table: "audit_logs"
    file: "internal/model/audit_log.go"
    business_purpose: "记录平台所有业务操作的审计日志，用于安全审计和合规检查"
    fields:
      - name: "id"
        type: "uint"
        business_meaning: "审计日志唯一标识"
        gorm: "primarykey;autoIncrement"
        json: "id"
        validation: ""
      - name: "user_id"
        type: "*uint"
        business_meaning: "执行操作的用户ID，可为空（系统操作）"
        gorm: "index"
        json: "user_id"
        validation: ""
      - name: "username"
        type: "string"
        business_meaning: "执行操作的用户名"
        gorm: "column:username;type:varchar(64)"
        json: "username"
        validation: ""
      - name: "action"
        type: "string"
        business_meaning: "操作动作类型（CREATE、UPDATE、DELETE、READ等）"
        gorm: "column:action;type:varchar(32);not null;index"
        json: "action"
        validation: "required,oneof=CREATE READ UPDATE DELETE LOGIN LOGOUT"
      - name: "module"
        type: "string"
        business_meaning: "操作所属功能模块（应用管理、用户管理、项目管理等）"
        gorm: "column:module;type:varchar(32);not null;index"
        json: "module"
        validation: "required"
      - name: "resource_type"
        type: "string"
        business_meaning: "操作的资源类型（USER、APP、PROJECT等）"
        gorm: "column:resource_type;type:varchar(32);index"
        json: "resource_type"
        validation: ""
      - name: "resource_id"
        type: "*uint"
        business_meaning: "操作的资源ID"
        gorm: "column:resource_id;index"
        json: "resource_id"
        validation: ""
      - name: "resource_name"
        type: "string"
        business_meaning: "操作的资源名称"
        gorm: "column:resource_name;type:varchar(64)"
        json: "resource_name"
        validation: ""
      - name: "description"
        type: "string"
        business_meaning: "操作的详细描述"
        gorm: "column:description;type:text"
        json: "description"
        validation: ""
      - name: "ip_address"
        type: "string"
        business_meaning: "客户端IP地址"
        gorm: "column:ip_address;type:varchar(45);index"
        json: "ip_address"
        validation: ""
      - name: "user_agent"
        type: "string"
        business_meaning: "客户端用户代理信息"
        gorm: "column:user_agent;type:varchar(255)"
        json: "user_agent"
        validation: ""
      - name: "request_method"
        type: "string"
        business_meaning: "HTTP请求方法"
        gorm: "column:request_method;type:varchar(10)"
        json: "request_method"
        validation: ""
      - name: "request_url"
        type: "string"
        business_meaning: "请求URL路径"
        gorm: "column:request_url;type:varchar(255)"
        json: "request_url"
        validation: ""
      - name: "request_params"
        type: "datatypes.JSON"
        business_meaning: "请求参数（脱敏处理）"
        gorm: "column:request_params;type:json"
        json: "request_params"
        validation: ""
      - name: "response_status"
        type: "*int"
        business_meaning: "HTTP响应状态码"
        gorm: "column:response_status;index"
        json: "response_status"
        validation: ""
      - name: "response_time"
        type: "*int"
        business_meaning: "响应时间（毫秒）"
        gorm: "column:response_time"
        json: "response_time"
        validation: ""
      - name: "success"
        type: "bool"
        business_meaning: "操作是否成功"
        gorm: "column:success;default:true;index"
        json: "success"
        validation: ""
      - name: "error_message"
        type: "string"
        business_meaning: "错误信息（如果操作失败）"
        gorm: "column:error_message;type:text"
        json: "error_message"
        validation: ""
      - name: "created_at"
        type: "time.Time"
        business_meaning: "操作发生时间"
        gorm: "column:created_at;not null;index"
        json: "created_at"
        validation: ""
    relationships:
      - type: "belongs_to"
        target: "User"
        business_rule: "审计日志关联操作用户，但支持匿名操作"
        foreign_key: "user_id"
    methods:
      - name: "IsSuccessful"
        purpose: "检查操作是否成功执行"
      - name: "GetDuration"
        purpose: "获取操作执行时长的人类可读格式"
      - name: "MaskSensitiveData"
        purpose: "脱敏处理敏感信息"
    indexes:
      - fields: ["user_id", "created_at"]
        type: "btree"
        purpose: "按用户查询操作历史的性能优化"
      - fields: ["module", "action", "created_at"]
        type: "btree"
        purpose: "按模块和操作类型查询的性能优化"
      - fields: ["created_at"]
        type: "btree"
        purpose: "时间范围查询优化"
      - fields: ["ip_address", "created_at"]
        type: "btree"
        purpose: "安全分析和IP追踪优化"
```

## 4. 接口定义

```yaml
interfaces:
  repository_interface:
    file: "internal/interface/repository/audit_log.go"
    interface: "AuditLogRepository"
    purpose: "审计日志数据访问操作的契约"
    methods:
      - name: "Create"
        signature: "Create(ctx context.Context, log *model.AuditLog) error"
        purpose: "创建新的审计日志记录"
      - name: "GetByID"
        signature: "GetByID(ctx context.Context, id uint) (*model.AuditLog, error)"
        purpose: "根据ID查询审计日志详情"
      - name: "List"
        signature: "List(ctx context.Context, filter *AuditLogFilter) ([]*model.AuditLog, int64, error)"
        purpose: "分页查询审计日志，支持多条件筛选"
      - name: "GetStatistics"
        signature: "GetStatistics(ctx context.Context, filter *StatisticsFilter) (*AuditLogStatistics, error)"
        purpose: "获取审计日志统计信息"
      - name: "Export"
        signature: "Export(ctx context.Context, filter *AuditLogFilter) ([]*model.AuditLog, error)"
        purpose: "导出审计日志数据（不分页）"
      - name: "DeleteOldLogs"
        signature: "DeleteOldLogs(ctx context.Context, beforeDate time.Time) (int64, error)"
        purpose: "清理过期日志数据（系统维护）"
    
  service_interface:
    file: "internal/interface/service/audit_log.go"
    interface: "AuditLogService"
    purpose: "审计日志业务逻辑操作的契约"
    methods:
      - name: "RecordLog"
        signature: "RecordLog(ctx context.Context, req *request.RecordAuditLogRequest) error"
        purpose: "记录业务操作的审计日志"
      - name: "GetAuditLog"
        signature: "GetAuditLog(ctx context.Context, id uint) (*response.AuditLogResponse, error)"
        purpose: "获取单个审计日志详情"
      - name: "ListAuditLogs"
        signature: "ListAuditLogs(ctx context.Context, req *request.ListAuditLogsRequest) (*response.ListAuditLogsResponse, error)"
        purpose: "查询审计日志列表，支持筛选和分页"
      - name: "GetStatistics"
        signature: "GetStatistics(ctx context.Context, req *request.AuditLogStatisticsRequest) (*response.AuditLogStatisticsResponse, error)"
        purpose: "获取审计日志统计分析数据"
      - name: "ExportLogs"
        signature: "ExportLogs(ctx context.Context, req *request.ExportAuditLogsRequest) ([]byte, string, error)"
        purpose: "导出审计日志到指定格式"

  dependency_injection:
    pattern: "在main.go中进行手动依赖注入"
    constructor_naming: "NewAuditLogRepository, NewAuditLogService, NewAuditLogController"
    interface_usage: "将接口作为构造函数参数传递"
```

## 5. Repository层

```yaml
repository:
  file: "internal/repository/audit_log.go"
  interface_file: "internal/interface/repository/audit_log.go"
  interface: "AuditLogRepository"
  purpose: "审计日志数据访问抽象层"
  methods:
    - name: "Create"
      purpose: "持久化新的审计日志记录"
      signature: "Create(ctx context.Context, log *model.AuditLog) error"
      business_logic: "验证必填字段，自动设置创建时间"
      error_handling: "处理数据库连接错误和约束违规"
    
    - name: "GetByID"
      purpose: "根据ID查询审计日志"
      signature: "GetByID(ctx context.Context, id uint) (*model.AuditLog, error)"
      error_handling: "处理记录不存在的情况"
    
    - name: "List"
      purpose: "分页查询审计日志，支持多条件筛选"
      signature: "List(ctx context.Context, filter *AuditLogFilter) ([]*model.AuditLog, int64, error)"
      pagination: "使用LIMIT和OFFSET实现分页"
      filtering: "支持用户ID、操作类型、模块、时间范围等条件筛选"
      
    - name: "GetStatistics"
      purpose: "获取审计日志统计信息"
      signature: "GetStatistics(ctx context.Context, filter *StatisticsFilter) (*AuditLogStatistics, error)"
      aggregation: "统计操作次数、成功率、热门用户和操作类型"
      
    - name: "Export"
      purpose: "导出审计日志数据"
      signature: "Export(ctx context.Context, filter *AuditLogFilter) ([]*model.AuditLog, error)"
      limit: "导出数据量限制，防止内存溢出"
      
    - name: "DeleteOldLogs"
      purpose: "清理过期的审计日志"
      signature: "DeleteOldLogs(ctx context.Context, beforeDate time.Time) (int64, error)"
      maintenance: "系统维护功能，定期清理历史数据"

  implementation_notes:
    - "所有方法必须使用WithContext(ctx)"
    - "使用一致的错误包装处理"
    - "记录关键操作日志"
    - "查询性能优化，使用适当的索引"
```

## 6. Service层

```yaml
service:
  file: "internal/service/audit_log.go"
  interface_file: "internal/interface/service/audit_log.go"
  interface: "AuditLogService"
  purpose: "审计日志业务逻辑编排"
  dependencies:
    - "AuditLogRepository"
    - "UserRepository"
    - "logger.Logger"
  
  methods:
    - name: "RecordLog"
      purpose: "记录业务操作的审计日志"
      signature: "RecordLog(ctx context.Context, req *request.RecordAuditLogRequest) error"
      business_rules:
        - "自动提取用户信息和请求上下文"
        - "敏感数据脱敏处理"
        - "必填字段验证"
      steps:
        - "验证输入参数"
        - "提取用户和请求信息"
        - "脱敏处理敏感数据"
        - "保存审计日志"
      error_handling: "使用errors.NewAppError包装错误"
    
    - name: "GetAuditLog"
      purpose: "获取审计日志详情，包含权限检查"
      signature: "GetAuditLog(ctx context.Context, id uint) (*response.AuditLogResponse, error)"
      authorization: "验证用户查看日志权限"
      steps:
        - "权限验证"
        - "查询日志记录"
        - "关联用户信息"
        - "返回格式化响应"
    
    - name: "ListAuditLogs"
      purpose: "查询审计日志列表，支持筛选"
      signature: "ListAuditLogs(ctx context.Context, req *request.ListAuditLogsRequest) (*response.ListAuditLogsResponse, error)"
      steps:
        - "参数验证和默认值设置"
        - "权限控制（用户只能查看相关日志）"
        - "构建查询条件"
        - "分页查询"
        - "格式化响应数据"
    
    - name: "GetStatistics"
      purpose: "获取审计日志统计分析"
      signature: "GetStatistics(ctx context.Context, req *request.AuditLogStatisticsRequest) (*response.AuditLogStatisticsResponse, error)"
      validation: "时间范围和统计维度验证"
      
    - name: "ExportLogs"
      purpose: "导出审计日志到指定格式"
      signature: "ExportLogs(ctx context.Context, req *request.ExportAuditLogsRequest) ([]byte, string, error)"
      formats: "支持CSV、Excel、JSON格式导出"
      security: "导出操作本身记录审计日志"

  logging_strategy:
    - "记录所有审计日志操作的开始和结果"
    - "使用结构化日志记录关键业务参数"
    - "错误日志包含足够的上下文信息"
    - "导出操作记录详细的操作日志"
```

## 7. Controller层

```yaml
controller:
  file: "internal/controller/audit_log.go"
  purpose: "审计日志HTTP请求处理"
  dependencies:
    - "AuditLogService"
    - "logger.Logger"
  
  common_methods:
    - name: "bindAndValidateRequest"
      purpose: "通用请求绑定和验证"
      reuse: "所有处理器使用"
  
  endpoints:
    - method: "GET"
      path: "/api/v1/audit-logs"
      handler: "ListAuditLogs"
      purpose: "查询审计日志列表"
      auth: true
      query_params: ["page", "page_size", "user_id", "action", "module", "resource_type", "start_time", "end_time", "ip_address"]
      response: "response.ListAuditLogsResponse"
      status_codes: [200, 400, 401, 403, 500]
      business_flow: "管理员查询平台操作记录"
      error_handling: "使用errors.HandleError处理所有错误"
    
    - method: "GET"
      path: "/api/v1/audit-logs/:id"
      handler: "GetAuditLog"
      purpose: "获取单个审计日志详情"
      auth: true
      params: ["id (路径参数)"]
      response: "response.AuditLogResponse"
      status_codes: [200, 401, 403, 404, 500]
      validation: "验证ID参数格式"
    
    - method: "GET"
      path: "/api/v1/audit-logs/statistics"
      handler: "GetStatistics"
      purpose: "获取审计日志统计信息"
      auth: true
      query_params: ["start_time", "end_time", "group_by"]
      response: "response.AuditLogStatisticsResponse"
      status_codes: [200, 400, 401, 403, 500]
      validation: "验证时间范围和分组参数"
    
    - method: "GET"
      path: "/api/v1/audit-logs/export"
      handler: "ExportLogs"
      purpose: "导出审计日志"
      auth: true
      query_params: ["format", "start_time", "end_time", "user_id", "action", "module"]
      response: "文件下载"
      status_codes: [200, 400, 401, 403, 500]
      download: "设置适当的Content-Type和文件名"

  middleware_usage:
    - "使用统一的认证中间件"
    - "应用错误处理中间件"
    - "使用日志中间件进行请求跟踪"
    - "应用i18n中间件"
    - "添加审计日志记录中间件"
```

## 8. Router集成

```yaml
router_configuration:
  file: "internal/router/router.go"
  integration_point: "SetupRouter函数"
  controller_struct: "Controllers"
  
  route_group_setup:
    path: "/api/v1/audit-logs"
    middleware: ["认证中间件（需要管理员权限）"]
    methods:
      - route: "GET /"
        handler: "controllers.AuditLogController.ListAuditLogs"
        purpose: "查询审计日志列表"
        auth_required: true
      - route: "GET /:id"
        handler: "controllers.AuditLogController.GetAuditLog"
        purpose: "获取审计日志详情"
        auth_required: true
      - route: "GET /statistics"
        handler: "controllers.AuditLogController.GetStatistics"
        purpose: "获取统计信息"
        auth_required: true
      - route: "GET /export"
        handler: "controllers.AuditLogController.ExportLogs"
        purpose: "导出审计日志"
        auth_required: true

  controller_struct_update:
    location: "router.Controllers结构体"
    addition: "AuditLogController *controller.AuditLogController"
    
  route_registration_example: |
    // 审计日志管理路由（需要管理员权限）
    auditLogGroup := protected.Group("/audit-logs")
    auditLogGroup.Use(middleware.RequireAdmin()) // 仅管理员可访问
    auditLogGroup.GET("", controllers.AuditLogController.ListAuditLogs)
    auditLogGroup.GET("/:id", controllers.AuditLogController.GetAuditLog)
    auditLogGroup.GET("/statistics", controllers.AuditLogController.GetStatistics)
    auditLogGroup.GET("/export", controllers.AuditLogController.ExportLogs)

  audit_middleware:
    description: "添加审计日志记录中间件，自动记录API操作"
    implementation: "middleware/audit.go"
    usage: "在需要记录的路由组上应用"
```

## 9. 数据传输对象

```yaml
dto_structures:
  request_file: "internal/dto/request/audit_log.go"
  response_file: "internal/dto/response/audit_log.go"
  
  requests:
    RecordAuditLogRequest:
      purpose: "记录审计日志的输入"
      validation_strategy: "使用Gin绑定标签"
      fields:
        - name: "user_id"
          type: "*uint"
          business_meaning: "操作用户ID"
          validation: ""
          json: "user_id"
        - name: "action"
          type: "string"
          business_meaning: "操作类型"
          validation: "required,oneof=CREATE READ UPDATE DELETE LOGIN LOGOUT"
          json: "action"
        - name: "module"
          type: "string"
          business_meaning: "功能模块"
          validation: "required"
          json: "module"
        - name: "resource_type"
          type: "string"
          business_meaning: "资源类型"
          validation: ""
          json: "resource_type"
        - name: "resource_id"
          type: "*uint"
          business_meaning: "资源ID"
          validation: ""
          json: "resource_id"
        - name: "resource_name"
          type: "string"
          business_meaning: "资源名称"
          validation: ""
          json: "resource_name"
        - name: "description"
          type: "string"
          business_meaning: "操作描述"
          validation: ""
          json: "description"
    
    ListAuditLogsRequest:
      purpose: "查询审计日志的参数"
      embedding: "嵌入通用分页结构"
      fields:
        - name: "Page"
          type: "int"
          default: 1
          json: "page"
          validation: "min=1"
        - name: "PageSize"
          type: "int"
          default: 20
          json: "page_size"
          validation: "min=1,max=100"
        - name: "UserID"
          type: "*uint"
          json: "user_id"
          purpose: "按用户ID筛选"
        - name: "Action"
          type: "string"
          json: "action"
          purpose: "按操作类型筛选"
        - name: "Module"
          type: "string"
          json: "module"
          purpose: "按模块筛选"
        - name: "ResourceType"
          type: "string"
          json: "resource_type"
          purpose: "按资源类型筛选"
        - name: "ResourceID"
          type: "*uint"
          json: "resource_id"
          purpose: "按资源ID筛选"
        - name: "StartTime"
          type: "string"
          json: "start_time"
          purpose: "开始时间"
          validation: "omitempty,datetime=2006-01-02T15:04:05Z07:00"
        - name: "EndTime"
          type: "string"
          json: "end_time"
          purpose: "结束时间"
          validation: "omitempty,datetime=2006-01-02T15:04:05Z07:00"
        - name: "IPAddress"
          type: "string"
          json: "ip_address"
          purpose: "按IP地址筛选"
        - name: "Success"
          type: "*bool"
          json: "success"
          purpose: "按操作结果筛选"
    
    AuditLogStatisticsRequest:
      purpose: "审计日志统计查询参数"
      fields:
        - name: "StartTime"
          type: "string"
          json: "start_time"
          validation: "omitempty,datetime=2006-01-02T15:04:05Z07:00"
        - name: "EndTime"
          type: "string"
          json: "end_time"
          validation: "omitempty,datetime=2006-01-02T15:04:05Z07:00"
        - name: "GroupBy"
          type: "string"
          json: "group_by"
          default: "day"
          validation: "omitempty,oneof=hour day week month"
    
    ExportAuditLogsRequest:
      purpose: "导出审计日志参数"
      fields:
        - name: "Format"
          type: "string"
          json: "format"
          default: "csv"
          validation: "omitempty,oneof=csv excel json"
        - name: "StartTime"
          type: "string"
          json: "start_time"
          validation: "omitempty,datetime=2006-01-02T15:04:05Z07:00"
        - name: "EndTime"
          type: "string"
          json: "end_time"
          validation: "omitempty,datetime=2006-01-02T15:04:05Z07:00"
        - name: "UserID"
          type: "*uint"
          json: "user_id"
        - name: "Action"
          type: "string"
          json: "action"
        - name: "Module"
          type: "string"
          json: "module"
  
  responses:
    AuditLogResponse:
      purpose: "审计日志详情的公开表示"
      data_mapping: "从model.AuditLog转换"
      fields:
        - name: "ID"
          type: "uint"
          json: "id"
          source: "model.ID"
        - name: "User"
          type: "*UserBasicInfo"
          json: "user"
          source: "关联查询用户信息"
        - name: "Action"
          type: "string"
          json: "action"
          source: "model.Action"
        - name: "Module"
          type: "string"
          json: "module"
          source: "model.Module"
        - name: "ResourceType"
          type: "string"
          json: "resource_type"
          source: "model.ResourceType"
        - name: "ResourceID"
          type: "*uint"
          json: "resource_id"
          source: "model.ResourceID"
        - name: "ResourceName"
          type: "string"
          json: "resource_name"
          source: "model.ResourceName"
        - name: "Description"
          type: "string"
          json: "description"
          source: "model.Description"
        - name: "IPAddress"
          type: "string"
          json: "ip_address"
          source: "model.IPAddress"
        - name: "UserAgent"
          type: "string"
          json: "user_agent"
          source: "model.UserAgent"
        - name: "RequestMethod"
          type: "string"
          json: "request_method"
          source: "model.RequestMethod"
        - name: "RequestURL"
          type: "string"
          json: "request_url"
          source: "model.RequestURL"
        - name: "RequestParams"
          type: "interface{}"
          json: "request_params"
          source: "model.RequestParams（脱敏后）"
        - name: "ResponseStatus"
          type: "*int"
          json: "response_status"
          source: "model.ResponseStatus"
        - name: "ResponseTime"
          type: "*int"
          json: "response_time"
          source: "model.ResponseTime"
        - name: "Success"
          type: "bool"
          json: "success"
          source: "model.Success"
        - name: "ErrorMessage"
          type: "string"
          json: "error_message"
          source: "model.ErrorMessage"
        - name: "CreatedAt"
          type: "string"
          json: "created_at"
          format: "time.RFC3339"
          source: "model.CreatedAt"
    
    ListAuditLogsResponse:
      purpose: "分页的审计日志集合"
      structure: "标准列表响应格式"
      fields:
        - name: "Items"
          type: "[]*AuditLogResponse"
          json: "items"
          source: "转换后的model列表"
        - name: "Pagination"
          type: "*response.PaginationInfo"
          json: "pagination"
          structure: "标准分页信息"
    
    AuditLogStatisticsResponse:
      purpose: "审计日志统计信息"
      fields:
        - name: "TotalOperations"
          type: "int64"
          json: "total_operations"
        - name: "SuccessOperations"
          type: "int64"
          json: "success_operations"
        - name: "FailedOperations"
          type: "int64"
          json: "failed_operations"
        - name: "SuccessRate"
          type: "float64"
          json: "success_rate"
        - name: "TopUsers"
          type: "[]*UserOperationStat"
          json: "top_users"
        - name: "TopActions"
          type: "[]*ActionStat"
          json: "top_actions"
        - name: "Timeline"
          type: "[]*TimelineStat"
          json: "timeline"

  conversion_methods:
    - "ToAuditLogResponse(model *model.AuditLog, user *model.User) *AuditLogResponse"
    - "ToAuditLogResponseList(models []*model.AuditLog, userMap map[uint]*model.User) []*AuditLogResponse"
    - "FromRecordAuditLogRequest(req *RecordAuditLogRequest) *model.AuditLog"
```

## 10. 实施任务清单

```yaml
tasks:
  phase_1_data_layer:
    - task: "创建AuditLog数据模型"
      file: "internal/model/audit_log.go"
      description: "定义GORM标签和关联关系的结构体"
      details:
        - "定义所有字段和GORM标签"
        - "实现业务逻辑方法（IsSuccessful、GetDuration、MaskSensitiveData）"
        - "定义关联关系"
        - "添加索引定义"
    
    - task: "在main.go中添加数据库迁移"
      file: "main.go"
      description: "在AutoMigrate中注册新模型"
      location: "数据库迁移部分"
  
  phase_2_interface_definitions:
    - task: "定义AuditLogRepository接口"
      file: "internal/interface/repository/audit_log.go"
      description: "创建数据访问接口契约"
      details:
        - "定义所有CRUD方法签名"
        - "包含适当的上下文和错误处理"
        - "为每个方法添加文档"
        - "定义筛选和统计相关的数据结构"
    
    - task: "定义AuditLogService接口"
      file: "internal/interface/service/audit_log.go" 
      description: "创建服务接口契约"
      details:
        - "定义所有业务方法签名"
        - "包含请求/响应DTO类型"
        - "添加详细文档"
  
  phase_3_repository_layer:
    - task: "实现AuditLogRepository"
      file: "internal/repository/audit_log.go"
      description: "使用GORM实现数据访问接口"
      details:
        - "实现所有接口方法"
        - "添加适当的错误处理和日志记录"
        - "所有数据库操作使用WithContext"
        - "包含构造函数：NewAuditLogRepository"
        - "实现高效的查询和索引使用"
        - "添加数据脱敏逻辑"
  
  phase_4_service_layer:
    - task: "实现AuditLogService"
      file: "internal/service/audit_log.go"
      description: "实现业务逻辑接口"
      details:
        - "实现所有接口方法"
        - "添加业务规则验证"
        - "处理DTO转换"
        - "添加结构化日志记录"
        - "使用统一错误处理"
        - "包含构造函数：NewAuditLogService"
        - "实现导出功能的各种格式支持"
        - "添加权限验证逻辑"
  
  phase_5_controller_layer:
    - task: "创建请求DTOs"
      file: "internal/dto/request/audit_log.go"
      description: "定义API请求数据结构"
      validation: "添加参数验证标签"
      details:
        - "RecordAuditLogRequest"
        - "ListAuditLogsRequest"
        - "AuditLogStatisticsRequest"
        - "ExportAuditLogsRequest"
    
    - task: "创建响应DTOs"
      file: "internal/dto/response/audit_log.go"
      description: "定义API响应数据结构"
      conversion: "实现模型到响应的转换方法"
      details:
        - "AuditLogResponse"
        - "ListAuditLogsResponse"
        - "AuditLogStatisticsResponse"
        - "相关统计数据结构"
    
    - task: "实现AuditLogController"
      file: "internal/controller/audit_log.go"
      description: "创建Gin处理器，包含适当的错误处理"
      details:
        - "实现所有HTTP处理器"
        - "使用通用的bindAndValidateRequest方法"
        - "应用统一响应格式"
        - "添加适当的HTTP状态码"
        - "包含构造函数：NewAuditLogController"
        - "实现文件下载处理"
  
  phase_6_middleware_implementation:
    - task: "创建审计日志中间件"
      file: "internal/middleware/audit.go"
      description: "自动记录API操作的中间件"
      details:
        - "提取请求信息（IP、User-Agent、URL等）"
        - "捕获响应状态和执行时间"
        - "异步记录审计日志，避免影响性能"
        - "处理敏感数据脱敏"
        - "支持白名单配置（某些操作不记录）"
    
    - task: "创建管理员权限中间件"
      file: "internal/middleware/admin.go"
      description: "验证管理员权限的中间件"
      details:
        - "检查用户角色是否为管理员"
        - "返回适当的错误响应"
  
  phase_7_router_integration:
    - task: "更新Controllers结构体"
      file: "internal/router/router.go"
      description: "在Controllers结构体中添加AuditLogController"
      location: "Controllers结构体定义"
    
    - task: "注册审计日志路由"
      file: "internal/router/router.go"
      description: "在SetupRouter函数中添加审计日志路由组"
      details:
        - "为audit-logs创建路由组"
        - "添加所有CRUD端点"
        - "应用认证和管理员权限中间件"
        - "应用审计日志记录中间件到其他业务路由"
        - "遵循现有路由模式"
  
  phase_8_dependency_injection:
    - task: "初始化AuditLogRepository"
      file: "main.go"
      description: "在main函数中添加数据访问层初始化"
      location: "Repository初始化部分"
      code_example: "auditLogRepo := repository.NewAuditLogRepository(db)"
    
    - task: "初始化AuditLogService"
      file: "main.go"
      description: "添加服务层初始化和依赖"
      location: "Service初始化部分"
      code_example: "auditLogService := service.NewAuditLogService(auditLogRepo, userRepo, zapLogger)"
    
    - task: "初始化AuditLogController"
      file: "main.go"
      description: "添加控制器初始化"
      location: "Controller初始化部分"
      code_example: "auditLogController := controller.NewAuditLogController(auditLogService, zapLogger)"
    
    - task: "更新Controllers结构体初始化"
      file: "main.go"
      description: "在路由设置中添加AuditLogController"
      location: "Router设置部分"
      code_example: |
        r := router.SetupRouter(&router.Controllers{
          UserController: userController,
          I18nController: i18nController,
          AuditLogController: auditLogController,
        }, cfg, zapLogger)
  
  phase_9_integration_setup:
    - task: "集成审计日志到现有业务操作"
      files: ["internal/service/user.go", "其他业务服务"]
      description: "在关键业务操作中添加审计日志记录"
      details:
        - "用户创建、更新、删除操作"
        - "登录、登出操作"
        - "权限变更操作"
        - "应用管理操作"
        - "项目管理操作"
        - "系统配置变更操作"
    
    - task: "配置审计日志自动清理"
      file: "internal/service/audit_log.go"
      description: "添加定期清理过期日志的功能"
      details:
        - "创建清理任务"
        - "配置保留期限（默认1年）"
        - "添加系统任务调度"
  
  phase_10_internationalization:
    - task: "添加中文翻译"
      file: "pkg/i18n/locales/zh-CN.yaml"
      description: "为审计日志相关消息添加中文翻译"
      keys:
        - "audit_log.list.success"
        - "audit_log.detail.not_found"
        - "audit_log.export.success"
        - "audit_log.statistics.success"
        - "audit_log.permission.denied"
    
    - task: "添加英文翻译"
      file: "pkg/i18n/locales/en-US.yaml"
      description: "为审计日志相关消息添加英文翻译"
  
  phase_11_testing:
    - task: "编写Repository单元测试"
      file: "internal/repository/audit_log_test.go"
      description: "测试所有数据访问方法"
      coverage: "覆盖正常和错误情况"
      details:
        - "测试CRUD操作"
        - "测试查询筛选功能"
        - "测试统计功能"
        - "测试分页功能"
    
    - task: "编写Service单元测试"
      file: "internal/service/audit_log_test.go"
      description: "使用mock测试业务逻辑"
      mocking: "模拟Repository依赖"
      details:
        - "测试业务规则验证"
        - "测试权限控制"
        - "测试数据转换"
        - "测试导出功能"
    
    - task: "编写Controller集成测试"
      file: "internal/controller/audit_log_test.go"
      description: "测试完整的HTTP工作流"
      setup: "设置测试数据库和路由"
      details:
        - "测试所有HTTP端点"
        - "测试参数验证"
        - "测试错误处理"
        - "测试权限控制"
    
    - task: "编写中间件测试"
      file: "internal/middleware/audit_test.go"
      description: "测试审计日志中间件功能"
      scenarios: "测试各种请求场景"
    
    - task: "编写API集成测试"
      file: "test/integration/audit_log_test.go"
      description: "端到端API测试"
      scenarios: "测试完整用户场景"
  
  phase_12_performance_optimization:
    - task: "数据库索引优化"
      description: "确保查询性能"
      details:
        - "验证索引效果"
        - "分析慢查询"
        - "优化复合索引"
    
    - task: "分页查询优化"
      description: "大数据量查询优化"
      details:
        - "使用游标分页（如需要）"
        - "查询结果缓存"
        - "限制导出数据量"
    
    - task: "异步日志记录"
      description: "提高日志记录性能"
      details:
        - "使用队列异步处理"
        - "批量插入优化"
        - "错误重试机制"
```

## 11. 质量验收标准

```yaml
acceptance_criteria:
  functionality:
    - "所有审计日志CRUD操作正常工作"
    - "多条件筛选查询功能正确"
    - "统计分析数据准确"
    - "导出功能支持多种格式"
    - "敏感数据正确脱敏"
    - "分页和排序功能正常"
    - "审计日志自动记录所有业务操作"
    - "权限控制正确实施"
  
  code_quality:
    - "所有函数都有适当的错误处理"
    - "代码遵循项目约定"
    - "日志记录实施恰当"
    - "输入验证防止无效数据"
    - "使用统一响应格式"
    - "Context正确传递给所有层"
    - "接口定义清晰，依赖注入正确"
  
  testing:
    - "单元测试覆盖率 > 80%"
    - "集成测试覆盖主要工作流"
    - "所有测试成功通过"
    - "错误场景得到测试"
    - "性能关键路径有基准测试"
  
  security:
    - "所有端点需要管理员认证"
    - "敏感数据正确脱敏"
    - "输入清理防止注入攻击"
    - "敏感数据不被记录到日志"
    - "使用HTTPS和安全头"
    - "导出功能有适当的权限控制"
  
  performance:
    - "查询响应时间 < 500ms"
    - "数据库查询已优化"
    - "适当的索引已建立"
    - "分页防止大结果集"
    - "导出功能有数据量限制"
    - "异步日志记录不影响业务性能"
  
  usability:
    - "API端点有适当的文档"
    - "错误消息清晰且本地化"
    - "统计信息格式易于理解"
    - "导出文件格式标准化"
  
  compliance:
    - "审计日志符合合规要求"
    - "数据保留策略正确实施"
    - "日志完整性得到保证"
    - "访问控制符合安全标准"

deployment_checklist:
  database:
    - "迁移脚本已测试"
    - "索引创建已验证"
    - "数据备份计划"
    - "性能监控设置"
  
  configuration:
    - "环境变量已文档化"
    - "配置验证"
    - "合理的默认值"
    - "日志保留期配置"
  
  monitoring:
    - "健康检查端点"
    - "关键指标监控"
    - "错误率告警"
    - "性能监控"
    - "磁盘空间监控（日志增长）"
  
  maintenance:
    - "日志清理任务调度"
    - "备份和恢复程序"
    - "性能调优指南"
    - "故障排除文档"
```

## 12. 实施注意事项

### 关键技术决策
1. **数据脱敏策略**: 在记录和导出时自动脱敏密码、密钥等敏感信息
2. **性能优化**: 使用异步队列记录日志，避免影响业务操作性能
3. **存储策略**: 考虑日志数据增长，实施自动清理机制
4. **权限控制**: 仅管理员可访问审计日志，确保安全合规

### 扩展性考虑
1. **大数据量支持**: 预留分表分库的接口设计
2. **实时分析**: 为未来的实时监控预留数据接口
3. **第三方集成**: 支持导出到外部审计系统
4. **多租户支持**: 数据隔离和权限隔离设计

### 运维要求
1. **监控告警**: 监控日志记录失败和存储空间使用
2. **备份策略**: 审计日志的备份和归档策略
3. **性能调优**: 定期分析查询性能和索引效果
4. **合规检查**: 定期验证日志完整性和访问控制

此实施计划提供了审计日志功能的完整技术实现路径，遵循现有项目架构模式，确保代码质量、安全性和可维护性。
