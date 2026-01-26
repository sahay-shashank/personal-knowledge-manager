# Personal Knowledge Manager 📙
A terminal-first note-taking application written in Go, focused on secure storage, fast retrieval, and clean filesystem design.
## Features 🚀

- [x] New Note Creation
- [x] Note Editing
- [x] Note Linking
- [x] Note Taging
- [x] Note Indexing for faster search
- [x] Note Encryption on disk
- [x] Multi-user system
- [x] Note Listing
- [ ] Multi-machine syncing system
- [ ] Terminal User Interface

## How-To Guide

### Build Binary ⚙️

- **STEP 1:** Clone the repo and change directory
    ```bash
    git clone https://github.com/sahay-shashank/personal-knowledge-manager.git
    cd personal-knowledge-manager/
    ```
-  **STEP 2:** The build process can be triggered using the Taskfile. Perform the below mentioned command:
    ```bash
    $ task build
    ```
- **STEP 3:** To add the binary to be accessible from all terminal, export the `$PATH` or add it to your profile (`~/.bashrc` or `~/.zshenv`) using the following command:
    - Get your working directory:
        ```bash
        $ echo $(pwd)
        ```
    - Add to profile or run in the terminal
        ```bash
        export PATH=$PATH:"<your working directory>/build/"
        ```
- **STEP 4:** All done!! Hurrah!! 🙌 🎉

### Download from the release page of github

👨‍💻 Working on it!! Yet to build the release pipeline.

### Download using your package manager

👨‍💻 Working on it!! Appreciate your patience.
