module github.com/tweag/trustix/contrib/trustix-nix

go 1.15

require (
	github.com/bakins/logrus-middleware v0.0.0-20180426214643-ce4c6f8deb07
	github.com/bakins/test-helpers v0.0.0-20141028124846-af83df64dc31 // indirect
	github.com/coreos/go-systemd v0.0.0-20190321100706-95778dfbb74e
	github.com/golang/protobuf v1.4.3
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.2.0 // indirect
	github.com/sirupsen/logrus v1.2.0
	github.com/spf13/cobra v1.0.0
	github.com/stretchr/testify v1.5.1
	github.com/tweag/trustix v0.0.0-20201216011910-cb45e22716fa
	github.com/ulikunitz/xz v0.5.10
	google.golang.org/grpc v1.35.0
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.1.0 // indirect
	google.golang.org/protobuf v1.25.0
)

replace github.com/tweag/trustix => ../../.
