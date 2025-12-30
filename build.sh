#!/bin/bash
set -e

echo dirname "$0"

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
KEYBOARD_DIR="$(dirname "$SCRIPT_DIR")"

ZMK_DIR="$KEYBOARD_DIR/zmk"
CONFIG_DIR="$SCRIPT_DIR/config"
OUTPUT_DIR="$SCRIPT_DIR/firmware/$(date +%Y%m%d)"

mkdir -p "$OUTPUT_DIR"

# Activate ZMK environment
cd "$ZMK_DIR"
source .venv/bin/activate

build_half() {
    local side=$1
    local pristine=${2:--p}
    echo "Building $side half..."
    west build $pristine -s app -b nice_nano -- \
        -DSHIELD="corne_$side nice_view_adapter nice_view" \
        -DZMK_CONFIG="$CONFIG_DIR"
    cp build/zephyr/zmk.uf2 "$OUTPUT_DIR/corne_$side.uf2"
    echo "Done: $OUTPUT_DIR/corne_$side.uf2"
}

build_reset() {
    echo "Building settings_reset firmware..."
    west build -p -s app -b nice_nano -- -DSHIELD=settings_reset
    cp build/zephyr/zmk.uf2 "$OUTPUT_DIR/settings_reset.uf2"
    echo "Done: $OUTPUT_DIR/settings_reset.uf2"
}

case "${1:-both}" in
    left)  build_half left "" ;;
    right) build_half right ;;
    both)  build_half left "" && build_half right ;;
    reset) build_reset ;;
    *)     echo "Usage: $0 [left|right|both|reset]" && exit 1 ;;
esac

echo ""
echo "Firmware ready in $OUTPUT_DIR"
echo "Flash by double-tap reset, then copy .uf2 to mounted drive."
