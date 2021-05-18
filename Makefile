all: build test lint doc format

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

doc:
	for pkg in packages/*; do \
	  bash -c "cd $$pkg && nix-shell --run 'make doc'"; \
	done

direnv-allow:
	find . -name .envrc -exec direnv allow {} \;

develop:
	hivemind
