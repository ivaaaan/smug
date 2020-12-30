VERSION_REGEX  := 's/(v[0-9\.]+)/$(version)/g'

build:
	go build -o smug *.go

test:
	go test .

coverage:
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out

release:
ifndef GITHUB_TOKEN
	$(error GITHUB_TOKEN is not defined)
endif
	sed -E -i.bak $(VERSION_REGEX) 'main.go'
	git commit -am '$(version)'
	git tag -a $(version) -m '$(version)'
	git push origin $(version)
	goreleaser --rm-dist
