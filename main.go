package main

import (
	"fmt"
)

func main() {
	dataFile := "data/brew-dependencies.txt"
	data, err := Asset(dataFile)
	if err != nil {
		panic(fmt.Sprintf("Data file [%s] not embedded!", dataFile))
	}

	dataAsString := string(data[:])
	fmt.Println("aaa")
	fmt.Println(dataAsString)
}
