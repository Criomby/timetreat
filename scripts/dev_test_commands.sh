#!/usr/bin/env bash

# Script to run and test basic commands.
# Takes a positional argument for the testdata log file to use,
# e.g. empty.log to use the log file in testdata/empty.log.
# The log file is then copied in a tmp dir for the tests.

set -e

project_dir="$(dirname "$(dirname "$(readlink -f "$0")")")"

cd $project_dir

POSITIONAL_ARGS=()

ARG_LOGFILE=""
ARG_LOGFILE_DEFAULT="9_regular.log"

while [[ $# -gt 0 ]]; do
  case $1 in
    -*|--*)
      echo "Unknown option: $1"
      exit 1
      ;;
    *)
      POSITIONAL_ARGS+=("$1")
      shift
      ;;
  esac
done

set -- "${POSITIONAL_ARGS[@]}"

if [[ ${#POSITIONAL_ARGS[@]} -gt 1 ]]; then
    echo "too many arguments"
    exit 1
fi

ARG_LOGFILE=""${POSITIONAL_ARGS[0]}""

# set TIMETREAT_LOG env var
if [[ "${ARG_LOGFILE}" = "" ]]; then
    export TIMETREAT_LOG="${project_dir}/testdata/${ARG_LOGFILE_DEFAULT}"
else
    export TIMETREAT_LOG="${project_dir}/testdata/${ARG_LOGFILE}"
fi

function echoCommandBanner() {
    echo "

#################################

    ╭────────────────────────────╮
    │ ${1}
    ╰────────────────────────────╯
>>>
    "
}

function echoDescription() {
    echo "
---------------------------------------
> ${1}
---------------------------------------
    "
}

# ╭────────────────────────────╮
# │            root            │
# ╰────────────────────────────╯
echoCommandBanner ROOT

echoDescription "show help output"
go run .

# ╭────────────────────────────╮
# │            list            │
# ╰────────────────────────────╯
echoCommandBanner LIST

echoDescription "show two entries with delta"
go run . list -d -n 2

echoDescription "show 9 entries with delta"
go run . ls -d
