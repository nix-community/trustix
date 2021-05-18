build: doc
	$(MAKE) -C schema
	$(MAKE) -C rpc
	$(MAKE) -C api
	bash -c 'cd ../trustix-python && nix-shell --run "make build"'

test:
	true  # noop, no tests for codegen

lint:
	@true  # noop, no linter for codegen

format:
	gofmt -d .

doc:
	protoc --doc_out=. --doc_opt=markdown,doc.md \
	-I=. -I=.. \
	*/*.proto
