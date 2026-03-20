//go:build mage

package main

import (
	"github.com/magefile/mage/sh"
)

var Default = Build

func Build() error {
	return sh.RunV("go", "build", "-o", "perc", "./cmd/perc")
}

func Run() error {
	return sh.RunV("go", "run", "./cmd/perc")
}

func Test() error {
	return sh.RunV("go", "test", "./...")
}

func TestRPN() error {
	return sh.RunV("go", "test", "./internal/rpn/...")
}

func Install() error {
	return sh.RunV("go", "install", "./cmd/perc")
}

func Repl() error {
	return sh.RunV("go", "run", "./cmd/perc", "--repl")
}
