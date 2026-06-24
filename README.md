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
Download the latest `pvm_*_windows_amd64.zip` from [Releases](https://github.com/medchakkir/pvm/releases), extract `pvm.exe`, and place it in a folder on your PATH.

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

PVM installs PHP versions into `~/.pvm/versions/` and maintains a set of `.bat` shims in `~/.pvm/shims/` (`php`, `php-cgi`, and `composer` when installed). That shims directory stays permanently on your PATH — switching versions simply rewrites the shims to point to a different PHP binary. No registry editing on every switch, no terminal restarts needed.

```
~/.pvm/
├── versions/
│   ├── 8.3.7-TS/     ← extracted PHP installs (+ composer.phar if requested)
│   └── 8.2.18-NTS/
├── shims/
│   ├── php.bat       ← always on PATH, rewrites on `pvm use`
│   ├── php-cgi.bat   ← FastCGI launcher
│   └── composer.bat  ← present when the active version has Composer
└── current           ← active version name
```

Install Composer alongside PHP with `pvm install --with-composer <version>`. PVM picks the right Composer line for the runtime (Composer 2.2 LTS for PHP < 7.2, latest stable otherwise).

---

## Managing Extensions

PVM can enable and disable PHP extensions by editing the active version's `php.ini` directly. No manual file editing required.

```powershell
pvm extensions list
```

Shows every extension available for the active PHP version, with one of four statuses:

| Status           | Meaning                                       |
| ---------------- | --------------------------------------------- |
| `[enabled]`      | Active in `php.ini`                           |
| `[disabled]`     | Present in `php.ini` but commented out        |
| `[available]`    | DLL exists in `ext/` but not yet in `php.ini` |
| `[missing file]` | Entry in `php.ini` but the DLL is gone        |

```powershell
pvm extensions enable curl
pvm extensions enable curl,gd,mbstring   # comma-separated for multiple
```

If the extension already has an entry in `php.ini`, PVM uncomments it. If it only exists as a `.dll` in the `ext/` folder, PVM adds the `extension=` line automatically. If no `php.ini` exists yet, PVM creates one from `php.ini-development`.

```powershell
pvm extensions disable xdebug
pvm extensions disable xdebug,opcache
```

Comments out the directive in `php.ini`. Works for both `extension=` and `zend_extension=` lines (e.g. Xdebug, OPcache).

> **Note:** Extensions are per-version. Run `pvm use <version>` first to switch to the version you want to configure.

---

## Thread Safe vs Non-Thread Safe

- **TS (Thread Safe)** — for use with Apache mod_php or multi-threaded SAPIs. Default.
- **NTS (Non-Thread Safe)** — for use with FastCGI / PHP-FPM / Nginx. Use `--nts` flag.

If you're running Laravel with `php artisan serve`, either works. If you're running behind Nginx or a FastCGI server, prefer NTS.

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
