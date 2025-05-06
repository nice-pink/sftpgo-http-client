package main

import (
	"flag"
	"os"

	"github.com/nice-pink/goutil/pkg/log"
)

func main() {
	config := flag.String("config", "", "Path to config file.")
	flag.Parse()

	log.Info("*** Start")
	log.Info(os.Args)
}
