package handlers

import (
	"context"
	"time"

	"github.com/openzipkin/zipkin-go"
	"google.golang.org/grpc"

	"bookinfo/bookdetails-service/global"
	"bookinfo/bookdetails-service/models"
	pb "bookinfo/pb/details"
	cpb "bookinfo/pb/comments"
)

// NewService returns a naïve, stateless implementation of Service.
func NewService() pb.BookDetailsServer {
	return bookdetailsService{}
}

type bookdetailsService struct{}

type SvcMiddleware func(pb.BookDetailsServer) pb.BookDetailsServer

// Detail implements Service.
func (s bookdetailsService) Detail(ctx context.Context, in *pb.DetailReq) (*pb.DetailResp, error) {
	span, newCtx := global.ZPTracer.StartSpanFromContext(
		ctx,
		"getting-book-info",
		zipkin.Tags(map[string]string{"func": "detail"}),
	)
	span.Tag("l", "v")
	span.Annotate(time.Now(), "detail svc--start...")
	defer func() {
		span.Annotate(time.Now(), "detail svc--end...")
		span.Finish()
	}()
	var resp pb.DetailResp
	{
		if in.Id == 0 {
			resp.Code = global.ERROR_PARAMS_ERROR.Code
			resp.Msg = global.ERROR_PARAMS_ERROR.Msg
			return &resp, nil
		}
		book := models.Books{}

		global.BOOK_DB.WarpRawScan(newCtx, &book, "select * from books where id = ?", in.Id)

		if book.ID == 0 {
			resp.Code = global.ERROR_RESOURCE_NOT_FOUND.Code
			resp.Msg = global.ERROR_RESOURCE_NOT_FOUND.Msg
			return &resp, nil
		}

		resp.Code = global.SUCCESS.Code
		resp.Msg = global.SUCCESS.Msg
		resp.Data = &pb.DetailRespData{
			Id:    int32(book.ID),
			Name:  book.Name,
			Intro: book.Intro,
		}
	}

	//comments
	{
		c, _ := global.NewGrpcClient(
			newCtx,
			span,
			global.Conf.Servers.BookComments.Grpc,
			func(ctx context.Context, conn *grpc.ClientConn) (resp interface{}, err error) {
				c := cpb.NewBookCommentsClient(conn)
				resp, err = c.Get(ctx, &cpb.GetReq{Id: 1})

				return
			},
			grpc.WithInsecure(),
			grpc.WithTimeout(10*time.Second),
		)
		res, err := c.Go()

		if err != nil {
			return &resp, nil
		}

		commentsResp := res.(*cpb.GetResp)
		if commentsResp.Code != global.SUCCESS.Code {
			return &resp, nil
		}
		resp.Data.Comments = commentsResp.Data
	}
	//resp = pb.DetailResp{
	//// Code:
	//// Msg:
	//// Data:
	//}
	return &resp, nil
}
