#!/bin/bash

# setup.sh - Automated setup script for go-bn128
# This script creates the project structure and verifies the installation

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m' # No Color

# Project settings
PROJECT_NAME="go-bn128"
PROJECT_URL="github.com/zacksfF/go-bn128"

# Print colored output
print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

print_info() {
    echo -e "${BLUE}→${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

print_header() {
    echo -e "\n${BOLD}${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${BOLD}${CYAN}  $1${NC}"
    echo -e "${BOLD}${CYAN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}\n"
}

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Main setup function
main() {
    print_header "BN128 Elliptic Curve Library Setup"
    
    echo -e "${BOLD}This script will:${NC}"
    echo "  1. Check prerequisites (Go, Make)"
    echo "  2. Create project directory structure"
    echo "  3. Initialize Go module"
    echo "  4. Verify the installation"
    echo "  5. Run initial tests"
    echo ""
    
    read -p "Continue with setup? (y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_warning "Setup cancelled."
        exit 0
    fi
    
    check_prerequisites
    create_directory_structure
    initialize_module
    create_gitignore
    verify_installation
    run_initial_tests
    show_next_steps
}

# Check prerequisites
check_prerequisites() {
    print_header "Checking Prerequisites"
    
    # Check Go
    if command_exists go; then
        GO_VERSION=$(go version | awk '{print $3}')
        print_success "Go is installed: $GO_VERSION"
        
        # Check Go version (need 1.21+)
        GO_MINOR=$(echo $GO_VERSION | sed 's/go1\.\([0-9]*\).*/\1/')
        if [ "$GO_MINOR" -lt 21 ]; then
            print_warning "Go 1.21+ recommended. You have $GO_VERSION"
        fi
    else
        print_error "Go is not installed"
        echo "  Please install Go from: https://golang.org/dl/"
        exit 1
    fi
    
    # Check Make
    if command_exists make; then
        MAKE_VERSION=$(make --version | head -n1)
        print_success "Make is installed: $MAKE_VERSION"
    else
        print_error "Make is not installed"
        echo "  macOS: xcode-select --install"
        echo "  Linux: sudo apt-get install build-essential"
        exit 1
    fi
    
    # Check Git (optional but recommended)
    if command_exists git; then
        GIT_VERSION=$(git --version)
        print_success "Git is installed: $GIT_VERSION"
    else
        print_warning "Git is not installed (recommended)"
    fi
    
    echo ""
}

# Create directory structure
create_directory_structure() {
    print_header "Creating Project Structure"
    
    # Check if directory exists
    if [ -d "$PROJECT_NAME" ]; then
        print_warning "Directory '$PROJECT_NAME' already exists"
        read -p "  Remove and recreate? (y/n) " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            rm -rf "$PROJECT_NAME"
            print_info "Removed existing directory"
        else
            print_error "Setup cannot continue with existing directory"
            exit 1
        fi
    fi
    
    # Create main directory
    mkdir -p "$PROJECT_NAME"
    cd "$PROJECT_NAME"
    print_success "Created directory: $PROJECT_NAME"
    
    # Create subdirectories
    mkdir -p examples/{zksnark,bls,pairing}
    print_success "Created examples directory"
    
    mkdir -p docs
    print_success "Created docs directory"
    
    mkdir -p .github/workflows
    print_success "Created .github/workflows directory"
    
    echo ""
}

# Initialize Go module
initialize_module() {
    print_header "Initializing Go Module"
    
    print_info "Creating go.mod..."
    cat > go.mod << EOF
module $PROJECT_URL

go 1.21

// Pure Go implementation - no external dependencies
EOF
    print_success "Created go.mod"
    
    print_info "Initializing Go module..."
    go mod tidy
    print_success "Go module initialized"
    
    echo ""
}

# Create .gitignore
create_gitignore() {
    print_header "Creating Configuration Files"
    
    print_info "Creating .gitignore..."
    cat > .gitignore << 'EOF'
# Binaries
*.test
*.out
*.prof

# Coverage
coverage.out
coverage.html

# Benchmarks
bench.txt

# IDE
.vscode/
.idea/
*.swp
*.swo

# OS
.DS_Store
Thumbs.db
.directory

# Build artifacts
/bin/
/dist/

# Go
go.work
go.work.sum
EOF
    print_success "Created .gitignore"
    
    echo ""
}

# Verify installation
verify_installation() {
    print_header "Verifying Installation"
    
    # Check if required files exist
    REQUIRED_FILES=("bn128.go" "bn128_test.go" "bn128_bench_test.go" "Makefile")
    MISSING_FILES=()
    
    for file in "${REQUIRED_FILES[@]}"; do
        if [ ! -f "$file" ]; then
            MISSING_FILES+=("$file")
        fi
    done
    
    if [ ${#MISSING_FILES[@]} -gt 0 ]; then
        print_warning "Missing required files:"
        for file in "${MISSING_FILES[@]}"; do
            echo "    - $file"
        done
        echo ""
        print_info "Please copy the following files to the project directory:"
        echo "    1. bn128.go (core implementation)"
        echo "    2. bn128_test.go (test suite)"
        echo "    3. bn128_bench_test.go (benchmarks)"
        echo "    4. Makefile (build automation)"
        echo ""
        print_warning "Setup will continue, but you'll need these files to run tests"
        return 1
    else
        print_success "All required files present"
        
        # Verify files are not empty
        for file in "${REQUIRED_FILES[@]}"; do
            if [ ! -s "$file" ]; then
                print_warning "$file is empty"
            fi
        done
    fi
    
    echo ""
}

# Run initial tests
run_initial_tests() {
    print_header "Running Initial Tests"
    
    # Check if files exist and are not empty
    if [ ! -s "bn128.go" ] || [ ! -s "bn128_test.go" ]; then
        print_warning "Core files not ready. Skipping tests."
        return 0
    fi
    
    print_info "Running quick check..."
    if make quick 2>&1 | tee /tmp/bn128_test.log; then
        print_success "Quick check passed!"
    else
        print_warning "Quick check failed or not ready yet"
        echo "  This is normal if you haven't copied the code files yet"
    fi
    
    echo ""
}

# Show next steps
show_next_steps() {
    print_header "Setup Complete!"
    
    echo -e "${BOLD} Project location:${NC}"
    echo "  $(pwd)"
    echo ""
    
    echo -e "${BOLD} Next steps:${NC}"
    echo "  1. Copy the source files:"
    echo "     - bn128.go"
    echo "     - bn128_test.go"
    echo "     - bn128_bench_test.go"
    echo "     - Makefile"
    echo ""
    echo "  2. Run tests:"
    echo "     ${CYAN}make quick${NC}      # Quick validation (30s)"
    echo "     ${CYAN}make test${NC}       # Full test suite (2-3m)"
    echo "     ${CYAN}make coverage${NC}   # Generate coverage report"
    echo ""
    echo "  3. Run benchmarks:"
    echo "     ${CYAN}make bench${NC}      # All benchmarks"
    echo "     ${CYAN}make bench-g1${NC}   # G1 operations only"
    echo "     ${CYAN}make bench-pairing${NC} # Pairing operations only"
    echo ""
    echo "  4. Full verification:"
    echo "     ${CYAN}make verify${NC}     # Complete pipeline (5-15m)"
    echo ""
    
    echo -e "${BOLD} Useful commands:${NC}"
    echo "  ${CYAN}make help${NC}          # Show all available commands"
    echo "  ${CYAN}make stats${NC}         # Project statistics"
    echo "  ${CYAN}make clean${NC}         # Clean generated files"
    echo ""
    
    echo -e "${BOLD}📚 Documentation:${NC}"
    echo "  - README.md         # Main documentation"
    echo "  - QUICKSTART.md     # Getting started guide"
    echo "  - RUN.md            # Detailed run instructions"
    echo "  - STRUCTURE.md      # Project structure"
    echo ""
    
    print_success "Project setup complete!"
    echo ""
}

# Run main function
main