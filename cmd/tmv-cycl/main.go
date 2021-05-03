// Copyright ©2021 The tomuvol-cyclone5-module Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command tmv-cycl builds kernel modules for a Cyclone-V board running
// a Linux-SoCFPGA-3.10-ltsi kernel.
// tmv-cycl makes sure the dedicated Docker image is built and ready; if not
// it will build it.
// tmv-cycl then builds the kernel module inside that Docker image.
//
// tmv-cycl expects the kernel module Makefile to point at a kernel location
// of KERNEL_LOCATION=/build/linux.
//
package main // import "github.com/sbinet-lpc/tomuvol-cyclone5-module/cmd/tmv-cycl"

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/sbinet-lpc/tomuvol-cyclone5-module/xbuild"
)

func main() {
	xmain()
}

func xmain() {
	log.SetPrefix("tmv-cycl: ")
	log.SetFlags(0)

	flag.Parse()

	if flag.NArg() <= 0 {
		log.Fatalf("missing path to directory holding kernel module")
	}

	for _, dir := range flag.Args() {
		err := build(dir)
		if err != nil {
			log.Fatalf("could not build %q: %+v", dir, err)
		}
	}
}

func build(dir string) error {
	err := xbuild.Docker()
	if err != nil {
		return fmt.Errorf("could not build docker x-compilation image: %w", err)
	}

	log.Printf("building %q...", dir)

	src, err := filepath.Abs(dir)
	if err != nil {
		return fmt.Errorf("could not build absolute path to sources: %w", err)
	}

	tmp, err := os.MkdirTemp("", "tmv-cycl-")
	if err != nil {
		return fmt.Errorf("could not create tmp dir: %w", err)
	}
	defer os.RemoveAll(tmp)

	buildScript := buildModuleScript
	if !buildModule(dir) {
		buildScript = buildBinScript
	}

	err = os.WriteFile(filepath.Join(tmp, "run.sh"), []byte(buildScript), 0644)
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

	log.Printf("building %q... [done]", dir)
	return nil
}

const buildModuleScript = `#!/bin/bash

set -e
set -x

/bin/cp -ra /build/src /build/module

cd /build/module
make

cd /build/src
/bin/cp /build/module/*ko /build/src/.
modinfo ./*ko
`

const buildBinScript = `#!/bin/bash

set -e
set -x

cd /build/src
make
`

func buildModule(dir string) bool {
	raw, err := os.ReadFile(filepath.Join(dir, "Makefile"))
	if err != nil {
		log.Printf("could not read Makefile from %q: %+v", dir, err)
		return false
	}

	return bytes.Contains(raw, []byte("KERNEL_LOCATION="))
}
