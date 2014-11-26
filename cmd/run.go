package main

import (
	"flag"
	"github.com/riywo/go8086"
)

func main() {
	debug := flag.Bool("d", false, "debug")
	trace := flag.Bool("t", false, "trace")
	prefix := flag.String("p", "", "path prefix")

	flag.Parse()
	go8086.Debug = *debug
	go8086.Trace = *trace
	go8086.MinixPathPrefix = *prefix

	file := flag.Args()[0]
	args := flag.Args()[0:]
	envs := []string{
		"USER=root",
		"HOME=/",
		"PAGER=more",
		"LOGNAME=root",
		"TERM=minix",
		"PATH=/usr/local/bin:/bin:/usr/bin",
		"SHELL=/bin/sh",
		"TZ=GMT0",
		"EDITOR=vi",
	}
	go8086.Run(file, args, envs)
}
