# corne-build

![](https://github.com/dhavalsavalia/corne-keyboard-layout/actions/workflows/build.yml/badge.svg)

My ZMK config for a 42-key Corne split keyboard with nice!view displays.

> Inspired by [thrly's corne layout](https://github.com/thrly/thrly-corne-zmk) — check out their excellent [blog post on split keyboard customization](https://thrly.com/blog/thoughts-on-customising-a-split-keyboard-layout/)

![Keymap](keymap.svg)

## Features

- QWERTY with balanced homerow mods (⌘⎇⌃⇧ or GACS)
- Sticky shift on right outer thumb key and left outer pinky
- Space morphs to underscore when shifted
- Nav layer with mouse keys and arrow navigation
- Horizontal number row with aligned F-keys and symbols

### Combos


| Keys      | Action          |
| ----------- | ----------------- |
| Q + W     | Escape          |
| G + H     | Caps Word       |
| T + Y     | Caps Lock       |
| E + R     | Dash            |
| W + E + R | Em-dash (`--`)  |
| U + I     | Underscore      |
| U + J     | Single quote    |
| I + K     | At symbol       |
| O + P     | Delete          |
| P + ;     | Volume up       |
| ; + /     | Volume down     |
| Q + A     | Brightness up   |
| A + Z     | Brightness down |

### Macros


| Trigger         | Output                     |
| ----------------- | ---------------------------- |
| Double-tap`>`   | `=>` (arrow function)      |
| W + E + R combo | `--` (em-dash with spaces) |

## Layers


| Layer | Access         |
| ------- | ---------------- |
| Base  | Default        |
| Nav   | Hold Space     |
| Sym   | Hold Backspace |
| Num   | Hold Return    |
| Utils | Hold Tab       |

## Build & Flash

Grab firmware artifacts from latest workflow run from ![this link](https://github.com/dhavalsavalia/corne-keyboard-layout/actions/workflows/build.yml).
It builds for each haves,

- Left: `corne_left nice_view_adapter nice_view-nice_nano-zmk.uf2`
- Right: `corne_right nice_view_adapter nice_view-nice_nano-zmk.uf2`

Flash by putting each half into bootloader mode (double-tap reset button) and copying the respective UF2 file to the mounted drive. The drive is usually named `NICENANO` on my Mac.
Once successfully flashed, nice!nano should eject itself. Do this for both halves.

## Credits

- [thrly](https://github.com/thrly/thrly-corne-zmk) for keymap inspiration
- [keymap-drawer](https://github.com/caksoylar/keymap-drawer) for visualization
- [zmk-helpers](https://github.com/urob/zmk-helpers) for cleaner keymap syntax
