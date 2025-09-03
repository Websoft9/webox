package main

import (
"api-service/pkg/security"
"fmt"
)

func main() {
password := "password123"
hash, err := security.HashPassword(password)
if err != nil {
fmt.Printf("Error hashing password: %v\n", err)
return
}

fmt.Printf("Password: %s\n", password)
fmt.Printf("Hash: %s\n", hash)

// 验证密码
isValid := security.CheckPasswordHash(password, hash)
fmt.Printf("Password valid: %v\n", isValid)

// 验证错误密码
isInvalid := security.CheckPasswordHash("wrongpassword", hash)
fmt.Printf("Wrong password valid: %v\n", isInvalid)
}
