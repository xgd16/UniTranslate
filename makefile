.PHONY: all run clean help

# 项目名称
APP = uniTranslate

GO_ROOT =$(shell go env GOROOT)
RACE = -ldflags="-s -w" -pgo=auto -o
GLOBAL_CONFIG = CGO_ENABLED=0

ifeq ($(OS),Windows_NT)
    IS_WINDOWS := 1
endif

BUILD_CMD = $(if $(IS_WINDOWS), \
    	SET ${GLOBAL_CONFIG}&SET GOOS=$(1)&SET GOARCH=$(2)&"$(GO_ROOT)\bin\go" build $(RACE) .\bin\$(1)_$(2)\$(APP)$(3) .\main.go, \
    	${GLOBAL_CONFIG} GOOS=$(1) GOARCH=$(2) $(GO_ROOT)/bin/go build $(RACE) ./bin/$(1)_$(2)/$(APP)$(3) ./main.go)

## linux: 编译打包linux
.PHONY: linux-amd64
linux-amd64:
	$(call BUILD_CMD,linux,amd64)
.PHONY: linux-arm64
linux-arm64:
	$(call BUILD_CMD,linux,arm64)

## win: 编译打包win
.PHONY: win-amd64
win-amd64:
	$(call BUILD_CMD,windows,amd64,.exe)
.PHONY: win-arm64
win-arm64:
	$(call BUILD_CMD,windows,arm64,.exe)

## mac: 编译打包mac
.PHONY: mac-amd64
mac-amd64:
	$(call BUILD_CMD,darwin,amd64)
.PHONY: mac-arm64
mac-arm64:
	$(call BUILD_CMD,darwin,arm64)

.PHONY: darwin-arm64-lib
darwin-arm64-lib:
	go build $(RACE) -buildmode=c-shared -o ./bin/${APP}darwin_arm64_lib.so ./main.go

## 编译win，linux，mac平台
.PHONY: all
all:win-amd64 win-arm64 linux-amd64 linux-arm64 mac-amd64 mac-arm64

run:
	@go run ./

.PHONY: tidy
tidy:
	@go mod tidy

.PHONY: update
update:
	@go get -u

initCfg:
	@cp ./config.bak.yaml ./config.yaml

## 清理二进制文件
clean:
	@rm -f ./bin/*

help:
	@echo "make - 格式化 Go 代码, 并编译生成二进制文件"
	@echo "make mac-amd64 - 编译 Go 代码, 生成mac-amd64的二进制文件"
	@echo "make linux-amd64 - 编译 Go 代码, 生成linux-amd64二进制文件"
	@echo "make win-amd64 - 编译 Go 代码, 生成windows-amd64二进制文件"
	@echo "make mac-arm64 - 编译 Go 代码, 生成mac-arm64的二进制文件"
	@echo "make linux-arm64 - 编译 Go 代码, 生成linux-arm64二进制文件"
	@echo "make win-arm64 - 编译 Go 代码, 生成windows-arm64二进制文件"
	@echo "make tidy - 执行go mod tidy"
	@echo "make run - 直接运行 Go 代码"
	@echo "make clean - 移除编译的二进制文件"
	@echo "make all - 编译多平台的二进制文件"
	@echo "make update - 更新 mod 扩展库"