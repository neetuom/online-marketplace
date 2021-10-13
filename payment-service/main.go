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
	"payment-service.com/payment-service/dbconn"
	pymtProto "payment-service.com/payment-service/proto/payment"
)

const (
	port = ":60062"
)

type Payment struct {
	Payments []*pymtProto.Payment
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

/* This function adds new Payment in the Payment repository */
func (pymt *Payment) AddPayment(ctx context.Context, req *pymtProto.Payment) (*pymtProto.PaymentResponse, error) {

	log.Printf(" AddPayment Method - Started")

	session := dbconn.SessionSetUp()

	payment := pymt.Payments

	log.Println(" Payment request :- ", req)

	orderDateTime := dateTime(req.OrderDate)
	invoiceDateTime := dateTime(req.InvoiceDate)

	fmt.Println(orderDateTime)
	fmt.Println(invoiceDateTime)

	err := session.Query("insert into payment(payment_id, order_id, order_date, invoice_number, invoice_date, payment_method, tracking_id, grand_total) values(?,?,?,?,?,?,?,?)", req.PaymentId, req.OrderId, orderDateTime, req.InvoiceNumber, invoiceDateTime, req.PaymentMethod, req.TrackingId, req.GrandTotal).Exec()
	if err != nil {
		log.Println(" Error while inserting payment record in the repository")
		log.Println(err)
	}

	payment = append(payment, req)

	response := &pymtProto.PaymentResponse{Flag: true, Pymts: payment}

	log.Printf(" AddPayment Method - Ended")

	return response, nil
}

/* This function updates payment in the payment repository */
func (pymt *Payment) UpdatePayment(ctx context.Context, req *pymtProto.PaymentUpdateRequest) (*pymtProto.PaymentUpdateResponse, error) {

	log.Printf(" UpdatePayment Method - Started")

	session := dbconn.SessionSetUp()

	//log.Println(" payment request :- ", req)

	err := session.Query("update payment set grand_total = ? where payment_id = ?", req.GrandTotal, req.PaymentId).Exec()
	if err != nil {
		log.Println(" Error while updating payment record in the repository")
		log.Println(err)
	}

	response := &pymtProto.PaymentUpdateResponse{Updated: true}

	log.Printf(" UpdatePayment Method - Ended")

	return response, nil
}

/* This function delete a payment from the payment repository */
func (pymt *Payment) DeletePayment(ctx context.Context, req *pymtProto.PaymentID) (*pymtProto.DeletePaymentResponse, error) {

	log.Printf(" DeletePayment Method - Started")

	session := dbconn.SessionSetUp()

	log.Println(" customer request :- ", req.PaymentId)

	err := session.Query("delete from payment where payment_id = ?", req.PaymentId).Exec()
	if err != nil {
		log.Println(" Error while deleting the payment record from repository")
		log.Println(err)
	}

	response := &pymtProto.DeletePaymentResponse{Deleted: true}

	log.Printf(" DeletePayment Method - Ended")

	return response, nil
}

/* This function fetch the list of payments */
func (pymt *Payment) GetPaymentList(ctx context.Context, req *pymtProto.PaymentRequest) (*pymtProto.PaymentResponse, error) {

	log.Printf(" GetPaymentList Method - Started")

	session := dbconn.SessionSetUp()

	payment := pymt.Payments

	m := map[string]interface{}{}

	/* Query execution */
	iter := session.Query("select * from payment").Iter()

	/* using iterator with mapScan() to get all the payment records */
	for iter.MapScan(m) {
		payment = append(payment, &pymtProto.Payment{
			PaymentId:     m["payment_id"].(string),
			OrderId:       m["order_id"].(string),
			OrderDate:     m["order_date"].(string),
			InvoiceNumber: m["invoice_number"].(string), // Payment table have int value,so converting it to go 32 bit integer.
			InvoiceDate:   m["invoice_date"].(string),
			PaymentMethod: m["payment_method"].(string),
			TrackingId:    m["tracking_id"].(string),
			GrandTotal:    m["grand_total"].(float64),
		})

		m = map[string]interface{}{}
	}

	response := &pymtProto.PaymentResponse{Flag: true, Pymts: payment}

	log.Printf(" GetPaymentList Method - Ended")

	return response, nil
}

/* This function fetch the payment list by paymentId */
func (pymt *Payment) GetPaymentByID(ctx context.Context, req *pymtProto.PaymentID) (*pymtProto.PaymentResponse, error) {

	log.Printf(" GetPaymentByID Method - Started")

	session := dbconn.SessionSetUp()

	payment := pymt.Payments

	m := map[string]interface{}{}

	iter := session.Query("select  * from payment where payment_id = ? ", req.PaymentId).Iter()

	for iter.MapScan(m) {
		payment = append(payment, &pymtProto.Payment{
			PaymentId:     m["payment_id"].(string),
			OrderId:       m["order_id"].(string),
			OrderDate:     m["order_date"].(string),
			InvoiceNumber: m["invoice_number"].(string), // Payment table have int value,so converting it to go 32 bit integer.
			InvoiceDate:   m["invoice_date"].(string),
			PaymentMethod: m["payment_method"].(string),
			TrackingId:    m["tracking_id"].(string),
			GrandTotal:    m["grand_total"].(float64),
		})

		log.Println("List of payment  :- ", payment, " Number of Payment :- ", int32(len(payment)))

		m = map[string]interface{}{}
	}

	response := &pymtProto.PaymentResponse{Flag: true, Pymts: payment}

	log.Printf(" GetPaymentByID Method - Ended")

	return response, nil
}

/* main function*/
func main() {
	log.Printf("payment-service - main function - started")

	// Set-up our gRPC server.
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()

	pymt := &Payment{}

	/* Register our service with the gRPC server, this will tie our implementation into the auto-generated interface code for our protobuf definition. */

	pymtProto.RegisterPaymentServiceServer(srv, pymt)

	log.Println("Running on port:", port)
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	log.Printf("payment-service - main function - ended")
}
