# `apt-container` &mdash; a cleaner way to build your Docker images

[![Go Reference](https://pkg.go.dev/badge/github.com/fsufitch/apt-container.svg)](https://pkg.go.dev/github.com/fsufitch/apt-container)

`apt-container` is a wrapper around `apt-get` to help keep your Dockerfiles/Containerfiles slim, while still using good layering practices for Apt.

### Example: Best-practices normal Dockerfile

```dockerfile
# Installing straight up
RUN apt-get update -q && \
    apt-get install -y -q git python3 figlet && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists

# Installing from a file
COPY mypackages.txt ./
RUN apt-get update -q && \
    apt-get install -y -q $(cat mypackages.txt) && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists
# Better hope the file is not CRLF or has no comments!
```

Using `apt-container` to reduce the Apt commands to one-liners:

```dockerfile
# Installing straight up
RUN apt-container install git python3 figlet

# Installing from a file
COPY mypackages.txt ./
RUN apt-container install -r mypackages.txt
```

Install me by using this command, or by building the binary from source.

```bash
go install github.com/fsufitch/apt-container@latest
```

## Usage
```
$ apt-container install --help
NAME:
   apt-container install

USAGE:
   apt-container install [command options] packages...

DESCRIPTION:
   container-friendly version of 'apt-get install' (see apt-get man pages)

OPTIONS:
   --interactive                    run interactively (include -y in apt-get commands) (default: false)
   --simulate, -s                   run in simulated mode (include -s in apt-get commands) (default: false)
   --no-update, -U                  skip running 'apt-get update' before installing (default: false)
   --keep-cache, -C                 keep package cache after install; don't run 'apt-get clean' (default: false)
   --keep-lists, -L                 keep package lists used for install (default: false)
   --apt-lists-dir value            directory that apt keeps its stuff (default: "/var/lib/apt/lists")
   --extra-options value, -o value  options to passthrough to the apt-get command
   --quiet value, -q value          how many '-q' to pass to apt-get (default: 1)
   --requirements FILE, -r FILE     read package list from FILE (one per line, # is a comment); mutually exclusive with packages as arguments
   --help, -h                       show help
```