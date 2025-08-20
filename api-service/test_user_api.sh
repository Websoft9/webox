#!/bin/bash

# Websoft9 用户管理API测试脚本
# 测试用户管理的所有接口功能

BASE_URL="http://localhost:8080/api/v1"
ADMIN_TOKEN=""

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}=== Websoft9 用户管理API测试 ===${NC}"

# 1. 用户注册
echo -e "\n${YELLOW}1. 测试用户注册${NC}"
REGISTER_RESPONSE=$(curl -s -X POST "${BASE_URL}/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "TestPass123"
  }')

echo "注册响应: $REGISTER_RESPONSE"

# 2. 用户登录
echo -e "\n${YELLOW}2. 测试用户登录${NC}"
LOGIN_RESPONSE=$(curl -s -X POST "${BASE_URL}/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "TestPass123"
  }')

echo "登录响应: $LOGIN_RESPONSE"

# 提取token
TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
echo "获取到Token: $TOKEN"

if [ -z "$TOKEN" ]; then
    echo -e "${RED}登录失败，无法获取Token${NC}"
    exit 1
fi

# 3. 获取用户资料
echo -e "\n${YELLOW}3. 测试获取用户资料${NC}"
PROFILE_RESPONSE=$(curl -s -X GET "${BASE_URL}/users/profile" \
  -H "Authorization: Bearer $TOKEN")

echo "用户资料响应: $PROFILE_RESPONSE"

# 4. 创建用户（需要管理员权限）
echo -e "\n${YELLOW}4. 测试创建用户${NC}"
CREATE_USER_RESPONSE=$(curl -s -X POST "${BASE_URL}/users" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "group_id": 1,
    "username": "newuser",
    "email": "newuser@example.com",
    "password": "NewPass123",
    "nickname": "新用户",
    "phone": "13800138001",
    "gender": 1,
    "signature": "这是我的个性签名",
    "timezone": "Asia/Shanghai",
    "language": "zh-CN",
    "role_ids": [2],
    "status": 1
  }')

echo "创建用户响应: $CREATE_USER_RESPONSE"

# 5. 获取用户列表
echo -e "\n${YELLOW}5. 测试获取用户列表${NC}"
USER_LIST_RESPONSE=$(curl -s -X GET "${BASE_URL}/users?page=1&page_size=10&keyword=test" \
  -H "Authorization: Bearer $TOKEN")

echo "用户列表响应: $USER_LIST_RESPONSE"

# 6. 获取用户详情
echo -e "\n${YELLOW}6. 测试获取用户详情${NC}"
USER_DETAIL_RESPONSE=$(curl -s -X GET "${BASE_URL}/users/1" \
  -H "Authorization: Bearer $TOKEN")

echo "用户详情响应: $USER_DETAIL_RESPONSE"

# 7. 更新用户信息
echo -e "\n${YELLOW}7. 测试更新用户信息${NC}"
UPDATE_USER_RESPONSE=$(curl -s -X PUT "${BASE_URL}/users/1" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "nickname": "更新的昵称",
    "signature": "更新的个性签名",
    "status": 1
  }')

echo "更新用户响应: $UPDATE_USER_RESPONSE"

# 8. 修改密码
echo -e "\n${YELLOW}8. 测试修改密码${NC}"
CHANGE_PASSWORD_RESPONSE=$(curl -s -X PUT "${BASE_URL}/users/1/password" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "old_password": "TestPass123",
    "new_password": "NewTestPass123",
    "confirm_password": "NewTestPass123"
  }')

echo "修改密码响应: $CHANGE_PASSWORD_RESPONSE"

# 9. 测试查询功能
echo -e "\n${YELLOW}9. 测试用户搜索和筛选${NC}"

# 按关键词搜索
SEARCH_RESPONSE=$(curl -s -X GET "${BASE_URL}/users?keyword=test&page=1&page_size=5" \
  -H "Authorization: Bearer $TOKEN")
echo "关键词搜索响应: $SEARCH_RESPONSE"

# 按状态筛选
STATUS_FILTER_RESPONSE=$(curl -s -X GET "${BASE_URL}/users?status=1&page=1&page_size=5" \
  -H "Authorization: Bearer $TOKEN")
echo "状态筛选响应: $STATUS_FILTER_RESPONSE"

# 按角色筛选
ROLE_FILTER_RESPONSE=$(curl -s -X GET "${BASE_URL}/users?role_id=1&page=1&page_size=5" \
  -H "Authorization: Bearer $TOKEN")
echo "角色筛选响应: $ROLE_FILTER_RESPONSE"

# 排序测试
SORT_RESPONSE=$(curl -s -X GET "${BASE_URL}/users?sort=created_at&order=desc&page=1&page_size=5" \
  -H "Authorization: Bearer $TOKEN")
echo "排序响应: $SORT_RESPONSE"

# 10. 删除用户（谨慎操作）
echo -e "\n${YELLOW}10. 测试删除用户${NC}"
# 注意：这里为了演示，我们删除一个不存在的用户ID，以避免真的删除用户
DELETE_USER_RESPONSE=$(curl -s -X DELETE "${BASE_URL}/users/999" \
  -H "Authorization: Bearer $TOKEN")

echo "删除用户响应: $DELETE_USER_RESPONSE"

echo -e "\n${GREEN}=== 用户管理API测试完成 ===${NC}"
echo -e "${YELLOW}注意：实际使用时请根据真实的用户ID和权限进行操作${NC}"
