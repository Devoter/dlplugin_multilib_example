# dlplugin_multilib_example

This example demonstrates a loading multiple dynamic libraries (C++ and Go) with the same interface into the single project.


Additionally, you may want to compare a solution with [hashicorp/go-plugin](https://github.com/hashicorp/go-plugin)-based examples. See [gRPC](./grpcplug/) and [net/rpc](./rpcplug/) links.


## Building

```sh
cd cppplug && make && cd ..
cd goplug && make && cd ..
make
```


## Running

```sh
./multilib -plug1 cppplug/libcppplug.so -plug2 goplug/libgoplug.so
# or
./multilib -plug1 goplug/libgoplug.so -plug2 cppplug/libcppplug.so
```


## Benchmarking

At first, [build](#building) plugins.

```sh
go test -benchmem -bench=.
```

## License

[LICENSE](./LICENSE)