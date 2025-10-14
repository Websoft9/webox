#!/bin/bash

# Websoft9 User Initialization Script
# Purpose: Create a new user in the Websoft9 system
# Usage: ./init_user.sh <username_or_email> <password> <role_code>

set -euo pipefail

# Script configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
API_SERVICE_DIR="$(dirname "$SCRIPT_DIR")"
CREATE_USER_TOOL_DIR="$API_SERVICE_DIR/scripts"
CREATE_USER_BINARY="$CREATE_USER_TOOL_DIR/create-user"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}" >&2
}

print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ $1${NC}"
}

# Function to show usage
show_usage() {
    cat << EOF
Usage: $0 <username_or_email> <password> <role_code>

Create a new user in the Websoft9 system.

Parameters:
  username_or_email  Username or email address for the new user
  password          Password for the new user (min 6 characters)
  role_code         Role code to assign to the user (e.g., admin, manager, developer)

Examples:
  $0 admin@example.com password123 admin
  $0 johndoe secretpass developer

Notes:
  - If username contains '@', it will be treated as an email address
  - If username doesn't contain '@', a default email will be generated
  - Role code must exist in the system and be active before creating the user
  - Password should be secure and meet system requirements

Common Role Codes:
  admin       - System Administrator
  manager     - Project Manager
  developer   - Developer
  viewer      - Viewer

To see available roles, check the database or use the web interface.
EOF
}

# Function to validate parameters
validate_parameters() {
    local username="$1"
    local password="$2"
    local role_code="$3"

    # Validate username/email
    if [[ -z "$username" ]]; then
        print_error "Username or email cannot be empty"
        return 1
    fi

    if [[ ${#username} -lt 3 ]]; then
        print_error "Username must be at least 3 characters long"
        return 1
    fi

    # If it contains @, validate as email
    if [[ "$username" == *"@"* ]]; then
        if [[ ! "$username" =~ ^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$ ]]; then
            print_error "Invalid email format"
            return 1
        fi
    fi

    # Validate password
    if [[ -z "$password" ]]; then
        print_error "Password cannot be empty"
        return 1
    fi

    if [[ ${#password} -lt 6 ]]; then
        print_error "Password must be at least 6 characters long"
        return 1
    fi

    # Validate role code
    if [[ -z "$role_code" ]]; then
        print_error "Role code cannot be empty"
        return 1
    fi

    return 0
}

# Function to build the create-user tool
build_create_user_tool() {
    print_info "Building create-user tool..."

    # Check if Go is available
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed or not in PATH"
        print_info "Please install Go 1.19 or later"
        return 1
    fi

    # Check if make is available
    if ! command -v make &> /dev/null; then
        print_error "Make is not installed or not in PATH"
        return 1
    fi

    # Build the tool using make
    cd "$API_SERVICE_DIR"
    if make build-create-user > /dev/null 2>&1; then
        print_success "Create-user tool built successfully"
        return 0
    else
        print_error "Failed to build create-user tool"
        print_info "You may need to run: make deps"
        return 1
    fi
}

# Function to check if binary exists and is executable
check_binary() {
    if [[ ! -f "$CREATE_USER_BINARY" ]]; then
        print_warning "Create-user binary not found, building..."
        if ! build_create_user_tool; then
            return 1
        fi
    elif [[ ! -x "$CREATE_USER_BINARY" ]]; then
        print_warning "Create-user binary exists but is not executable, rebuilding..."
        if ! build_create_user_tool; then
            return 1
        fi
    fi
    return 0
}

# Function to create user
create_user() {
    local username="$1"
    local password="$2"
    local role_code="$3"

    print_info "Creating user with username/email: $username"
    print_info "Assigning role code: $role_code"

    # Run the create-user tool
    if "$CREATE_USER_BINARY" -username="$username" -password="$password" -role-code="$role_code"; then
        print_success "User created successfully!"
        return 0
    else
        print_error "Failed to create user"
        return 1
    fi
}

# Main function
main() {
    print_info "Websoft9 User Initialization Script"
    print_info "======================================"

    # Check if help is requested
    if [[ $# -eq 0 ]] || [[ "$1" == "-h" ]] || [[ "$1" == "--help" ]]; then
        show_usage
        exit 0
    fi

    # Check parameter count
    if [[ $# -ne 3 ]]; then
        print_error "Invalid number of parameters"
        echo
        show_usage
        exit 1
    fi

    local username="$1"
    local password="$2"
    local role_code="$3"

    # Validate parameters
    if ! validate_parameters "$username" "$password" "$role_code"; then
        echo
        show_usage
        exit 1
    fi

    # Check and build binary if needed
    if ! check_binary; then
        exit 1
    fi

    # Create user
    if create_user "$username" "$password" "$role_code"; then
        echo
        print_success "User initialization completed successfully!"
        print_info "User '$username' has been created with role code '$role_code'"
        exit 0
    else
        echo
        print_error "User initialization failed!"
        print_info "Please check the error messages above and try again"
        exit 1
    fi
}

# Trap to handle script interruption
trap 'print_error "Script interrupted by user"; exit 130' INT TERM

# Run main function with all arguments
main "$@"
