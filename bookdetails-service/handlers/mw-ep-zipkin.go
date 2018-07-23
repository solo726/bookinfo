package handlers

import (
	"time"
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/openzipkin/zipkin-go/model"
	"github.com/openzipkin/zipkin-go/propagation/b3"
	zipkingo "github.com/openzipkin/zipkin-go"

	"bookinfo/bookdetails-service/global"
	"google.golang.org/grpc/metadata"
)

func ZipkinEndpointMiddleware() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			zipkinTracer := global.ZPTracer
			//zipkinTracer, err := global.NewZipkinTracer()
			if err != nil {
				global.Logger.Error("zipkinTracer create failed,", err)
				return next(ctx, request)
			}
			var sc model.SpanContext
			md, _ := metadata.FromIncomingContext(ctx)

			sc = zipkinTracer.Extract(b3.ExtractGRPC(&md))

			//http 请求没有sc,需要完善 for kong zipkin
			//if len(md) == 0 {
			//	md = metadata.MD{}
			//	sc = ctx.Value("b3-http-con").(model.SpanContext)
			//}

			span := zipkinTracer.StartSpan("in service book-details", zipkingo.Parent(sc))

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
