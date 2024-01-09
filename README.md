# SSEx

SSEx is a terminal user interface (TUI) that allows you to execute commands on a remote server through SSH. 🪡

<p align="center">
  <img src="./assets/demo.gif") width="500"/>
</p>

## Installation

You can install the appropriate binary for your operating system by visiting the [Release page]().

**Note**:  
If you're on macOS, you may need to run:

```sh
xattr -c ./ssex_Darwin_x86_64.tar.gz
```

to (to avoid "unknown developer" warning)

## Pre-requisites

SSEx asssume you have setup SSH keys on your local machine and the remote server. \
For now SSEx looks for private key at `~/.ssh/id_rsa`.

### WIP:

- Add support for custom SSH keys
- Save connection details for future use
