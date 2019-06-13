VERSION := $(shell cat VERSION)

build:
	@echo Building $(VERSION) for linux, darwin and windows
	
	CGO_ENABLED=0 GOOS=linux go build -o api-connector-bulk-${VERSION}.linux -a -installsuffix cgo
	CGO_ENABLED=0 GOOS=darwin go build -o api-connector-bulk-${VERSION}.darwin -a -installsuffix cgo
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o api-connector-bulk-${VERSION}.exe -a -installsuffix cgo

clean:
	rm -f api-connector-bulk
	rm -f api-connector-bulk-${VERSION}*

release: clean build
	@echo Releasing $(VERSION)
	
	shasum -a 256 api-connector-bulk-${VERSION}.linux > api-connector-bulk-${VERSION}.linux.sha
	shasum -a 256 api-connector-bulk-${VERSION}.darwin > api-connector-bulk-${VERSION}.darwin.sha
	shasum -a 256 api-connector-bulk-${VERSION}.exe > api-connector-bulk-${VERSION}.exe.sha

	shasum -c api-connector-bulk-${VERSION}.linux.sha
	shasum -c api-connector-bulk-${VERSION}.darwin.sha
	shasum -c api-connector-bulk-${VERSION}.exe.sha