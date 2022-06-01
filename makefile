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

define packconfig
	@echo "package config ......."
	mkdir package
	mkdir package/conf
	mkdir package/devaccess
	mv $(BINARY_NAME) ./package/
	cp $(CONFPATH) package/conf/
	cp $(DEVCONFPATH) package/devaccess/
	cp -r $(FRONTENDSOURCE) package/
	tar -zcvf oyster-iot.tar.gz package
	rm -rf package
endef
all: test build
build:
	$(GOBUILD) -o $(BINARY_NAME) -v
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
	$(GOBUILD) -o $(BINARY_NAME) -v
	$(call packconfig)
	@echo "package complete!"

package-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME) -v
	$(call packconfig)
	@echo "package complete!"

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v
docker-build:
	docker run --rm -it -v "$(GOPATH)":/go -w /go/src/bitbucket.org/rsohlich/makepost golang:latest go build -o "$(BINARY_UNIX)" -v