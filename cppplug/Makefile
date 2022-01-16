CPP=g++
NAME=libcppplug
CFLAGS=-Wall -fPIC -shared
CPPFLAGS=-std=c++2a
PREFIX=/usr/local

$(NAME).so: device.o device.hpp $(NAME).h
	$(CPP) $(CFLAGS) $(CPPFLAGS) -o $(NAME).so main.cpp device.o device.hpp $(NAME).h

device.o: device.cpp
	$(CPP) $(CFLAGS) $(CPPFLAGS) -c device.cpp

clean:
	rm -f $(NAME).so device.o

install:
	mkdir -p $(PREFIX)/{lib,include}
	cp $(NAME).so $(PREFIX)/lib
	cp $(NAME).h $(PREFIX)/include
