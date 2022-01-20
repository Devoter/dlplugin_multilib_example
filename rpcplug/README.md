# rpcplug

This example uses [hashicorp/go-plugin](https://github.com/hashicorp/go-plugin) based on [net/rpc](https://pkg.go.dev/net/rpc).

## Building

```sh
make
```

## Running

```sh
./rpclib -plug="./device-go-rpc"
```


## Benchmarking

```sh
go test -benchmem -bench=.
```


## License

[LICENSE](../LICENSE)