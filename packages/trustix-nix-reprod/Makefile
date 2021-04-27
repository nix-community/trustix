build:
	bash -c "cd trustix_nix_reprod/js && nix-shell --run 'npm run build'"
	true  # noop, no build step for python

test:
	@true

lint:
	true  # noop, no linter

format:
	black --quiet --check .
