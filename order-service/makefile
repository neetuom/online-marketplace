# echo order.proto
grpc-proto:
	protoc -I. --go_out=plugins=grpc:. proto/order/order.proto


#grpc-gateway
grpc-gateway:
	protoc -I/usr/local/include -I. \
	-I/Users/arshita/WorkspaceGO/src \
	--grpc-gateway_out=logtostderr=true,grpc_api_configuration=proto/order/order.yml:. \
	proto/order/order.proto	


# Generate swagger definitions
grpc-swagger:
	protoc -I/usr/local/include -I. \
	-I/Users/arshita/WorkspaceGO/src \
	--swagger_out=logtostderr=true,grpc_api_configuration=proto/order/order.yml:. \
	proto/order/order.proto
   
# Command to create the module
mod:
	go mod init order-service.com/order-service
# echo Dockerfile
build:
	docker build -f Dockerfile -t order-service:1 .
# run
run:
	docker run -p 50052:50051 -e MICRO_SERVER_ADDRESS=:50051 order-service:1
