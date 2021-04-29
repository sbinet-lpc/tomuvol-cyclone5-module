.PHONY: all

build: linux-socfpga
	mkdir -p ./modules
	DOCKER_BUILDKIT=1 \
	 docker build -t cyclone5 -o modules --progress=plain .

linux-socfpga:
	git clone --branch socfpga-3.10-ltsi \
		git@github.com:sbinet-lpc/tomuvol-linux-socfpga \
		./linux-socfpga

clean:
	/bin/rm -fr ./linux-socfpga ./modules

