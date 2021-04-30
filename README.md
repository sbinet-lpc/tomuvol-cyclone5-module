# tomuvol-cyclone5-module

`tomuvol-cyclone5-module` is a small repository to build kernel modules for a Cyclone-V SocKit board.

This repository builds a `socfpga-3.10-ltsi` kernel and builds kernel modules:

- `hello world` module
- `eda-irq` module adapted from [zhemao/interrupt_example](https://github.com/zhemao/interrupt_example)

This is for a Cyclone-V board like:

```
$> uname -a
Linux cyclone5 3.10.31-ltsi-05172-g28bac3e #1 SMP Tue Oct 18 16:02:00 EDT 2016 armv7l GNU/Linux
```

## Installation

```
$> go get github.com/sbinet-lpc/tomuvol-cyclone5-module/cmd/tmv-cycl
$> tmv-cycl ./hello
tmv-cycl: building docker image "tomuvol-cyclone5-3.10"...
Cloning into './linux-socfpga'...
remote: warning: multi-pack bitmap is missing required reverse index
remote: Enumerating objects: 6125779, done.
remote: Counting objects: 100% (6125779/6125779), done.
remote: Compressing objects: 100% (914709/914709), done.
Receiving objects: 100% (6125779/6125779), 1.32 GiB | 19.25 MiB/s, done.
remote: Total 6125779 (delta 5171045), reused 6122122 (delta 5167388), pack-reused 0
Resolving deltas: 100% (5171045/5171045), done.
Updating files: 100% (43679/43679), done.
Sending build context to Docker daemon  2.131GB
Step 1/10 : from ubuntu:12.04
 ---> 5b117edd0b76
Step 2/10 : run apt-get update -y
[...]
  INSTALL drivers/usb/gadget/libcomposite.ko
  DEPMOD  3.10.31-ltsi-05172-g28bac3e
Removing intermediate container 0180f4c55ffb
 ---> 67fc43264117
Step 10/10 : workdir /build
 ---> Running in 283a421f731f
Removing intermediate container 283a421f731f
 ---> 7097ec311889
Successfully built 7097ec311889
Successfully tagged tomuvol-cyclone5-3.10:latest
tmv-cycl: building docker image "tomuvol-cyclone5-3.10"... [done]
+ /bin/cp -ra /build/src /build/module
+ cd /build/module
+ make
make ARCH=arm SUBARCH=arm CROSS_COMPILE=arm-linux-gnueabihf- -C /build/linux M=/build/module modules
make[1]: Entering directory `/build/linux'
  CC [M]  /build/module/bjr.o
  Building modules, stage 2.
  MODPOST 1 modules
  CC      /build/module/bjr.mod.o
  LD [M]  /build/module/bjr.ko
make[1]: Leaving directory `/build/linux'
+ cd /build/src
+ /bin/cp /build/module/bjr.ko /build/src/.
+ modinfo ./bjr.ko
filename:       ./bjr.ko
license:        Dual BSD/GPL
depends:        
vermagic:       3.10.31-ltsi-05172-g28bac3e SMP mod_unload ARMv7 p2v8 
```
