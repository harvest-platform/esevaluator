PROG_NAME := "esevaluator"
GIT_VERSION := $(shell git log -1 --pretty=format:"%h (%ci)" .)

setup: install tls

clean:
	go clean ./...

install:
	go get github.com/mitchellh/gox

build:
	go build -o $(GOPATH)/bin/esevaluator ./cmd/esevaluator

dist-build:
	mkdir -p bin

	gox -output="bin/esevaluator-{{.OS}}.{{.Arch}}" \
		-os="linux windows darwin" \
		-arch="amd64" \
		./cmd/esevaluator > /dev/null

dist-zip:
	cd dist && zip $(PROG_NAME)-darwin-amd64.zip darwin-amd64/*
	cd dist && zip $(PROG_NAME)-linux-amd64.zip linux-amd64/*
	cd dist && zip $(PROG_NAME)-windows-amd64.zip windows-amd64/*

dist: dist-build dist-zip

tls:
	@if [ ! -a cert.pem ]; then \
		echo >&2 'Creating self-signed TLS certs.'; \
		go run $(GOROOT)/src/crypto/tls/generate_cert.go --host localhost; \
	fi
