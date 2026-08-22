// SPDX-License-Identifier: MIT
// Copyright (c) 2026 Paul Buetow

// Package internal provides version information for the gt application.
//
// This package contains internal constants that are used across the gt project.
// It is not intended for direct use by external code and may change without notice.
//
// # Package Location
//
// The internal package is located at internal/version.go and contains:
//   - Version: Current application version string
//
// # Version Format
//
// The version string follows semantic versioning (SemVer) format:
//   - Major.Minor.Patch (e.g., "v0.4.2")
//   - Pre-release versions may include suffixes like "-beta", "-rc1", etc.
//   - Build metadata may be appended for development builds
//
// # Usage in Code
//
// To access the version from the main command:
//
//	import "github.com/snonux/gt/internal"
//
//	func main() {
//	    fmt.Println("gt version", internal.Version)
//	}
//
// # Version History
//
// Current: v0.5.2
//
// See the git repository for complete version history and release notes.
package internal

// Version is the current version of the gt application.
//
// This constant is defined at build time and can be overridden during builds:
//
//	go build -ldflags="-X 'github.com/snonux/gt/internal.Version=v0.5.2-20260524'"
//
// The version is used for the "gt version" command-line subcommand.
//
// Example output:
//
//	$ gt version
//	v0.5.2
const Version = "v0.5.2"
