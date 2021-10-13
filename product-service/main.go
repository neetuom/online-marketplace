package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"product-service.com/product-service/models"
	prodProto "product-service.com/product-service/proto/product"
)

const (
	port = ":60060"
)

type Product struct {
	Products []*prodProto.Product
}

/* This function adds new Product in the Product repository */
func (prod *Product) AddProduct(ctx context.Context, req *prodProto.Product) (*prodProto.ProductResponse, error) {

	log.Printf(" AddProduct Method - Started")

	session := models.SessionSetUp()

	product := prod.Products

	//log.Println(" product request :- ", req)

	err := session.Query("insert into product(product_id,name,size,colour,price,quantity,description) values(?,?,?,?,?,?,?)", req.ProductId, req.Name, req.Size, req.Colour, req.Price, req.Quantity, req.Description).Exec()
	if err != nil {
		log.Println(" Error while inserting product record in the repository")
		log.Println(err)
	}

	product = append(product, req)

	response := &prodProto.ProductResponse{Flag: true, Prods: product}

	log.Printf(" AddProduct Method - Ended")

	return response, nil
}

/* This function updates Product in the Product repository */
func (prod *Product) UpdateProduct(ctx context.Context, req *prodProto.UpdateRequest) (*prodProto.UpdateResponse, error) {

	log.Printf(" UpdateProduct Method - Started")

	session := models.SessionSetUp()

	//log.Println(" product request :- ", req)

	err := session.Query("update product set price = ? where product_id = ?", req.Price, req.ProductID).Exec()
	if err != nil {
		log.Println(" Error while updating product record in the repository")
		log.Println(err)
	}

	response := &prodProto.UpdateResponse{Updated: true}

	log.Printf(" UpdateProduct Method - Ended")

	return response, nil
}

/* This function delete a Product from the Product repository */
func (prod *Product) DeleteProduct(ctx context.Context, req *prodProto.ProductID) (*prodProto.DeleteResponse, error) {

	log.Printf(" DeleteProduct Method - Started")

	session := models.SessionSetUp()

	log.Println(" product request :- ", req.ProductID)

	err := session.Query("delete from product where product_id = ?", req.ProductID).Exec()
	if err != nil {
		log.Println(" Error while deleting the product record from repository")
		log.Println(err)
	}

	response := &prodProto.DeleteResponse{Deleted: true}

	log.Printf(" DeleteProduct Method - Ended")

	return response, nil
}

/* This function fetch the list of Products */
func (prod *Product) GetProductList(ctx context.Context, req *prodProto.ProductRequest) (*prodProto.ProductResponse, error) {

	log.Printf(" GetProductList Method - Started")

	session := models.SessionSetUp()

	product := prod.Products

	m := map[string]interface{}{}

	iter := session.Query("select * from Product").Iter()

	for iter.MapScan(m) {
		product = append(product, &prodProto.Product{
			ProductId:   m["product_id"].(string),
			Name:        m["name"].(string),
			Size:        m["size"].(string),
			Colour:      m["colour"].(string),
			Price:       m["price"].(float64),
			Quantity:    int32(m["quantity"].(int)), // Product table have int value,so converting it to go 32 bit integer.
			Description: m["description"].(string),
		})
		//log.Println("List of product  :- ", product, " Number of Product :- ", int32(len(product)))
		//log.Printf(m["product_id"].(string), m["name"].(string), m["size"].(string), m["colour"].(string), m["price"].(float64), int32(m["quantity"].(int)), m["description"].(string))
		m = map[string]interface{}{}
	}

	response := &prodProto.ProductResponse{Flag: true, Prods: product}

	log.Printf(" GetProductList Method - Ended")

	return response, nil
}

/* This function fetch the Product list by ProductID */
func (prod *Product) GetProductByID(ctx context.Context, req *prodProto.ProductID) (*prodProto.ProductResponse, error) {

	log.Printf(" GetProductByID Method - Started")

	session := models.SessionSetUp()

	product := prod.Products

	m := map[string]interface{}{}

	iter := session.Query("select  * from product where product_id = ? ", req.ProductID).Iter()

	for iter.MapScan(m) {
		product = append(product, &prodProto.Product{
			ProductId:   m["product_id"].(string),
			Name:        m["name"].(string),
			Size:        m["size"].(string),
			Colour:      m["colour"].(string),
			Price:       m["price"].(float64),
			Quantity:    int32(m["quantity"].(int)), // Product table have int value,so converting it to go 32 bit integer.
			Description: m["description"].(string),
		})
		log.Println("List of product  :- ", product, " Number of Product :- ", int32(len(product)))
		//log.Printf(m["product_id"].(string), m["name"].(string), m["size"].(string), m["colour"].(string), m["price"].(float64), int32(m["quantity"].(int)), m["description"].(string))
		m = map[string]interface{}{}
	}

	response := &prodProto.ProductResponse{Flag: true, Prods: product}

	log.Printf(" GetProductByID Method - Ended")

	return response, nil
}

/* main function*/
func main() {
	log.Printf("product-service - main function - started")

	// Set-up our gRPC server.
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()

	prod := &Product{}

	/* Register our service with the gRPC server, this will tie our implementation into the auto-generated interface code for our protobuf definition. */
	prodProto.RegisterProductServiceServer(srv, prod)

	/* Register reflection service on gRPC server. */
	reflection.Register(srv)

	log.Println("Running on port:", port)
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	log.Printf("product-service - main function - ended")
}
