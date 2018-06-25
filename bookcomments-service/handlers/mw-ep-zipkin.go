package handlers

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/openzipkin/zipkin-go/model"
	"google.golang.org/grpc/metadata"
	"github.com/openzipkin/zipkin-go/propagation/b3"
	zipkingo "github.com/openzipkin/zipkin-go"

	"bookinfo/bookcomments-service/global"
)

func ZipkinEndpointMiddleware() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			zipkinTracer := global.ZPTracer
			var sc model.SpanContext
			md, _ := metadata.FromIncomingContext(ctx)
			sc = zipkinTracer.Extract(b3.ExtractGRPC(&md))

			span := zipkinTracer.StartSpan("in service book-comments get comments", zipkingo.Parent(sc))
			span.Annotate(time.Now(), "in endpoint")

			defer func() {
				span.Annotate(time.Now(), "out endpoint")
				span.Finish()
			}()

			ctx = zipkingo.NewContext(ctx, span)
			return next(ctx, request)
		}
	}
}
