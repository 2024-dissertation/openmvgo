#!/bin/bash
# ====================
# Entrypoint for OpenMVG / OpenMVS dockerfile
# Handles initialization, secrets, and graceful startup
# Sam Laister 2025
# ====================
set -e

# Colors for logging
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')] INFO:${NC} $1"
}

warn() {
    echo -e "${YELLOW}[$(date +'%Y-%m-%d %H:%M:%S')] WARN:${NC} $1"
}

error() {
    echo -e "${RED}[$(date +'%Y-%m-%d %H:%M:%S')] ERROR:${NC} $1"
}

graceful_shutdown() {
    # TODO graceful shutdown if needed
    log "Graceful shutdown completed"
    exit 0
}

# Set up signal handlers
trap graceful_shutdown SIGTERM SIGINT

# Pre-flight checks
preflight_checks() {
    log "Running pre-flight checks..."
    
    # Check if required binaries exist
    local required_bins=("openmvgo")
    for bin in "${required_bins[@]}"; do
        if ! command -v "$bin" >/dev/null 2>&1; then
            error "Required binary not found: $bin"
            exit 1
        fi
    done
    
    # Check OpenMVG/OpenMVS binaries
    local openmvg_bins=("openMVG_main_SfMInit_ImageListing" "openMVG_main_ComputeFeatures")
    for bin in "${openmvg_bins[@]}"; do
        if ! command -v "$bin" >/dev/null 2>&1; then
            warn "OpenMVG binary not found: $bin"
        fi
    done
    
    log "Pre-flight checks completed successfully"
}

# Main execution
main() {
    log "Starting application with PID: $$"
    log "Running as user: $(whoami)"
    log "Working directory: $(pwd)"
    log "Environment: ${GIN_MODE:-development}"
    
    # Run initialization steps
    preflight_checks
    
    # Continue to Dockerfile
    log "Executing command: $*"
    exec "$@"
}

# Execute main function with all arguments
main "$@"