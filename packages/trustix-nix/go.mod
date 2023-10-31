module github.com/nix-community/trustix/packages/trustix-nix

go 1.18

require (
	connectrpc.com/connect v1.12.0
	github.com/bakins/logrus-middleware v0.0.0-20180426214643-ce4c6f8deb07
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf
	github.com/nix-community/go-nix v0.0.0-20231012070617-9b176785e54d
	github.com/nix-community/trustix/packages/trustix v0.0.0-20231027092553-b0e71501e6f6
	github.com/nix-community/trustix/packages/trustix-proto v0.0.0-20231027042222-1fba619f3548
	github.com/prometheus/client_golang v1.17.0
	github.com/sirupsen/logrus v1.9.3
	github.com/spf13/cobra v1.7.0
	github.com/stretchr/testify v1.8.4
	github.com/ulikunitz/xz v0.5.11
)

require (
	github.com/bakins/test-helpers v0.0.0-20141028124846-af83df64dc31 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/matttproud/golang_protobuf_extensions/v2 v2.0.0 // indirect
	github.com/nix-community/trustix/packages/unixtransport v0.0.0-20231027042222-1fba619f3548 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.5.0 // indirect
	github.com/prometheus/common v0.45.0 // indirect
	github.com/prometheus/procfs v0.12.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/sys v0.13.0 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/nix-community/trustix/packages/go-lib => ../go-lib
	github.com/nix-community/trustix/packages/trustix => ../trustix
	github.com/nix-community/trustix/packages/trustix-proto => ../trustix-proto
	github.com/nix-community/trustix/packages/unixtransport => ../unixtransport
)
