package global

import (
	"os"
	"time"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/reporter"
	"log"
)

func newZPTracer() *zipkin.Tracer {
	var (
		err         error
		serviceName = Conf.Zipkin.ServiceName
		ZPReporter  = newZPReporter()
	)
	zEP, _ := zipkin.NewEndpoint(serviceName, "localhost"+Conf.GrpcServer.Addr, )
	ZipkinTracer, err := zipkin.NewTracer(
		ZPReporter, zipkin.WithLocalEndpoint(zEP),
	)
	if err != nil {
		log.Fatal("zipkin report error:", err)
		os.Exit(1)
	}
	return ZipkinTracer
}

func newZPReporter() reporter.Reporter {
	return zipkinhttp.NewReporter(
		Conf.Zipkin.Addr+"/api/v2/spans",
		zipkinhttp.Timeout(time.Duration(Conf.Zipkin.Reporter.Timeout)*time.Second),
		zipkinhttp.BatchSize(Conf.Zipkin.Reporter.BatchSize),
		zipkinhttp.BatchInterval(time.Duration(Conf.Zipkin.Reporter.BatchInterval)*time.Second),
	)
}