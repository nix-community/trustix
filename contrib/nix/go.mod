module github.com/tweag/trustix/contrib/nix

go 1.15

require (
	github.com/bakins/logrus-middleware v0.0.0-20180426214643-ce4c6f8deb07
	github.com/coreos/go-systemd v0.0.0-20190321100706-95778dfbb74e
	github.com/golang/protobuf v1.4.3
	github.com/sirupsen/logrus v1.2.0
	github.com/spf13/cobra v1.0.0
	github.com/stretchr/testify v1.4.0
	github.com/tweag/trustix v0.0.0-20201216011910-cb45e22716fa
	google.golang.org/grpc v1.33.1
	google.golang.org/protobuf v1.25.0
)

replace github.com/tweag/trustix => ../../.
