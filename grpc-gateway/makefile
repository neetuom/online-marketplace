# Command to create the module
mod:
	go mod init grpc-gateway.com/grpc-gateway
# echo Dockerfile
build:
	docker build -f Dockerfile -t grpc-gateway:1 .

# build plugin
build-plugin:
	go build -buildmode=plugin -o ./plugin/grpc-gateway-plugin.so

# run
run:
	docker run -p 8081:8081 -e MICRO_SERVER_ADDRESS=:8081 grpc-gateway:1

# docker compose remove all images

docker-rm-all:
	docker-compose -f docker-compose-online-marketplace.yml down --rmi all

#echo "generating static files" 
statik-file-gen:
	statik -src=proto/swagger
	mv statik gateway/.
