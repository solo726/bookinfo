// Code generated by truss. DO NOT EDIT.
// Rerunning truss will overwrite this file.
// Version: 7dc4d5d85c
// Version Date: 2018年 5月28日 星期一 22时12分59秒 UTC

// Package grpc provides a gRPC client for the BookDetails service.
package grpc

import (
	"context"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"

	// This Service
	"bookinfo/bookdetails-service/svc"
	pb "bookinfo/pb/details"
)

// New returns an service backed by a gRPC client connection. It is the
// responsibility of the caller to dial, and later close, the connection.
func New(conn *grpc.ClientConn, options ...ClientOption) (pb.BookDetailsServer, error) {
	var cc clientConfig

	for _, f := range options {
		err := f(&cc)
		if err != nil {
			return nil, errors.Wrap(err, "cannot apply option")
		}
	}

	clientOptions := []grpctransport.ClientOption{
		grpctransport.ClientBefore(
			contextValuesToGRPCMetadata(cc.headers)),
	}
	var detailEndpoint endpoint.Endpoint
	{
		detailEndpoint = grpctransport.NewClient(
			conn,
			"details.BookDetails",
			"Detail",
			EncodeGRPCDetailRequest,
			DecodeGRPCDetailResponse,
			pb.DetailResp{},
			clientOptions...,
		).Endpoint()
	}

	return svc.Endpoints{
		DetailEndpoint: detailEndpoint,
	}, nil
}

// GRPC Client Decode

// DecodeGRPCDetailResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC detail reply to a user-domain detail response. Primarily useful in a client.
func DecodeGRPCDetailResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.DetailResp)
	return reply, nil
}

// GRPC Client Encode

// EncodeGRPCDetailRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain detail request to a gRPC detail request. Primarily useful in a client.
func EncodeGRPCDetailRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.DetailReq)
	return req, nil
}

type clientConfig struct {
	headers []string
}

// ClientOption is a function that modifies the client config
type ClientOption func(*clientConfig) error

func CtxValuesToSend(keys ...string) ClientOption {
	return func(o *clientConfig) error {
		o.headers = keys
		return nil
	}
}

func contextValuesToGRPCMetadata(keys []string) grpctransport.ClientRequestFunc {
	return func(ctx context.Context, md *metadata.MD) context.Context {
		var pairs []string
		for _, k := range keys {
			if v, ok := ctx.Value(k).(string); ok {
				pairs = append(pairs, k, v)
			}
		}

		if pairs != nil {
			*md = metadata.Join(*md, metadata.Pairs(pairs...))
		}

		return ctx
	}
}
