## 财务
#protoc --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. ./finance/*.proto ./live/*.proto


## 直播
protoc --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. ./live/*.proto

## 财务
protoc --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. ./finance/*.proto

##协议生成后，请务必执行脚本
go run remove_omitempty.go