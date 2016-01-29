EXECUTABLE_NAME=do-package-tree
GO_DEV_FLAGS=-v #-race

all: generate-data-file compile build

build-brew-packages:
	./generate-package-list.sh > data/brew-dependencies.txt

generate-data-file:
	go-bindata -o data.go data/...

refresh-data: build-brew-packages generate-data-file

compile:
	go test $(GO_DEV_FLAGS) && go vet && golint

build:
	GOOS=darwin go build -o $(EXECUTABLE_NAME)_darwin *.go
	GOOS=linux go build -o $(EXECUTABLE_NAME)_linux *.go

run:
	time ./do-package-tree
