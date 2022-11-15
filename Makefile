PROJECT := etms
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

run-plugin:
	go build -gcflags="-m -m" -buildmode=plugin -o=../test/bin/plugins/excel.so ./excel/main.go
	go build -gcflags="-m -m" -buildmode=plugin -o=../test/bin/plugins/word.so ./word/main.go
	go build -gcflags="-m -m" -buildmode=plugin -o=../test/bin/plugins/gob-stream.so ./gob-stream/main.go
	go build -gcflags="-m -m" -buildmode=plugin -o=./bin/plugins/sonic-json.so ./plugins/sonic-json/main.go
	go build -gcflags="-m -m" -buildmode=plugin -o=./bin/plugins/mysql.so ./plugins/mysql/main.go
	go build -gcflags="-m -m" -buildmode=plugin -o=./bin/plugins/xid.so ./plugins/xid/main.go
	go build -gcflags="-m -m" -buildmode=plugin -o=./bin/plugins/snowflake.so ./plugins/snowflake/main.go
	go build -gcflags="-m -m" -buildmode=plugin -o=./bin/plugins/redis.so ./plugins/redis/main.go
	go build -gcflags="-m -m" -buildmode=plugin -o=./bin/plugins/local-fs.so ./plugins/local-fs/main.go
	go build -gcflags="-m -m" -buildmode=plugin -o=./bin/plugins/pgsql.so ./plugins/pgsql/main.go
	go build -gcflags="-m -m" -buildmode=plugin -o=./bin/plugins/dameng.so ./plugins/dameng/main.go
	go build -gcflags="-m -m" -buildmode=plugin -o=./bin/plugins/mxtong-sms.so ./plugins/mxtong-sms/main.go
