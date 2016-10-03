VERSION=0.6

clean:
	rm -rf $(PWD)/.go
	rm -rf out

goGet: clean
	GOPATH=$(PWD)/.go go get github.com/op/go-logging

go-bootstrapper.darwin.amd64: src/main.go goGet
	GOPATH=$(PWD)/.go GOOS=darwin GOARCH=386 go build -o out/go-bootstrapper-$(VERSION).darwin.amd64 src/main.go

go-bootstrapper.darwin.386: src/main.go goGet
	GOPATH=$(PWD)/.go GOOS=darwin GOARCH=amd64 go build -o out/go-bootstrapper-$(VERSION).darwin.x64   src/main.go

go-bootstrapper.linux.amd64: src/main.go goGet
	GOPATH=$(PWD)/.go GOOS=linux GOARCH=386 go build -o out/go-bootstrapper-$(VERSION).linux.amd64 src/main.go
go-bootstrapper.linux.386: src/main.go goGet
	GOPATH=$(PWD)/.go GOOS=linux GOARCH=amd64 go build -o out/go-bootstrapper-$(VERSION).linux.x64   src/main.go

go-bootstrapper.windows.amd64.exe: src/main.go goGet
	GOPATH=$(PWD)/.go GOOS=windows GOARCH=386 go build -o out/go-bootstrapper-$(VERSION).windows.amd64.exe src/main.go
go-bootstrapper.windows.386.exe: src/main.go goGet
	GOPATH=$(PWD)/.go GOOS=windows GOARCH=amd64 go build -o out/go-bootstrapper-$(VERSION).windows.x64.exe   src/main.go

all: clean go-bootstrapper.darwin.amd64 go-bootstrapper.darwin.386 go-bootstrapper.linux.amd64 go-bootstrapper.linux.386 go-bootstrapper.windows.amd64.exe go-bootstrapper.windows.386.exe
