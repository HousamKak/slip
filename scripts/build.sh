#!/bin/bash
# Build script for Slip - Cross-platform compilation
# Usage: ./scripts/build.sh [version]

set -e

# Configuration
VERSION="${1:-dev}"
BUILD_DIR="build"
BINARY_NAME="slip"
PACKAGE="slip/cmd/slip"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Print colored messages
info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Clean build directory
clean() {
    info "Cleaning build directory..."
    rm -rf "$BUILD_DIR"
    mkdir -p "$BUILD_DIR"
    success "Build directory cleaned"
}

# Get git commit hash (short)
get_commit_hash() {
    if git rev-parse --short HEAD >/dev/null 2>&1; then
        git rev-parse --short HEAD
    else
        echo "unknown"
    fi
}

# Get build timestamp
get_timestamp() {
    date -u '+%Y-%m-%d_%H:%M:%S_UTC'
}

# Build for a specific platform
build_platform() {
    local os=$1
    local arch=$2
    local output_name="${BINARY_NAME}"

    if [ "$os" = "windows" ]; then
        output_name="${BINARY_NAME}.exe"
    fi

    local output_path="${BUILD_DIR}/${BINARY_NAME}-${os}-${arch}"
    if [ "$os" = "windows" ]; then
        output_path="${output_path}.exe"
    fi

    info "Building for ${os}/${arch}..."

    # Build with version info and optimizations
    GOOS=$os GOARCH=$arch go build \
        -ldflags="-s -w -X main.Version=${VERSION} -X main.Commit=$(get_commit_hash) -X main.BuildTime=$(get_timestamp)" \
        -o "$output_path" \
        "$PACKAGE"

    if [ $? -eq 0 ]; then
        local size=$(du -h "$output_path" | cut -f1)
        success "Built ${os}/${arch} (${size})"

        # Create SHA256 checksum
        if command -v sha256sum >/dev/null 2>&1; then
            (cd "$BUILD_DIR" && sha256sum "$(basename "$output_path")" > "$(basename "$output_path").sha256")
        elif command -v shasum >/dev/null 2>&1; then
            (cd "$BUILD_DIR" && shasum -a 256 "$(basename "$output_path")" > "$(basename "$output_path").sha256")
        fi
    else
        error "Failed to build ${os}/${arch}"
        return 1
    fi
}

# Build all platforms
build_all() {
    info "Starting cross-platform build for version: ${VERSION}"
    echo ""

    # Linux
    build_platform linux amd64
    build_platform linux arm64

    # macOS
    build_platform darwin amd64
    build_platform darwin arm64

    # Windows
    build_platform windows amd64

    echo ""
    success "All builds completed!"
}

# Create release archive
create_archives() {
    info "Creating release archives..."

    cd "$BUILD_DIR"

    for file in slip-*; do
        if [[ ! "$file" =~ \.sha256$ ]]; then
            local archive_name="${file}.tar.gz"
            if [[ "$file" == *.exe ]]; then
                archive_name="${file%.exe}.zip"
                zip -q "$archive_name" "$file" "${file}.sha256" 2>/dev/null || true
            else
                tar -czf "$archive_name" "$file" "${file}.sha256" 2>/dev/null || true
            fi

            if [ -f "$archive_name" ]; then
                success "Created $archive_name"
            fi
        fi
    done

    cd ..
}

# Show build summary
summary() {
    echo ""
    info "Build Summary:"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo "Version:    ${VERSION}"
    echo "Commit:     $(get_commit_hash)"
    echo "Build Time: $(get_timestamp)"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""
    ls -lh "$BUILD_DIR"
    echo ""
    success "Build artifacts are in: $BUILD_DIR/"
}

# Main execution
main() {
    info "Slip Build Script v1.0"
    echo ""

    # Check if Go is installed
    if ! command -v go >/dev/null 2>&1; then
        error "Go is not installed or not in PATH"
        exit 1
    fi

    local go_version=$(go version)
    info "Using: $go_version"
    echo ""

    clean
    build_all
    create_archives
    summary
}

# Run main function
main "$@"
