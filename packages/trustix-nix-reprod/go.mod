module github.com/nix-community/trustix/packages/trustix-nix-reprod

go 1.18

replace (
	github.com/nix-community/trustix/packages/go-lib => ../go-lib
	github.com/nix-community/trustix/packages/trustix => ../trustix
	github.com/nix-community/trustix/packages/trustix-nix => ../trustix-nix
	github.com/nix-community/trustix/packages/trustix-proto => ../trustix-proto
	github.com/nix-community/trustix/packages/unixtransport => ../unixtransport
)

require (
	github.com/BurntSushi/toml v1.2.0
	github.com/adrg/xdg v0.4.0
	github.com/bakins/logrus-middleware v0.0.0-20180426214643-ce4c6f8deb07
	github.com/bufbuild/connect-go v0.5.0
	github.com/buger/jsonparser v1.1.1
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf
	github.com/hashicorp/golang-lru v0.5.4
	github.com/kyleconroy/sqlc v1.15.0
	github.com/mattn/go-sqlite3 v1.14.15
	github.com/nix-community/go-nix v0.0.0-20220906172053-6b0185c1190b
	github.com/nix-community/trustix/packages/go-lib v0.0.0-20220906074512-b3b86d141df1
	github.com/nix-community/trustix/packages/trustix v0.0.0-20220906074512-b3b86d141df1
	github.com/nix-community/trustix/packages/trustix-nix v0.0.0-20220906074512-b3b86d141df1
	github.com/nix-community/trustix/packages/trustix-proto v0.0.0-20220906074512-b3b86d141df1
	github.com/pbnjay/memory v0.0.0-20210728143218-7b4eea64cf58
	github.com/pressly/goose v2.7.0+incompatible
	github.com/pressly/goose/v3 v3.7.0
	github.com/sirupsen/logrus v1.9.0
	github.com/spf13/cobra v1.5.0
	golang.org/x/net v0.0.0-20220826154423-83b083e8dc8b
	google.golang.org/protobuf v1.28.1
)

require (
	github.com/ClickHouse/clickhouse-go v1.5.4 // indirect
	github.com/antlr/antlr4/runtime/Go/antlr v0.0.0-20220626175859-9abda183db8e // indirect
	github.com/bytecodealliance/wasmtime-go v0.39.0 // indirect
	github.com/cloudflare/golz4 v0.0.0-20150217214814-ef862a3cdc58 // indirect
	github.com/cznic/mathutil v0.0.0-20181122101859-297441e03548 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/denisenkom/go-mssqldb v0.12.2 // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/golang-sql/civil v0.0.0-20220223132316-b832511892a9 // indirect
	github.com/golang-sql/sqlexp v0.1.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/hashicorp/go-uuid v1.0.3 // indirect
	github.com/inconshreveable/mousetrap v1.0.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/klauspost/cpuid/v2 v2.1.1 // indirect
	github.com/lib/pq v1.10.6 // indirect
	github.com/minio/sha256-simd v1.0.0 // indirect
	github.com/mr-tron/base58 v1.2.0 // indirect
	github.com/multiformats/go-multihash v0.2.1 // indirect
	github.com/multiformats/go-varint v0.0.6 // indirect
	github.com/nix-community/trustix/packages/unixtransport v0.0.0-20220906074512-b3b86d141df1 // indirect
	github.com/pganalyze/pg_query_go/v2 v2.1.2 // indirect
	github.com/pingcap/errors v0.11.5-0.20210425183316-da1aaba5fb63 // indirect
	github.com/pingcap/log v0.0.0-20210906054005-afc726e70354 // indirect
	github.com/pingcap/tidb/parser v0.0.0-20220725134311-c80026e61f00 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20220927061507-ef77025ab5aa // indirect
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/ziutek/mymysql v1.5.4 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	go.uber.org/zap v1.19.1 // indirect
	golang.org/x/crypto v0.0.0-20220926161630-eccd6366d1be // indirect
	golang.org/x/exp v0.0.0-20220927162542-c76eaa363f9d // indirect
	golang.org/x/sys v0.0.0-20220927170352-d9d178bc13c6 // indirect
	golang.org/x/text v0.3.7 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	lukechampine.com/blake3 v1.1.7 // indirect
	modernc.org/cc/v3 v3.39.0 // indirect
	modernc.org/libc v1.20.0 // indirect
	modernc.org/sqlite v1.19.1 // indirect
)
