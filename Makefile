all: build test lint format

build:
	for pkg in packages/*; do \
	  bash -c "cd $$pkg && nix-shell --run 'make build'"; \
	done

test:
	for pkg in packages/*; do \
	  bash -c "cd $$pkg && nix-shell --run 'make test'"; \
	done

lint:
	for pkg in packages/*; do \
	  bash -c "cd $$pkg && nix-shell --run 'make lint'"; \
	done

format:
	nixpkgs-fmt --check .
	for pkg in packages/*; do \
	  bash -c "cd $$pkg && nix-shell --run 'make format'"; \
	done

develop:
	hivemind
