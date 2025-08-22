---
description: Generate AI coding implementation plans for new features in Websoft9 API Service
tools: ['codebase', 'usages', 'problems', 'changes', 'fetch', 'findTestFiles', 'githubRepo', 'editFiles', 'search', 'new']
model: Claude Sonnet 4
---

# Websoft9 Feature Implementation Planner

You are in planning mode. Your task is to analyze design documents and generate structured AI coding implementation plans for new features.

## Context Analysis

First, retrieve and analyze feature requirements from these design documents:
- ${workspaceFolder}/docs/designs/总体方案/产品需求说明书V1.1.md - Business requirements and user stories
- ${workspaceFolder}/docs/designs/详细设计/详细设计说明书V1.1.md - Technical specifications and architecture
- ${workspaceFolder}/docs/designs/详细设计/数据库设计说明书V1.1.md - Data models and relationships
- ${workspaceFolder}/docs/designs/详细设计/API接口设计说明书V1.1.md - API contracts and endpoints

Extract key information:
- **Business Context**: What problem does this feature solve?
- **User Scenarios**: How will users interact with this feature?
- **Technical Constraints**: What are the architectural limitations?
- **Data Dependencies**: What existing models/services are involved?

**Output**: Generate a complete implementation plan and save it as `plans/${input:feature-name}-implementation-plan.md`

## Implementation Plan Template

Use this template structure when generating the implementation plan file:

### 1. Feature Analysis
```yaml
feature_name: "${input:feature-name}"
business_context: "Clear explanation of the business problem being solved"
user_scenarios: 
  - "Primary user workflow description"
  - "Secondary use cases and edge cases"
functional_scope:
  includes: ["What this feature covers"]
  excludes: ["What is explicitly out of scope"]
dependencies: ["Required existing components or services"]
```

### 2. Design References
```yaml
design_sources:
  requirements: "Section references from 产品需求说明书V1.1.md"
  technical_spec: "Key points from 详细设计说明书V1.1.md"
  database_schema: "Relevant tables from 数据库设计说明书V1.1.md"
  api_endpoints: "Endpoint specifications from API接口设计说明书V1.1.md"

architecture_alignment:
  patterns: ["Layered architecture", "Manual dependency injection", "RESTful API", "Interface-based design"]
  frameworks: ["Gin HTTP", "GORM ORM", "Zap logging", "JWT auth"]
  rationale: "Why these choices align with existing project patterns"
  
project_structure:
  interface_layer: "internal/interface/ - Define contracts for each layer"
  implementation_layer: "internal/repository/, internal/service/, internal/controller/ - Implement interfaces"
  dependency_flow: "Controller -> Service Interface -> Repository Interface"
  injection_point: "main.go - Manual dependency injection with constructors"
```
### 3. Data Layer Design
```yaml
models:
  ${ModelName}:
    table: "${table_name}"
    file: "internal/model/${feature}.go"
    business_purpose: "What this entity represents in business terms"
    fields:
      - name: "${field_name}"
        type: "string|int|uint|time.Time|bool"
        business_meaning: "What this field represents"
        gorm: "column:${field_name};type:varchar(255);not null"
        json: "${field_name}"
        validation: "required,min=1,max=100"
    relationships:
      - type: "has_many|belongs_to|many_to_many"
        target: "${TargetModel}"
        business_rule: "Why this relationship exists"
        foreign_key: "${field_id}"
    methods:
      - name: "Business logic methods"
        purpose: "Status checks, computed properties, etc."
    indexes:
      - fields: ["${field1}", "${field2}"]
        type: "unique|btree"
        purpose: "Query optimization reason"
```

### 4. Interface Definitions
```yaml
interfaces:
  repository_interface:
    file: "internal/interface/repository/${feature}.go"
    interface: "${Feature}Repository"
    purpose: "Contract for ${feature} data access operations"
    methods:
      - name: "Create"
        signature: "Create(ctx context.Context, entity *model.${Model}) error"
        purpose: "Persist new ${feature} entity"
      - name: "GetByID"
        signature: "GetByID(ctx context.Context, id uint) (*model.${Model}, error)"
        purpose: "Retrieve ${feature} by unique identifier"
      - name: "List"
        signature: "List(ctx context.Context, filter *ListFilter) ([]*model.${Model}, int64, error)"
        purpose: "Query ${feature} entities with filtering and pagination"
      - name: "Update"
        signature: "Update(ctx context.Context, entity *model.${Model}) error"
        purpose: "Modify existing ${feature} entity"
      - name: "Delete"
        signature: "Delete(ctx context.Context, id uint) error"
        purpose: "Remove ${feature} entity (soft delete)"
    
  service_interface:
    file: "internal/interface/service/${feature}.go"
    interface: "${Feature}Service"
    purpose: "Contract for ${feature} business logic operations"
    methods:
      - name: "Create${Model}"
        signature: "Create${Model}(ctx context.Context, req *request.Create${Model}Request) (*response.${Model}Response, error)"
        purpose: "Create new ${feature} with business rule validation"
      - name: "Get${Model}"
        signature: "Get${Model}(ctx context.Context, id uint) (*response.${Model}Response, error)"
        purpose: "Retrieve ${feature} with permission checks"
      - name: "List${Model}s"
        signature: "List${Model}s(ctx context.Context, req *request.List${Model}Request) (*response.List${Model}Response, error)"
        purpose: "Query ${feature} list with filtering"
      - name: "Update${Model}"
        signature: "Update${Model}(ctx context.Context, id uint, req *request.Update${Model}Request) (*response.${Model}Response, error)"
        purpose: "Update ${feature} with business validation"
      - name: "Delete${Model}"
        signature: "Delete${Model}(ctx context.Context, id uint) error"
        purpose: "Delete ${feature} with cascade handling"

  dependency_injection:
    pattern: "Manual dependency injection in main.go"
    constructor_naming: "New${Feature}Repository, New${Feature}Service, New${Feature}Controller"
    interface_usage: "Pass interfaces as constructor parameters"
```

### 5. Repository Layer
```yaml
repository:
  file: "internal/repository/${feature}.go"
  interface_file: "internal/interface/repository/${feature}.go"
  interface: "${Feature}Repository"
  purpose: "Data access abstraction for ${feature} operations"
  methods:
    - name: "Create"
      purpose: "Persist new ${feature} entity"
      signature: "Create(ctx context.Context, entity *model.${Model}) error"
      business_logic: "Validation and constraints applied"
      error_handling: "Handle database errors and constraint violations"
    
    - name: "GetByID"
      purpose: "Retrieve ${feature} by unique identifier"
      signature: "GetByID(ctx context.Context, id uint) (*model.${Model}, error)"
      error_handling: "Handle record not found cases"
    
    - name: "List"
      purpose: "Query ${feature} entities with filtering and pagination"
      signature: "List(ctx context.Context, filter *ListFilter) ([]*model.${Model}, int64, error)"
      pagination: "Use LIMIT and OFFSET for pagination"
      
    - name: "Update"
      purpose: "Modify existing ${feature} entity"
      signature: "Update(ctx context.Context, entity *model.${Model}) error"
      concurrency: "Handle concurrent updates with optimistic locking"
      
    - name: "Delete"
      purpose: "Remove ${feature} entity (soft delete)"
      signature: "Delete(ctx context.Context, id uint) error"
      soft_delete: "Use GORM's DeletedAt field"

  implementation_notes:
    - "All methods must use WithContext(ctx)"
    - "Use consistent error wrapping"
    - "Log critical operations"
```
### 6. Service Layer
```yaml
service:
  file: "internal/service/${feature}.go"
  interface_file: "internal/interface/service/${feature}.go"
  interface: "${Feature}Service"
  purpose: "Business logic orchestration for ${feature} operations"
  dependencies:
    - "${feature}Repository"
    - "logger.Logger"
    - "Other required services"
  
  methods:
    - name: "Create${Model}"
      purpose: "Create new ${feature} with business rule validation"
      signature: "Create${Model}(ctx context.Context, req *request.Create${Model}Request) (*response.${Model}Response, error)"
      business_rules:
        - "Specific validation rules"
        - "Business constraint checks"
        - "Permission validation"
      steps:
        - "Validate input parameters"
        - "Check business rules"
        - "Save via repository"
        - "Return formatted response"
      error_handling: "Use errors.NewAppError for wrapping"
    
    - name: "Get${Model}"
      purpose: "Retrieve ${feature} with permission checks"
      signature: "Get${Model}(ctx context.Context, id uint) (*response.${Model}Response, error)"
      authorization: "User access control logic"
      steps:
        - "Query by ID"
        - "Verify user permissions"
        - "Return formatted response"
    
    - name: "List${Model}s"
      purpose: "Query ${feature} list with filtering"
      signature: "List${Model}s(ctx context.Context, req *request.List${Model}Request) (*response.List${Model}Response, error)"
      steps:
        - "Apply user-specific filters"
        - "Paginate results"
        - "Return formatted response"
    
    - name: "Update${Model}"
      purpose: "Update ${feature} with business validation"
      signature: "Update${Model}(ctx context.Context, id uint, req *request.Update${Model}Request) (*response.${Model}Response, error)"
      validation: "Check update permissions and business rules"
    
    - name: "Delete${Model}"
      purpose: "Delete ${feature} with cascade handling"
      signature: "Delete${Model}(ctx context.Context, id uint) error"
      cascade_logic: "Handle related data cleanup"

  logging_strategy:
    - "Log all business operation starts and results"
    - "Use structured logging for key business parameters"
    - "Error logs include sufficient context"
```

### 7. Controller Layer
```yaml
controller:
  file: "internal/controller/${feature}.go"
  purpose: "HTTP request handling for ${feature} operations"
  dependencies:
    - "${feature}Service"
    - "logger.Logger"
  
  common_methods:
    - name: "bindAndValidateRequest"
      purpose: "Generic request binding and validation"
      reuse: "Used by all handlers"
  
  endpoints:
    - method: "POST"
      path: "/api/v1/${resource}"
      handler: "Create${Model}"
      purpose: "Create new ${feature}"
      auth: true
      request: "request.Create${Model}Request"
      response: "response.${Model}Response"
      status_codes: [201, 400, 401, 500]
      business_flow: "User submits ${feature} creation request"
      error_handling: "Use errors.HandleError for all errors"
    
    - method: "GET"
      path: "/api/v1/${resource}/:id"
      handler: "Get${Model}"
      purpose: "Retrieve specific ${feature}"
      auth: true
      params: ["id (path parameter)"]
      response: "response.${Model}Response"
      status_codes: [200, 401, 404, 500]
      validation: "Validate ID parameter format"
    
    - method: "GET"
      path: "/api/v1/${resource}"
      handler: "List${Model}s"
      purpose: "Query ${feature} list"
      auth: true
      query_params: ["page", "size", "filter"]
      response: "response.List${Model}Response"
      status_codes: [200, 401, 500]
      pagination: "Support pagination parameters"
    
    - method: "PUT"
      path: "/api/v1/${resource}/:id"
      handler: "Update${Model}"
      purpose: "Modify existing ${feature}"
      auth: true
      request: "request.Update${Model}Request"
      response: "response.${Model}Response"
      status_codes: [200, 400, 401, 404, 500]
      concurrency: "Handle concurrent update conflicts"
    
    - method: "DELETE"
      path: "/api/v1/${resource}/:id"
      handler: "Delete${Model}"
      purpose: "Remove ${feature}"
      auth: true
      status_codes: [204, 401, 404, 500]
      confirmation: "May require confirmation parameter"

  middleware_usage:
    - "Use unified auth middleware"
    - "Apply error handling middleware"
    - "Use logging middleware for request tracking"
    - "Apply i18n middleware"
```
### 8. Router Integration
```yaml
router_configuration:
  file: "internal/router/router.go"
  integration_point: "SetupRouter function"
  controller_struct: "Controllers"
  
  route_group_setup:
    path: "/api/v1/${resource}"
    middleware: ["Auth middleware for protected routes"]
    methods:
      - route: "POST /"
        handler: "controllers.${Feature}Controller.Create${Model}"
        purpose: "Create new ${feature}"
        auth_required: true
      - route: "GET /:id"
        handler: "controllers.${Feature}Controller.Get${Model}"
        purpose: "Get ${feature} by ID"
        auth_required: true
      - route: "GET /"
        handler: "controllers.${Feature}Controller.List${Model}s"
        purpose: "List ${feature}s with pagination"
        auth_required: true
      - route: "PUT /:id"
        handler: "controllers.${Feature}Controller.Update${Model}"
        purpose: "Update existing ${feature}"
        auth_required: true
      - route: "DELETE /:id"
        handler: "controllers.${Feature}Controller.Delete${Model}"
        purpose: "Delete ${feature}"
        auth_required: true

  controller_struct_update:
    location: "router.Controllers struct"
    addition: "${Feature}Controller *controller.${Feature}Controller"
    
  route_registration_example: |
    // ${Feature} management routes (protected)
    ${feature}Group := protected.Group("/${resource}")
    ${feature}Group.POST("", controllers.${Feature}Controller.Create${Model})
    ${feature}Group.GET("/:id", controllers.${Feature}Controller.Get${Model})
    ${feature}Group.GET("", controllers.${Feature}Controller.List${Model}s)
    ${feature}Group.PUT("/:id", controllers.${Feature}Controller.Update${Model})
    ${feature}Group.DELETE("/:id", controllers.${Feature}Controller.Delete${Model})
```

### 9. Data Transfer Objects
```yaml
dto_structures:
  request_file: "internal/dto/request/${feature}.go"
  response_file: "internal/dto/response/${feature}.go"
  
  requests:
    Create${Model}Request:
      purpose: "Input for ${feature} creation"
      validation_strategy: "Use Gin binding tags"
      fields:
        - name: "${field_name}"
          type: "string"
          business_meaning: "What this input represents"
          validation: "required,min=1,max=100"
          json: "${field_name}"
          i18n_key: "${feature}.${field_name}.required"
    
    Update${Model}Request:
      purpose: "Input for ${feature} modification"
      partial_update: "Support partial field updates"
      fields:
        - name: "${field_name}"
          type: "string"
          validation: "omitempty,min=1,max=100"
          json: "${field_name}"
          update_behavior: "Only update when provided"
    
    List${Model}Request:
      purpose: "Query parameters for ${feature} listing"
      embedding: "Embed common pagination structure"
      fields:
        - name: "Page"
          type: "int"
          default: 1
          json: "page"
          validation: "min=1"
        - name: "Size"
          type: "int"
          default: 20
          json: "size"
          validation: "min=1,max=100"
        - name: "Search"
          type: "string"
          json: "search"
          purpose: "Fuzzy search keywords"
  
  responses:
    ${Model}Response:
      purpose: "Public representation of ${feature}"
      data_mapping: "Converted from model.${Model}"
      fields:
        - name: "ID"
          type: "uint"
          json: "id"
          source: "model.ID"
        - name: "${field_name}"
          type: "string"
          json: "${field_name}"
          source: "model.${FieldName}"
        - name: "CreatedAt"
          type: "string"
          json: "created_at"
          format: "time.RFC3339"
          source: "model.CreatedAt"
        - name: "UpdatedAt"
          type: "string"
          json: "updated_at"
          format: "time.RFC3339"
          source: "model.UpdatedAt"
    
    List${Model}Response:
      purpose: "Paginated ${feature} collection"
      structure: "Standard list response format"
      fields:
        - name: "Items"
          type: "[]*${Model}Response"
          json: "items"
          source: "Converted model list"
        - name: "Pagination"
          type: "*response.PaginationInfo"
          json: "pagination"
          structure: "Standard pagination info"

  conversion_methods:
    - "To${Model}Response(model *model.${Model}) *${Model}Response"
    - "To${Model}ResponseList(models []*model.${Model}) []*${Model}Response"
    - "FromCreate${Model}Request(req *Create${Model}Request) *model.${Model}"
```

### 10. Implementation Task Checklist
```yaml
tasks:
  phase_1_data_layer:
    - task: "Create ${Model} data model"
      file: "internal/model/${feature}.go"
      description: "Define struct with GORM tags and relationships"
      details:
        - "Define all fields with GORM tags"
        - "Implement business logic methods"
        - "Define associations"
        - "Add index definitions"
    
    - task: "Add database migration in main.go"
      file: "main.go"
      description: "Register new model in AutoMigrate"
      location: "Database migration section"
  
  phase_2_interface_definitions:
    - task: "Define ${Feature}Repository interface"
      file: "internal/interface/repository/${feature}.go"
      description: "Create repository interface contract"
      details:
        - "Define all CRUD method signatures"
        - "Include proper context and error handling"
        - "Add documentation for each method"
    
    - task: "Define ${Feature}Service interface"
      file: "internal/interface/service/${feature}.go" 
      description: "Create service interface contract"
      details:
        - "Define all business method signatures"
        - "Include request/response DTO types"
        - "Add comprehensive documentation"
  
  phase_3_repository_layer:
    - task: "Define ${Feature}Repository interface"
      file: "internal/interface/repository/${feature}.go"
      description: "Create repository interface with CRUD methods"
      methods: ["Create", "GetByID", "List", "Update", "Delete"]
    
    - task: "Implement ${Feature}Repository"
      file: "internal/repository/${feature}.go"
      description: "Implement repository interface using GORM"
      details:
        - "Implement all interface methods"
        - "Add proper error handling and logging"
        - "Use WithContext for all database operations"
        - "Include constructor function: New${Feature}Repository"
  
  phase_4_service_layer:
    - task: "Define ${Feature}Service interface"
      file: "internal/interface/service/${feature}.go"
      description: "Create service interface with business methods"
      methods: ["Create${Model}", "Get${Model}", "List${Model}s", "Update${Model}", "Delete${Model}"]
    
    - task: "Implement ${Feature}Service"
      file: "internal/service/${feature}.go"
      description: "Implement service interface with business logic"
      details:
        - "Implement all interface methods"
        - "Add business rule validation"
        - "Handle DTO conversions"
        - "Add structured logging"
        - "Use unified error handling"
        - "Include constructor function: New${Feature}Service"
  
  phase_5_controller_layer:
    - task: "Create request DTOs"
      file: "internal/dto/request/${feature}.go"
      description: "Define API request data structures"
      validation: "Add parameter validation tags"
    
    - task: "Create response DTOs"
      file: "internal/dto/response/${feature}.go"
      description: "Define API response data structures"
      conversion: "Implement model to response conversion methods"
    
    - task: "Implement ${Feature}Controller"
      file: "internal/controller/${feature}.go"
      description: "Create Gin handlers with proper error handling"
      details:
        - "Implement all HTTP handlers"
        - "Use common bindAndValidateRequest method"
        - "Apply unified response format"
        - "Add appropriate HTTP status codes"
        - "Include constructor function: New${Feature}Controller"
  
  phase_6_router_integration:
    - task: "Update Controllers struct"
      file: "internal/router/router.go"
      description: "Add ${Feature}Controller to Controllers struct"
      location: "Controllers struct definition"
    
    - task: "Register ${feature} routes"
      file: "internal/router/router.go"
      description: "Add ${feature} route group in SetupRouter function"
      details:
        - "Create route group for ${resource}"
        - "Add all CRUD endpoints"
        - "Apply authentication middleware"
        - "Follow existing routing patterns"
  
  phase_7_dependency_injection:
  phase_7_dependency_injection:
    - task: "Initialize ${Feature}Repository"
      file: "main.go"
      description: "Add repository initialization in main function"
      location: "Repository initialization section"
      code_example: "${feature}Repo := repository.New${Feature}Repository(db)"
    
    - task: "Initialize ${Feature}Service"
      file: "main.go"
      description: "Add service initialization with dependencies"
      location: "Service initialization section"
      code_example: "${feature}Service := service.New${Feature}Service(${feature}Repo, jwtAuth, zapLogger)"
    
    - task: "Initialize ${Feature}Controller"
      file: "main.go"
      description: "Add controller initialization"
      location: "Controller initialization section"
      code_example: "${feature}Controller := controller.New${Feature}Controller(${feature}Service, zapLogger)"
    
    - task: "Update Controllers struct initialization"
      file: "main.go"
      description: "Add ${Feature}Controller to router Controllers"
      location: "Router setup section"
      code_example: |
        r := router.SetupRouter(&router.Controllers{
          UserController: userController,
          I18nController: i18nController,
          ${Feature}Controller: ${feature}Controller,
        }, cfg, zapLogger)
  
  phase_8_internationalization:
    - task: "Add Chinese translations"
      file: "pkg/i18n/locales/zh-CN.yaml"
      description: "Add Chinese translations for ${feature} related messages"
    
    - task: "Add English translations"
      file: "pkg/i18n/locales/en-US.yaml"
      description: "Add English translations for ${feature} related messages"
  
  phase_9_testing:
    - task: "Write Repository unit tests"
      file: "internal/repository/${feature}_test.go"
      description: "Test all data access methods"
      coverage: "Cover normal and error cases"
    
    - task: "Write Service unit tests"
      file: "internal/service/${feature}_test.go"
      description: "Test business logic with mocks"
      mocking: "Mock repository dependencies"
    
    - task: "Write Controller integration tests"
      file: "internal/controller/${feature}_test.go"
      description: "Test complete HTTP workflows"
      setup: "Setup test database and routes"
    
    - task: "Write API integration tests"
      file: "test/integration/${feature}_test.go"
      description: "End-to-end API testing"
      scenarios: "Test complete user scenarios"
```

### 11. Quality Acceptance Criteria
```yaml
acceptance_criteria:
  functionality:
    - "All CRUD operations work correctly"
    - "Business rules are properly enforced"
    - "Error handling covers edge cases"
    - "API responses match specifications"
    - "Pagination and search functionality works"
    - "Soft delete mechanism implemented correctly"
  
  code_quality:
    - "All functions have proper error handling"
    - "Code follows project conventions"
    - "Logging is implemented appropriately"
    - "Input validation prevents invalid data"
    - "Unified response format is used"
    - "Context is properly passed through all layers"
  
  testing:
    - "Unit test coverage > 80%"
    - "Integration tests cover main workflows"
    - "All tests pass successfully"
    - "Error scenarios are tested"
    - "Performance-critical paths have benchmarks"
  
  security:
    - "All endpoints require authentication"
    - "Authorization rules are enforced"
    - "Input sanitization prevents injection attacks"
    - "Sensitive data is not logged"
    - "HTTPS and security headers are used"
  
  documentation:
    - "API endpoints have proper comments"
    - "Business logic is clearly documented"
    - "Error codes and messages are documented"
    - "Internationalization messages are complete"
  
  performance:
    - "Database queries are optimized"
    - "Proper indexes are established"
    - "Pagination prevents large result sets"
    - "Caching strategy (if applicable)"

deployment_checklist:
  database:
    - "Migration scripts tested"
    - "Index creation verified"
    - "Data backup plan"
  
  configuration:
    - "Environment variables documented"
    - "Configuration validation"
    - "Reasonable default values"
  
  monitoring:
    - "Health check endpoints"
    - "Key metrics monitoring"
    - "Error rate alerting"
```

## Usage Instructions

Follow these steps when generating implementation plans:

### Preparation Phase
1. **Analyze design documents** - First understand requirements and business context
2. **Identify data dependencies** - Determine relationships with existing models
3. **Extract API contracts** - Get specific endpoint definitions from design docs
4. **Assess technical constraints** - Consider existing architecture and performance requirements

### Plan Generation Process
1. **Replace template variables** - Substitute ${feature}, ${Model}, etc. with actual values
2. **Specify implementation details** - Provide concrete field names, validation rules, business logic
3. **Ensure consistency** - Follow existing project naming and structure conventions
4. **Completeness check** - Ensure all required components and tests are covered

### Output Requirements
- Generated plan should be detailed enough for autonomous implementation
- Include specific file paths, method signatures, and implementation details
- Provide clear task breakdown and acceptance criteria
- Consider internationalization, testing, and documentation requirements

### Reference Implementation
Refer to existing user management feature as a pattern:
- Interfaces: `internal/interface/repository/user.go`, `internal/interface/service/user.go`
- Model: `internal/model/user.go`
- Repository: `internal/repository/user.go` (implements UserRepository interface)
- Service: `internal/service/user.go` (implements UserService interface)
- Controller: `internal/controller/user.go`
- DTOs: `internal/dto/request/user.go`, `internal/dto/response/user.go`
- Router integration: `internal/router/router.go` (Controllers struct and route registration)
- Dependency injection: `main.go` (manual initialization and wiring)

### Key Implementation Notes
- **Interface First**: Always define interfaces before implementations
- **Constructor Pattern**: Use New${Component}Name functions for dependency injection  
- **Layer Dependencies**: Controller depends on Service interface, Service depends on Repository interface
- **Router Integration**: Update both Controllers struct and route registration in SetupRouter
- **Manual DI**: Follow the existing pattern in main.go for dependency initialization

Maintain consistency with existing code style and architectural patterns.