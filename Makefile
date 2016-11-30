PROG_NAME := "esevaluator"

setup: install tls

clean:
	go clean ./...

install:
	go get github.com/mitchellh/gox

build:
	go build -o $(GOPATH)/bin/$(PROG_NAME) ./cmd/esevaluator

dist-build:
	mkdir -p dist

	gox -output="./dist/{{.OS}}-{{.Arch}}/$(PROG_NAME)" \
		-os="linux windows darwin" \
		-arch="amd64" \
		./cmd/esevaluator > /dev/null

docker:
	docker build -t harvest-platform/$(PROG_NAME) .

dist-zip:
	cd dist && zip $(PROG_NAME)-darwin-amd64.zip darwin-amd64/*
	cd dist && zip $(PROG_NAME)-linux-amd64.zip linux-amd64/*
	cd dist && zip $(PROG_NAME)-windows-amd64.zip windows-amd64/*

dist: dist-build dist-zip docker

tls:
	@if [ ! -a cert.pem ]; then \
		echo >&2 'Creating self-signed TLS certs.'; \
		go run $(shell go env GOROOT)/src/crypto/tls/generate_cert.go --host localhost; \
	fi
