package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"

	"grpc-gateway.com/grpc-gateway/gateway"
)

var (
	productEndpoint = flag.String("Product_Service_endpoint", "product-service:60060", "product.service address")
	/*productEndpoint = flag.String("Product_Service_endpoint", "localhost:60060", "product.service address")*/
	customerEndpoint = flag.String("Customer_Service_endpoint", "customer-service:60061", "customer.service address")
	/*customerEndpoint = flag.String("Customer_Service_endpoint", "localhost:60061", "customer.service address")*/
	paymentEndpoint = flag.String("Payment_Service_endpoint", "payment-service:60062", "payment.service address")
	/*paymentEndpoint = flag.String("Payment_Service_endpoint", "localhost:60062", "payment.service address")*/
	orderEndpoint = flag.String("Order_Service_endpoint", "order-service:60063", "order.service address")
	/*orderEndpoint = flag.String("Order_Service_endpoint", "localhost:60063", "order.service address")*/
	port = flag.Int("p", 8081, "port of the service")
)

/* Create and start the server using that handler */
func main() {

	log.Printf("grpc-gateway - main  method  - STARTED")

	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux, err := gateway.New(ctx, *productEndpoint, *customerEndpoint, *paymentEndpoint, *orderEndpoint)
	if err != nil {
		log.Printf("Setting up the gateway: %s", err.Error())
		return
	}

	srvAddr := fmt.Sprintf(":%d", *port)

	s := &http.Server{
		Addr:    srvAddr,
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		log.Printf("Shutting down the http server")
		if err := s.Shutdown(context.Background()); err != nil {
			log.Printf("Failed to shutdown http server: %v", err)
		}
	}()

	log.Printf("Starting listening at %s", srvAddr)
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("Failed to listen and serve: %v", err)
	}

	log.Printf("grpc-gateway - main  method  - ENDED")
}
