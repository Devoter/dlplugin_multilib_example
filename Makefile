NAME=multilib

all: $(NAME)

clean:
	rm -f $(NAME)

$(NAME): main.go papi/device_plugin.go device/device.go cerror/cerror.go
	go build -o $(NAME) main.go
