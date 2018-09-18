VERSION=2.3
GOFLAGS=

clean:
	rm -rf out

goGet:
	go get -u ./...

go-bootstrapper.darwin.amd64: goGet
	GOOS=darwin GOARCH=amd64 go build $(GOFLAGS) -o out/go-bootstrapper-$(VERSION).darwin.amd64 .
go-bootstrapper.darwin.386: goGet
	GOOS=darwin GOARCH=386 go build $(GOFLAGS) -o out/go-bootstrapper-$(VERSION).darwin.386 .

go-bootstrapper.linux.amd64: goGet
	GOOS=linux GOARCH=amd64 go build $(GOFLAGS) -o out/go-bootstrapper-$(VERSION).linux.amd64 .
go-bootstrapper.linux.386: goGet
	GOOS=linux GOARCH=386 go build $(GOFLAGS) -o out/go-bootstrapper-$(VERSION).linux.386 .

go-bootstrapper.windows.amd64: goGet
	GOOS=windows GOARCH=amd64 go build $(GOFLAGS) -o out/go-bootstrapper-$(VERSION).windows.amd64.exe .
go-bootstrapper.windows.386: goGet
	GOOS=windows GOARCH=386 go build $(GOFLAGS) -o out/go-bootstrapper-$(VERSION).windows.386.exe .

all: clean go-bootstrapper.darwin.amd64 go-bootstrapper.darwin.386 go-bootstrapper.linux.amd64 go-bootstrapper.linux.386 go-bootstrapper.windows.amd64 go-bootstrapper.windows.386
