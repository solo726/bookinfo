package server

import (
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"

	// 3d Party
	"google.golang.org/grpc"

	// This Service
	"bookinfo/bookcomments-service/handlers"
	"bookinfo/bookcomments-service/svc"
	pb "bookinfo/pb/comments"
	"bookinfo/bookcomments-service/global"
)


func NewEndpoints() svc.Endpoints {
	// Business domain.
	var service pb.BookCommentsServer
	{
		service = handlers.NewService()
		// Wrap Service with middlewares. See handlers/middlewares.go
		service = handlers.WrapService(service)
	}

	// Endpoint domain.
	var (
		postEndpoint = svc.MakePostEndpoint(service)
		getEndpoint  = svc.MakeGetEndpoint(service)
	)

	endpoints := svc.Endpoints{
		PostEndpoint: postEndpoint,
		GetEndpoint:  getEndpoint,
	}

	// Wrap selected Endpoints with middlewares. See handlers/middlewares.go
	endpoints = handlers.WrapEndpoints(endpoints)

	return endpoints
}

// Run starts a new http server, gRPC server, and a debug server with the
// passed config and logger
func Run() {
	endpoints := NewEndpoints()

	// Mechanical domain.
	errc := make(chan error)

	// Interrupt handler.
	go handlers.InterruptHandler(errc)

	// Debug listener.
	go func() {
		log.Println("transport", "debug", "addr", global.Conf.DebugServer.Addr)

		errc <- http.ListenAndServe(global.Conf.DebugServer.Addr, nil)
	}()
	//go func() {
	//	log.Println("transport", "debug", "addr", cfg.DebugAddr)
	//
	//	m := http.NewServeMux()
	//	m.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	//	m.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	//	m.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	//	m.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	//	m.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
	//
	//	errc <- http.ListenAndServe(cfg.DebugAddr, m)
	//}()

	// HTTP transport.
	go func() {
		log.Println("transport", "HTTP", "addr", global.Conf.HttpServer.Addr)
		h := svc.MakeHTTPHandler(endpoints)
		errc <- http.ListenAndServe(global.Conf.HttpServer.Addr, h)
	}()

	// gRPC transport.
	go func() {
		log.Println("transport", "gRPC", "addr", global.Conf.GrpcServer.Addr)
		ln, err := net.Listen("tcp", global.Conf.GrpcServer.Addr)
		if err != nil {
			errc <- err
			return
		}

		srv := svc.MakeGRPCServer(endpoints)
		s := grpc.NewServer()
		pb.RegisterBookCommentsServer(s, srv)

		errc <- s.Serve(ln)
	}()

	// Run!
	log.Println("exit", <-errc)
}
