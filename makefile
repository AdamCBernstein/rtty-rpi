UNAME := $(shell uname)

# Determin platform type. Looking for Rpi-like platform
platform:
ifeq ("$(arch)", "armv6l")
	make -f makefile.nopi
else
	make -f makefile.arm
endif

all:
	platform

clean:
	make -f makefile.nopi clean
