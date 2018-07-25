package tracing

import (
	"fmt"
	"io"
	"time"

	ot "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/transport/zipkin"
	zk "github.com/uber/jaeger-client-go/zipkin"
	"go.uber.org/zap"

	"github.com/yejiayu/go-cita/log"
)

var (
	httpTimeout = 5 * time.Second
	sampler     = jaeger.NewConstSampler(true)
	poolSpans   = jaeger.TracerOptions.PoolSpans(false)
	logger      = spanLogger{}
)

/* TODO:
 *   - Support only tracing when trace context information is already present (mixer)
 *   - Support tracing for some percentage of requests (pilot)
 */

type holder struct {
	closer io.Closer
	tracer ot.Tracer
}

// indirection for testing
type newZipkin func(url string, options ...zipkin.HTTPOption) (*zipkin.HTTPTransport, error)

// Configure initializes Istio's tracing subsystem.
//
// You typically call this once at process startup.
// Once this call returns, the tracing system is ready to accept data.
func Configure(serviceName string, url string) (io.Closer, error) {
	return configure(serviceName, url, zipkin.NewHTTPTransport)
}

func configure(serviceName, url string, nz newZipkin) (io.Closer, error) {
	reporters := make([]jaeger.Reporter, 0, 3)

	trans, err := nz(url, zipkin.HTTPLogger(logger), zipkin.HTTPTimeout(httpTimeout))
	if err != nil {
		return nil, fmt.Errorf("could not build zipkin reporter: %v", err)
	}
	reporters = append(reporters, jaeger.NewRemoteReporter(trans))

	var rep jaeger.Reporter
	if len(reporters) == 0 {
		// leave the default NoopTracer in place since there's no place for tracing to go...
		return holder{}, nil
	} else if len(reporters) == 1 {
		rep = reporters[0]
	} else {
		rep = jaeger.NewCompositeReporter(reporters...)
	}

	var tracer ot.Tracer
	var closer io.Closer

	zipkinPropagator := zk.NewZipkinB3HTTPHeaderPropagator()
	injector := jaeger.TracerOptions.Injector(ot.HTTPHeaders, zipkinPropagator)
	extractor := jaeger.TracerOptions.Extractor(ot.HTTPHeaders, zipkinPropagator)
	tracer, closer = jaeger.NewTracer(serviceName, sampler, rep, poolSpans, injector, extractor)

	// NOTE: global side effect!
	ot.SetGlobalTracer(tracer)

	return holder{
		closer: closer,
		tracer: tracer,
	}, nil
}

func (h holder) Close() error {
	if ot.GlobalTracer() == h.tracer {
		ot.SetGlobalTracer(ot.NoopTracer{})
	}

	var err error
	if h.closer != nil {
		err = h.closer.Close()
	}

	return err
}

type spanLogger struct{}

// Report implements the Report() method of jaeger.Reporter
func (spanLogger) Report(span *jaeger.Span) {
	log.Info("Reporting span",
		zap.String("operation", span.OperationName()),
		zap.String("span", span.String()))
}

// Close implements the Close() method of jaeger.Reporter.
func (spanLogger) Close() {}

// Error implements the Error() method of log.Logger.
func (spanLogger) Error(msg string) {
	log.Error(msg)
}

// Infof implements the Infof() method of log.Logger.
func (spanLogger) Infof(msg string, args ...interface{}) {
	log.Infof(msg, args...)
}
