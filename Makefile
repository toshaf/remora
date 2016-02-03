
PKG=github.com/toshaf/remora

main:
	@GOBIN=./bin go install $(PKG)/...

race:
	@GOBIN=./bin go install -race $(PKG)/...

bin:
	@mkdir bin

.phony:

test: .phony
	@go test -race ./...

bench: .phony
	@go test -bench=. ./...

fmt:
	@go fmt ./...

run: race
	@./bin/server bin/client

runabs: main
	@./bin/server `pwd`/bin/client

clean:
	@-rm -rf bin
