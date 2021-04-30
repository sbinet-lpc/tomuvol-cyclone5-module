CC=arm-linux-gnueabihf-gcc
CFLAGS=-Wall -O3

all: c-eda-irq-watch

%: %.c
	$(CC) $(CFLAGS) $< $(LDFLAGS) -o $@

