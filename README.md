# Smug - tmux session manager

Inspired by [tmuxinator](https://github.com/tmuxinator/tmuxinator) and [tmuxp](https://github.com/tmux-python/tmuxp).

Smug automates your [tmux](https://github.com/tmux/tmux) workflow. You can create a single configuration file, and smug will create all required windows and panes from it.

![gif](https://raw.githubusercontent.com/ivaaaan/gifs/master/smug.gif)

Configuration used in this GIF can be found [here](#example-2).

## Usage

`smug <command> <project>[:window name] [-w window name]`.

### Examples

To start/stop a project and all windows, run:

```console
xyz@localhost:~$ smug start project

xyz@localhost:~$ smug stop project
```

When you already have a running session, and you want to create only some windows from the configuration file, you can do something like this:

```console
xyz@localhost:~$ smug start project:window1

xyz@localhost:~$ smug start project:window1,window2

xyz@localhost:~$ smug start project -w window1

xyz@localhost:~$ smug start project -w window1 -w window2

xyz@localhost:~$ smug stop project:window1

xyz@localhost:~$ smug stop project -w window1 -w window2
```

## Installation

#### Prerequisite Tools

* [Git](https://git-scm.com/)
* [Go (we test it with the last 2 major versions)](https://golang.org/dl/)

#### Fetch from GitHub

The easiest is to clone Smug from GitHub and install it using `go-cli`:

```console
xyz@localhost:~$ cd /tmp
xyz@localhost:/tmp$ git clone https://github.com/ivaaaan/smug.git
xyz@localhost:/tmp$ cd smug
xyz@localhost:/tmp/smug$ go install
```

## Configuration

Configuration files stored in the `~/.config/smug` directory in the `YAML` format, e.g `~/.config/smug/your_project.yml`.

### Examples

#### Example 1


```yaml
session: blog

root: ~/Developer/blog

before_start:
  - docker-compose -f my-microservices/docker-compose.yml up -d # my-microservices/docker-compose.yml is a relative to `root`

stop:
  - docker stop $(docker ps -q)

windows:
  - name: code
    root: blog # a relative path to root
    manual: true # you can start this window only manually, using the -w arg
    commands:
      - docker-compose start
    panes:
      - type: horizontal
        root: .
        commands:
          - docker-compose exec php /bin/sh
          - clear

  - name: infrastructure
    root: ~/Developer/blog/my-microservices
    panes:
      - type: horizontal
        root: .
        commands:
          - docker-compose up -d
          - docker-compose exec php /bin/sh
          - clear
```
#### Example 2

```yaml
session: blog

root: ~/Code/blog

before_start:
  - docker-compose up -d

stop:
  - docker-compose stop

windows:
  - name: code
    commands:
      - vim app/dependencies.php
    panes:
      - type: horizontal
        commands:
          - make run-tests
  - name: ssh
    commands:
      - ssh -i ~/keys/blog.pem ubuntu@127.0.0.1
```
