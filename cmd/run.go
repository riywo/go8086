package main

import (
	"github.com/riywo/go8086"
	"io/ioutil"
	"os"
)

func main() {
	file := os.Args[1]
	args := os.Args[1:]
	bs, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	go8086.Run(bs, args)
}
