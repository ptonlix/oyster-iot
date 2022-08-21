# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=oyster-iot
BINARY_UNIX=$(BINARY_NAME)_unix
CONFPATH=./conf/app.conf
DEVCONFPATH=./devaccess/config.ini
FRONTENDSOURCE=static

APP_NAME=${BINARY_NAME}
APP_VERSION= 0.0.8v
BUILD_VERSION=$(shell git log -1 --oneline)
BUILD_TIME=$(shell date )
GIT_REVISION=$(shell git rev-parse --short HEAD)
GIT_BRANCH=$(shell git name-rev --name-only HEAD)
GO_VERSION=$(shell go version)
VERSIONINFO = "-s -X 'main.AppName=${APP_NAME}' \
            -X 'main.AppVersion=${APP_VERSION}' \
            -X 'main.BuildVersion=${BUILD_VERSION}' \
            -X 'main.BuildTime=${BUILD_TIME}' \
            -X 'main.GitRevision=${GIT_REVISION}' \
            -X 'main.GitBranch=${GIT_BRANCH}' \
            -X 'main.GoVersion=${GO_VERSION}'" \

define packconfig
	@echo "package config ......."
	mkdir package
	mkdir package/conf
	mkdir package/devaccess
	mkdir package/log
	mv $(BINARY_NAME) ./package/
	cp $(CONFPATH) package/conf/
	cp $(DEVCONFPATH) package/devaccess/
	cp -r $(FRONTENDSOURCE) package/
	tar -zcvf oyster-iot-$(APP_VERSION).tar.gz package
	rm -rf package
endef
all: test build
build:
	$(GOBUILD) -ldflags $(VERSIONINFO) -o $(BINARY_NAME) -v
test:
	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)
beego:
	bee run

package: 
	$(GOBUILD) -ldflags $(VERSIONINFO) -o $(BINARY_NAME) -v 
	$(call packconfig)
	@echo "package complete!"

package-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags $(VERSIONINFO) -o $(BINARY_NAME) -v
	$(call packconfig)
	@echo "package complete!"

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v
docker-build:
	docker run --rm -it -v "$(GOPATH)":/go -w /go/src/bitbucket.org/rsohlich/makepost golang:latest go build -o "$(BINARY_UNIX)" -v