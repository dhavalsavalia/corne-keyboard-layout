# Corne ZMK Keymap

Custom ZMK keymap for Corne keyboard with nice!nano v2 + nice!view displays.

![Keymap](keymap.svg)

## Layers

| # | Name  | Access    |
|---|-------|-----------|
| 0 | BASE  | Default   |
| 1 | NUM   | Hold BSPC |
| 2 | SYM   | Hold SPACE|
| 3 | FUN   | Hold DEL  |
| 4 | MEDIA | Hold ESC  |

## Building

```bash
./build.sh        # both halves
./build.sh left   # left only
./build.sh right  # right only
```

Output goes to `./firmware/`

## Flashing

1. Double-tap reset button (enters bootloader, mounts as USB drive)
2. Copy `.uf2` file to the mounted drive
3. Keyboard auto-reboots with new firmware
4. Repeat for other half (right first, then left)
