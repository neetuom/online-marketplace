package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"grpc-gateway.com/grpc-gateway/gateway"
)

func init() {
	fmt.Println("krakend-grpc-gateway plugin loaded!!!")
}

//var GRPCRegisterer = registerer("grpc-gateway")
var ClientRegisterer = registerer("grpc-gateway")

type registerer string

func (r registerer) RegisterClients(f func(
	name string,
	handler func(context.Context, map[string]interface{}) (http.Handler, error),
)) {
	f(string(r), func(ctx context.Context, extra map[string]interface{}) (http.Handler, error) {
		fmt.Println("Plugin :: RegisterClients method get called")
		cfg := parse(extra)
		if cfg == nil {
			return nil, errors.New("wrong config")
		}
		if cfg.name != string(r) {
			return nil, fmt.Errorf("unknown register %s", cfg.name)
		}

		fmt.Println("cfg ::::->", cfg)
		return gateway.New(ctx, cfg.productEndpoint, cfg.customerEndpoint, cfg.paymentEndpoint, cfg.orderEndpoint)
	})
}

func parse(extra map[string]interface{}) *opts {
	name, ok := extra["name"].(string)
	if !ok {
		return nil
	}

	rawEs, ok := extra["endpoints"]
	if !ok {
		return nil
	}
	es, ok := rawEs.([]interface{})
	if !ok || len(es) < 4 {
		return nil
	}
	endpoints := make([]string, len(es))
	for i, e := range es {
		fmt.Println("e.(string) :::::->", e.(string))
		endpoints[i] = e.(string)
	}

	return &opts{
		name:             name,
		productEndpoint:  endpoints[0],
		customerEndpoint: endpoints[1],
		paymentEndpoint:  endpoints[2],
		orderEndpoint:    endpoints[3],
	}
}

type opts struct {
	name             string
	productEndpoint  string
	customerEndpoint string
	paymentEndpoint  string
	orderEndpoint    string
}

func main() {}
