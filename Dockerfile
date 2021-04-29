from ubuntu:12.04 as build-stage

run apt-get update -y

run apt-get install -y gcc-arm-linux-gnueabihf \
	bc binutils bison \
	coreutils \
	diffutils \
	flex \
	git \
	gcc gcc-multilib \
	make \
	socat \
	texinfo \
	unzip \
	;

run apt-get install -y u-boot-tools

add ./linux-socfpga      /build/linux-3.10-ltsi
add ./eda-irq            /build/eda-irq
add ./hello              /build/hello

workdir /build/linux-3.10-ltsi

run make FOO=1 CROSS_COMPILE=arm-linux-gnueabihf- ARCH=arm socfpga_defconfig
#add ./eda-kernel-config  /build/linux-3.10-ltsi/.config
#run make CROSS_COMPILE=arm-linux-gnueabihf- ARCH=arm oldconfig
run make CROSS_COMPILE=arm-linux-gnueabihf- ARCH=arm -j8 uImage LOADADDR=0x8000
run make CROSS_COMPILE=arm-linux-gnueabihf- ARCH=arm -j8 modules
run make ARCH=arm INSTALL_MOD_PATH=/build/mnt modules_install

workdir /build/hello
run make

workdir /build/eda-irq
run make

run mkdir -p /build/out
run /bin/cp  /build/eda-irq/*.ko /build/out/.
run /bin/cp  /build/hello/*.ko /build/out/.
run modinfo /build/out/*ko

####

FROM scratch AS export-stage
COPY --from=build-stage /build/out/*.ko /


