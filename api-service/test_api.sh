#!/bin/bash

echo "Starting Websoft9 Web Service..."
./api-service &
SERVER_PID=$!

# 等待服务启动
sleep 3

echo "Testing health endpoint..."
curl -s http://localhost:8080/health

echo -e "\n\nTesting user registration..."
curl -s -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"password123"}'

echo -e "\n\nTesting user login..."
curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}'

echo -e "\n\nStopping server..."
kill $SERVER_PID 2>/dev/null || true
wait $SERVER_PID 2>/dev/null || true

echo "Test completed!"