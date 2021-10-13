package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"4d63.com/tz"
	"google.golang.org/grpc"
	"order-service.com/order-service/dbconn"
	oderProto "order-service.com/order-service/proto/order"
)

const (
	port = ":60063"
)

type Order struct {
	Orders []*oderProto.Order
}

/* This function calculate today's date & time */
func dateTime(dt string) string {

	ts := time.Now()

	timeStr := strconv.Itoa(ts.Hour()) + ":" + strconv.Itoa(ts.Minute()) + ":" + strconv.Itoa(ts.Second())

	fmt.Println(timeStr)

	loc, _ := tz.LoadLocation("Asia/Kolkata")

	const longform = "2006-01-02 15:04:05 MST"

	datetime := dt + " " + timeStr + " " + "IST"

	t, _ := time.ParseInLocation(longform, datetime, loc)

	return t.String()
}

/* This function adds new order in the orders repository */
func (oder *Order) AddOrder(ctx context.Context, req *oderProto.Order) (*oderProto.OrderResponse, error) {

	log.Printf(" AddOrder Method - Started")

	session := dbconn.SessionSetUp()

	order := oder.Orders

	log.Println(" Order request :- ", req)

	dateCreatedDateTime := dateTime(req.DateCreated)

	fmt.Println(dateCreatedDateTime)

	err := session.Query("insert into orders(order_id, product_id, cust_id, amount, date_created) values(?,?,?,?,?)", req.OrderId, req.ProductId, req.CustId, req.Amount, dateCreatedDateTime).Exec()
	if err != nil {
		log.Println(" Error while inserting record in the orders repository")
		log.Println(err)
	}

	order = append(order, req)

	response := &oderProto.OrderResponse{Flag: true, Oders: order}

	log.Printf(" AddOrder Method - Ended")

	return response, nil
}

/* This function updates order in the orders repository */
func (oder *Order) UpdateOrder(ctx context.Context, req *oderProto.OrderUpdateRequest) (*oderProto.OrderUpdateResponse, error) {

	log.Printf(" UpdateOrder Method - Started")

	session := dbconn.SessionSetUp()

	//log.Println(" order request :- ", req)

	err := session.Query("update orders set amount = ? where order_id = ?", req.Amount, req.OrderId).Exec()
	if err != nil {
		log.Println(" Error while updating order record in the repository")
		log.Println(err)
	}

	response := &oderProto.OrderUpdateResponse{Updated: true}

	log.Printf(" UpdateOrder Method - Ended")

	return response, nil
}

/* This function delete a order from the orders repository */
func (oder *Order) DeleteOrder(ctx context.Context, req *oderProto.OrderID) (*oderProto.DeleteOrderResponse, error) {

	log.Printf(" DeleteOrder Method - Started")

	session := dbconn.SessionSetUp()

	log.Println(" order request :- ", req.OrderId)

	err := session.Query("delete from orders where order_id = ?", req.OrderId).Exec()
	if err != nil {
		log.Println(" Error while deleting the order record from repository")
		log.Println(err)
	}

	response := &oderProto.DeleteOrderResponse{Deleted: true}

	log.Printf(" DeleteOrder Method - Ended")

	return response, nil
}

/* This function fetch the list of orders */
func (oder *Order) GetOrderList(ctx context.Context, req *oderProto.OrderRequest) (*oderProto.OrderResponse, error) {

	log.Printf(" GetOrderList Method - Started")

	session := dbconn.SessionSetUp()

	order := oder.Orders

	m := map[string]interface{}{}

	/* Query execution */
	iter := session.Query("select * from orders").Iter()

	/* using iterator with mapScan() to get all the order records */
	for iter.MapScan(m) {
		order = append(order, &oderProto.Order{
			OrderId:     m["order_id"].(string),
			ProductId:   m["product_id"].(string),
			CustId:      m["cust_id"].(string),
			Amount:      m["amount"].(float64),
			DateCreated: m["date_created"].(string),
		})

		m = map[string]interface{}{}
	}

	response := &oderProto.OrderResponse{Flag: true, Oders: order}

	log.Printf(" GetOrderList Method - Ended")

	return response, nil
}

/* This function fetch the order list by orderId */
func (oder *Order) GetOrderByID(ctx context.Context, req *oderProto.OrderID) (*oderProto.OrderResponse, error) {

	log.Printf(" GetOrderByID Method - Started")

	session := dbconn.SessionSetUp()

	order := oder.Orders

	m := map[string]interface{}{}

	iter := session.Query("select  * from orders where order_id = ? ", req.OrderId).Iter()

	for iter.MapScan(m) {
		order = append(order, &oderProto.Order{
			OrderId:     m["order_id"].(string),
			ProductId:   m["product_id"].(string),
			CustId:      m["cust_id"].(string),
			Amount:      m["amount"].(float64),
			DateCreated: m["date_created"].(string),
		})

		log.Println("List of order  :- ", order, " Number of Order :- ", int32(len(order)))

		m = map[string]interface{}{}
	}

	response := &oderProto.OrderResponse{Flag: true, Oders: order}

	log.Printf(" GetOrderByID Method - Ended")

	return response, nil
}

/* main function*/
func main() {
	log.Printf("order-service - main function - started")

	// Set-up our gRPC server.
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()

	oder := &Order{}

	/* Register our service with the gRPC server, this will tie our implementation into the auto-generated interface code for our protobuf definition. */

	oderProto.RegisterOrderServiceServer(srv, oder)

	log.Println("Running on port:", port)
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	log.Printf("order-service - main function - ended")
}
