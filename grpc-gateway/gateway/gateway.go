package gateway

import (
	"context"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/rakyll/statik/fs"
	"google.golang.org/grpc"

	_ "grpc-gateway.com/grpc-gateway/gateway/statik"
	customerGW "grpc-gateway.com/grpc-gateway/proto/customer"
	orderGW "grpc-gateway.com/grpc-gateway/proto/order"
	paymentGW "grpc-gateway.com/grpc-gateway/proto/payment"
	productGW "grpc-gateway.com/grpc-gateway/proto/product"
)

/* Creating the gateway http.Handler */
func New(ctx context.Context, productEndpoint, customerEndpoint, paymentEndpoint, orderEndpoint string) (http.Handler, error) {
	log.Printf("Gateway handler New method  - STARTED")

	gw := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	// product.microservice address
	if err := productGW.RegisterProductServiceHandlerFromEndpoint(ctx, gw, productEndpoint, opts); err != nil {
		return nil, err
	}

	// customer.microservice address
	if err := customerGW.RegisterCustomerServiceHandlerFromEndpoint(ctx, gw, customerEndpoint, opts); err != nil {
		return nil, err
	}

	// payment.microservice address
	if err := paymentGW.RegisterPaymentServiceHandlerFromEndpoint(ctx, gw, paymentEndpoint, opts); err != nil {
		return nil, err
	}

	// order.microservice address
	if err := orderGW.RegisterOrderServiceHandlerFromEndpoint(ctx, gw, orderEndpoint, opts); err != nil {
		return nil, err
	}

	statikFS, err := fs.New()
	if err != nil {
		return nil, err
	}
	mux := http.NewServeMux()
	mux.Handle("/swagger/", http.StripPrefix("/swagger/", http.FileServer(statikFS)))
	mux.Handle("/", gw)

	log.Printf("Gateway handler New method  - ENDED")

	return mux, nil
}
