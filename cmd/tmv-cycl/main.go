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
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
		err := buildModule(dir)
		if err != nil {
			log.Fatalf("could not build kernel module %q: %+v", dir, err)
		}
	}
}

func buildModule(dir string) error {
	err := buildDocker()
	if err != nil {
		return fmt.Errorf("could not build docker x-compilation image: %w", err)
	}

	dir, err = filepath.Abs(dir)
	if err != nil {
		return fmt.Errorf("could not build absolute path to kernel module: %w", err)
	}

	log.Printf("building kernel module from %q...", dir)

	tmp, err := os.MkdirTemp("", "tmv-cycl-")
	if err != nil {
		return fmt.Errorf("could not create tmp dir: %w", err)
	}
	defer os.RemoveAll(tmp)

	err = os.WriteFile(filepath.Join(tmp, "run.sh"), []byte(buildScript), 0644)
	if err != nil {
		return fmt.Errorf("could not create build script: %w", err)
	}

	cmd := exec.Command(
		"docker", "run", "--rm", "-t",
		"-v", dir+":/build/src",
		"-v", tmp+":/build/x",
		dockerImageName,
		"/bin/sh", "/build/x/run.sh",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("could not run kernel module build script: %w", err)
	}

	log.Printf("building kernel module from %q... [done]", dir)
	return nil
}

const dockerImageName = "tomuvol-cyclone5-3.10"

func buildDocker() error {
	if !hasDocker() {
		return fmt.Errorf("docker not installed or unavailable")
	}

	if hasDockerImage(dockerImageName) {
		return nil
	}

	log.Printf("building docker image %q...", dockerImageName)

	dir, err := os.MkdirTemp("", "tmv-cycl-")
	if err != nil {
		return fmt.Errorf("could not create tmp dir: %w", err)
	}
	defer os.RemoveAll(dir)

	cmd := exec.Command(
		"git", "clone", "--branch", "socfpga-3.10-ltsi",
		"git@github.com:sbinet-lpc/tomuvol-linux-socfpga",
		"./linux-socfpga",
	)
	cmd.Dir = dir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("could not fetch linux-socfpga sources: %w", err)
	}

	fname := filepath.Join(dir, "Dockerfile")
	err = os.WriteFile(fname, []byte(dockerImage), 0644)
	if err != nil {
		return fmt.Errorf("could not create Dockerfile: %w", err)
	}

	cmd = exec.Command("docker", "build", "-t", dockerImageName, ".")
	cmd.Dir = dir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf(
			"could not build docker image %q: %w", dockerImageName, err,
		)
	}

	log.Printf("building docker image %q... [done]", dockerImageName)
	return nil
}

func hasDocker() bool {
	cmd := exec.Command("docker", "info")
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	return cmd.Run() == nil
}

func hasDockerImage(name string) bool {
	buf := new(bytes.Buffer)
	cmd := exec.Command("docker", "images", name)
	cmd.Stdout = buf
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Printf("could not run cmd %q: %+v", cmd.Args, err)
		return false
	}

	sc := bufio.NewScanner(buf)
	if !sc.Scan() {
		return false
	}
	if !strings.HasPrefix(sc.Text(), "REPOSITORY") {
		return false
	}

	for sc.Scan() {
		txt := sc.Text()
		if strings.HasPrefix(txt, name+" ") {
			return true
		}
	}
	return false
}

const dockerImage = `
from ubuntu:12.04

run apt-get update -y
run apt-get install -y \
	bc binutils bison \
	coreutils \
	diffutils \
	flex \
	git \
	gcc gcc-multilib gcc-arm-linux-gnueabihf \
	make \
	socat \
	texinfo \
	u-boot-tools unzip \
	;

add ./linux-socfpga /build/linux

workdir /build/linux

run make CROSS_COMPILE=arm-linux-gnueabihf- ARCH=arm socfpga_defconfig
run make CROSS_COMPILE=arm-linux-gnueabihf- ARCH=arm -j8 uImage LOADADDR=0x8000
run make CROSS_COMPILE=arm-linux-gnueabihf- ARCH=arm -j8 modules
run make ARCH=arm INSTALL_MOD_PATH=/build/mnt modules_install

workdir /build
`

const buildScript = `#!/bin/bash

set -e
set -x

/bin/cp -ra /build/src /build/module

cd /build/module
make

cd /build/src
/bin/cp /build/module/*ko /build/src/.
modinfo ./*ko
`
