# Corne ZMK Keymap

Custom ZMK keymap for Corne keyboard with nice!nano v2 + nice!view displays.

![Keymap](keymap.svg)

## Layers


| # | Name  | Access     |
| --- | ------- | ------------ |
| 0 | BASE  | Default    |
| 1 | NUM   | Hold BSPC  |
| 2 | SYM   | Hold SPACE |
| 3 | FUN   | Hold DEL   |
| 4 | MEDIA | Hold ESC   |

## Building

```bash
./build.sh        # both halves
./build.sh left   # left only
./build.sh right  # right only
./build.sh reset  # settings_reset firmware
```

Output goes to `./firmware/`

## Resetting (Clean Flash)

Use when keyboard has pairing issues or needs a fresh start:

1. Remove old "Corne" from Bluetooth devices, pair as new
2. Build reset firmware: `./build.sh reset`
3. **Left half reset**: Plug in left → double-tap reset → drag `settings_reset.uf2` → unplug
4. **Right half reset**: Plug in right → double-tap reset → drag `settings_reset.uf2` → unplug
5. Wait a few seconds
6. Flash normal firmware (see below)

## Flashing

1. Build firmware: `./build.sh`
2. **Left half**: Plug in left → double-tap reset → drag `corne_left.uf2` → unplug
3. **Right half**: Plug in right → double-tap reset → drag `corne_right.uf2` → unplug
4. Wait a few seconds
5. Test typing on both halves and check layers activate correctly
6. If pairing issues occur, perform a reset (see above)

> **Troubleshooting**: If one half doesn't type after connecting, repeat all steps. After flashing settings_reset to both halves, tap reset once on each half simultaneously to force them to pair with each other.
