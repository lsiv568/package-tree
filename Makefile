GO_FLAGS=-v -race

all: compile

build-brew-packages:
	./generate-package-list.sh > data/brew-dependencies.txt

generate-data-file:
	go-bindata -o data.go data/...

refresh-data: build-brew-packages generate-data-file

compile:
	go build $(GO_FLAGS) -o do-package-tree  *.go && go test $(GO_FLAGS) && go vet && golint

dependencies:
	go get -u github.com/golang/lint/golint
	go get -u github.com/jteeuwen/go-bindata/...

