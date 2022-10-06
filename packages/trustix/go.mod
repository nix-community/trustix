module github.com/nix-community/trustix/packages/trustix

go 1.18

require (
	github.com/BurntSushi/toml v1.2.0
	github.com/bufbuild/connect-go v1.0.0
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf
	github.com/dop251/goja v0.0.0-20221003171542-5ea1285e6c91
	github.com/hashicorp/go-memdb v1.3.3
	github.com/lazyledger/smt v0.2.0
	github.com/nix-community/trustix/packages/go-lib v0.0.0-20221006094532-116e5118302b
	github.com/nix-community/trustix/packages/trustix-proto v0.0.0-20221006094532-116e5118302b
	github.com/nix-community/trustix/packages/unixtransport v0.0.0-20221006094532-116e5118302b
	github.com/sirupsen/logrus v1.9.0
	github.com/spf13/cobra v1.5.0
	github.com/stretchr/testify v1.8.0
	go.etcd.io/bbolt v1.3.6
	go.uber.org/multierr v1.8.0
	golang.org/x/net v0.0.0-20221004154528-8021a29435af
	google.golang.org/protobuf v1.28.1
)

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dlclark/regexp2 v1.7.0 // indirect
	github.com/go-sourcemap/sourcemap v2.1.3+incompatible // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-uuid v1.0.3 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/inconshreveable/mousetrap v1.0.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	golang.org/x/exp v0.0.0-20221004215720-b9f4876ce741 // indirect
	golang.org/x/sys v0.0.0-20220928140112-f11e5e49a4ec // indirect
	golang.org/x/text v0.3.7 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/nix-community/trustix/packages/go-lib => ../go-lib
	github.com/nix-community/trustix/packages/trustix-proto => ../trustix-proto
	github.com/nix-community/trustix/packages/unixtransport => ../unixtransport
)
