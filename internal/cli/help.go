package cli

import "fmt"

func globalHelp() {
	fmt.Print(`
╔════════════════════════════════════════════════════════════════════════════╗
║                                                                            ║
║  pkm - A Personal Knowledge Manager based on Zettelkasten                  ║
║                                                                            ║
║  A CLI tool for building interconnected notes that form a second brain     ║
║                                                                            ║
╚════════════════════════════════════════════════════════════════════════════╝

USAGE:
  pkm [--user <username>] [--storeDirectory <path>] <command> [arguments]

FLAGS:
  --user <username>            Username (required for note operations)
  --storeDirectory <path>      Storage directory (default: ~/.pkm)

COMMANDS:`)

	// Display commands with descriptions
	commands := []struct {
		name string
		desc string
	}{
		{"user", "Create users, change passwords, manage accounts"},
		{"note", "Create, read, edit, and delete encrypted notes"},
		{"link", "Create and manage links between notes for knowledge discovery"},
		{"tag", "Organize notes with tags for categorization and search"},
		{"search", "Search notes by keywords or tags"},
		{"help", "Show detailed help for a command"},
	}

	for _, cmd := range commands {
		fmt.Printf("  %-15s %s\n", cmd.name, cmd.desc)
	}

	fmt.Print(`
QUICK START:

  1. Initialize a new user:
     $ pkm user init alice
     
  2. Create a note:
     $ pkm --user alice note new "My first note"
     
  3. Link notes together:
     $ pkm --user alice link add <source-id> <target-id>
     
  4. Search notes:
     $ pkm --user alice search keyword "keyword"

DETAILED HELP:

  Get help for a specific command:
    $ pkm <command> help
    $ pkm user help
    $ pkm note help
    $ pkm link help
    $ pkm tag help
    $ pkm search help

USER COMMANDS:

  pkm user init <username>
    Create a new user account with password
    
  pkm user passwd <username>
    Change password for existing user

NOTE COMMANDS:

  pkm --user <username> note new <title>
    Create a new note (opens $EDITOR for content)
    
  pkm --user <username> note edit <note-id>
    Edit existing note content
    
  pkm --user <username> note get <note-id>
    Display note content in terminal
    
  pkm --user <username> note delete <note-id>
    Delete a note permanently
    
  pkm --user <username> note list
    List all notes with IDs and titles

LINK COMMANDS:

  pkm --user <username> link add <source-id> <target-id>
    Create a link from source note to target note
    
  pkm --user <username> link remove <source-id> <target-id>
    Remove a link between two notes
    
  pkm --user <username> link list <note-id>
    Show all outgoing links from a note

TAG COMMANDS:

  pkm --user <username> tag add <note-id> <tag1,tag2,...>
    Add one or more tags to a note
    
  pkm --user <username> tag remove <note-id> <tag1,tag2,...>
    Remove tags from a note
    
  pkm --user <username> tag list <note-id>
    Show all tags on a note

SEARCH COMMANDS:

  pkm --user <username> search keyword <term1> [term2] ...
    Find notes by searching keywords in title and content
    
  pkm --user <username> search tag <tag1> [tag2] ...
    Find notes by tag (returns notes with all specified tags)

COMMON WORKFLOWS:

  Building a Zettelkasten:
    $ pkm user init alice
    $ pkm --user alice note new "Graph Theory Basics"
    $ pkm --user alice note new "BFS Algorithm"
    $ pkm --user alice link add <graph-theory-id> <bfs-id>
    $ pkm --user alice tag add <graph-theory-id> "algorithms,graphs"
    $ pkm --user alice search tag "algorithms"

  Exploring your knowledge:
    $ pkm --user alice note list
    $ pkm --user alice note get <note-id>
    $ pkm --user alice link list <note-id>
    $ pkm --user alice search keyword "recursion"

DATA & SECURITY:

  Storage:
    ~/.pkm/
    ├── .crypt              (encrypted user keys - NEVER commit to git)
    ├── <username>/
    │   ├── <note-id>.pkm   (encrypted notes)
    │   └── .index.pkm      (encrypted search index)

  Encryption:
    ✓ AES-256-GCM encryption
    ✓ PBKDF2 key derivation (100,000 iterations)
    ✓ Random salts and nonces per user and operation
    ✓ Zero plaintext storage

  Git Sync:
    - Encrypted notes are safe to push to git
    - NEVER commit .crypt file to git
    - .crypt stays local-only on each machine
    - Different machines can share .crypt via secure transfer

TROUBLESHOOTING:

  "User not found"
    → First time? Create user: pkm user init <username>

  "Decrypt failed (wrong password?)"
    → Password incorrect
    → To reset: pkm user passwd <username>

  "Note file not found"
    → Check note ID: pkm --user <username> note list
    → IDs are long UUIDs (not incremental)

  "Editor not opening"
    → Set $EDITOR: export EDITOR=nano
    → Default: vi (vim)

For more information:
  GitHub: https://github.com/sahay-shashank/personal-knowledge-manager
`)
}

// Help messages for each command
func (noteCmd *NoteCommand) Help() string {
	return `
NOTE OPERATIONS

USAGE:
  pkm --user <username> note <subcommand> [arguments]

SUBCOMMANDS:
  new <title>              Create a new note (opens $EDITOR)
  edit <note-id>           Edit an existing note
  get <note-id>            Display note content
  delete <note-id>         Delete a note
  list                     List all notes
  help                      Show this help message

EXAMPLES:
  $ pkm --user alice note new "Graph Theory"
  $ pkm --user alice note list
  $ pkm --user alice note get 550e8400-e29b
  $ pkm --user alice note edit 550e8400-e29b
  $ pkm --user alice note delete 550e8400-e29b

NOTES:
  • Editors: Uses $EDITOR environment variable (default: vi)
  • Format: Notes are stored as JSON with encryption
  • Links: Add links using 'link add' command
  • Tags: Add tags using 'tag add' command
`
}

func (userCmd *UserCommand) Help() string {
	return `
USER MANAGEMENT

USAGE:
  pkm user <subcommand> <username>

SUBCOMMANDS:
  init <username>         Create a new user account
  passwd <username>       Change password for existing user
  help                    Show this help message

EXAMPLES:
  $ pkm user init alice
  $ pkm user init bob
  $ pkm user passwd alice

SECURITY:
  • Passwords: Prompted interactively (never passed as argument)
  • Confirmation: Password change requires current password
  • Encryption: All user keys are encrypted in .crypt file
`
}

func (linkCmd *LinkCommand) Help() string {
	return `
LINK MANAGEMENT

USAGE:
  pkm --user <username> link <subcommand> [arguments]

SUBCOMMANDS:
  add <source-id> <target-id>    Create link from source to target
  remove <source-id> <target-id> Remove link between notes
  list <note-id>                  List all links from a note
  help                            Show this help message

EXAMPLES:
  $ pkm --user alice link add 550e8400-e29b 6ba7b810-9dad
  $ pkm --user alice link list 550e8400-e29b
  $ pkm --user alice link remove 550e8400-e29b 6ba7b810-9dad

ABOUT LINKS:
  • Directional: A→B is different from B→A
  • Backlinks: Automatically tracked for discovery
  • No cycles: Links create knowledge graph, not circular
  • UUID-based: Use full note IDs for accuracy
`
}

func (tagCmd *TagCommand) Help() string {
	return `
TAG MANAGEMENT

USAGE:
  pkm --user <username> tag <subcommand> [arguments]

SUBCOMMANDS:
  add <note-id> <tags>           Add tags to a note (comma-separated)
  remove <note-id> <tags>        Remove tags from a note
  list <note-id>                 List all tags on a note
  help                           Show this help message

EXAMPLES:
  $ pkm --user alice tag add 550e8400-e29b "learning,graphs"
  $ pkm --user alice tag list 550e8400-e29b
  $ pkm --user alice tag remove 550e8400-e29b "learning"

TAG GUIDELINES:
  • Format: lowercase, hyphen-separated (e.g., machine-learning)
  • Comma-separated: Use commas without spaces
  • Case-insensitive: Automatically converted to lowercase
  • Unique per note: Duplicate tags are rejected
  • Searchable: Search by tag with 'search tag' command
`
}

func (searchCmd *SearchCommand) Help() string {
	return `
NOTE SEARCH

USAGE:
  pkm --user <username> search <type> <terms...>

SEARCH TYPES:
  keyword <term1> [term2] ...   Search by keywords in title/content
  tag <tag1> [tag2] ...         Search by tags (intersection)

EXAMPLES:
  $ pkm --user alice search keyword "graph theory"
  $ pkm --user alice search keyword recursion
  $ pkm --user alice search tag learning
  $ pkm --user alice search tag productivity algorithms

SEARCH BEHAVIOR:
  • Keywords: Case-insensitive substring match in title and content
  • Multiple keywords: AND logic (all must be present)
  • Tags: Case-insensitive exact match (intersection if multiple)
  • Results: Returns note IDs and titles
  • Index: Uses built-in keyword/tag index for speed
`
}
