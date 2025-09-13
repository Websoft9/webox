#!/bin/bash

# 添加辅助函数到结构体和测试函数之间
awk '/func TestGetUserProfile\(t \*testing.T\)/{
  print "// Helper function to convert boolean to string";
  print "func boolToString(b bool) string {";
  print "if b {";
  print "return \"true\"";
  print "}";
  print "return \"false\"";
  print "}";
  print "";
  print $0;
  next;
}
{print}' ./internal/service/user_profile_test.go > temp.go && mv temp.go ./internal/service/user_profile_test.go
