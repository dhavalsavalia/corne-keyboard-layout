# corne-build

My ZMK config for a 42-key Corne split keyboard with nice!view displays.

> Inspired by [thrly's corne layout](https://github.com/thrly/thrly-corne-zmk) — check out their excellent [blog post on split keyboard customization](https://thrly.com/blog/thoughts-on-customising-a-split-keyboard-layout/)

![Keymap](keymap.svg)

## Features

- QWERTY with balanced homerow mods (⌘-⎇-⌃-⇧)
- Sticky shift on outer thumb keys
- Space morphs to underscore when shifted
- Nav layer with mouse keys and arrow navigation
- Horizontal number row with mirrored homerow mods

### Combos

| Keys | Action |
|------|--------|
| Q + W | Escape |
| G + H | Caps Word |
| T + Y | Caps Lock |
| E + R | Dash |
| W + E + R | Em-dash (` -- `) |
| U + I | Underscore |
| U + J | Single quote |
| I + K | At symbol |
| O + P | Delete |
| P + ; | Volume up |
| ; + / | Volume down |
| Q + A | Brightness up |
| A + Z | Brightness down |

### Macros

| Trigger | Output |
|---------|--------|
| Double-tap `>` | `=>` (arrow function) |
| W + E + R combo | ` -- ` (em-dash with spaces) |

## Layers

| Layer | Access |
|-------|--------|
| Base | Default |
| Nav | Hold Space |
| Sym | Hold Backspace |
| Num | Hold Return |
| Utils | Hold Tab |

## Build & Flash

```bash
./build.sh          # build both halves
kbflash             # flash via TUI
```

## Credits

- [thrly](https://github.com/thrly/thrly-corne-zmk) for keymap inspiration
- [keymap-drawer](https://github.com/caksoylar/keymap-drawer) for visualization
