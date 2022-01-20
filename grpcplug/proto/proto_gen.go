//go:generate protoc --proto_path=. --go_out=. --go_opt=paths=source_relative --go_opt=Mdevice.proto=github.com/Devoter/dlplugin_multilib_example/rpcplug/proto device.proto device.proto
//go:generate protoc --proto_path=. --go-grpc_out=. --go-grpc_opt=paths=source_relative device.proto device.proto
package proto
