#!/bin/bash

# MCP Configuration Demo Script
# This script demonstrates the MCP configuration functionality

set -e

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 Websoft9 MCP Configuration Demo${NC}"
echo "========================================="

# 1. Show the generated MCP configuration
echo -e "\n${GREEN}📋 Generated MCP Configuration:${NC}"
if [ -f "api-service/configs/mcp.json" ]; then
    echo "File: api-service/configs/mcp.json"
    cat api-service/configs/mcp.json | head -20
    echo "..."
else
    echo "MCP configuration not found. Run './scripts/setup-mcp.sh' first."
    exit 1
fi

# 2. Show the VSCode configuration
echo -e "\n${GREEN}🔧 VSCode Configuration:${NC}"
if [ -f ".vscode/settings.json" ]; then
    echo "File: .vscode/settings.json"
    cat .vscode/settings.json
else
    echo "VSCode configuration not found."
fi

# 3. Show available MCP servers in YAML configuration
echo -e "\n${GREEN}⚙️ API Service MCP Configuration (YAML):${NC}"
echo "File: api-service/configs/config.yaml"
grep -A 20 "^mcp:" api-service/configs/config.yaml || echo "MCP section not found in config.yaml"

# 4. Show agent MCP configuration
echo -e "\n${GREEN}🤖 Agent MCP Configuration (YAML):${NC}"
echo "File: websoft9-agent/configs/agent.yaml"
grep -A 10 "^mcp:" websoft9-agent/configs/agent.yaml || echo "MCP section not found in agent.yaml"

# 5. Test the build
echo -e "\n${GREEN}🏗️ Testing Build:${NC}"
cd api-service
echo "Building API service..."
if make build > /dev/null 2>&1; then
    echo "✅ API service builds successfully"
else
    echo "❌ API service build failed"
fi

cd ../websoft9-agent
echo "Building agent..."
if make build > /dev/null 2>&1; then
    echo "✅ Agent builds successfully"
else
    echo "❌ Agent build failed"
fi

cd ..

# 6. Run tests
echo -e "\n${GREEN}🧪 Running Tests:${NC}"
cd api-service
if go test ./pkg/mcp/... > /dev/null 2>&1; then
    echo "✅ MCP package tests pass"
else
    echo "❌ MCP package tests failed"
fi

cd ..

echo -e "\n${GREEN}🎉 Demo Complete!${NC}"
echo ""
echo "Next steps to use MCP:"
echo "1. Open this project in VSCode"
echo "2. Install GitHub Copilot extension"
echo "3. Start using @workspace in Copilot Chat for enhanced context"
echo "4. Start the API service: cd api-service && make run"
echo "5. Use the MCP management APIs at http://localhost:8080/api/v1/mcp/"