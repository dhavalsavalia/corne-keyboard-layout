# Faithful port of josean's QMK Corne layout to ZMK

Date: 2026-06-14

## Goal

Port [josean's QMK Corne layout](https://raw.githubusercontent.com/josean-dev/dev-environment-files/refs/heads/main/qmk/crkbd_rev1_3x6_josean.json)
to this ZMK Corne **faithfully** — same layers, keycodes, and beginner-friendly
philosophy — as a clean baseline. Personal customizations (combos, homerow mods)
come in a later, separate pass.

## Toolchain (already updated, prerequisite)

Pinned to latest ZMK stable, verified compatible with this keymap:

- `config/west.yml`: `zmk` and `zmk-helpers` → `revision: v0.3.0`
- `.github/workflows/build.yml`: `build-user-config.yml@v0.3.0`
- `.github/workflows/draw-keymap.yml`: `keymap-drawer==0.23.0`

## Decisions

- **Syntax:** zmk-helpers stays as a dependency (`west.yml`) and via includes
  (`helper.h`, `42.h`). The 3 layers are written as plain devicetree `layer { bindings }`
  nodes — the same proven style as the current file. Helpers are unused by this
  faithful port but retained for the upcoming customization pass (combos, homerow mods).
- **RGB keys:** josean's RAISE layer has 6 RGB underglow keys. Nice!View has no
  underglow, so these map to `&none` (layout-faithful, hardware-honest).
- **Scope:** full clean replace of `config/corne.keymap`. The current 5-layer design
  (homerow mods, 13 combos, mouse keys) is removed; it stays recoverable in git history.
- **Layer names:** `BASE` / `LOWER` / `RAISE` (canonical QMK terminology, faithful to
  josean's beginner philosophy). `#define BASE 0`, `LOWER 1`, `RAISE 2`.
- **Behaviors:** stock `&mt` / `&mo` defaults — no behavior tuning, to stay simple.

## Files changed

| File | Change |
|------|--------|
| `config/corne.keymap` | Replace entirely with 3-layer faithful port |
| `config/corne.conf` | Remove `CONFIG_ZMK_POINTING` (no mouse keys in this design) |
| `keymap/corne.yaml` | Regenerate from the new keymap so the drawn SVG matches |

## Thumb cluster (josean's exact mapping, left → right)

```
LH2:&kp LGUI   LH1:&mo LOWER   LH0:&mt LALT RET
RH0:&mt LC(LS(LA(LGUI))) SPACE   RH1:&kp BSPC   RH2:&mo RAISE
```

The right-inner thumb is HYPER/Space: hold = Ctrl+Shift+Alt+Gui (ZMK has no single
HYPER alias, so it is nested), tap = Space. Left-inner thumb is Alt/Enter mod-tap.

## Layer bindings (42-key, row-major: 12 top, 12 mid, 12 bottom, 6 thumb)

### BASE (layer 0)

```
&kp TAB    &kp Q  &kp W  &kp E  &kp R  &kp T      &kp Y  &kp U  &kp I      &kp O    &kp P     &kp ESC
&kp LCTRL  &kp A  &kp S  &kp D  &kp F  &kp G      &kp H  &kp J  &kp K      &kp L    &kp SEMI  &kp SQT
&kp LSHFT  &kp Z  &kp X  &kp C  &kp V  &kp B      &kp N  &kp M  &kp COMMA  &kp DOT  &kp FSLH  &kp RSHFT
               &kp LGUI  &mo LOWER  &mt LALT RET      &mt LC(LS(LA(LGUI))) SPACE  &kp BSPC  &mo RAISE
```

### LOWER (layer 1 — symbols + numbers)

```
&trans  &kp EXCL  &kp AT  &kp HASH  &kp DLLR  &kp PRCNT    &kp CARET  &kp AMPS   &kp ASTRK  &kp LPAR  &kp RPAR  &kp BSLH
&trans  &kp N1    &kp N2  &kp N3    &kp N4    &kp N5       &kp MINUS  &kp EQUAL  &kp GRAVE  &kp LBKT  &kp RBKT  &kp PIPE
&trans  &kp N6    &kp N7  &kp N8    &kp N9    &kp N0       &kp UNDER  &kp PLUS   &kp TILDE  &kp LBRC  &kp RBRC  &trans
                          &trans  &trans  &trans              &trans  &trans  &trans
```

### RAISE (layer 2 — F-keys, media, nav; RGB → `&none`)

```
&bootloader  &kp F1      &kp F2      &kp F3        &kp F4        &kp F5      &kp F6    &kp F7    &kp F8  &kp F9     &kp F10  &none
&trans       &kp C_PREV  &kp C_NEXT  &kp C_VOL_DN  &kp C_VOL_UP  &kp C_PP    &kp LEFT  &kp DOWN  &kp UP  &kp RIGHT  &none    &none
&trans       &none       &none       &none         &none         &none       &none     &none     &none   &none      &none    &trans
                          &trans  &trans  &trans              &trans  &trans  &trans
```

`&bootloader` (top-left) is josean's `QK_BOOT`. RAISE bottom-left held josean's 6 RGB
keys, now `&none`. Right-side `KC_NO` positions are also `&none`.

## Keycode translation notes (QMK → ZMK)

Non-obvious mappings used above: `KC_SCLN`→`SEMI`, `KC_QUOT`→`SQT`, `KC_SLSH`→`FSLH`,
`KC_BSLS`→`BSLH`, `KC_GRV`→`GRAVE`, `KC_LBRC/RBRC`→`LBKT/RBKT`, `KC_LCBR/RCBR`→`LBRC/RBRC`,
`KC_ENT`→`RET`, `KC_MPRV/MNXT`→`C_PREV/C_NEXT`, `KC_VOLD/VOLU`→`C_VOL_DN/C_VOL_UP`,
`KC_MPLY`→`C_PP`. `MO(n)`→`&mo`, `MT(mod,kc)`→`&mt mod kc`, `KC_TRNS`→`&trans`,
`KC_NO`→`&none`, `QK_BOOT`→`&bootloader`.

## Verification

- Build succeeds via GitHub Actions (push to `config/**` triggers `build.yml`).
- All 42 josean keycodes accounted for per layer; only RGB downgraded (hardware limit).
- keymap-drawer regenerates `keymap/corne.yaml` and the SVG reflects 3 layers.

## Out of scope (future customization pass)

Homerow mods, combos, additional layers (NAV/UTILS), mouse keys, behavior tuning.
