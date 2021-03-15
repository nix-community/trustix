module github.com/tweag/trustix/packages/trustix-nix

go 1.15

require (
	github.com/bakins/logrus-middleware v0.0.0-20180426214643-ce4c6f8deb07
	github.com/bakins/test-helpers v0.0.0-20141028124846-af83df64dc31 // indirect
	github.com/coreos/go-systemd v0.0.0-20190321100706-95778dfbb74e
	github.com/dgraph-io/badger/v2 v2.2007.2 // indirect
	github.com/golang/protobuf v1.5.0
	github.com/sirupsen/logrus v1.2.0
	github.com/spf13/cobra v1.0.1-0.20201006035406-b97b5ead31f7
	github.com/stretchr/testify v1.7.0
	github.com/tweag/trustix/packages/trustix v0.0.0-20201216011910-cb45e22716fa
	github.com/tweag/trustix/packages/trustix-proto v0.0.0-00010101000000-000000000000
	github.com/ulikunitz/xz v0.5.10
	golang.org/x/crypto v0.0.0-20200728195943-123391ffb6de // indirect
	google.golang.org/grpc v1.36.0
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.1.0 // indirect
	google.golang.org/protobuf v1.26.0
)

replace (
	github.com/tweag/trustix/packages/trustix => ../trustix
	github.com/tweag/trustix/packages/trustix-proto => ../trustix-proto
)
