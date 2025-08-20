#!/bin/bash

# Websoft9 API Service User Management Test Script
# 测试用户管理相关的API接口

# 设置基础URL
BASE_URL="http://localhost:8080/api/v1"
AUTH_TOKEN=""

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}=== Websoft9 用户管理API测试 ===${NC}"

# 测试健康检查
echo -e "\n${YELLOW}1. 测试健康检查${NC}"
response=$(curl -s -w "HTTPSTATUS:%{http_code}" -X GET "${BASE_URL%/api/v1}/health")
http_code=$(echo $response | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')
body=$(echo $response | sed -e 's/HTTPSTATUS\:.*//g')

if [ "$http_code" -eq 200 ]; then
    echo -e "${GREEN}✓ 健康检查通过${NC}"
    echo "$body" | jq .
else
    echo -e "${RED}✗ 健康检查失败 (HTTP: $http_code)${NC}"
    echo "$body"
fi

# 测试用户注册
echo -e "\n${YELLOW}2. 测试用户注册${NC}"
register_response=$(curl -s -w "HTTPSTATUS:%{http_code}" -X POST "${BASE_URL}/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser001",
    "email": "testuser001@example.com",
    "password": "password123"
  }')

http_code=$(echo $register_response | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')
body=$(echo $register_response | sed -e 's/HTTPSTATUS\:.*//g')

if [ "$http_code" -eq 200 ]; then
    echo -e "${GREEN}✓ 用户注册成功${NC}"
    echo "$body" | jq .
else
    echo -e "${RED}✗ 用户注册失败 (HTTP: $http_code)${NC}"
    echo "$body"
fi

# 测试用户登录
echo -e "\n${YELLOW}3. 测试用户登录${NC}"
login_response=$(curl -s -w "HTTPSTATUS:%{http_code}" -X POST "${BASE_URL}/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser001",
    "password": "password123"
  }')

http_code=$(echo $login_response | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')
body=$(echo $login_response | sed -e 's/HTTPSTATUS\:.*//g')

if [ "$http_code" -eq 200 ]; then
    echo -e "${GREEN}✓ 用户登录成功${NC}"
    echo "$body" | jq .
    # 提取token用于后续请求
    AUTH_TOKEN=$(echo "$body" | jq -r '.data.token')
    echo -e "${GREEN}Token: $AUTH_TOKEN${NC}"
else
    echo -e "${RED}✗ 用户登录失败 (HTTP: $http_code)${NC}"
    echo "$body"
    exit 1
fi

# 测试获取用户列表（需要认证）
echo -e "\n${YELLOW}4. 测试获取用户列表${NC}"
if [ -n "$AUTH_TOKEN" ]; then
    users_response=$(curl -s -w "HTTPSTATUS:%{http_code}" -X GET "${BASE_URL}/users?page=1&page_size=10" \
      -H "Authorization: Bearer $AUTH_TOKEN")
    
    http_code=$(echo $users_response | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')
    body=$(echo $users_response | sed -e 's/HTTPSTATUS\:.*//g')
    
    if [ "$http_code" -eq 200 ]; then
        echo -e "${GREEN}✓ 获取用户列表成功${NC}"
        echo "$body" | jq .
    else
        echo -e "${RED}✗ 获取用户列表失败 (HTTP: $http_code)${NC}"
        echo "$body"
    fi
else
    echo -e "${RED}✗ 无法获取用户列表，缺少认证token${NC}"
fi

# 测试创建用户（需要认证和权限）
echo -e "\n${YELLOW}5. 测试创建用户${NC}"
if [ -n "$AUTH_TOKEN" ]; then
    create_response=$(curl -s -w "HTTPSTATUS:%{http_code}" -X POST "${BASE_URL}/users" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $AUTH_TOKEN" \
      -d '{
        "username": "newuser002",
        "email": "newuser002@example.com",
        "password": "password123",
        "nickname": "新用户",
        "group_id": 1
      }')
    
    http_code=$(echo $create_response | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')
    body=$(echo $create_response | sed -e 's/HTTPSTATUS\:.*//g')
    
    if [ "$http_code" -eq 200 ]; then
        echo -e "${GREEN}✓ 创建用户成功${NC}"
        echo "$body" | jq .
        NEW_USER_ID=$(echo "$body" | jq -r '.data.id')
        echo -e "${GREEN}新用户ID: $NEW_USER_ID${NC}"
    else
        echo -e "${RED}✗ 创建用户失败 (HTTP: $http_code)${NC}"
        echo "$body"
    fi
else
    echo -e "${RED}✗ 无法创建用户，缺少认证token${NC}"
fi

# 测试获取用户详情
echo -e "\n${YELLOW}6. 测试获取用户详情${NC}"
if [ -n "$AUTH_TOKEN" ] && [ -n "$NEW_USER_ID" ] && [ "$NEW_USER_ID" != "null" ]; then
    detail_response=$(curl -s -w "HTTPSTATUS:%{http_code}" -X GET "${BASE_URL}/users/$NEW_USER_ID" \
      -H "Authorization: Bearer $AUTH_TOKEN")
    
    http_code=$(echo $detail_response | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')
    body=$(echo $detail_response | sed -e 's/HTTPSTATUS\:.*//g')
    
    if [ "$http_code" -eq 200 ]; then
        echo -e "${GREEN}✓ 获取用户详情成功${NC}"
        echo "$body" | jq .
    else
        echo -e "${RED}✗ 获取用户详情失败 (HTTP: $http_code)${NC}"
        echo "$body"
    fi
else
    echo -e "${YELLOW}跳过用户详情测试（用户创建失败或无token）${NC}"
fi

# 测试更新用户
echo -e "\n${YELLOW}7. 测试更新用户${NC}"
if [ -n "$AUTH_TOKEN" ] && [ -n "$NEW_USER_ID" ] && [ "$NEW_USER_ID" != "null" ]; then
    update_response=$(curl -s -w "HTTPSTATUS:%{http_code}" -X PUT "${BASE_URL}/users/$NEW_USER_ID" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $AUTH_TOKEN" \
      -d '{
        "nickname": "更新的昵称",
        "phone": "13888888888"
      }')
    
    http_code=$(echo $update_response | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')
    body=$(echo $update_response | sed -e 's/HTTPSTATUS\:.*//g')
    
    if [ "$http_code" -eq 200 ]; then
        echo -e "${GREEN}✓ 更新用户成功${NC}"
        echo "$body" | jq .
    else
        echo -e "${RED}✗ 更新用户失败 (HTTP: $http_code)${NC}"
        echo "$body"
    fi
else
    echo -e "${YELLOW}跳过用户更新测试（用户创建失败或无token）${NC}"
fi

# 测试修改密码
echo -e "\n${YELLOW}8. 测试修改密码${NC}"
if [ -n "$AUTH_TOKEN" ] && [ -n "$NEW_USER_ID" ] && [ "$NEW_USER_ID" != "null" ]; then
    password_response=$(curl -s -w "HTTPSTATUS:%{http_code}" -X PUT "${BASE_URL}/users/$NEW_USER_ID/password" \
      -H "Content-Type: application/json" \
      -H "Authorization: Bearer $AUTH_TOKEN" \
      -d '{
        "old_password": "password123",
        "new_password": "newpassword123"
      }')
    
    http_code=$(echo $password_response | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')
    body=$(echo $password_response | sed -e 's/HTTPSTATUS\:.*//g')
    
    if [ "$http_code" -eq 200 ]; then
        echo -e "${GREEN}✓ 修改密码成功${NC}"
        echo "$body" | jq .
    else
        echo -e "${RED}✗ 修改密码失败 (HTTP: $http_code)${NC}"
        echo "$body"
    fi
else
    echo -e "${YELLOW}跳过密码修改测试（用户创建失败或无token）${NC}"
fi

# 测试删除用户
echo -e "\n${YELLOW}9. 测试删除用户${NC}"
if [ -n "$AUTH_TOKEN" ] && [ -n "$NEW_USER_ID" ] && [ "$NEW_USER_ID" != "null" ]; then
    delete_response=$(curl -s -w "HTTPSTATUS:%{http_code}" -X DELETE "${BASE_URL}/users/$NEW_USER_ID" \
      -H "Authorization: Bearer $AUTH_TOKEN")
    
    http_code=$(echo $delete_response | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')
    body=$(echo $delete_response | sed -e 's/HTTPSTATUS\:.*//g')
    
    if [ "$http_code" -eq 200 ]; then
        echo -e "${GREEN}✓ 删除用户成功${NC}"
        echo "$body" | jq .
    else
        echo -e "${RED}✗ 删除用户失败 (HTTP: $http_code)${NC}"
        echo "$body"
    fi
else
    echo -e "${YELLOW}跳过用户删除测试（用户创建失败或无token）${NC}"
fi

echo -e "\n${YELLOW}=== 测试完成 ===${NC}"
