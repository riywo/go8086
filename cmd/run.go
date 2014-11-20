package main

import (
	"flag"
	"github.com/riywo/go8086"
	"io/ioutil"
)

func main() {
	debug := flag.Bool("d", false, "debug")

	flag.Parse()

	file := flag.Args()[0]
	args := flag.Args()[0:]
	bs, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	go8086.Debug = *debug
	go8086.Run(bs, args)
}
