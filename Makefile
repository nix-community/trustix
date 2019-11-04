.PHONY: contract run-geth run-contract all

all: contract build

contract:
	abigen --sol contracts/registry.sol --pkg registry --type NarRegistry --out registry/registry.go

build:
	go build

nix:
	vgo2nix

test:
	./trustix

# All commands prefixed with run- are meant to be implementing some kind of watch-mode for development
run-build:
	reflex -r \.go$$ make build

run-geth:
	geth --ipcpath $(PWD)/geth.sock --dev

run-contract:
	reflex -r \.sol$$ make contract

run-test:
	reflex -g trustix make test
