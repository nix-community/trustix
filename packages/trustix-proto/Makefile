all: go py

go:
	$(MAKE) -C schema
	$(MAKE) -C proto
	$(MAKE) -C api

py:
	bash -c 'cd ../trustix-python && nix-shell --run make'
