VERSION := $(shell cat VERSION)

build:
	@echo Building $(VERSION) for darwin, linux, and windows-64
	
	CGO_ENABLED=0 GOOS=darwin go build -o api-connector-bulk-${VERSION}.darwin -a -installsuffix cgo
	CGO_ENABLED=0 GOOS=linux go build -o api-connector-bulk-${VERSION}.linux -a -installsuffix cgo
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o api-connector-bulk-${VERSION}.exe -a -installsuffix cgo

clean:
	rm -f api-connector-bulk*

release: clean build
	@echo Releasing $(VERSION)
	
	shasum -a 256 api-connector-bulk-${VERSION}.darwin > api-connector-bulk-${VERSION}.darwin.sha
	shasum -a 256 api-connector-bulk-${VERSION}.linux > api-connector-bulk-${VERSION}.linux.sha
	shasum -a 256 api-connector-bulk-${VERSION}.exe > api-connector-bulk-${VERSION}.exe.sha

	shasum -c api-connector-bulk-${VERSION}.darwin.sha
	shasum -c api-connector-bulk-${VERSION}.linux.sha
	shasum -c api-connector-bulk-${VERSION}.exe.sha

docker:
	@echo Building docker container for ${VERSION}, this expects a release in GitHub!

	docker build --build-arg VERSION=${VERSION} --tag=cybergrx/api-connector-bulk:${VERSION} -f Dockerfile .
	docker tag cybergrx/api-connector-bulk:${VERSION} cybergrx/api-connector-bulk:latest

	docker push cybergrx/api-connector-bulk:${VERSION}
	docker push cybergrx/api-connector-bulk:latest
	