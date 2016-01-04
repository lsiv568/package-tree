EXECUTABLE_NAME=do-package-tree
GO_DEV_FLAGS=-v -race

all: dependencies generate-data-file compile build

build-brew-packages:
	./generate-package-list.sh > data/brew-dependencies.txt

generate-data-file:
	go-bindata -o data.go data/...

refresh-data: build-brew-packages generate-data-file

compile:
	go test $(GO_DEV_FLAGS) && go vet && golint

dependencies:
	go get -u github.com/golang/lint/golint
	go get -u github.com/jteeuwen/go-bindata/...

build:
	go build -o $(EXECUTABLE_NAME) *.go
