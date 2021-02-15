# Determin platform type. Looking for Rpi-like platform
platform:
ifeq ("$(shell arch)", "armv6l")
	echo makefile.arm
	make -f makefile.arm
else
	echo makefile.nopi
	make -f makefile.nopi
endif

all:
	platform

clean:
	make -f makefile.nopi clean
