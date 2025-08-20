# Websoft9 用户管理API测试脚本
# 测试用户管理的所有接口功能

$BASE_URL = "http://localhost:8080/api/v1"
$ADMIN_TOKEN = ""

Write-Host "=== Websoft9 用户管理API测试 ===" -ForegroundColor Green

# 1. 用户注册
Write-Host "`n1. 测试用户注册" -ForegroundColor Yellow
$registerBody = @{
    username = "testuser"
    email = "test@example.com"
    password = "TestPass123"
} | ConvertTo-Json

try {
    $registerResponse = Invoke-RestMethod -Uri "$BASE_URL/auth/register" -Method Post -Body $registerBody -ContentType "application/json"
    Write-Host "注册响应: $($registerResponse | ConvertTo-Json)"
} catch {
    Write-Host "注册失败: $($_.Exception.Message)" -ForegroundColor Red
}

# 2. 用户登录
Write-Host "`n2. 测试用户登录" -ForegroundColor Yellow
$loginBody = @{
    username = "testuser"
    password = "TestPass123"
} | ConvertTo-Json

try {
    $loginResponse = Invoke-RestMethod -Uri "$BASE_URL/auth/login" -Method Post -Body $loginBody -ContentType "application/json"
    Write-Host "登录响应: $($loginResponse | ConvertTo-Json)"
    
    $TOKEN = $loginResponse.data.token
    Write-Host "获取到Token: $TOKEN"
    
    if (-not $TOKEN) {
        Write-Host "登录失败，无法获取Token" -ForegroundColor Red
        exit 1
    }
} catch {
    Write-Host "登录失败: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

# 3. 获取用户资料
Write-Host "`n3. 测试获取用户资料" -ForegroundColor Yellow
$headers = @{
    "Authorization" = "Bearer $TOKEN"
}

try {
    $profileResponse = Invoke-RestMethod -Uri "$BASE_URL/users/profile" -Method Get -Headers $headers
    Write-Host "用户资料响应: $($profileResponse | ConvertTo-Json)"
} catch {
    Write-Host "获取用户资料失败: $($_.Exception.Message)" -ForegroundColor Red
}

# 4. 创建用户（需要管理员权限）
Write-Host "`n4. 测试创建用户" -ForegroundColor Yellow
$createUserBody = @{
    group_id = 1
    username = "newuser"
    email = "newuser@example.com"
    password = "NewPass123"
    nickname = "新用户"
    phone = "13800138001"
    gender = 1
    signature = "这是我的个性签名"
    timezone = "Asia/Shanghai"
    language = "zh-CN"
    role_ids = @(2)
    status = 1
} | ConvertTo-Json

try {
    $createUserResponse = Invoke-RestMethod -Uri "$BASE_URL/users" -Method Post -Body $createUserBody -ContentType "application/json" -Headers $headers
    Write-Host "创建用户响应: $($createUserResponse | ConvertTo-Json)"
} catch {
    Write-Host "创建用户失败: $($_.Exception.Message)" -ForegroundColor Red
}

# 5. 获取用户列表
Write-Host "`n5. 测试获取用户列表" -ForegroundColor Yellow
try {
    $userListResponse = Invoke-RestMethod -Uri "$BASE_URL/users?page=1&page_size=10&keyword=test" -Method Get -Headers $headers
    Write-Host "用户列表响应: $($userListResponse | ConvertTo-Json -Depth 3)"
} catch {
    Write-Host "获取用户列表失败: $($_.Exception.Message)" -ForegroundColor Red
}

# 6. 获取用户详情
Write-Host "`n6. 测试获取用户详情" -ForegroundColor Yellow
try {
    $userDetailResponse = Invoke-RestMethod -Uri "$BASE_URL/users/1" -Method Get -Headers $headers
    Write-Host "用户详情响应: $($userDetailResponse | ConvertTo-Json -Depth 2)"
} catch {
    Write-Host "获取用户详情失败: $($_.Exception.Message)" -ForegroundColor Red
}

# 7. 更新用户信息
Write-Host "`n7. 测试更新用户信息" -ForegroundColor Yellow
$updateUserBody = @{
    nickname = "更新的昵称"
    signature = "更新的个性签名"
    status = 1
} | ConvertTo-Json

try {
    $updateUserResponse = Invoke-RestMethod -Uri "$BASE_URL/users/1" -Method Put -Body $updateUserBody -ContentType "application/json" -Headers $headers
    Write-Host "更新用户响应: $($updateUserResponse | ConvertTo-Json)"
} catch {
    Write-Host "更新用户失败: $($_.Exception.Message)" -ForegroundColor Red
}

# 8. 修改密码
Write-Host "`n8. 测试修改密码" -ForegroundColor Yellow
$changePasswordBody = @{
    old_password = "TestPass123"
    new_password = "NewTestPass123"
    confirm_password = "NewTestPass123"
} | ConvertTo-Json

try {
    $changePasswordResponse = Invoke-RestMethod -Uri "$BASE_URL/users/1/password" -Method Put -Body $changePasswordBody -ContentType "application/json" -Headers $headers
    Write-Host "修改密码响应: $($changePasswordResponse | ConvertTo-Json)"
} catch {
    Write-Host "修改密码失败: $($_.Exception.Message)" -ForegroundColor Red
}

# 9. 测试查询功能
Write-Host "`n9. 测试用户搜索和筛选" -ForegroundColor Yellow

# 按关键词搜索
try {
    $searchResponse = Invoke-RestMethod -Uri "$BASE_URL/users?keyword=test&page=1&page_size=5" -Method Get -Headers $headers
    Write-Host "关键词搜索响应: $($searchResponse | ConvertTo-Json -Depth 3)"
} catch {
    Write-Host "关键词搜索失败: $($_.Exception.Message)" -ForegroundColor Red
}

# 按状态筛选
try {
    $statusFilterResponse = Invoke-RestMethod -Uri "$BASE_URL/users?status=1&page=1&page_size=5" -Method Get -Headers $headers
    Write-Host "状态筛选响应: $($statusFilterResponse | ConvertTo-Json -Depth 3)"
} catch {
    Write-Host "状态筛选失败: $($_.Exception.Message)" -ForegroundColor Red
}

# 按角色筛选
try {
    $roleFilterResponse = Invoke-RestMethod -Uri "$BASE_URL/users?role_id=1&page=1&page_size=5" -Method Get -Headers $headers
    Write-Host "角色筛选响应: $($roleFilterResponse | ConvertTo-Json -Depth 3)"
} catch {
    Write-Host "角色筛选失败: $($_.Exception.Message)" -ForegroundColor Red
}

# 排序测试
try {
    $sortResponse = Invoke-RestMethod -Uri "$BASE_URL/users?sort=created_at&order=desc&page=1&page_size=5" -Method Get -Headers $headers
    Write-Host "排序响应: $($sortResponse | ConvertTo-Json -Depth 3)"
} catch {
    Write-Host "排序失败: $($_.Exception.Message)" -ForegroundColor Red
}

# 10. 删除用户（谨慎操作）
Write-Host "`n10. 测试删除用户" -ForegroundColor Yellow
# 注意：这里为了演示，我们删除一个不存在的用户ID，以避免真的删除用户
try {
    $deleteUserResponse = Invoke-RestMethod -Uri "$BASE_URL/users/999" -Method Delete -Headers $headers
    Write-Host "删除用户响应: $($deleteUserResponse | ConvertTo-Json)"
} catch {
    Write-Host "删除用户失败: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host "`n=== 用户管理API测试完成 ===" -ForegroundColor Green
Write-Host "注意：实际使用时请根据真实的用户ID和权限进行操作" -ForegroundColor Yellow
