module github.com/nix-community/trustix/packages/trustix

go 1.15

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/coreos/go-systemd v0.0.0-20190321100706-95778dfbb74e
	github.com/dop251/goja v0.0.0-20210427212725-462d53687b0d
	github.com/golang/protobuf v1.5.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.3.0
	github.com/hashicorp/go-memdb v1.3.0
	github.com/hashicorp/go-uuid v1.0.1
	github.com/lazyledger/smt v0.0.0-20200827143353-42131aab296f
	github.com/nix-community/trustix/packages/go-lib v0.0.0-00010101000000-000000000000
	github.com/nix-community/trustix/packages/trustix-proto v0.0.0-00010101000000-000000000000
	github.com/sirupsen/logrus v1.2.0
	github.com/spf13/cobra v1.0.1-0.20201006035406-b97b5ead31f7
	github.com/stretchr/testify v1.8.0
	go.etcd.io/bbolt v1.3.5
	go.uber.org/multierr v1.6.0
	google.golang.org/grpc v1.36.0
)

replace github.com/nix-community/trustix/packages/trustix-proto => ../trustix-proto

replace github.com/nix-community/trustix/packages/go-lib => ../go-lib
