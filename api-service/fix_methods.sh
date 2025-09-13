#!/bin/bash

# 修复方法名
sed -i 's/func (m \*UserProfileMockLogger) UserProfileInfoContext/func (m \*UserProfileMockLogger) InfoContext/g' ./internal/service/user_profile_test.go
sed -i 's/func (m \*UserProfileMockLogger) UserProfileWarnContext/func (m \*UserProfileMockLogger) WarnContext/g' ./internal/service/user_profile_test.go
sed -i 's/func (m \*UserProfileMockLogger) UserProfileErrorContext/func (m \*UserProfileMockLogger) ErrorContext/g' ./internal/service/user_profile_test.go
sed -i 's/func (m \*UserProfileMockLogger) UserProfileDebug/func (m \*UserProfileMockLogger) Debug/g' ./internal/service/user_profile_test.go
sed -i 's/func (m \*UserProfileMockLogger) UserProfileInfo/func (m \*UserProfileMockLogger) Info/g' ./internal/service/user_profile_test.go
sed -i 's/func (m \*UserProfileMockLogger) UserProfileWarn/func (m \*UserProfileMockLogger) Warn/g' ./internal/service/user_profile_test.go
sed -i 's/func (m \*UserProfileMockLogger) UserProfileError/func (m \*UserProfileMockLogger) Error/g' ./internal/service/user_profile_test.go
sed -i 's/func (m \*UserProfileMockLogger) UserProfileFatal/func (m \*UserProfileMockLogger) Fatal/g' ./internal/service/user_profile_test.go

# 检查是否存在并添加DebugContext方法
if ! grep -q "func (m \*UserProfileMockLogger) DebugContext" ./internal/service/user_profile_test.go; then
  sed -i '/func (m \*UserProfileMockLogger) Debug/a\\
// DebugContext mocks the DebugContext method\\
func (m *UserProfileMockLogger) DebugContext(ctx context.Context, msg string, fields ...logger.Field) {\\
\tm.Called(ctx, msg, fields)\\
}' ./internal/service/user_profile_test.go
fi

# 添加IsDebugEnabled等方法
if ! grep -q "func (m \*UserProfileMockLogger) IsDebugEnabled" ./internal/service/user_profile_test.go; then
  sed -i '/func (m \*UserProfileMockLogger) Fatal/a\\
// IsDebugEnabled mocks the IsDebugEnabled method\\
func (m *UserProfileMockLogger) IsDebugEnabled() bool {\\
\treturn m.Called().Bool(0)\\
}\\
\\
// IsInfoEnabled mocks the IsInfoEnabled method\\
func (m *UserProfileMockLogger) IsInfoEnabled() bool {\\
\treturn m.Called().Bool(0)\\
}\\
\\
// IsWarnEnabled mocks the IsWarnEnabled method\\
func (m *UserProfileMockLogger) IsWarnEnabled() bool {\\
\treturn m.Called().Bool(0)\\
}\\
\\
// IsErrorEnabled mocks the IsErrorEnabled method\\
func (m *UserProfileMockLogger) IsErrorEnabled() bool {\\
\treturn m.Called().Bool(0)\\
}' ./internal/service/user_profile_test.go
fi

# 添加boolToString函数
if ! grep -q "func boolToString" ./internal/service/user_profile_test.go; then
  sed -i '/func (m \*UserProfileMockLogger) IsErrorEnabled/a\\
// Helper function to convert boolean to string\\
func boolToString(b bool) string {\\
\tif b {\\
\t\treturn "true"\\
\t}\\
\treturn "false"\\
}' ./internal/service/user_profile_test.go
fi

