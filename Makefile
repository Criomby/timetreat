PROJECT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))

.PHONY: build
build:
	mkdir -p build
	go build -o ./build

.PHONY: install
install: build
	cp ./build/timetreat ~/.local/bin/timetreat

.PHONY: symlink
symlink: build
	ln -sf $(PROJECT_DIR)build/timetreat ~/.local/bin/timetreat

.PHONY: leaks
leaks:
	gitleaks --no-banner dir
	trufflehog git --log-level=-1 --fail file://.
