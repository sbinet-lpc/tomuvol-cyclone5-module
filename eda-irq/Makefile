KERNEL_LOCATION=/build/linux
ARMMAKE=make ARCH=arm SUBARCH=arm CROSS_COMPILE=arm-linux-gnueabihf-

obj-m := eda_irq.o

eda_irq.ko: eda_irq.c
	$(ARMMAKE) -C $(KERNEL_LOCATION) M=$(PWD) modules

clean:
	rm -f *.ko *.o *.mod.c *.symvers *.order
