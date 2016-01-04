package main

type Client interface {
	Write(string) (int, error)
}
