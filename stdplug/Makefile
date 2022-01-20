NAME=stdlib

all: $(NAME)

clean:
	rm -f $(NAME)

$(NAME): main.go papi/device_plugin.go ../device/device.go
	go build -o $(NAME) .