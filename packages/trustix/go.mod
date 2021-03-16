module github.com/tweag/trustix/packages/trustix

go 1.15

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/Shopify/go-lua v0.0.0-20191113154418-05ce435a9edd
	github.com/coreos/go-systemd v0.0.0-20190321100706-95778dfbb74e
	github.com/dgraph-io/badger/v2 v2.2007.2
	github.com/golang/protobuf v1.5.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.3.0
	github.com/hashicorp/go-memdb v1.3.0
	github.com/lazyledger/smt v0.0.0-20200827143353-42131aab296f
	github.com/sirupsen/logrus v1.2.0
	github.com/spf13/cobra v1.0.1-0.20201006035406-b97b5ead31f7
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stretchr/testify v1.7.0
	github.com/tweag/trustix/packages/trustix-proto v0.0.0-00010101000000-000000000000
	github.com/ugorji/go v1.1.4 // indirect
	golang.org/x/crypto v0.0.0-20200728195943-123391ffb6de // indirect
	golang.org/x/net v0.0.0-20210119194325-5f4716e94777
	google.golang.org/grpc v1.36.0
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.1.0
	google.golang.org/protobuf v1.26.0
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace github.com/tweag/trustix/packages/trustix-proto => ../trustix-proto
