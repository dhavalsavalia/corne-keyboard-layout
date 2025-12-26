# Corne Local Build

Local ZMK firmware builds for Corne keyboard with nice!nano v2 + nice!view displays.

## Setup (One-Time)

### Dependencies

```bash
brew install cmake ninja gperf ccache dtc libmagic
```

> Should be managed through [chezmoi](https://github.com/dhavalsavalia/dotfiles-chezmoi) already.

### ZMK Toolchain

```bash
cd /Users/dhavalsavalia/projects/keyboard
git clone https://github.com/zmkfirmware/zmk.git
cd zmk

python3 -m venv .venv
source .venv/bin/activate

pip install west
west init -l app/
west update
west zephyr-export
west packages pip --install
west sdk install
```

## Building

```bash
./build.sh        # both halves
./build.sh left   # left only
./build.sh right  # right only
```

Output goes to `./firmware/`

## Flashing

1. First Right Half then Left. Plug right half while its on.
2. Double-tap reset button on keyboard half (enters bootloader, mounts as USB drive)
3. Copy `.uf2` file to the mounted drive
4. Keyboard auto-reboots with new firmware
5. Repeat for other half

## Directory Structure

```
/Users/dhavalsavalia/projects/keyboard/
├── zmk/                      # ZMK firmware (upstream)
├── miryoku_zmk/              # Keymap config (mods branch)
│   ├── config/
│   │   ├── corne.keymap      # Build entry point
│   │   └── corne.conf        # Hardware config (display, sleep)
│   ├── miryoku/
│   │   └── custom_config.h   # Keymap customizations
│   └── keymap.yaml           # Visual reference (not used in build)
└── corne-build/              # This directory
    ├── build.sh
    ├── README.md
    └── firmware/             # Build output
```

## Config Files


| File                      | Purpose                                   |
| --------------------------- | ------------------------------------------- |
| `miryoku/custom_config.h` | Keymap: layout, layers, mods, combos      |
| `config/corne.conf`       | Hardware: display, sleep, battery widget  |
| `keymap.yaml`             | Visual documentation only (keymap-drawer) |

## Notes

- Board name is `nice_nano` (not `nice_nano_v2` - changed in 2024)
- Shield includes `nice_view_adapter nice_view` for display support
- Left half builds without `-p`, right half uses `-p` (pristine) for clean rebuild
