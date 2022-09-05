module github.com/nix-community/trustix/packages/trustix-nix-reprod

go 1.18

replace (
	github.com/nix-community/trustix/packages/go-lib => ../go-lib
	github.com/nix-community/trustix/packages/trustix => ../trustix
	github.com/nix-community/trustix/packages/trustix-proto => ../trustix-proto
	github.com/nix-community/trustix/packages/unixtransport => ../unixtransport
)

require (
	github.com/adrg/xdg v0.4.0
	github.com/bufbuild/connect-go v0.4.0
	github.com/buger/jsonparser v1.1.1
	github.com/hashicorp/golang-lru v0.5.4
	github.com/kyleconroy/sqlc v1.15.0
	github.com/nix-community/go-nix v0.0.0-20220822154651-3df711b31eb2
	github.com/nix-community/trustix/packages/go-lib v0.0.0-00010101000000-000000000000
	github.com/nix-community/trustix/packages/trustix v0.0.0-20220831055858-ad6617ff041f
	github.com/nix-community/trustix/packages/trustix-proto v0.0.0-20220831055858-ad6617ff041f
	github.com/pbnjay/memory v0.0.0-20210728143218-7b4eea64cf58
	github.com/pressly/goose v2.7.0+incompatible
	github.com/pressly/goose/v3 v3.7.0
	github.com/sirupsen/logrus v1.9.0
	github.com/spf13/cobra v1.5.0
	google.golang.org/protobuf v1.28.1
	modernc.org/sqlite v1.18.1
)

require (
	github.com/ClickHouse/clickhouse-go v1.5.4 // indirect
	github.com/antlr/antlr4/runtime/Go/antlr v1.4.10 // indirect
	github.com/benbjohnson/clock v1.3.0 // indirect
	github.com/bytecodealliance/wasmtime-go v0.40.0 // indirect
	github.com/cloudflare/golz4 v0.0.0-20150217214814-ef862a3cdc58 // indirect
	github.com/cznic/mathutil v0.0.0-20181122101859-297441e03548 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/denisenkom/go-mssqldb v0.12.2 // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/golang-sql/civil v0.0.0-20220223132316-b832511892a9 // indirect
	github.com/golang-sql/sqlexp v0.1.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51 // indirect
	github.com/klauspost/cpuid/v2 v2.1.1 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/lib/pq v1.10.6 // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	github.com/mattn/go-sqlite3 v1.14.15 // indirect
	github.com/minio/sha256-simd v1.0.0 // indirect
	github.com/mr-tron/base58 v1.2.0 // indirect
	github.com/multiformats/go-multihash v0.2.1 // indirect
	github.com/multiformats/go-varint v0.0.6 // indirect
	github.com/nix-community/trustix/packages/unixtransport v0.0.0-00010101000000-000000000000 // indirect
	github.com/pganalyze/pg_query_go/v2 v2.1.2 // indirect
	github.com/pingcap/errors v0.11.5-0.20210425183316-da1aaba5fb63 // indirect
	github.com/pingcap/log v1.1.0 // indirect
	github.com/pingcap/tidb/parser v0.0.0-20220905075856-6b8cf9d5b29b // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20200410134404-eec4a21b6bb0 // indirect
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/ziutek/mymysql v1.5.4 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	go.uber.org/zap v1.23.0 // indirect
	golang.org/x/crypto v0.0.0-20220829220503-c86fa9a7ed90 // indirect
	golang.org/x/exp v0.0.0-20220827204233-334a2380cb91 // indirect
	golang.org/x/mod v0.6.0-dev.0.20220419223038-86c51ed26bb4 // indirect
	golang.org/x/sys v0.0.0-20220829200755-d48e67d00261 // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/tools v0.1.12 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	lukechampine.com/blake3 v1.1.7 // indirect
	lukechampine.com/uint128 v1.2.0 // indirect
	modernc.org/cc/v3 v3.38.0 // indirect
	modernc.org/ccgo/v3 v3.16.9 // indirect
	modernc.org/libc v1.18.0 // indirect
	modernc.org/mathutil v1.5.0 // indirect
	modernc.org/memory v1.3.0 // indirect
	modernc.org/opt v0.1.3 // indirect
	modernc.org/strutil v1.1.3 // indirect
	modernc.org/token v1.0.1 // indirect
)
