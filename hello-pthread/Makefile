CC=arm-linux-gnueabihf-gcc
CFLAGS=-Wall -O3
LDFLAGS=-lpthread

all: hello-pthread

%: %.c
	$(CC) $(CFLAGS) $< $(LDFLAGS) -o $@

