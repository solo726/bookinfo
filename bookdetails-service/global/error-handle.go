package global

import (
	"net/http"
	"runtime/debug"
	"context"

	"google.golang.org/grpc"
	"log"
	"github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware"
)

type HttpServerPanicHandler struct {
	http.Handler
}

func (h HttpServerPanicHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if e := recover(); e != nil {
			Logger.Errorln(e, string(debug.Stack()))
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Keep calm! Try again after prayer."))
		}
	}()
	h.Handler.ServeHTTP(w, r)
}

func GrpcServerPanicHandlerOptions() []grpc.ServerOption{
	var recoveryOpts = []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(func(p interface{}) error {
			Logger.Errorln("this is panic")
			Logger.Errorln(p, string(debug.Stack()))
			return nil
		}),
	}

	var serverOpts = []grpc.ServerOption{

		grpc_middleware.WithUnaryServerChain(
			grpc_recovery.UnaryServerInterceptor(recoveryOpts...),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_recovery.StreamServerInterceptor(recoveryOpts...),
		),
	}

	return serverOpts
}

func loadGrpcOpts() []grpc.ServerOption {
	var unaryInterceptor grpc.UnaryServerInterceptor
	var streamServerInterceptor = func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		defer func() {
			log.Println("ab")
			if e := recover(); e != nil {
				Logger.Errorln(info, string(debug.Stack()))
			}
		}()
		// 继续处理请求
		return handler(srv, ss)
	}

	unaryInterceptor = func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		defer func() {
			log.Println("ac")
			if e := recover(); e != nil {
				Logger.Errorln(info, string(debug.Stack()), resp, err)
				return
			}
		}()
		// 继续处理请求
		return handler(ctx, req)
	}

	var opts []grpc.ServerOption
	opts = append(opts, grpc.UnaryInterceptor(unaryInterceptor), grpc.StreamInterceptor(streamServerInterceptor))
	return opts
}
