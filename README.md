# PVM — PHP Version Manager for Windows

> Install, switch, and manage multiple PHP versions on Windows with a single command.

There is no official PHP version manager, and existing community tools are weakest on Windows. PVM fills that gap — built specifically for Windows developers who need to switch PHP versions between projects without the manual hassle.

---

## Installation

**One-liner (PowerShell):**

```powershell
irm https://github.com/medchakkir/pvm/releases/latest/download/install.ps1 | iex
```

**Manual:**
Download the latest `pvm-v*-windows-amd64.zip` from [Releases](https://github.com/medchakkir/pvm/releases), extract `pvm.exe`, and place it in a folder on your PATH.

---

## Commands

```powershell
# See what's available
pvm list-remote                  # all available PHP versions (TS + NTS)
pvm list-remote --ts             # Thread Safe builds only
pvm list-remote --nts            # Non-Thread Safe builds only

# Install a version
pvm install 8.3                  # install latest PHP 8.3 (Thread Safe)
pvm install 8.3.7                # install exact version
pvm install --nts 8.3            # install Non-Thread Safe build
pvm install --with-composer 8.3  # install PHP 8.3 and Composer

# Switch versions
pvm use 8.2                      # switch to PHP 8.2 (auto-detects TS/NTS)
pvm use --nts 8.2                # switch to NTS build specifically
pvm use                          # switch using .pvmrc from current or parent directory

# Pin project version
pvm init                         # write .pvmrc with the currently active version
pvm init 8.3                    # write .pvmrc with a specific installed version

# Manage installs
pvm list                         # show installed versions + active one
pvm current                      # show active version and its path
pvm uninstall 8.1                # remove PHP 8.1 (prompts for confirmation)

# Extensions
pvm extensions list              # show all extensions and their status
pvm extensions enable curl       # enable a single extension
pvm extensions enable curl,gd,mbstring  # enable multiple at once
pvm extensions disable xdebug   # disable an extension

# PATH setup
pvm bin                          # print the shims path to add to your PATH
```

---

## How It Works

PVM stores PHP versions in `~/.pvm/versions/` and maintains `.bat` shims in `~/.pvm/shims/` that stay on your PATH. Switching versions simply rewrites the shims to point to a different binary—no registry edits, no terminal restarts.

With `pvm install --with-composer <version>`, Composer is installed alongside PHP.

---

## Per-project pinning with .pvmrc

Run `pvm init` in a project directory to write the active version to `.pvmrc`. Then `pvm use` with no arguments will walk up the current directory tree, find the nearest `.pvmrc`, and activate the pinned version.

---

## Managing Extensions

```powershell
pvm extensions list                    # show all extensions and their status
pvm extensions enable curl             # enable an extension
pvm extensions enable curl,gd,mbstring # enable multiple
pvm extensions disable xdebug          # disable an extension
```

Extension statuses shown by `pvm extensions list`:

- `[enabled]` — the extension is active in `php.ini`
- `[disabled]` — the extension is present in `php.ini` but commented out
- `[available]` — the DLL exists in `ext/` but no entry is present in `php.ini`
- `[missing file]` — `php.ini` contains the extension entry, but the DLL is missing

PVM edits `php.ini` directly. It handles both `extension=` and `zend_extension=` directives (e.g., Xdebug, OPcache). If no `php.ini` exists yet, PVM creates one from `php.ini-development`.

---

## Thread Safe vs Non-Thread Safe

- **TS (Thread Safe)** — Default. Use for Apache mod_php or multi-threaded SAPIs.
- **NTS (Non-Thread Safe)** — Use `--nts` for FastCGI / PHP-FPM / Nginx environments.

---

## Requirements

- Windows 10 or later (x64)
- PowerShell 5.1+ (for the installer)
- Internet connection (for downloading PHP versions)

---

## Contributing

Contributions are very welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) for how to get started.

---

## License

This project is licensed under the [MIT License](LICENSE).
