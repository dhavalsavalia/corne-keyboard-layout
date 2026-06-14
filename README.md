# corne-build

![](https://github.com/dhavalsavalia/corne-keyboard-layout/actions/workflows/build.yml/badge.svg)

My ZMK config for a 42-key Corne split keyboard with nice!view displays.

A ZMK port of [josean's QMK Corne layout](https://github.com/josean-dev/dev-environment-files/blob/main/qmk/crkbd_rev1_3x6_josean.json) — same three layers and keycodes — with GACS homerow mods layered on top.

![Keymap](keymap.svg)

## Features

- QWERTY base layer with **GACS homerow mods**: pinky→index = GUI/ALT/CTRL/SHIFT, mirrored on both hands
- **Hyper** (Ctrl+Shift+Alt+Cmd) on the inner-index keys `G` and `H`
- Robust positional hold-taps (opposite-hand only, hold-trigger-on-release) so fast typing doesn't misfire mods
- Right-hand mods replicated on the Lower layer, so `mod+num` (e.g. aerospace `alt+1`) works while numbers are held
- Three layers: Base, Lower (symbols + numbers), Raise (F-keys, media, arrows)
- Left-inner thumb mod-tap: tap Enter / hold Alt
- Right-inner thumb mod-tap: tap Space / hold Hyper
- Momentary layer switches on the remaining inner/outer thumbs

## Layers

| Layer | Access                  | Purpose                       |
| ----- | ----------------------- | ----------------------------- |
| Base  | Default                 | QWERTY + GACS homerow mods    |
| Lower | Hold left-inner thumb   | Symbols + numbers             |
| Raise | Hold right-outer thumb  | F-keys, media, arrow nav      |

### Thumb cluster

```
Left:   Cmd      Lower(hold)    Alt/Enter
Right:  Hyper/Space   Backspace    Raise(hold)
```

## Build & Flash

Grab firmware artifacts from latest workflow run from ![this link](https://github.com/dhavalsavalia/corne-keyboard-layout/actions/workflows/build.yml).
It builds for each haves,

- Left: `corne_left nice_view_adapter nice_view-nice_nano-zmk.uf2`
- Right: `corne_right nice_view_adapter nice_view-nice_nano-zmk.uf2`

Flash by putting each half into bootloader mode (double-tap reset button) and copying the respective UF2 file to the mounted drive. The drive is usually named `NICENANO` on my Mac.
Once successfully flashed, nice!nano should eject itself. Do this for both halves.

## Credits

- [josean](https://github.com/josean-dev/dev-environment-files) for the original QMK layout this ports
- [keymap-drawer](https://github.com/caksoylar/keymap-drawer) for visualization
- [zmk-helpers](https://github.com/urob/zmk-helpers) for cleaner keymap syntax
