VERSION = $(shell git describe --tags --abbrev=0)

version:
	@echo $(VERSION)

build:
	go build -ldflags "-X=main.version=$(VERSION)" -gcflags "all=-trimpath=$(GOPATH)"

test:
	go test

coverage:
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out

release:
ifndef GITHUB_TOKEN
	$(error GITHUB_TOKEN is not defined)
endif
	git tag -a $(version) -m '$(version)'
	git push origin $(version)
	VERSION=$(version) goreleaser --rm-dist
