# Personal Knowledge Manager (pkm) 📙

**pkm** is a terminal-first, encrypted personal knowledge manager inspired by the **Zettelkasten** method.
It helps you build a *second brain* using small, interconnected notes — all stored securely on disk and fully controlled by you.

> ✨ Think: Obsidian-style linking + Git-friendly storage + strong encryption — without leaving the terminal.

## Why pkm?

* 🧠 **Zettelkasten-inspired** — notes are atomic, linkable, and discoverable
* 🔐 **Security first** — notes are encrypted at rest (AES-256-GCM)
* ⚡ **Fast & local** — no cloud, no latency, no vendor lock-in
* 🗂 **Clean filesystem layout** — works naturally with Git
* 🧑 **Multi-user support** — separate encrypted knowledge bases per user

## Features 🚀

### Core

* ✅ Create, edit, list, and delete notes
* ✅ Link notes to build a knowledge graph
* ✅ Tag notes for categorization
* ✅ Full-text keyword search
* ✅ Tag-based search
* ✅ Note indexing for fast queries

### Security

* ✅ AES-256-GCM encryption on disk
* ✅ PBKDF2 key derivation (100k iterations)
* ✅ Per-user encryption keys
* ✅ Zero plaintext storage

### System

* ✅ Multi-user support
* ✅ Git-safe encrypted notes
* ⏳ Multi-machine sync (planned)
* ⏳ Terminal UI (TUI) (planned)

## Installation ⚙️

### Option 1: Download from GitHub Releases (Recommended)

Prebuilt binaries are available for major platforms on the **GitHub Releases** page.

```bash
# Example (Linux x86_64)
curl -LO https://github.com/<username>/personal-knowledge-manager/releases/latest/download/pkm-linux-amd64
chmod +x pkm-linux-amd64
mv pkm-linux-amd64 /usr/local/bin/pkm
```

---

### Option 2: Build from source

#### Prerequisites

* Go (latest stable)
* `task` (Taskfile runner)

```bash
git clone https://github.com/<username>/personal-knowledge-manager.git
cd personal-knowledge-manager
task build
```

Add the binary to your `PATH`:

```bash
export PATH="$PATH:$(pwd)/build"
```

---

### Option 3: Package Managers (Planned)

* Homebrew
* AUR
* Scoop

> Contributions welcome 🙂

## Quick Start ⚡

```bash
# 1. Create a user
pkm user init alice

# 2. Create a note
pkm --user alice note new "My first note"

# 3. List notes
pkm --user alice note list

# 4. Link notes
pkm --user alice link add <source-id> <target-id>

# 5. Search notes
pkm --user alice search keyword "zettelkasten"
```

## Data Layout 📁

```text
~/.pkm/
├── .crypt               # Encrypted user keys (DO NOT commit)
├── alice/
│   ├── <note-id>.pkm    # Encrypted notes
│   └── .index.pkm       # Encrypted search index
```

### Git Usage

✔ Safe to commit:

* User directories (`<username>/`)
* Encrypted `.pkm` files

❌ Never commit:

* `.crypt`

## Security Model 🔐

* AES-256-GCM authenticated encryption
* PBKDF2 key derivation
* Random salts and nonces per user
* Passwords are never stored
* No plaintext ever written to disk

> Even if your repo is public, your notes remain private.

## Philosophy 🧩

pkm is designed around these principles:

* **Small notes > big documents**
* **Links create insight**
* **Local-first > cloud-first**
* **You own your data**

## Roadmap 🛣

* 🔄 Multi-machine sync
* 🖥 Terminal UI (TUI)
* 📊 Graph visualization
* 🔌 Plugin system

## Contributing 🤝

Contributions, ideas, and feedback are welcome.
Open an issue or submit a PR — even small improvements matter.
