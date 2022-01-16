NAME=libgoplug
GOFLAGS=-buildmode=c-shared
PREFIX=/usr/local

all: $(NAME).so

clean:
	rm -f $(NAME).so $(NAME).h

install: $(NAME).so $(NAME).h
	mkdir -p $(PREFIX)/{bin,include}
	cp $(NAME).so $(PREFIX)/bin
	cp $(NAME).h $(PREFIX)/include

$(NAME).so: main.go
	go build $(GOFLAGS) -o $(NAME).so .
