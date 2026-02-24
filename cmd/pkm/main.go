package main

import (
	cli "github.com/sahay-shashank/personal-knowledge-manager/internal/cli"
)

var VERSION = "dev"
var COMMIT = "unknown"

func main() {
	cli.NewCli(VERSION, COMMIT)
}
