package utils

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	conf "face-recognition-svc/app/config"

	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

var Tracer = opentracing.Tracer(nil)

func InitJaeger(c *conf.Config) (opentracing.Tracer, io.Closer, error) {
	cfg := config.Configuration{
		ServiceName: c.Jaeger.ServiceName,
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeRateLimiting,
			Param: float64(c.Jaeger.TracePerSecond), // 100 traces per second
		},
		Reporter: &config.ReporterConfig{
			LocalAgentHostPort: c.Jaeger.Host + ":" + c.Jaeger.Port,
		},
	}

	var err error
	var closer io.Closer
	Tracer, closer, err = cfg.NewTracer(config.Logger(jaeger.StdLogger))
	logrus.Printf("JAEGER RUNNING ON %s:%s\n", c.Jaeger.Host, c.Jaeger.Port)
	return Tracer, closer, err
}

func StartSpanFromRequest(tracer opentracing.Tracer, r *http.Request, funcDesc string) opentracing.Span {
	spanCtx, _ := Extract(tracer, r)
	return tracer.StartSpan(funcDesc, ext.RPCServerOption(spanCtx))
}

func StartSpan(e echo.Context, funcDesc string) (context.Context, opentracing.Span) {
	span := StartSpanFromRequest(Tracer, e.Request(), funcDesc)
	ctx := opentracing.ContextWithSpan(e.Request().Context(), span)
	span.SetTag("Route", e.Request().RequestURI)
	span.SetTag("Method", e.Request().Method)
	span.SetTag("Host", e.Request().Host)
	span.SetTag("Headers", e.Request().Header)
	return ctx, span
}

func SpanFromContext(ctx context.Context, funcDesc string) (opentracing.Span, context.Context) {
	span, ctx := opentracing.StartSpanFromContext(ctx, funcDesc)
	return span, ctx
}

func LogEvent(span opentracing.Span, desc string, event any) {
	if str, ok := event.(string); ok {
		// If event is a string, log it directly
		span.LogFields(log.Object(desc, str))
	} else {
		// If event is not a string, marshal it to JSON and log
		jsonData, err := json.Marshal(event)
		if err != nil {
			// If marshalling fails, log the error and event as a fallback
			span.LogFields(
				log.Object("error", err.Error()),
				log.Object(desc, "error marshalling event"),
			)
		} else {
			span.LogFields(
				log.Object(desc, string(jsonData)),
			)
		}
	}
}

func LogEventError(span opentracing.Span, err error) {
	span.SetTag("error", true)
	span.LogFields(log.String("error", err.Error()))
}

func Inject(span opentracing.Span, request *http.Request) error {
	return span.Tracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(request.Header))
}

func Extract(tracer opentracing.Tracer, r *http.Request) (opentracing.SpanContext, error) {
	return tracer.Extract(
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(r.Header))
}
