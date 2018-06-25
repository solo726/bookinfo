package global

import (
	"context"
	"github.com/openzipkin/zipkin-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"github.com/openzipkin/zipkin-go/propagation/b3"
)

func NewGrpcClient(ctx context.Context, zipkinSpan zipkin.Span, grpcAddr string, f ExecFunc, opts ...grpc.DialOption) (grpcClient, error) {

	var c grpcClient
	md := metadata.New(make(map[string]string))
	b3.InjectGRPC(&md)(zipkinSpan.Context())

	ctx = metadata.NewOutgoingContext(
		ctx,
		md,
	)

	conn, err := grpc.DialContext(
		ctx,
		grpcAddr,
		opts...,
	)

	if err != nil {
		return c, err
	}

	c.ctx = ctx
	c.Conn = conn
	c.Func = f

	return c, nil
}

type ExecFunc func(ctx context.Context, conn *grpc.ClientConn) (interface{}, error)

type grpcClient struct {
	ctx  context.Context
	Conn *grpc.ClientConn
	Func ExecFunc
}

func (this grpcClient) Go() (interface{}, error) {
	defer this.Conn.Close()
	return this.Func(this.ctx, this.Conn)
}

func (this grpcClient) Close() {
	this.Conn.Close()
}
