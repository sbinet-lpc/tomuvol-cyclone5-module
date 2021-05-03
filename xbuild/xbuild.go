// Copyright ©2021 The tomuvol-cyclone5-module Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xbuild

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const ImageName = "tomuvol-cyclone5-3.10"

func Docker() error {
	if !HasDocker() {
		return fmt.Errorf("docker not installed or unavailable")
	}

	if HasDockerImage(ImageName) {
		return nil
	}

	log.Printf("building docker image %q...", ImageName)

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

	cmd = exec.Command("docker", "build", "-t", ImageName, ".")
	cmd.Dir = dir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf(
			"could not build docker image %q: %w", ImageName, err,
		)
	}

	log.Printf("building docker image %q... [done]", ImageName)
	return nil
}

func HasDocker() bool {
	cmd := exec.Command("docker", "info")
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	return cmd.Run() == nil
}

func HasDockerImage(name string) bool {
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
