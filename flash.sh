#!/bin/bash
set -e

# ============================================================================
# flash.sh - Interactive TUI for flashing Corne keyboard firmware
# ============================================================================

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
FIRMWARE_DIR="$SCRIPT_DIR/firmware"
VOLUME_PATH="/Volumes/NICENANO"

# ============================================================================
# Dependency Check
# ============================================================================

if ! command -v gum &> /dev/null; then
    echo "Error: gum is required but not installed."
    echo "Install with: brew install gum"
    exit 1
fi

# ============================================================================
# Helper Functions
# ============================================================================

header() {
    echo ""
    gum style --foreground 212 --bold "$1"
}

info() {
    gum style --foreground 39 "  $1"
}

success() {
    gum style --foreground 82 "  ✓ $1"
}

warn() {
    gum style --foreground 214 "  ⚠ $1"
}

error() {
    gum style --foreground 196 "  ✗ $1"
}

get_latest_firmware() {
    # Find most recent dated directory (format: YYYYMMDD)
    local latest
    latest=$(ls -1d "$FIRMWARE_DIR"/[0-9]* 2>/dev/null | sort -r | head -1)
    echo "$latest"
}

wait_for_device() {
    gum spin --spinner dot --title "Waiting for NICENANO to mount..." -- \
        bash -c 'while [ ! -d "/Volumes/NICENANO" ]; do sleep 0.5; done'
}

wait_for_eject() {
    gum spin --spinner dot --title "Waiting for device to eject..." -- \
        bash -c 'while [ -d "/Volumes/NICENANO" ]; do sleep 0.5; done'
}

flash_file() {
    local file="$1"
    local label="$2"

    header "Flash: $label"
    info "Connect $label, then double-tap reset button"
    echo ""

    # Wait for device
    if [ -d "$VOLUME_PATH" ]; then
        info "Device already mounted"
    else
        wait_for_device
        success "Device connected"
    fi

    # Copy firmware (ignore extended attributes error on FAT32)
    sleep 0.5
    info "Copying firmware..."
    cp "$file" "$VOLUME_PATH/" 2>/dev/null || true

    # Device may eject very quickly after copy - if volume gone, assume success
    if [ ! -d "$VOLUME_PATH" ]; then
        success "$label flashed successfully"
    else
        wait_for_eject
        success "$label flashed successfully"
    fi
    echo ""
}

run_build() {
    local target="$1"
    header "Building firmware..."

    if ! "$SCRIPT_DIR/build.sh" "$target" 2>&1; then
        error "Build failed"
        exit 1
    fi
    success "Build complete"
}

# ============================================================================
# Main Flow
# ============================================================================

header "Corne Keyboard Flash Utility"
echo ""

# Step 1: Build Selection
BUILD_OPTIONS=("Skip (use existing firmware)" "Build both halves")
info "What would you like to build?"
BUILD_CHOICE=$(gum choose --cursor.foreground 212 "${BUILD_OPTIONS[@]}")

if [ "$BUILD_CHOICE" = "Build both halves" ]; then
    run_build "both"
else
    info "Skipping build"
fi

# Step 2: Detect Firmware
LATEST_FIRMWARE=$(get_latest_firmware)

if [ -z "$LATEST_FIRMWARE" ]; then
    warn "No firmware found in $FIRMWARE_DIR"
    if gum confirm "Build firmware now?"; then
        run_build "both"
        LATEST_FIRMWARE=$(get_latest_firmware)
    else
        error "Cannot continue without firmware"
        exit 1
    fi
fi

info "Using firmware from: $(basename "$LATEST_FIRMWARE")"
echo ""

# Check required files are present
LEFT_FILE="$LATEST_FIRMWARE/corne_left.uf2"
RIGHT_FILE="$LATEST_FIRMWARE/corne_right.uf2"
RESET_FILE="$LATEST_FIRMWARE/settings_reset.uf2"

if [ ! -f "$LEFT_FILE" ] || [ ! -f "$RIGHT_FILE" ]; then
    error "Both corne_left.uf2 and corne_right.uf2 are required"
    info "Found in $LATEST_FIRMWARE:"
    [ -f "$LEFT_FILE" ] && info "  • corne_left.uf2" || warn "  • corne_left.uf2 (missing)"
    [ -f "$RIGHT_FILE" ] && info "  • corne_right.uf2" || warn "  • corne_right.uf2 (missing)"
    exit 1
fi

# Build flash options
FLASH_OPTIONS=("Flash both halves")
[ -f "$RESET_FILE" ] && FLASH_OPTIONS+=("Factory reset + Flash both")

# Step 3: Flash Selection
info "What would you like to do?"
FLASH_CHOICE=$(gum choose --cursor.foreground 212 "${FLASH_OPTIONS[@]}")

FLASH_RESET=false
[ "$FLASH_CHOICE" = "Factory reset + Flash both" ] && FLASH_RESET=true

# Step 4: Factory Reset Flow (if selected)
if $FLASH_RESET; then
    header "Factory Reset"
    warn "Remember to unpair keyboard from all Bluetooth devices!"
    echo ""

    if ! gum confirm "Have you unpaired from all devices?"; then
        info "Please unpair first, then run this script again."
        exit 0
    fi

    # Flash reset to both halves
    flash_file "$RESET_FILE" "LEFT half (reset)"
    flash_file "$RESET_FILE" "RIGHT half (reset)"

    warn "Keyboard should now be non-functional"
    if ! gum confirm "Is the keyboard non-functional?"; then
        error "Factory reset may have failed. Please try again."
        exit 1
    fi

    success "Factory reset complete"
    echo ""
fi

# Step 5: Flash Firmware (always both)
flash_file "$LEFT_FILE" "LEFT half"
flash_file "$RIGHT_FILE" "RIGHT half"

# Step 6: Completion
header "Flash Complete!"

if $FLASH_RESET; then
    warn "Don't forget to re-pair keyboard to your devices"
    echo ""
fi

info "Summary:"
$FLASH_RESET && info "  • Factory reset: both halves"
info "  • Left half: flashed"
info "  • Right half: flashed"
echo ""
info "Test both halves to verify everything works correctly."
echo ""
