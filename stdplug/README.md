# stdplug

This example uses the [plugin](https://pkg.go.dev/plugin) package.

## Building

```sh
cd plug
make
cd ..
make
```

## Running

```sh
./stdlib -plug="./plug/stdplug.so"
```


## Benchmarking

```sh
go test -benchmem -bench=.
```


## License

[LICENSE](../LICENSE)