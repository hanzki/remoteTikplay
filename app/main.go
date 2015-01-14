package main

import (
	"code.google.com/p/gcfg"
	"fmt"
	"os"
)

type Config struct {
	Basic struct {
		Message string
	}
}

func main() {
	var cfg Config
	err := gcfg.ReadFileInto(&cfg, "config.gcfg")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(cfg.Basic.Message)
}
