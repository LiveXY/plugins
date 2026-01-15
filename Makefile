PROJECT := plugins
BIN_PATH := $(CURDIR)/bin/

OS := $(if $(GOOS),$(GOOS),linux)
ARCH := $(if $(GOARCH),$(GOARCH),amd64)

RUN_PATH := ./
PARENT_PATH := ../
EXTENSION :=
ifeq ($(shell go env GOOS), windows)
	RUN_PATH :=
	PARENT_PATH := ..\\
	EXTENSION := .exe
endif

update-depend:
	go get -t -u -v ./... && go mod tidy

govulncheck:
	~/go/bin/govulncheck ./...

gosec:
	~/go/bin/gosec ./...

build-plugins:
	#go build -ldflags="-s -w" -buildmode=plugin -o=../../test/bin/plugins/opengaussb.so ./opengaussb/main.go
	go build -ldflags="-s -w" -buildmode=plugin -o=../../test/bin/plugins/excel.so ./excel/main.go
	go build -ldflags="-s -w" -buildmode=plugin -o=../../test/bin/plugins/word.so ./word/main.go
	go build -ldflags="-s -w" -buildmode=plugin -o=../../test/bin/plugins/gob-stream.so ./gob-stream/main.go
	go build -ldflags="-s -w" -buildmode=plugin -o=../../test/bin/plugins/sonic-json.so ./sonic-json/main.go
	go build -ldflags="-s -w" -buildmode=plugin -o=../../test/bin/plugins/go-json.so ./go-json/main.go
	go build -ldflags="-s -w" -buildmode=plugin -o=../../test/bin/plugins/mysql.so ./mysql/main.go
	go build -ldflags="-s -w" -buildmode=plugin -o=../../test/bin/plugins/xid.so ./xid/main.go
	go build -ldflags="-s -w" -buildmode=plugin -o=../../test/bin/plugins/snowflake.so ./snowflake/main.go
	go build -ldflags="-s -w" -buildmode=plugin -o=../../test/bin/plugins/redis.so ./redis/main.go
	go build -ldflags="-s -w" -buildmode=plugin -o=../../test/bin/plugins/local-fs.so ./local-fs/main.go
	go build -ldflags="-s -w" -buildmode=plugin -o=../../test/bin/plugins/pgsql.so ./pgsql/main.go
	go build -ldflags="-s -w" -buildmode=plugin -o=../../test/bin/plugins/dameng.so ./dameng/main.go
	go build -ldflags="-s -w" -buildmode=plugin -o=../../test/bin/plugins/mxtong-sms.so ./mxtong-sms/main.go
	go build -ldflags="-s -w" -buildmode=plugin -o=../../test/bin/plugins/chrome2pdf.so ./chrome2pdf/main.go
