package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"4d63.com/tz"
	"customer-service.com/customer-service/dbconn"
	custProto "customer-service.com/customer-service/proto/customer"
	"google.golang.org/grpc"
)

const (
	port = ":60061"
)

type Customer struct {
	Customers []*custProto.Customer
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

/* This function adds new Customer in the Customer repository */
func (prod *Customer) AddCustomer(ctx context.Context, req *custProto.Customer) (*custProto.CustomerResponse, error) {

	log.Printf(" AddCustomer Method - Started")

	session := dbconn.SessionSetUp()

	customer := prod.Customers

	log.Println(" Customer request :- ", req)

	dateTime := dateTime(req.RegisteredOn)

	fmt.Println(dateTime)

	err := session.Query("insert into customer(cust_id,fullname,address,mobile,email,registered_on) values(?,?,?,?,?,?)", req.CustId, req.Fullname, req.Address, req.Mobile, req.Email, dateTime).Exec()
	if err != nil {
		log.Println(" Error while inserting product record in the repository")
		log.Println(err)
	}

	customer = append(customer, req)

	response := &custProto.CustomerResponse{Flag: true, Custs: customer}

	log.Printf(" AddCustomer Method - Ended")

	return response, nil
}

/* This function updates Customer in the Customer repository */
func (prod *Customer) UpdateCustomer(ctx context.Context, req *custProto.CustomerUpdateRequest) (*custProto.CustomerUpdateResponse, error) {

	log.Printf(" UpdateProduct Method - Started")

	session := dbconn.SessionSetUp()

	//log.Println(" product request :- ", req)

	err := session.Query("update customer set mobile = ? where cust_id = ?", req.Mobile, req.CustId).Exec()
	if err != nil {
		log.Println(" Error while updating product record in the repository")
		log.Println(err)
	}

	response := &custProto.CustomerUpdateResponse{Updated: true}

	log.Printf(" UpdateProduct Method - Ended")

	return response, nil
}

/* This function delete a Customer from the Customer repository */
func (prod *Customer) DeleteCustomer(ctx context.Context, req *custProto.CustomerID) (*custProto.DeleteCustomerResponse, error) {

	log.Printf(" DeleteCustomer Method - Started")

	session := dbconn.SessionSetUp()

	log.Println(" customer request :- ", req.CustId)

	err := session.Query("delete from customer where cust_id = ?", req.CustId).Exec()
	if err != nil {
		log.Println(" Error while deleting the product record from repository")
		log.Println(err)
	}

	response := &custProto.DeleteCustomerResponse{Deleted: true}

	log.Printf(" DeleteCustomer Method - Ended")

	return response, nil
}

/* This function fetch the list of Customers */
func (prod *Customer) GetCustomerList(ctx context.Context, req *custProto.CustomerRequest) (*custProto.CustomerResponse, error) {

	log.Printf(" GetCustomerList Method - Started")

	session := dbconn.SessionSetUp()

	customer := prod.Customers

	m := map[string]interface{}{}

	/* Query execution */
	iter := session.Query("select * from customer").Iter()

	/* using iterator with mapScan() to get all the customer records */
	for iter.MapScan(m) {
		customer = append(customer, &custProto.Customer{
			CustId:       m["cust_id"].(string),
			Fullname:     m["fullname"].(string),
			Address:      m["address"].(string),
			Mobile:       m["mobile"].(int64), // Customer table have int value,so converting it to go 32 bit integer.
			Email:        m["email"].(string),
			RegisteredOn: m["registered_on"].(string),
		})

		m = map[string]interface{}{}
	}

	response := &custProto.CustomerResponse{Flag: true, Custs: customer}

	log.Printf(" GetCustomerList Method - Ended")

	return response, nil
}

/* This function fetch the Customer list by CustomerID */
func (prod *Customer) GetCustomerByID(ctx context.Context, req *custProto.CustomerID) (*custProto.CustomerResponse, error) {

	log.Printf(" GetProductByID Method - Started")

	session := dbconn.SessionSetUp()

	customer := prod.Customers

	m := map[string]interface{}{}

	iter := session.Query("select  * from customer where cust_id = ? ", req.CustId).Iter()

	for iter.MapScan(m) {
		customer = append(customer, &custProto.Customer{
			CustId:       m["cust_id"].(string),
			Fullname:     m["fullname"].(string),
			Address:      m["address"].(string),
			Mobile:       m["mobile"].(int64), // Customer table have int value,so converting it to go 32 bit integer.
			Email:        m["email"].(string),
			RegisteredOn: m["registered_on"].(string),
		})

		log.Println("List of customer  :- ", customer, " Number of Customer :- ", int32(len(customer)))

		m = map[string]interface{}{}
	}

	response := &custProto.CustomerResponse{Flag: true, Custs: customer}

	log.Printf(" GetProductByID Method - Ended")

	return response, nil
}

/* main function*/
func main() {
	log.Printf("customer-service - main function - started")

	// Set-up our gRPC server.
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()

	cust := &Customer{}

	/* Register our service with the gRPC server, this will tie our implementation into the auto-generated interface code for our protobuf definition. */

	custProto.RegisterCustomerServiceServer(srv, cust)

	log.Println("Running on port:", port)
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	log.Printf("customer-service - main function - ended")
}
