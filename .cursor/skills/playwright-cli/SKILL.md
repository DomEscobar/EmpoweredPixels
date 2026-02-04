---
name: playwright-cli
description: Automates browser interactions using the Playwright CLI (npx playwright-cli). Use when automating browsers, running playwright-cli commands, testing web UIs, or when the user mentions Playwright CLI, browser automation, or playwright-cli.
---

# Browser Automation with playwright-cli

Install or run via: `npx skills add https://github.com/microsoft/playwright --skill playwright-cli` or use `npx playwright-cli` directly.

## Quick start

```bash
playwright-cli open https://playwright.dev
playwright-cli click e15
playwright-cli type "page.click"
playwright-cli press Enter
```

## Core workflow

1. **Navigate**: `playwright-cli open https://example.com`
2. **Snapshot** to get element refs (e.g. `e3`, `e5`).
3. **Interact** using refs from the snapshot.
4. **Re-snapshot** after significant page changes.

## Commands

### Core

| Command | Example |
|--------|---------|
| Open/close | `playwright-cli open https://example.com/` |
| | `playwright-cli close` |
| Type / fill | `playwright-cli type "search query"` |
| | `playwright-cli fill e5 "user@example.com"` |
| Click | `playwright-cli click e3` |
| | `playwright-cli dblclick e7` |
| Other | `playwright-cli drag e2 e8` |
| | `playwright-cli hover e4` |
| | `playwright-cli select e9 "option-value"` |
| | `playwright-cli upload ./document.pdf` |
| | `playwright-cli check e12` / `playwright-cli uncheck e12` |
| Snapshot / eval | `playwright-cli snapshot` |
| | `playwright-cli eval "document.title"` |
| | `playwright-cli eval "el => el.textContent" e5` |
| Dialogs | `playwright-cli dialog-accept` |
| | `playwright-cli dialog-accept "confirmation text"` |
| | `playwright-cli dialog-dismiss` |
| Viewport | `playwright-cli resize 1920 1080` |

### Navigation

- `playwright-cli go-back`
- `playwright-cli go-forward`
- `playwright-cli reload`

### Keyboard

- `playwright-cli press Enter`
- `playwright-cli press ArrowDown`
- `playwright-cli keydown Shift` / `playwright-cli keyup Shift`

### Mouse

- `playwright-cli mousemove 150 300`
- `playwright-cli mousedown` / `playwright-cli mousedown right`
- `playwright-cli mouseup` / `playwright-cli mouseup right`
- `playwright-cli mousewheel 0 100`

### Save as

- `playwright-cli screenshot` or `playwright-cli screenshot e5`
- `playwright-cli pdf`

### Tabs

- `playwright-cli tab-list`
- `playwright-cli tab-new` / `playwright-cli tab-new https://example.com/page`
- `playwright-cli tab-close` / `playwright-cli tab-close 2`
- `playwright-cli tab-select 0`

### Storage (auth state)

- `playwright-cli state-save` / `playwright-cli state-save auth.json`
- `playwright-cli state-load auth.json`

### DevTools

- `playwright-cli console` / `playwright-cli console warning`
- `playwright-cli network`
- `playwright-cli run-code "async page => await page.context().grantPermissions(['geolocation'])"`
- `playwright-cli tracing-start` / `playwright-cli tracing-stop`
- `playwright-cli video-start` / `playwright-cli video-stop video.webm`

### Configuration

```bash
playwright-cli config --config my-config.json
playwright-cli config --headed --in-memory --browser=firefox
playwright-cli --session=mysession config my-config.json
playwright-cli open --config=my-config.json
```

### Sessions

```bash
playwright-cli --session=mysession open example.com
playwright-cli --session=mysession click e6
playwright-cli session-list
playwright-cli session-stop mysession
playwright-cli session-stop-all
playwright-cli session-delete
playwright-cli session-delete mysession
```

## Examples

### Form submission

```bash
playwright-cli open https://example.com/form
playwright-cli snapshot
playwright-cli fill e1 "user@example.com"
playwright-cli fill e2 "password123"
playwright-cli click e3
playwright-cli snapshot
```

### Multi-tab

```bash
playwright-cli open https://example.com
playwright-cli tab-new https://example.com/other
playwright-cli tab-list
playwright-cli tab-select 0
playwright-cli snapshot
```

### Debugging (console, network, tracing)

```bash
playwright-cli open https://example.com
playwright-cli click e4
playwright-cli fill e7 "test"
playwright-cli console
playwright-cli network
```

```bash
playwright-cli open https://example.com
playwright-cli tracing-start
playwright-cli click e4
playwright-cli fill e7 "test"
playwright-cli tracing-stop
```

### Auth state reuse

```bash
# Login and save
playwright-cli open https://app.example.com/login
playwright-cli snapshot
playwright-cli fill e1 "user@example.com"
playwright-cli fill e2 "password123"
playwright-cli click e3
playwright-cli state-save auth.json

# Restore and skip login
playwright-cli state-load auth.json
playwright-cli open https://app.example.com/dashboard
```

## Specific topics

For deeper coverage see the Playwright docs or repo:

- Request mocking
- Running custom Playwright code (`run-code`)
- Session management (named sessions)
- Storage state (cookies, localStorage)
- Test generation
- Tracing and video recording
