# Keymap Styling Split Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Split keymap.yaml into separate content (corne.yaml) and styling (config.yaml) files in a new `keymap/` directory, adopting thrly's pink accent styling.

**Architecture:** Create `keymap/` directory with two files - `corne.yaml` for keymap content (layers, combos) and `config.yaml` for draw/parse configuration with thrly's styling. Update GitHub workflow and sync-keymap skill to use new paths.

**Tech Stack:** keymap-drawer, ZMK, GitHub Actions

---

### Task 1: Create keymap directory

**Files:**
- Create: `keymap/` directory

**Step 1: Create the directory**

```bash
mkdir -p keymap
```

**Step 2: Verify**

```bash
ls -la keymap
```
Expected: Empty directory exists

**Step 3: Commit**

```bash
git add keymap/.gitkeep 2>/dev/null || true
```
Note: Directory will be committed with files in next tasks.

---

### Task 2: Create keymap/config.yaml with thrly styling

**Files:**
- Create: `keymap/config.yaml`

**Step 1: Create config.yaml with thrly's styling**

```yaml
draw_config:
  key_w: 60
  key_h: 56
  split_gap: 30.0
  combo_w: 28
  combo_h: 26
  key_rx: 6.0
  key_ry: 6.0
  dark_mode: false
  n_columns: 1
  separate_combo_diagrams: false
  combo_diagrams_scale: 2
  inner_pad_w: 2.0
  inner_pad_h: 2.0
  outer_pad_w: 30.0
  outer_pad_h: 56.0
  line_spacing: 1.2
  arc_radius: 6.0
  append_colon_to_layer_header: false
  small_pad: 2.0
  legend_rel_x: 0.0
  legend_rel_y: 0.0
  draw_key_sides: false
  key_side_pars:
    rel_x: 0.0
    rel_y: 4.0
    rel_w: 12.0
    rel_h: 12.0
    rx: 4.0
    ry: 4.0
  svg_extra_style: |-
    /* font and background color specifications */
    svg.keymap {
      font-family: Ubuntu Mono, Inconsolata, Consolas, Liberation Mono, Menlo, monospace;
      font-size: 12px;
      font-weight: bold;
      text-rendering: optimizeLegibility;
    }

    /* default key styling */
    rect {
        fill: #f6f8fa;
        stroke: #c9cccf;
        stroke-width: 1;
    }

    /* color accent for combo boxes */
    rect.combo {
        fill: #fdb2e7;
    }

    /* color accent for held keys */
    rect.held, rect.combo.held {
        fill: #ee61bd;
    }

    /* color accent for ghost (optional) keys */
    rect.ghost, rect.combo.ghost {
        fill: #ddd;
    }

    text {
        text-anchor: middle;
        dominant-baseline: middle;
    }

    /* styling for layer labels */
    text.label {
        font-weight: bold;
        font-size: 16px;
        text-anchor: start;
        stroke: white;
        fill: #ee61bd;
        stroke-width: 2;
        paint-order: stroke;
    }

    /* styling for combo tap, and key hold/shifted label text */
    text.combo, text.hold, text.shifted {
        font-size: 11px;
    }

    text.hold {
        text-anchor: middle;
        dominant-baseline: auto;
        fill: #ee61bd;
    }

    text.shifted {
        text-anchor: middle;
        dominant-baseline: hanging;
        opacity: 0.7;
    }

    /* styling for hold/shifted label text in combo box */
    text.combo.hold, text.combo.shifted {
        font-size: 8px;
    }

    /* styling for combo dendrons */
    path {
        stroke-width: 1;
        stroke: gray;
        fill: none;
    }
    path.combo {
      stroke-dasharray: 4, 4;
      stroke-opacity: 0.3;
    }
  footer_text: ""
  shrink_wide_legends: 0
  style_layer_activators: true
  glyph_tap_size: 16
  glyph_hold_size: 12
  glyph_shifted_size: 10
  glyphs: {}

parse_config:
  preprocess: true
  skip_binding_parsing: false
  raw_binding_map:
    "&mmv MOVE_LEFT": MOUSE ←
    "&mmv MOVE_RIGHT": MOUSE →
    "&mmv MOVE_UP": MOUSE ↑
    "&mmv MOVE_DOWN": MOUSE ↓
    "&mkp LCLK": LClk
    "&mkp MB1": LClk
    "&mkp MB2": RClk
    "&mkp MB3": Mid Clk
    "&mkp MB4": BACK
    "&mkp MB5": FORWARD
    "&msc SCRL_DOWN": SCROLL ↓
    "&msc SCRL_UP": SCROLL ↑
    "&bootloader": BOOT
    "&studio_unlock": UNLOCK
    "&caps_word": ⇪ WORD

  sticky_label: sticky
  toggle_label: toggle
  tap_toggle_label: tap-toggle
  trans_legend:
    t: ▽
    type: trans

  modifier_fn_map:
    left_ctrl: Ctl
    right_ctrl: Ctl
    left_shift: Sft
    right_shift: Sft
    left_alt: Alt
    right_alt: AltGr
    left_gui: Gui
    right_gui: Gui
    keycode_combiner: "{mods}+{key}"
    mod_combiner: "{mod_1}+{mod_2}"

  zmk_keycode_map:
    AMPERSAND: "&"
    AMPS: "&"
    APOS: "'"
    APOSTROPHE: "'"
    ASTERISK: "*"
    ASTRK: "*"
    AT: "@"
    AT_SIGN: "@"
    BACKSLASH: \
    BSLH: \
    BACKSPACE: ⌫
    CAPSLOCK: ⇪ LOCK
    CARET: ^
    COLON: ":"
    COMMA: ","
    DLLR: $
    DOLLAR: $
    DOT: .
    DOUBLE_QUOTES: '"'
    DQT: '"'
    EQUAL: "="
    EXCL: "!"
    EXCLAMATION: "!"
    FSLH: /
    GRAVE: "`"
    GREATER_THAN: ">"
    GT: ">"
    HASH: "#"
    LBKT: "["
    LBRC: "{"
    LEFT_BRACE: "{"
    LEFT_BRACKET: "["
    LEFT_PARENTHESIS: (
    LESS_THAN: <
    LPAR: (
    LT: <
    MINUS: "-"
    NON_US_BACKSLASH: \
    NON_US_BSLH: "|"
    NON_US_HASH: "#"
    NUHS: "#"
    PERCENT: "%"
    PERIOD: .
    PIPE: "|"
    PIPE2: "|"
    PLUS: +
    POUND: "#"
    PRCNT: "%"
    QMARK: "?"
    QUESTION: "?"
    RBKT: "]"
    RBRC: "}"
    RIGHT_BRACE: "}"
    RIGHT_BRACKET: "]"
    RIGHT_PARENTHESIS: )
    RPAR: )
    SEMI: ;
    SEMICOLON: ;
    SINGLE_QUOTE: "'"
    SLASH: /
    SPACE: ⎵
    SQT: "'"
    STAR: "*"
    TILDE: "~"
    TILDE2: "~"
    UNDER: _
    UNDERSCORE: _
    LCTRL: ⌃
    LALT: ⎇
    LGUI: ⌘
    LSHFT: ⇧
    RCTRL: ⌃
    LEFT_CONTROL: ⌃
    LEFT_ALT: ⎇
    LEFT_GUI: ⌘
    LMETA: ⌘
    LEFT_SHFT: ⇧
    LEFT_SHIFT: ⇧
    RALT: ⎇
    RGUI: ⌘
    RSHFT: ⇧
    RIGHT_CONTROL: ⌃
    RIGHT_ALT: ⎇
    RIGHT_GUI: ⌘
    RIGHT_SHFT: ⇧
    RIGHT_SHIFT: ⇧
    LEFT: ←
    RIGHT: →
    UP: ↑
    DOWN: ↓
    LEFT_ARROW: ←
    RIGHT_ARROW: →
    UP_ARROW: ↑
    DOWN_ARROW: ↓
    TAB: ↹
    RETURN: ⏎
    RET: ⏎
    ESCAPE: ESC
    C_PREVIOUS: ⏮
    C_NEXT: ⏭
    C_PLAY_PAUSE: ⏯
    C_MUTE: MUTE
    C_VOLUME_DOWN: |-
      VOL
      DOWN
    C_VOLUME_UP: |-
      VOL
      UP
    C_BRIGHTNESS_DEC: ☀↓
    C_BRIGHTNESS_INC: ☀↑
    DELETE: ⌦
    PAGE_UP: PG UP
    PAGE_DOWN: PG DN
```

**Step 2: Verify YAML syntax**

```bash
python3 -c "import yaml; yaml.safe_load(open('keymap/config.yaml'))" && echo "Valid YAML"
```
Expected: "Valid YAML"

---

### Task 3: Create keymap/corne.yaml with keymap content

**Files:**
- Create: `keymap/corne.yaml`
- Reference: `keymap.yaml` (current file to extract from)

**Step 1: Create corne.yaml with layout, layers, and combos**

Extract from current `keymap.yaml`:
- `layout:` block
- `layers:` block (all 5 layers)
- `combos:` block

Remove:
- `draw_config:` block (now in config.yaml)
- Comments can be preserved or simplified

**Step 2: Verify YAML syntax**

```bash
python3 -c "import yaml; yaml.safe_load(open('keymap/corne.yaml'))" && echo "Valid YAML"
```
Expected: "Valid YAML"

---

### Task 4: Test new keymap-drawer setup

**Files:**
- Read: `keymap/config.yaml`, `keymap/corne.yaml`
- Create: `keymap.svg` (regenerate)

**Step 1: Generate SVG with new setup**

```bash
keymap -c keymap/config.yaml draw keymap/corne.yaml -o keymap.svg
```
Expected: SVG generated successfully

**Step 2: Verify SVG was created**

```bash
ls -la keymap.svg
```
Expected: File exists with recent timestamp

**Step 3: Commit new keymap structure**

```bash
git add keymap/config.yaml keymap/corne.yaml keymap.svg
git commit -m "refactor: split keymap into content and styling files

- Create keymap/corne.yaml for layers and combos
- Create keymap/config.yaml with thrly-inspired pink styling
- Ubuntu Mono font, pink accents (#fdb2e7, #ee61bd)"
```

---

### Task 5: Update GitHub workflow

**Files:**
- Modify: `.github/workflows/draw-keymap.yml`

**Step 1: Update workflow to use new paths**

Replace the generate/draw steps:

```yaml
name: Update drawing

on:
  push:
    branches:
      - main
    paths:
      - 'config/*'
      - 'keymap/*'

  workflow_dispatch:

jobs:
  update_drawing:
    name: Update drawing
    runs-on: ubuntu-latest
    permissions:
      contents: write
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-python@v6
        with:
          python-version: 3.13
      - name: Install pipx
        uses: CfirTsabari/actions-pipx@v1
      - name: Install keymap-drawer
        run: pipx install keymap-drawer
      - name: Generate SVG
        run: keymap -c keymap/config.yaml draw keymap/corne.yaml > keymap.svg
      - name: Commit and push changes
        run: |
          git config --global user.name "github-actions[bot]"
          git config --global user.email "41898282+github-actions[bot]@users.noreply.github.com"

          git add keymap.svg
          git commit -m "update keymap drawing" || true
          git push origin main
```

**Step 2: Commit workflow changes**

```bash
git add .github/workflows/draw-keymap.yml
git commit -m "ci: update workflow for new keymap directory structure"
```

---

### Task 6: Update sync-keymap skill

**Files:**
- Modify: `.claude/skills/sync-keymap.md`

**Step 1: Update paths in skill**

Change all references:
- `keymap.yaml` → `keymap/corne.yaml`
- Add note about config.yaml for styling

Key changes:
- Line 8: `keymap/corne.yaml`
- Line 86-88: Update draw command
- Any other path references

**Step 2: Update the draw command**

```bash
keymap -c keymap/config.yaml draw keymap/corne.yaml -o keymap.svg
```

**Step 3: Commit skill update**

```bash
git add .claude/skills/sync-keymap.md
git commit -m "docs: update sync-keymap skill for new keymap paths"
```

---

### Task 7: Delete old keymap.yaml

**Files:**
- Delete: `keymap.yaml`

**Step 1: Remove old file**

```bash
git rm keymap.yaml
```

**Step 2: Commit deletion**

```bash
git commit -m "chore: remove old keymap.yaml (replaced by keymap/)"
```

---

### Task 8: Final verification

**Step 1: Verify directory structure**

```bash
ls -la keymap/
```
Expected:
```
config.yaml
corne.yaml
```

**Step 2: Regenerate and verify SVG**

```bash
keymap -c keymap/config.yaml draw keymap/corne.yaml -o keymap.svg
```
Expected: Success, SVG has pink styling

**Step 3: Check git status**

```bash
git status
```
Expected: Clean working tree

**Step 4: Push all changes**

```bash
git push origin main
```
