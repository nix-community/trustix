.PHONY: contract run-contract all

all: build

build:
	go build

nix:
	vgo2nix

test:
	./trustix

# All commands prefixed with run- are meant to be implementing some kind of watch-mode for development
run-build:
	reflex -r \.go$$ make build

run-test:
	reflex -g trustix make test
