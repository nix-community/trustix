build: doc
	protoc --go_out=. --go_opt=paths=source_relative \
		--connect-go_out=. --connect-go_opt=paths=source_relative \
		-I=. \
		api/api.proto

	protoc --go_out=. --go_opt=paths=source_relative \
		--connect-go_out=. --connect-go_opt=paths=source_relative \
		-I=. \
		rpc/rpc.proto

	protoc --go_out=. --go_opt=paths=source_relative \
		schema/*.proto

doc:
	protoc --doc_out=. --doc_opt=markdown,doc.md -I=. \
	*/*.proto
