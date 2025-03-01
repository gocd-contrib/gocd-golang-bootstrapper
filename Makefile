VERSION=2.26
GOFLAGS=

clean:
	rm -rf out

goGet:
	go get -u ./...

go-bootstrapper.darwin.amd64: goGet
	GOOS=darwin GOARCH=amd64 go build $(GOFLAGS) -o out/go-bootstrapper-$(VERSION).darwin.amd64 .
go-bootstrapper.darwin.arm64: goGet
	GOOS=darwin GOARCH=arm64 go build $(GOFLAGS) -o out/go-bootstrapper-$(VERSION).darwin.arm64 .

go-bootstrapper.linux.amd64: goGet
	GOOS=linux GOARCH=amd64 go build $(GOFLAGS) -o out/go-bootstrapper-$(VERSION).linux.amd64 .
go-bootstrapper.linux.arm64: goGet
	GOOS=linux GOARCH=arm64 go build $(GOFLAGS) -o out/go-bootstrapper-$(VERSION).linux.arm64 .

go-bootstrapper.windows.amd64: goGet
	GOOS=windows GOARCH=amd64 go build $(GOFLAGS) -o out/go-bootstrapper-$(VERSION).windows.amd64.exe .
go-bootstrapper.windows.arm64: goGet
	GOOS=windows GOARCH=arm64 go build $(GOFLAGS) -o out/go-bootstrapper-$(VERSION).windows.arm64.exe .

summary: goGet
	pushd out; sha256sum *; popd

all: clean go-bootstrapper.darwin.amd64 go-bootstrapper.darwin.arm64 go-bootstrapper.linux.amd64 go-bootstrapper.linux.arm64 go-bootstrapper.windows.amd64 go-bootstrapper.windows.arm64 summary
