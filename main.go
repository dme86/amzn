package main

import (
	"github.com/dme86/amzn/pkg"
)

func main() {
	bucket.Init()
	bucket.RootCmd.Execute()
}
