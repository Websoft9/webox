#!/bin/bash

# 添加缺少的import
sed -i '9a\ "api-service/internal/interface/repository"' ./internal/service/user_profile_test.go

# 替换所有MockLogger为UserProfileMockLogger
sed -i 's/func (m \*MockLogger)/func (m *UserProfileMockLogger)/g' ./internal/service/user_profile_test.go

# 检查修复后的文件
grep -n "UserProfileMockLogger" ./internal/service/user_profile_test.go
grep -n "repository" ./internal/service/user_profile_test.go
