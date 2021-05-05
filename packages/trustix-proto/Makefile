build:
	$(MAKE) -C schema
	$(MAKE) -C protocol
	$(MAKE) -C rpc
	$(MAKE) -C api
	bash -c 'cd ../trustix-python && nix-shell --run "make build"'

test:
	true  # noop, no tests for codegen

lint:
	@true  # noop, no linter for codegen

format:
	gofmt -d .
