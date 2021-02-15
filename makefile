UNAME := $(shell uname)

platform:
ifeq ("$(UNAME)", "Darwin")
	make -f makefile.nopi
else
	make -f makefile.arm
endif

all:
	platform

clean:
	make -f makefile.nopi clean
