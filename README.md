<div align="center">
    <img
        height="400"
        alt="timetreat-logo"
        src="assets/timetreat_blister.png"
    />
    <p><b>An easy-to-use and compatible time-tracking command line tool.</b></p>
</div>

<br>

# Setup

## Installation

### Release

**TODO**

Using the installer script:

- Linux & macOS:
```sh
curl -sS https://raw.githubusercontent.com/Criomby/timetreat/refs/heads/main/scripts/install_from_release.sh | sh
```

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

**See below for typical workflow examples.**

## Commands

> Use "timetreat [command] --help" for more information about a command.

- **start**

*Start a new tracking entry.*

```shell
# without a project name start now
timetreat start

# with project
timetreat start -p subscript

# long project name
timetreat start -p "project code completion"

# start tracking afterwards (e.g. it is now 9 a.m. but started at 8:30)
timetreat start -t 08:30

# round start time to 15 mins
# e.g. it's now 8:16
timetreat start -r 15m  # starts at 8:15
# or round to 30 mins
timetreat start -r 15m  # starts at 8:30
```

- **stop**

*Stop the currently running task.*

```shell
timetreat stop

# stop in 10 mins (assume it's 11 a.m.)
timetreat stop -t 11:10

# round stop time to 15 mins
# e.g. if it's 11:23, will stop at 11:30
timetreat stop -r 15m

# add/append description
timetreat stop -d "some additional detail"
```

- **current**

*Show current task information.*

```shell
# get full current task info
timetreat current

# alias to cur
timetreat cur

# get only project or only duration of current task
timetreat cur --project
timetreat cur --duration
# same as:
timetreat cur -p
timetreat cur -d
# etc.
```

- **list**

*List activities in log.*

```shell
timetreat list

# is aliased to ls
timetreat ls

# show each entry's duration
timetreat ls -d
```

- **last**

*Get last used projects.*

```shell
timetreat last
```

- **export**

*Export log in various formats to a separate file.*

```shell
# exports to csv relative to log file by default
timetreat export

# export to specific dir
timetreat export -d ~/Downloads
```

- **check**

*Verify the integrity of the log file.*

```shell
timetreat check
```

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

- :white_check_mark: Make the location & name of the log file configurable via env var
- :white_check_mark: Export log as csv
- :white_check_mark: Start/stop from a specific time
- :white_check_mark: Round start/stop time (e.g. `stop --round 15m`)
- Ask to replace project name and/or description if provided but already set
- Add optional tags to entries
- Generate customizable reports
  - Cli tables
  - Markdown
  - Html

If you'd like to request a feature or have feedback just open a new issue or comment on an existing one.

<br>

# Attributions

This project was heavily inspired by [`bartib`](https://github.com/nikolassv/bartib) by nikolassv.

I used his project for work to keep track of what I had been working on and how long my working days were becoming and it did a great job at it. However, one of the main drawbacks of the program were the non-standard logging format which hasn't changed despite various PRs and for which I ended up writing a Python script to automate the conversion to JSON entries for archiving and further processing. This is the main reason why I created timetreat and I came up with some more quality-of-life improvements along the way.

## Other Projects

| Project    | Description |
| -------- | ------- |
| [taskwarrior](https://github.com/GothenburgBitFactory/taskwarrior) | <ul><li>Focus on task management</li><li>Powerful & (very) complex cli</li><li>Highly customizable</li><li>C++</li></ul> |
| [Watson](https://github.com/jazzband/Watson) | <ul><li>Edit entries and generate reports</li><li>Simple interface</li><li>Python</li><li>unmaintained (last commit three years ago)</li></ul> |
| [timetrap](https://github.com/samg/timetrap) | <ul><li>Manage entries in timesheets</li><li>Uses natural language</li><li>Various export formats</li><li>Ruby</li></ul> |
| [timetrace](https://github.com/dominikbraun/timetrace)  | <ul><li>Simple and basic feature set</li><li>Only json exports</li><li>Go</li><li>unmaintained (last commit two years ago)</li></ul> |
