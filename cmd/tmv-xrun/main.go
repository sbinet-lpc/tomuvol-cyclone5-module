// Copyright ©2021 The tomuvol-cyclone5-module Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command tmv-xrun executes a given command in a cross-compiled environment.
package main // import "github.com/sbinet-lpc/tomuvol-cyclone5-module/cmd/tmv-xrun"

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sbinet-lpc/tomuvol-cyclone5-module/xbuild"
)

func main() {
	xmain()
}

func xmain() {
	log.SetPrefix("tmv-xrun: ")
	log.SetFlags(0)

	dir := flag.String("dir", ".", "path to directory to mount")

	flag.Parse()

	if flag.NArg() <= 0 {
		flag.Usage()
		log.Fatalf("missing command to execute")
	}

	err := build(*dir, flag.Args())
	if err != nil {
		log.Fatalf("could not run command %q: %+v", flag.Args(), err)
	}
}

func build(dir string, args []string) error {
	err := xbuild.Docker()
	if err != nil {
		return fmt.Errorf("could not build docker image: %w", err)
	}

	src, err := filepath.Abs(dir)
	if err != nil {
		return fmt.Errorf("could not build absolute path to sources: %w", err)
	}

	tmp, err := os.MkdirTemp("", "tmv-cycl-")
	if err != nil {
		return fmt.Errorf("could not create tmp dir: %w", err)
	}
	defer os.RemoveAll(tmp)

	bld := new(strings.Builder)
	fmt.Fprintf(bld, `#!/bin/bash

set -e
set -x

cd /build/src
%s
`, strings.Join(args, " "),
	)

	err = os.WriteFile(filepath.Join(tmp, "run.sh"), []byte(bld.String()), 0644)
	if err != nil {
		return fmt.Errorf("could not create build script: %w", err)
	}

	cmd := exec.Command(
		"docker", "run", "--rm", "-t",
		"-v", src+":/build/src",
		"-v", tmp+":/build/x",
		xbuild.ImageName,
		"/bin/sh", "/build/x/run.sh",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("could not build %q: %w", dir, err)
	}

	return nil
}
