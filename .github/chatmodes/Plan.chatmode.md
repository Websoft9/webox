---
description: Generate AI coding implementation plans for new features
tools: ['codebase', 'usages', 'problems', 'changes', 'fetch', 'findTestFiles', 'githubRepo', 'editFiles', 'search', 'new']
model: Claude Sonnet 4
---

# Feature Implementation Planning

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

architecture_decisions:
  patterns: ["Clean Architecture layers", "Dependency injection", "RESTful API"]
  frameworks: ["Gin for HTTP", "GORM for ORM", "Wire for DI"]
  rationale: "Why these choices align with project standards"
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
    indexes:
      - fields: ["${field1}", "${field2}"]
        type: "unique|btree"
        purpose: "Query optimization reason"
```

### 4. Repository Layer
```yaml
repository:
  file: "internal/repository/${feature}_repository.go"
  interface: "${Feature}Repository"
  purpose: "Data access abstraction for ${feature} operations"
  methods:
    - name: "Create"
      purpose: "Persist new ${feature} entity"
      params: ["ctx context.Context", "entity *model.${Model}"]
      returns: ["*model.${Model}", "error"]
      business_logic: "Validation and constraints applied"
    
    - name: "GetByID"
      purpose: "Retrieve ${feature} by unique identifier"
      params: ["ctx context.Context", "id uint"]
      returns: ["*model.${Model}", "error"]
    
    - name: "List"
      purpose: "Query ${feature} entities with filtering and pagination"
      params: ["ctx context.Context", "filter *${Filter}"]
      returns: ["[]*model.${Model}", "int64", "error"]
      
    - name: "Update"
      purpose: "Modify existing ${feature} entity"
      params: ["ctx context.Context", "entity *model.${Model}"]
      returns: ["error"]
      
    - name: "Delete"
      purpose: "Remove ${feature} entity (soft delete)"
      params: ["ctx context.Context", "id uint"]
      returns: ["error"]
```
### 5. Service Layer
```yaml
service:
  file: "internal/service/${feature}_service.go"
  interface: "${Feature}Service"
  purpose: "Business logic orchestration for ${feature} operations"
  methods:
    - name: "Create${Model}"
      purpose: "Create new ${feature} with business rule validation"
      input: "Create${Model}Request"
      output: "${Model}Response"
      business_rules:
        - "Specific validation rule"
        - "Business constraint check"
      steps:
        - "Validate input parameters"
        - "Check business rules"
        - "Save via repository"
        - "Return formatted response"
    
    - name: "Get${Model}"
      purpose: "Retrieve ${feature} with permission checks"
      input: "uint (id)"
      output: "${Model}Response"
      authorization: "User access control logic"
      steps:
        - "Query by ID"
        - "Verify user permissions"
        - "Return formatted response"
    
    - name: "List${Model}s"
      purpose: "Query ${feature} list with filtering"
      input: "List${Model}Request"
      output: "List${Model}Response"
      steps:
        - "Apply user-specific filters"
        - "Paginate results"
        - "Return formatted response"
```

### 6. API Layer
```yaml
controller:
  file: "internal/controller/${feature}_controller.go"
  purpose: "HTTP request handling for ${feature} operations"
  endpoints:
    - method: "POST"
      path: "/api/v1/${resource}"
      handler: "Create${Model}"
      purpose: "Create new ${feature}"
      auth: true
      request: "Create${Model}Request"
      response: "${Model}Response"
      status: [201, 400, 401, 500]
      business_flow: "User submits ${feature} creation request"
    
    - method: "GET"
      path: "/api/v1/${resource}/:id"
      handler: "Get${Model}"
      purpose: "Retrieve specific ${feature}"
      auth: true
      response: "${Model}Response"
      status: [200, 401, 404, 500]
    
    - method: "GET"
      path: "/api/v1/${resource}"
      handler: "List${Model}s"
      purpose: "Query ${feature} list"
      auth: true
      params: ["page", "size", "filter"]
      response: "List${Model}Response"
      status: [200, 401, 500]
    
    - method: "PUT"
      path: "/api/v1/${resource}/:id"
      handler: "Update${Model}"
      purpose: "Modify existing ${feature}"
      auth: true
      request: "Update${Model}Request"
      response: "${Model}Response"
      status: [200, 400, 401, 404, 500]
    
    - method: "DELETE"
      path: "/api/v1/${resource}/:id"
      handler: "Delete${Model}"
      purpose: "Remove ${feature}"
      auth: true
      status: [204, 401, 404, 500]
```
### 7. Data Transfer Objects
```yaml
dto_structures:
  requests:
    Create${Model}Request:
      purpose: "Input for ${feature} creation"
      fields:
        - name: "${field_name}"
          type: "string"
          business_meaning: "What this input represents"
          validation: "required,min=1,max=100"
          json: "${field_name}"
    
    Update${Model}Request:
      purpose: "Input for ${feature} modification"
      fields:
        - name: "${field_name}"
          type: "string"
          validation: "omitempty,min=1,max=100"
          json: "${field_name}"
    
    List${Model}Request:
      purpose: "Query parameters for ${feature} listing"
      fields:
        - name: "Page"
          type: "int"
          default: 1
          json: "page"
        - name: "Size"
          type: "int"
          default: 20
          json: "size"
  
  responses:
    ${Model}Response:
      purpose: "Public representation of ${feature}"
      fields:
        - name: "ID"
          type: "uint"
          json: "id"
        - name: "${field_name}"
          type: "string"
          json: "${field_name}"
    
    List${Model}Response:
      purpose: "Paginated ${feature} collection"
      fields:
        - name: "Items"
          type: "[]*${Model}Response"
          json: "items"
        - name: "Pagination"
          type: "*response.PaginationInfo"
          json: "pagination"
```

### 8. Implementation Tasks
```yaml
tasks:
  phase_1_foundation:
    - task: "Create ${Model} data model"
      file: "internal/model/${feature}.go"
      description: "Define struct with GORM tags and relationships"
    
    - task: "Add database migration"
      file: "internal/database/migrate.go"
      description: "Register model for auto-migration"
  
  phase_2_data_access:
    - task: "Define ${Feature}Repository interface"
      file: "internal/repository/${feature}_repository.go"
      description: "Create repository interface with CRUD methods"
    
    - task: "Implement repository with GORM"
      file: "internal/repository/${feature}_repository.go"
      description: "Implement all interface methods with error handling"
  
  phase_3_business_logic:
    - task: "Define ${Feature}Service interface"
      file: "internal/service/${feature}_service.go"
      description: "Create service interface with business methods"
    
    - task: "Implement service with validation"
      file: "internal/service/${feature}_service.go"
      description: "Implement business logic with validation and authorization"
  
  phase_4_api_endpoints:
    - task: "Create request/response DTOs"
      file: "internal/controller/${feature}_controller.go"
      description: "Define API data structures"
    
    - task: "Implement HTTP handlers"
      file: "internal/controller/${feature}_controller.go"
      description: "Create Gin handlers with proper error handling"
  
  phase_5_integration:
    - task: "Register routes"
      file: "internal/router/router.go"
      description: "Add ${feature} route group"
    
    - task: "Update dependency injection"
      file: "internal/service/services.go"
      description: "Wire up repository and service providers"
  
  phase_6_testing:
    - task: "Write unit tests"
      files: ["*_test.go"]
      description: "Test all layers with mocks"
    
    - task: "Write integration tests"
      file: "test/integration/${feature}_test.go"
      description: "Test complete workflows"
```

### 9. Quality Validation
```yaml
acceptance_criteria:
  functionality:
    - "All CRUD operations work correctly"
    - "Business rules are enforced"
    - "Error handling covers edge cases"
    - "API responses match specifications"
  
  code_quality:
    - "All functions have proper error handling"
    - "Code follows project conventions"
    - "Logging is implemented appropriately"
    - "Input validation prevents invalid data"
  
  testing:
    - "Unit test coverage > 80%"
    - "Integration tests cover main workflows"
    - "All tests pass successfully"
    - "Error scenarios are tested"
  
  security:
    - "Authentication is required for all endpoints"
    - "Authorization rules are enforced"
    - "Input sanitization prevents injection attacks"
    - "Sensitive data is not logged"
```

## Instructions

Reference the workspace coding standards in `.github/copilot-instructions.md` for project-specific guidelines. If this file is not present, follow the documented coding conventions for Go, Clean Architecture, and RESTful API design as outlined in your project documentation.

When generating the implementation plan:

1. **Analyze design documents** first to understand requirements
2. **Extract concrete details** for data models, API contracts, and business rules
3. **Replace template variables** (${feature}, ${Model}, etc.) with actual values
4. **Create comprehensive plan** covering all implementation phases
5. **Save plan file** to `plans/` directory with descriptive filename

The generated plan should be detailed enough for another developer or AI to implement the feature autonomously.