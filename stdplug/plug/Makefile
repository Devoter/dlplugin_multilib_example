NAME=stdplug
GOFLAGS=-buildmode=plugin

all: $(NAME).so

clean:
	rm -f $(NAME).so

$(NAME).so: main.go ../../device/device.go
	go build $(GOFLAGS) -o $(NAME).so .
