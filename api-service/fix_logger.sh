#!/bin/bash

# 查找UserProfileMockLogger结构体定义的结束位置
line_num=$(grep -n "}" ./internal/service/user_profile_test.go | grep -A1 "type UserProfileMockLogger struct" | tail -1 | cut -d: -f1)

# 在结构体定义后添加缺失的方法
sed -i "${line_num}a\\
// DebugContext mocks the DebugContext method\\
func (m *UserProfileMockLogger) DebugContext(ctx context.Context, msg string, fields ...logger.Field) {\\
\\tm.Called(ctx, msg, fields)\\
}\\
\\
// IsDebugEnabled mocks the IsDebugEnabled method\\
func (m *UserProfileMockLogger) IsDebugEnabled() bool {\\
\\treturn m.Called().Bool(0)\\
}\\
\\
// IsInfoEnabled mocks the IsInfoEnabled method\\
func (m *UserProfileMockLogger) IsInfoEnabled() bool {\\
\\treturn m.Called().Bool(0)\\
}\\
\\
// IsWarnEnabled mocks the IsWarnEnabled method\\
func (m *UserProfileMockLogger) IsWarnEnabled() bool {\\
\\treturn m.Called().Bool(0)\\
}\\
\\
// IsErrorEnabled mocks the IsErrorEnabled method\\
func (m *UserProfileMockLogger) IsErrorEnabled() bool {\\
\\treturn m.Called().Bool(0)\\
}\\
" ./internal/service/user_profile_test.go

# 添加辅助函数boolToString
echo '
// Helper function to convert boolean to string
func boolToString(b bool) string {
if b {
return "true"
}
return "false"
}' >> ./internal/service/user_profile_test.go

# 确保在每个测试用例中添加对新方法的mock
sed -i 's/mockLogger.On("InfoContext", mock.Anything, mock.Anything, mock.Anything).Return()/mockLogger.On("InfoContext", mock.Anything, mock.Anything, mock.Anything).Return()\n\t\t\tmockLogger.On("DebugContext", mock.Anything, mock.Anything, mock.Anything).Return()\n\t\t\tmockLogger.On("IsDebugEnabled").Return(true)\n\t\t\tmockLogger.On("IsInfoEnabled").Return(true)\n\t\t\tmockLogger.On("IsWarnEnabled").Return(true)\n\t\t\tmockLogger.On("IsErrorEnabled").Return(true)/g' ./internal/service/user_profile_test.go
