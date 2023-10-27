module github.com/nix-community/trustix/packages/trustix-nix-r13y

go 1.18

replace (
	github.com/nix-community/trustix/packages/go-lib => ../go-lib
	github.com/nix-community/trustix/packages/trustix => ../trustix
	github.com/nix-community/trustix/packages/trustix-nix => ../trustix-nix
	github.com/nix-community/trustix/packages/trustix-proto => ../trustix-proto
	github.com/nix-community/trustix/packages/unixtransport => ../unixtransport
)

require (
	github.com/BurntSushi/toml v1.3.2
	github.com/adrg/xdg v0.4.0
	github.com/bakins/logrus-middleware v0.0.0-20180426214643-ce4c6f8deb07
	github.com/bufbuild/connect-go v1.10.0
	github.com/buger/jsonparser v1.1.1
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf
	github.com/hashicorp/golang-lru v1.0.2
	github.com/kyleconroy/sqlc v1.19.1
	github.com/mattn/go-sqlite3 v1.14.17
	github.com/nix-community/go-nix v0.0.0-20231012070617-9b176785e54d
	github.com/nix-community/trustix/packages/go-lib v0.0.0-20231027042222-1fba619f3548
	github.com/nix-community/trustix/packages/trustix v0.0.0-20231027042222-1fba619f3548
	github.com/nix-community/trustix/packages/trustix-nix v0.0.0-20231027042222-1fba619f3548
	github.com/nix-community/trustix/packages/trustix-proto v0.0.0-20231027042222-1fba619f3548
	github.com/pbnjay/memory v0.0.0-20210728143218-7b4eea64cf58
	github.com/pressly/goose v2.7.0+incompatible
	github.com/pressly/goose/v3 v3.15.1
	github.com/prometheus/client_golang v1.17.0
	github.com/sirupsen/logrus v1.9.3
	github.com/spf13/cobra v1.7.0
	golang.org/x/net v0.17.0
	google.golang.org/protobuf v1.31.0
)

require (
	github.com/ClickHouse/clickhouse-go v1.5.4 // indirect
	github.com/antlr/antlr4/runtime/Go/antlr v1.4.10 // indirect
	github.com/antlr/antlr4/runtime/Go/antlr/v4 v4.0.0-20230321174746-8dcc6526cfb1 // indirect
	github.com/benbjohnson/clock v1.3.5 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bytecodealliance/wasmtime-go v1.0.0 // indirect
	github.com/bytecodealliance/wasmtime-go/v8 v8.0.0 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/cloudflare/golz4 v0.0.0-20150217214814-ef862a3cdc58 // indirect
	github.com/cubicdaiya/gonp v1.0.4 // indirect
	github.com/cznic/mathutil v0.0.0-20181122101859-297441e03548 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/denisenkom/go-mssqldb v0.12.3 // indirect
	github.com/go-sql-driver/mysql v1.7.1 // indirect
	github.com/golang-sql/civil v0.0.0-20220223132316-b832511892a9 // indirect
	github.com/golang-sql/sqlexp v0.1.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/cel-go v0.18.1 // indirect
	github.com/hashicorp/go-uuid v1.0.3 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.4.3 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/klauspost/cpuid/v2 v2.2.5 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/matttproud/golang_protobuf_extensions/v2 v2.0.0 // indirect
	github.com/minio/sha256-simd v1.0.1 // indirect
	github.com/mr-tron/base58 v1.2.0 // indirect
	github.com/multiformats/go-multihash v0.2.3 // indirect
	github.com/multiformats/go-varint v0.0.7 // indirect
	github.com/nix-community/trustix/packages/unixtransport v0.0.0-20231027042222-1fba619f3548 // indirect
	github.com/pganalyze/pg_query_go/v2 v2.2.0 // indirect
	github.com/pganalyze/pg_query_go/v4 v4.2.3 // indirect
	github.com/pingcap/errors v0.11.5-0.20210425183316-da1aaba5fb63 // indirect
	github.com/pingcap/failpoint v0.0.0-20220801062533-2eaa32854a6c // indirect
	github.com/pingcap/log v1.1.0 // indirect
	github.com/pingcap/tidb/parser v0.0.0-20231013125129-93a834a6bf8d // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_model v0.5.0 // indirect
	github.com/prometheus/common v0.45.0 // indirect
	github.com/prometheus/procfs v0.12.0 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	github.com/riza-io/grpc-go v0.2.0 // indirect
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/stoewer/go-strcase v1.3.0 // indirect
	github.com/ziutek/mymysql v1.5.4 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.26.0 // indirect
	golang.org/x/crypto v0.14.0 // indirect
	golang.org/x/exp v0.0.0-20231006140011-7918f672742d // indirect
	golang.org/x/sync v0.4.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	google.golang.org/genproto v0.0.0-20231012201019-e917dd12ba7a // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20231016165738-49dd2c1f3d0b // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231016165738-49dd2c1f3d0b // indirect
	google.golang.org/grpc v1.59.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	lukechampine.com/blake3 v1.2.1 // indirect
	modernc.org/cc/v3 v3.41.0 // indirect
	modernc.org/libc v1.24.1 // indirect
	modernc.org/sqlite v1.26.0 // indirect
)
