package main

import (
	"github.com/minhngoc274/genomic-system/genomic-service/bootstrap"
	"go.uber.org/fx"
)

// @title API
// @version 1.0.0
func main() {
	fx.New(bootstrap.All()).Run()
}
