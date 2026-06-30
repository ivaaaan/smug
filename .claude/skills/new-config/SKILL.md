---
name: new-config
description: >-
  Interactively scaffold a new smug tmux session configuration (a YAML
  template). Use when the user wants to create, scaffold, or set up a new
  smug project/config file, or asks for "smug new config", "a smug template",
  or "/new-config". Walks the user through session name, root, windows, panes,
  commands, env, and hooks, then writes a valid smug YAML file.
---

# Create a new smug configuration

Scaffold a new smug session config by **interviewing the user** with the
`AskUserQuestion` tool, then writing a valid YAML file in the smug format.

Do not write the file until you have asked the questions below. Ask in small
batches, infer sensible defaults, and only ask what you can't infer.

## 1. Gather requirements (interactive)

Ask with `AskUserQuestion`. Group related questions into a single call when you
can. Always offer a clear default and let the user override via "Other".

1. **Session name** — the tmux session name (`session:`). Suggest the current
   directory's basename as the default.
2. **Where to save** — one of:
   - `.smug.yml` in the current working directory (used by default by `smug`).
   - `~/.config/smug/<session>.yml` (loadable as `smug start <session>`).
   - A custom path the user types.
3. **Root directory** — the session `root:` (where windows start). Default `.`
   for a cwd config, or the cwd's absolute path for a `~/.config/smug` config.
4. **Windows** — ask how many windows and, for each, the name and what it runs
   (e.g. "editor", "server", "tests"). Keep it light; one window with a
   command is a fine minimal config. Only dig into panes/layout if the user
   wants a split window.
5. **Optional extras** — only ask if the user signals they want them:
   - `attach: true` (auto-attach after creation).
   - `env:` variables.
   - `before_start:` / `stop:` commands (e.g. `docker-compose up -d`).
   - `attach_hook:` / `detach_hook:`.

If the user says "just give me a basic one" or similar, skip straight to a
minimal config (session + root + one window) — don't force the full interview.

## 2. Smug config schema (reference)

Generate YAML matching this structure. Omit empty/optional fields.

```yaml
session: my_project        # tmux session name (required)
root: ~/code/my_project    # base dir for the session
attach: true               # optional: auto-attach after creation (default false)

env:                       # optional: environment variables
  FOO: bar

before_start:              # optional: runs once before the session is created
  - docker-compose up -d
stop:                      # optional: runs once before the session is killed
  - docker-compose stop

attach_hook: echo hi       # optional: runs each time first client attaches
detach_hook: echo bye      # optional: runs each time last client detaches

windows:
  - name: code             # window name
    root: .                # optional: path relative to session root
    selected: true         # optional: focus this window on start
    manual: false          # optional: only start via `-w` when true
    layout: main-vertical  # optional: tmux layout (e.g. tiled, main-horizontal)
    commands:              # commands run in the window's first pane
      - $EDITOR .
    panes:                 # optional: extra split panes
      - type: horizontal   # horizontal | vertical
        root: .            # optional
        commands:
          - npm run dev
```

Notes:
- `session` is the only truly required field. A useful minimal config is
  `session` + `root` + one `windows` entry with a `commands` list.
- `${var}` in the YAML is substituted from custom settings passed on the CLI
  (`smug start project var=value`) or the environment — only use it if the
  user asks for parameterization.
- Quote nothing unnecessarily; keep it idiomatic, readable YAML.

## 3. Write and confirm

1. Write the file to the chosen path with the `Write` tool. If saving to
   `~/.config/smug/`, expand `~` to the home directory and ensure the directory
   exists (it normally does — smug creates it).
2. Show the generated YAML to the user.
3. Tell them how to use it:
   - cwd `.smug.yml`: `smug start`
   - `~/.config/smug/<session>.yml`: `smug start <session>`
   - custom path: `smug start -f <path>`
4. Offer to validate it. If a `smug` binary is on `PATH`, run
   `smug print -f <path>` (parses the config and prints the tmux commands
   without starting a session). Inside this repo you can instead run
   `go run . print -f <path>`. If it errors, fix the YAML and retry.

Keep the result minimal and tailored to what the user actually asked for — no
speculative windows, panes, or hooks they didn't request.
