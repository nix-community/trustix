.PHONY: all build nix test run-build run-test

all: build

build:
	go build

nix:
	vgo2nix

test:
	./dev/test

# All commands prefixed with run- are meant to be implementing some kind of watch-mode for development
run-build:
	reflex -r \.go$$ make build

run-test:
	reflex -g trustix make test

run-mysql:
	./dev/mysql
