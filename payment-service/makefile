# echo payment.proto
grpc-proto:
	protoc -I. --go_out=plugins=grpc:. proto/payment/payment.proto


#grpc-gateway
grpc-gateway:
	protoc -I/usr/local/include -I. \
	-I/Users/arshita/WorkspaceGO/src \
	--grpc-gateway_out=logtostderr=true,grpc_api_configuration=proto/payment/payment.yml:. \
	proto/payment/payment.proto	


# Generate swagger definitions
grpc-swagger:
	protoc -I/usr/local/include -I. \
	-I/Users/arshita/WorkspaceGO/src \
	--swagger_out=logtostderr=true,grpc_api_configuration=proto/payment/payment.yml:. \
	proto/payment/payment.proto
   
# Command to create the module
mod:
	go mod init payment-service.com/payment-service
# echo Dockerfile
build:
	docker build -f Dockerfile -t payment-service:1 .
# run
run:
	docker run -p 50052:50051 -e MICRO_SERVER_ADDRESS=:50051 payment-service:1
