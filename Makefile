BUILDDIR = ./build
BUILDFLAGS =

.PHONY: test lint

APPS = nsq_to_consumer

$(BUILDDIR)/%:
	@mkdir -p $(dir $@)
	go build ${BUILDFLAGS} -o $@ ./

$(APPS): %: $(BUILDDIR)/%

test:
	go test -v -race -cover -coverprofile=coverage.txt -covermod=atomic ./...

lint:
	golangci-lint cache clean
	golangci-lint run --tests=false ./...

clean:
	rm -rf $(BUILDDIR)

tidy:
	go mod tidy
