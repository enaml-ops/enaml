package main

import (
	"io/ioutil"
	"os"

	"github.com/enaml-ops/enaml/generators"
)

func main() {
	packagename := os.Args[1]
	filename := os.Args[2]
	outputDir := os.Args[3]
	b, _ := ioutil.ReadFile(filename)
	generators.Generate(packagename, b, outputDir)
}
