# dlplugin_multilib_example

This example demonstrates a loading multiple dynamic libraries (C++ and Go) with the same interface into the single project.

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


## License

[LICENSE](./LICENSE)