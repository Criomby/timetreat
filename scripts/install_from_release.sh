#!/usr/bin/env bash

# This script gets the latest release from GitHub and installs it to ~/.local/bin

timetreat_ascii_logo=" _   _                _                  _
| | (_)              | |                | |
| |_ _ _ __ ___   ___| |_ _ __ ___  __ _| |_
| __| | '_ ' _ \ / _ \ __| '__/ _ \/ _' | __|
| |_| | | | | | |  __/ |_| | |  __/ (_| | |_
 \__|_|_| |_| |_|\___|\__|_|  \___|\__,_|\__|"

set -e

latest_release_version="$(curl -s --show-headers https://github.com/Criomby/timetreat/releases/latest | grep "location:" | sed "s/.*\///" | sed "s/\r//")"
latest_semver="${latest_release_version#*v}"
download_url="https://github.com/Criomby/timetreat/releases/download/${latest_release_version}"

detect_platform() {
  platform="$(uname -s | tr '[:upper:]' '[:lower:]')"
  case "${platform}" in
    msys_nt*) platform="pc-windows-msvc" ;;
    cygwin_nt*) platform="pc-windows-msvc";;
    mingw*) platform="pc-windows-msvc" ;;
  esac
  # no installer support on Windows
  if [ "{$platform}" = "pc-windows-msvc" ]; then
    echo "Error: Please download & install the exe manually on Windows."
    exit 1
  fi
  printf '%s' "${platform}"
}

detect_arch() {
  arch="$(uname -m | tr '[:upper:]' '[:lower:]')"
  case "${arch}" in
    x86_64) arch="amd64" ;;
    armv*) arch="arm" ;;
  esac
  # `uname -m` in some cases mis-reports 32-bit OS as 64-bit, so double check
  # 32-bit targets will have to be built from source
  if [ "${arch}" = "amd64" ] && [ "$(getconf LONG_BIT)" -eq 32 ]; then
    arch=i686
  elif [ "${arch}" = "arm64" ] && [ "$(getconf LONG_BIT)" -eq 32 ]; then
    arch=arm
  fi
  # adjust for armv* 64 bit targets
  if [ "${arch}" = "arm" ] && [ "$(getconf LONG_BIT)" -eq 64 ]; then
    arch=arm64
  fi
  printf '%s' "${arch}"
}

confirm() {
  if [ -z "${FORCE-}" ]; then
    printf "%s " "$* [y/N]"
    set +e
    read -r yn </dev/tty
    rc=$?
    set -e
    if [ $rc -ne 0 ]; then
      echo "Error reading from prompt (please re-run with the '--yes' option)"
      exit 1
    fi
    if [ "$yn" != "y" ] && [ "$yn" != "yes" ]; then
      echo 'Aborting (please answer "yes" to continue)'
      exit 1
    fi
  fi
}

echo -e "${timetreat_ascii_logo}

Install script
"

if [ "${latest_release_version}" != "" ]; then
    echo -e "Latest release version: ${latest_release_version}\n"
else
    echo "Error: Could not get latest release version from repo"
    exit 1
fi

platform=$(detect_platform)
arch=$(detect_arch)

# currently builds are available for
targets=(
    "linux/amd64"
    "linux/arm64"
    "darwin/arm64"
    "windows/amd64"
)
found=0
for target in "${targets[@]}"; do
    if [ "${platform}/${arch}" = "${target}" ]; then
        found=1
        break
    fi
done
if [ "${found}" -eq 0 ]; then
    echo "Error: no prebuilt binary available for ${platform}/${arch}\nPlease build from source"
    exit 1
fi

download_asset_name="timetreat_${latest_release_version:1}_${platform}_${arch}.tar.gz"
asset_download_url="${download_url}/${download_asset_name}"

confirm "Install?"

echo -e "\nDownloading from ${download_url}\n"

temp_dir="$(mktemp -d)"
cd ${temp_dir}

echo -e "${download_asset_name}"
curl -LO ${asset_download_url}
echo
echo -e "${download_asset_name}.sha256"
curl -LO ${asset_download_url}.sha256
echo
shasum -c ${download_asset_name}.sha256
echo
tar -xzf ${download_asset_name}
/bin/cp -f timetreat ~/.local/bin/

rm -r ${temp_dir}

echo -e "Timetreat installed successfully!"
