NAME=multilib

all: $(NAME)

clean:
	rm -f $(NAME)

$(NAME): dlprog/main.go papi/device_plugin.go device/device.go cerror/cerror.go
	go build -o $(NAME) ./dlprog/...
