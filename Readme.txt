
1). Kong API Gateway set up using below docker-compose file

    nodejs-kong-gateway/docker-compose-kong-gateway-cassandra.yml   # Cassandra

    nodejs-kong-gateway/docker-compose-kong-gateway-postgres.yml    # PostgreSQL

    Examples - a. gin-restapi  b.nodejs-kong-gateway  c. online-marketplace


2). ##########################################################################################################################  
      Kong API Gateway  -> grpc service mapping
    ##########################################################################################################################
    
     # Service mapping

     curl -XPOST localhost:8001/services \
     --data name=grpc \
     --data protocol=grpc \
     --data host=product-service \
     --data port=60060  

     # getproductlist route mapping


     curl -XPOST localhost:8001/services/grpc/routes \
     --data protocols=grpc \
     --data paths=/product.ProductService/GetProductList \
     --data name=getproductlist



   # Issue a gRPC request to the “GetProductList” method:

   $ grpcurl -d '{}' \
     -H 'kong-debug: 1' -plaintext \
     localhost:9080 product.ProductService/GetProductList



3). ##########################################################################################################################  
      Kong API Gateway(nodejs-kong-gateway)     -> Golang gRPC-Gateway (online-marketplace) service and route mapping
    ##########################################################################################################################
 
     Product -
     --------
     1). getproductlist

     # Service mapping

     curl -i -X POST --url http://localhost:8001/services/ --data 'name=grpc-gateway' --data 'url=http://grpc-gateway:8081'      

     # getproductlist route mapping

     curl -i -X POST --url http://localhost:8001/services/grpc-gateway/routes \
                     --data 'hosts[]=grpc-gateway' \
                     --data 'paths[]=/v1/product/getproductlist' \
                     --data 'strip_path=false' \
                     --data 'methods[]=GET'

    # Check getproductlist API response at 8000

    curl -i -X GET \
    --url http://localhost:8000/v1/product/getproductlist \
    --header 'Host: grpc-gateway'


4). #########################################################
            gin-restapi mapping to kong api gateway
    #########################################################
    

    curl -i -X POST --url http://localhost:8001/services/ --data 'name=gin-restapi' --data 'url=http://gin-restapi:8082' 

    # http://localhost:8080/books/1 -route 1

    curl -i -X POST --url http://localhost:8001/services/gin-restapi/routes \
                     --data 'hosts[]=gin-restapi' \
                     --data 'paths[]=/books/1' \
                     --data 'strip_path=false' \
                     --data 'methods[]=GET'


    curl -i -X GET \
    --url http://localhost:8000/books/1 \
    --header 'Host: gin-restapi'     

   # http://localhost:8080/books -route 2

   curl -i -X POST --url http://localhost:8001/services/gin-restapi/routes \
                     --data 'hosts[]=gin-restapi' \
                     --data 'paths[]=/books' \
                     --data 'strip_path=false' \
                     --data 'methods[]=GET'


    curl -i -X GET \
    --url http://localhost:8000/books \
    --header 'Host: gin-restapi'   
      


5). ############################################
         kong admin command
    ############################################

    # To Fetch a service details from Kong Api Gateway

      curl http://localhost:8001/services


    # To Delete a service from Kong Api Gateway

      http://localhost:8001/services/grpc-gateway


    # To Fetch routes details from Kong Api Gateway

      curl http://localhost:8001/routes


    # To Delete a routes from Kong Api Gateway

      http://localhost:8001/routes/5344a373-0237-4e86-99aa-77d979a3bc65

