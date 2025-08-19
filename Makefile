PROJECT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))

.PHONY: install
install:
	go build
	cp ./timetreat ~/.local/bin/timetreat

.PHONY: symlink
symlink:
	go build
	ln -sf $(PROJECT_DIR)timetreat ~/.local/bin/timetreat

.PHONY: leaks
leaks:
	gitleaks --no-banner dir
	trufflehog git --log-level=-1 --fail file://.

