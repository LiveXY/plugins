module github.com/livexy/plugins

go 1.25

require (
	gitee.com/chunanyong/dm v1.8.22
	github.com/bwmarrin/snowflake v0.3.0
	github.com/bytedance/sonic v1.15.0
	github.com/chromedp/cdproto v0.0.0-20250803210736-d308e07a266d
	github.com/chromedp/chromedp v0.14.2
	github.com/emirpasic/gods v1.18.1
	github.com/goccy/go-json v0.10.5
	github.com/golang/snappy v1.0.0
	github.com/gonfva/docxlib v0.0.0-20210517191039-d8f39cecf1ad
	github.com/livexy/linq v1.1.5
	github.com/livexy/pkg v1.1.4
	github.com/livexy/plugin v1.1.0
	github.com/livexy/plugins/opengaussb/opengauss v0.0.0-20260118060101-6a140b0a683a
	github.com/redis/go-redis/v9 v9.17.2
	github.com/rs/xid v1.6.0
	github.com/thoas/go-funk v0.9.3
	github.com/xuri/excelize/v2 v2.10.0
	go.uber.org/zap v1.27.1
	golang.org/x/text v0.33.0
	gorm.io/driver/mysql v1.6.0
	gorm.io/driver/postgres v1.6.0
	gorm.io/gorm v1.31.1
	gorm.io/plugin/dbresolver v1.6.2
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	gitee.com/opengauss/openGauss-connector-go-pq v1.0.7 // indirect
	github.com/bytedance/gopkg v0.1.3 // indirect
	github.com/bytedance/sonic/loader v0.5.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/chromedp/sysutil v1.1.0 // indirect
	github.com/cloudwego/base64x v0.1.6 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-json-experiment/json v0.0.0-20251027170946-4849db3c2f7e // indirect
	github.com/go-sql-driver/mysql v1.9.3 // indirect
	github.com/gobwas/httphead v0.1.0 // indirect
	github.com/gobwas/pool v0.2.1 // indirect
	github.com/gobwas/ws v1.4.0 // indirect
	github.com/golang/glog v1.2.5 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgx/v5 v5.8.0 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/klauspost/cpuid/v2 v2.3.0 // indirect
	github.com/puzpuzpuz/xsync/v4 v4.3.0 // indirect
	github.com/richardlehane/mscfb v1.0.6 // indirect
	github.com/richardlehane/msoleps v1.0.6 // indirect
	github.com/tiendc/go-deepcopy v1.7.2 // indirect
	github.com/tjfoc/gmsm v1.4.1 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	github.com/xuri/efp v0.0.1 // indirect
	github.com/xuri/nfp v0.0.2-0.20250530014748-2ddeb826f9a9 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/arch v0.23.0 // indirect
	golang.org/x/crypto v0.47.0 // indirect
	golang.org/x/net v0.49.0 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/sys v0.40.0 // indirect
	golang.org/x/xerrors v0.0.0-20240903120638-7835f813f4da // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
)

replace github.com/livexy/plugins/opengaussb/opengauss => ./opengaussb/opengauss
