// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

//go:build mage

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/magefile/mage/sh"
)

// Project description for build output
const binaryName = "gt"

// Default is the default target when no target is specified.
var Default = Build

// Build builds the perc binary.
func Build() error {
	fmt.Println("Building gt...")
	return sh.RunV("go", "build", "-o", binaryName, "./cmd/gt")
}

// Run runs the perc binary.
func Run() error {
	return sh.RunV("go", "run", "./cmd/gt")
}

// Test runs all tests.
func Test() error {
	fmt.Println("Running all tests...")
	return sh.RunV("go", "test", "./...")
}

// TestRPN runs tests for the RPN package.
func TestRPN() error {
	fmt.Println("Running RPN tests...")
	return sh.RunV("go", "test", "./internal/rpn/...")
}

// RPN runs tests for the RPN package (alias for TestRPN).
func RPN() error {
	return TestRPN()
}

// Install installs the perc binary to GOPATH/bin.
func Install() error {
	fmt.Println("Installing gt...")
	return sh.RunV("go", "install", "./cmd/gt")
}

// Repl starts the REPL mode.
func Repl() error {
	return sh.RunV("go", "run", "./cmd/gt", "--repl")
}

// Uninstall removes the perc binary from GOPATH/bin.
func Uninstall() error {
	fmt.Println("Uninstalling gt...")
	binPath := filepath.Join(getGOPATH(), "bin", binaryName)
	return os.Remove(binPath)
}

// getGOPATH returns the GOPATH environment variable.
func getGOPATH() string {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = filepath.Join(os.Getenv("HOME"), "go")
	}
	return gopath
}
