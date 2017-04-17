package main

import "os"

type AAA interface {
	hello(a string) (b int, err error)
	B
	os.FileInfo
}

type B interface {
	test()
}
