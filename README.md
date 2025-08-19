<div align="center">
    <img
        height="400"
        alt="timetreat-logo"
        src="https://github.com/user-attachments/assets/982a8376-9976-4b3d-83a7-9530a276782c"
    />
    <p><b>An easy-to-use and compatible time-tracking command line tool.</b></p>
</div>

<br>

# Setup

## Installation

### Release

**TODO**

Using the installer script:

- Linux: `curl | sh`
- macOS: `curl | sh`

Make sure `~/.local/bin` is in your PATH.

Or manually download the latest release from the releases section.

The installer script does not support Windows but I'd suggest you use WSL2 anyways and install the linux version.

### Source

```
git clone https://github.com/Criomby/timetreat.git
cd timetreat
make install

# or symlink to build in local repo for fast rebuilding of latest branch version
# make symlink
```

## Config

The only config option you need to decide on is the location of the log file.

The default creates/uses the log file at `~/timetreat.log`.

If you want to set a different file path and/or name you can set the environment variable `TIMETREAT_LOG`
with an absolute path to the file, e.g. `/home/user/timetreat/mytasks.log`. The parent directories must
already exist, the file will be created with the `start` command if it does not exist yet.

If you want to use multiple different log files see advanced usage below.

<br>

# How To

## Basics

1. Fist-time setup: Use the `start` command to create the log file. It will ask you to create a log file at the default or configured location.
2. Start a new entry with `start`. You don't have to specify a project name and description right away. You can leave it empty and when you `stop` the tracking later you can optionally provide the project name and an description or append to the existing description if already present. The use case for that is that I often found myself initially wanting to work on something but ended up working on something else or wanting to add additional items to the description I didn't think about when starting the task and before I `stop` the tracking without having to change existing entries or using multiple commands.

Typical workflow examples:

```shell
TODO
```

## Commands

> Use "timetreat [command] --help" for more information about a command.

- **start**

*Start a new tracking entry.*

TODO

- **stop**

*Stop the currently running task.*

TODO

- **current**

*Show current task information.*

TODO

- **list**

*List activities in log.*

TODO

- **last**

*Get last used projects.*

TODO

- **export**

*Export log in various formats to a separate file.*

TODO

- **check**

*Verify the integrity of the log file.*

TODO

<br>

## Advanced Usage

- Multiple log files

If you want to use different log files for different activities, you can pass the env var per command to set a config file, e.g. `TIMETREAT_LOG="/home/user/project_cocoa.log" timetreat start`.

If you frequently switch between different log files you can define a shell function to quickly switch between projects and set the right config file, e.g.

```shell
#!/usr/bin/env bash
function timetreat_wrapper() {
  if [[ "$1" -eq "cocoa" ]]; then
    TIMETREAT_LOG="/home/user/project_cocoa.log" timetreat ${@:2}
  elif [[ "$1" -eq "hardcopy" ]]; then
    TIMETREAT_LOG="/home/user/project_hardcopy.log" timetreat ${@:2}
  else
    echo "unknown log file for '$1'"
  fi
}
```

Then use the function to log to the cocoa log file like this: `timetreat_wrapper cocoa start -d "harvesting beans"`

- [Starship prompt integration](https://starship.rs/)

You can have information about your current task displayed in your starship prompt with this custom module:

```
[custom.timetreat]
command = "timetreat current --project --duration"
shell = ["sh", "--norc"]
when = "timetreat current | grep -q -v 'no task running'"
symbol = " "
```

This will show output in the format ` timetreat - 1h15m46s` at the end of your prompt. If no task is running, it will hide the prompt module completely.

Customize this according to your preferences, e.g. showing only project name.

<br>

# Roadmap

*In no particular order.*

- :white_check_mark: Make the location & name of the log file configurable via env var
- :white_check_mark: Export log as csv
- Ask to replace project name and/or description if provided but already set
- Generate customizable reports
  - Cli tables
  - Markdown
  - Html
- Start/stop from a specific time
- Round start/stop time (e.g. `stop --round 15m`)

If you'd like to request a feature or have feedback just open a new issue or comment on an existing one.

<br>

# Attributions

This project was heavily inspired by [`bartib`](https://github.com/nikolassv/bartib) by nikolassv.

I used his project for work to keep track of what I had been working on and how long my working days were becoming and it did a great job at it. However, one of the main drawbacks of the program were the non-standard logging format which hasn't changed despite various PRs and for which I ended up writing a Python script to automate the conversion to JSON entries for archiving and further processing. This is the main reason why I created timetreat and I came up with some more quality-of-life improvements along the way.

As I am no professional in Go (yet) I made use of Google Gemini in addition to the official docs as a wiki to ask various things about how Go works and how I can implement certain functions. However, it made quite severe mistakes along the way and it is currently difficult to get working code out of it. This led me to gaining a pretty good understanding of the language and the algorithms used while debugging, especially for reading & parsing files in buffers from the end of a file and the Go module system.

## Other Projects

| Project    | Status |
| -------- | ------- |
| [taskwarrior](https://github.com/GothenburgBitFactory/taskwarrior) | C++ |
| [Watson](https://github.com/jazzband/Watson) | Python, unmaintained (last commit three years ago) |
| [timetrap](https://github.com/samg/timetrap) | Ruby |
| [timetrace](https://github.com/dominikbraun/timetrace?tab=readme-ov-file#generate-a-report-beta)  | Go, unmaintained (last commit two years ago) |
