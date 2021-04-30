// Copyright ©2021 The tomuvol-cyclone5-module Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"syscall"
)

func main() {
	log.SetPrefix("eda-irq: ")
	log.SetFlags(0)

	flag.Parse()

	xmain()
}

const (
	sysfsFile   = "/sys/bus/platform/drivers/eda-irq/eda_irq"
	numSwitches = 4
	numKeys     = 4
)

func xmain() {
	irq, err := watch(sysfsFile, func(p []byte) {
		log.Printf("irq: %q", p)
	})
	if err != nil {
		log.Fatalf("irq: %+v", err)
	}
	defer irq.close()

	select {} // FIXME(sbinet): listen on an irq.done channel
}

func watch(fname string, cbk func(p []byte)) (*register, error) {
	f, err := os.OpenFile(fname, os.O_RDWR|os.O_SYNC, 0666)
	if err != nil {
		return nil, fmt.Errorf("could not open IRQ: %w", err)
	}

	irq := &register{
		f:   f,
		buf: make([]byte, 32*1024),
		cbk: cbk,
	}

	evt := syscall.EpollEvent{
		Events: syscall.EPOLLIN | (syscall.EPOLLET & 0xffffffff) | syscall.EPOLLPRI,
	}

	fd := int(irq.f.Fd())
	epoll.cbks[fd] = irq

	err = syscall.SetNonblock(fd, true)
	if err != nil {
		return nil, fmt.Errorf("could not set non-blocking mode: %w", err)
	}

	evt.Fd = int32(fd)

	err = syscall.EpollCtl(epoll.fd, syscall.EPOLL_CTL_ADD, fd, &evt)
	if err != nil {
		return nil, fmt.Errorf("could not add IRQ: %w", err)
	}

	return irq, nil
}

type register struct {
	f       *os.File
	initial bool
	buf     []byte
	cbk     func(p []byte)
}

func (reg *register) close() {
	if reg.f == nil {
		return
	}

	defer reg.f.Close()

	fd := int(reg.f.Fd())

	delete(epoll.cbks, fd)

	err := syscall.EpollCtl(epoll.fd, syscall.EPOLL_CTL_DEL, fd, nil)
	if err != nil {
		log.Printf("could not remove IRQ: %+v", err)
		return
	}

	err = syscall.SetNonblock(fd, false)
	if err != nil {
		log.Printf("could not set blocking mode: %+v", err)
	}

	err = reg.f.Close()
	if err != nil {
		log.Printf("could not close register: %+v", err)
	}

	reg.f = nil
}

var epoll struct {
	fd   int
	cbks map[int]*register
}

func init() {
	epoll.cbks = make(map[int]*register)

	var err error
	epoll.fd, err = syscall.EpollCreate1(0)
	if err != nil {
		log.Fatalf("could not create epoll FD: %+v", err)
	}

	go func() {
		var evts [64]syscall.EpollEvent

		for {
			n, err := syscall.EpollWait(epoll.fd, evts[:], -1)
			if err != nil {
				if err == syscall.EAGAIN {
					continue
				}
				log.Panicf("epoll wait error: %+v", err)
			}
			log.Printf("\n\n\n")
			log.Printf("epoll: evts=%d", n)

			for i := 0; i < n; i++ {
				fd := int(evts[i].Fd)
				log.Printf("epoll: i=%d, fd=%d", i, fd)
				reg, ok := epoll.cbks[fd]
				if !ok {
					continue
				}
				switch {
				case reg.initial:
					reg.initial = false
				default:
					nn, err := syscall.Read(fd, reg.buf)
					if err != nil {
						log.Printf("epoll: could not read register: %+v", err)
						continue
					}
					if nn <= 0 {
						continue
					}
					log.Printf("epoll: i=%d, fd=%d ==> n=%d", i, fd, nn)
					reg.cbk(reg.buf[:nn])
				}
			}
		}
	}()
}
