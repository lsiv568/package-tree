EXECUTABLE_NAME=do-package-tree
SUPPORTED_GOOS=darwin linux windows freebsd
GO_DEV_FLAGS=-v #-race

all: generate-data-file compile build

clean:
	rm -f $(EXECUTABLE_NAME)*

build-brew-packages:
	./generate-package-list.sh > data/brew-dependencies.txt

generate-data-file:
	go-bindata -o data.go data/...

refresh-data: build-brew-packages generate-data-file

compile: clean
	go test $(GO_DEV_FLAGS) && go vet && golint && go build  -o $(EXECUTABLE_NAME)_`go env GOOS` *.go

build: clean
	$(foreach goos,$(SUPPORTED_GOOS), GOOS=$(goos) go build -o $(EXECUTABLE_NAME)_$(goos) *.go;)

run:
	time ./do-package-tree
