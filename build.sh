#!/bin/sh -e

cd "$(dirname "$0")"

# Color definitions
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
WHITE='\033[1;37m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# Function to print colored messages
print_info() {
    printf "${BLUE}${BOLD}[BUILD]${NC} ${CYAN}%s${NC}\n" "$1"
}

print_success() {
    printf "${GREEN}${BOLD}[✓]${NC} %s\n" "$1"
}

print_error() {
    printf "${RED}${BOLD}[✗]${NC} ${RED}%s${NC}\n" "$1"
}

print_header() {
    printf "\n${PURPLE}${BOLD}━━━ %s ━━━${NC}\n\n" "$1"
}

# Check if go is available
if ! command -v go > /dev/null 2>&1; then
    print_error "Go not found in PATH"
    printf "${RED}Please install Go from ${BOLD}https://go.dev${NC}\n"
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}')
printf "${GREEN}${BOLD}[✓]${NC} Found ${BOLD}%s${NC}\n" "$GO_VERSION"

# Download dependencies
print_info "Downloading dependencies..."
go mod download
print_success "Dependencies downloaded"

# Create bin directory if it doesn't exist
mkdir -p bin

print_header "Building Cathedral"

# Build cathedral CLI binary
print_info "Building cathedral binary..."
if go build -o bin/cathedral ./cmd/cathedral; then
    print_success "Built cathedral successfully"
else
    print_error "Failed to build cathedral"
    exit 1
fi

# Build cathedral-web binary
print_info "Building cathedral-web binary..."
if go build -o bin/cathedral-web ./cmd/cathedral-web; then
    print_success "Built cathedral-web successfully"
else
    print_error "Failed to build cathedral-web"
    exit 1
fi

printf "\n"
printf "${GREEN}${BOLD}Build Complete!${NC}\n"
printf "${GREEN}Binaries built at:${NC}\n"
printf "${GREEN}  • ${BOLD}./bin/cathedral${NC}\n"
printf "${GREEN}  • ${BOLD}./bin/cathedral-web${NC}\n"
