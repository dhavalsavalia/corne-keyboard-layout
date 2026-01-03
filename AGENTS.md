# Agent Instructions

## Build Commands
```bash
./build.sh        # Build both halves
./build.sh left   # Left half only
./build.sh right  # Right half only
```
Output: `./firmware/YYYYMMDD/corne_{left,right}.uf2`

## Code Style
- **Bash**: Use `set -e` for error checking, quote variables, prefer `${var:-default}` syntax
- **YAML**: 2-space indentation, no trailing whitespace, document layers with comments
- **Devicetree**: Keep consistent with ZMK patterns, use uppercase for constants (BASE, NUM, etc.)
- **Naming**: Descriptive layer names, snake_case for files

## Workflow
```bash
bd ready              # Find work
bd update <id> --status in_progress  # Claim
# ...work...
bd close <id>         # Complete
```

## Session Completion (MANDATORY)
```bash
git add <files> && git commit -m "..."
bd sync && git push && git status
```
Work NOT complete until `git push` succeeds. Always push before ending session.
