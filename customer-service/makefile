# echo customer.proto
grpc-proto:
	protoc -I. --go_out=plugins=grpc:. proto/customer/customer.proto


#grpc-gateway
grpc-gateway:
	protoc -I/usr/local/include -I. \
	-I/Users/arshita/WorkspaceGO/src \
	--grpc-gateway_out=logtostderr=true,grpc_api_configuration=proto/customer/customer.yml:. \
	proto/customer/customer.proto	


# Generate swagger definitions
grpc-swagger:
	protoc -I/usr/local/include -I. \
	-I/Users/arshita/WorkspaceGO/src \
	--swagger_out=logtostderr=true,grpc_api_configuration=proto/customer/customer.yml:. \
	proto/customer/customer.proto
   
# Command to create the module
mod:
	go mod init customer-service.com/customer-service
# echo Dockerfile
build:
	docker build -f Dockerfile -t customer-service:1 .
# run
run:
	docker run -p 50052:50051 -e MICRO_SERVER_ADDRESS=:50051 customer-service:1
