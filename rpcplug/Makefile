NAME=rpclib
LIB_NAME=device-go-grpc

all: $(NAME) $(LIB_NAME)

clean:
	rm -f $(NAME) $(LIB_NAME)

$(NAME): main.go ../device/device.go ./shared/grpc.go ./shared/interface.go \
		./proto/device_grpc.pb.go ./proto/device.pb.go ./proto/proto_gen.go
	go build -o $(NAME) main.go

$(LIB_NAME): plug/main.go ../device/device.go ./shared/grpc.go ./shared/interface.go \
		./proto/device_grpc.pb.go ./proto/device.pb.go ./proto/proto_gen.go
	go build -o $(LIB_NAME) ./plug

generate:
	go generate ./proto
