build:
	go build -o smug *.go

test:
	go test .

coverage:
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out

release:
	rm -rf dist
	git tag -a $(version) -m '$(version)'
	git push origin $(version)
	goreleaser
