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
$> git clone git@github.com:sbinet-lpc/tomuvol-cyclone5-module
$> cd ./tomuvol-cyclone5-module
$> make
git clone --branch socfpga-3.10-ltsi \
	git@github.com:sbinet-lpc/tomuvol-linux-socfpga \
	./linux-socfpga
Cloning into './linux-socfpga'...
remote: warning: multi-pack bitmap is missing required reverse index
remote: Enumerating objects: 6125779, done.
remote: Counting objects: 100% (6125779/6125779), done.
remote: Compressing objects: 100% (914709/914709), done.
Receiving objects: 100% (6125779/6125779), 1.32 GiB | 16.54 MiB/s, done.
remote: Total 6125779 (delta 5171045), reused 6122122 (delta 5167388), pack-reused 0
Resolving deltas: 100% (5171045/5171045), done.
Updating files: 100% (43679/43679), done.
mkdir -p ./modules
DOCKER_BUILDKIT=1 \
 docker build -t cyclone5 -o out --progress=plain .
#1 [internal] load build definition from Dockerfile
#1 sha256:3379666a65f474986239e6bc8000979147025acca9e9e07937e151f353f23750
#1 transferring dockerfile: 38B done
#1 DONE 0.0s

[...]
#24 [build-stage 20/20] RUN modinfo /build/out/*ko
#24 sha256:252d4eda0b43b2b46390dcce0f3ffbfe25c7ef12b3de1ff1c3e0c97060a33f8c
#24 0.635 filename:       /build/out/bjr.ko
#24 0.635 license:        Dual BSD/GPL
#24 0.635 depends:        
#24 0.635 vermagic:       3.10.31-ltsi-05172-g28bac3e SMP mod_unload ARMv7 p2v8 
#24 0.635 filename:       /build/out/eda_irq.ko
#24 0.635 license:        Dual BSD/GPL
#24 0.635 depends:        
#24 0.635 vermagic:       3.10.31-ltsi-05172-g28bac3e SMP mod_unload ARMv7 p2v8 
#24 DONE 0.7s

#25 [export-stage 1/1] COPY --from=build-stage /build/out/*.ko /
#25 sha256:3410d2eb1142ff88dc48c7e6257cb55cfcf4bf524c19fad910a8844773d42869
#25 CACHED

#26 exporting to client
#26 sha256:b60a1292d407630dbb741f28ab6ea4ce3cca872ac28eeee56f4e66a182eca4bc
#26 copying files 105.02kB done
#26 DONE 0.0s
```

