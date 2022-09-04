build: doc
	$(MAKE) -C schema
	$(MAKE) -C rpc
	$(MAKE) -C api

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
