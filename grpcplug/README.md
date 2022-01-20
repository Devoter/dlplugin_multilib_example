# grpcplug

This example uses [hashicorp/go-plugin](https://github.com/hashicorp/go-plugin) based on [gRPC](https://grpc.io/).

## Building

```sh
make
```

## Running

```sh
./grpclib -plug="./device-go-grpc"
```


## Benchmarking

```sh
go test -benchmem -bench=.
```


## License

[LICENSE](../LICENSE)