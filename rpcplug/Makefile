NAME=rpclib
LIB_NAME=device-go-rpc

all: $(NAME) $(LIB_NAME)

clean:
	rm -f $(NAME) $(LIB_NAME)

$(NAME): main.go ../device/device.go ./shared/rpc.go ./shared/interface.go ./proto/device_rpc.go
	go build -o $(NAME) main.go

$(LIB_NAME): plug/main.go ../device/device.go ./shared/rpc.go ./shared/interface.go ./proto/device_rpc.go
	go build -o $(LIB_NAME) ./plug
