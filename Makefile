.PHONY: build test clean release

build:
	go build -o pvm.exe .

test:
	go test ./...

clean:
	rm -f pvm.exe

release:
	goreleaser release --clean
