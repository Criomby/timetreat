#!/usr/bin/env sh

# Builds all targets into ./build

set -e

cd "$(dirname "$(dirname "$(readlink -f "$0")")")"

mkdir -p build

targets=(
    "linux/amd64"
    "linux/arm64"
    "darwin/arm64"
    "windows/amd64"
)

ARG_ARCHIVE="false"
while [[ $# -gt 0 ]]; do
  case $1 in
    --archive)
    ARG_ARCHIVE="true"
      shift
      ;;
    -*|--*)
      echo "Unknown option: $1"
      exit 1
      ;;
    *)
      echo "Unknown argument: $1"
      exit 1
      ;;
  esac
done

for target in "${targets[@]}"; do
    os="${target%%/*}"
    arch="${target#*/}"
    echo "building $os/$arch"
    if [ "${os}" = "windows" ]; then
        binary_name="timetreat.exe"
    else
        binary_name="timetreat"
    fi
    GOOS="${os}" GOARCH="${arch}" go build -o ./build/${binary_name}
    if [ "${ARG_ARCHIVE}" = "true" ]; then
        archive_name="${binary_name}_${os}_${arch}.tar.gz"
        cd build
        tar -czf ${archive_name} ${binary_name}
        rm ${binary_name}
        shasum -a 256 ${archive_name} > ${archive_name}.sha256
        cd - > /dev/null
    fi
    echo "DONE"
done
