// Package tracing provide dependency injection definitions.
package tracing

import (
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"

	"sber_cloud/tw/container"
	"sber_cloud/tw/definition/config"
)

// DefTracing definition name.
const DefTracing = "tracing"

// Tracer alias
type Tracer = opentracing.Tracer

// Definition init func.
func init() {
	container.Register(func(builder *container.Builder, params map[string]interface{}) error {
		var (
			subProcess string
			ok         bool
		)
		if subProcess, ok = params["sub_process"].(string); !ok || len(subProcess) == 0 {
			subProcess = "default"
		}
		return builder.Add(container.Def{
			Name: DefTracing,
			Build: func(container container.Container) (_ interface{}, err error) {
				var conf config.Config
				if err = container.Fill(config.DefConfig, &conf); err != nil {
					return nil, err
				}

				if !conf.GetBool("tracing.enabled") {
					return opentracing.NoopTracer{}, nil
				}

				var jg = jaegercfg.Configuration{
					ServiceName: conf.GetString("service"),
					Sampler: &jaegercfg.SamplerConfig{
						Type:  jaeger.SamplerTypeConst,
						Param: 1,
					},
					Reporter: &jaegercfg.ReporterConfig{
						LocalAgentHostPort: conf.GetString("tracing.url"),
					},
				}

				var tracer opentracing.Tracer
				if tracer, _, err = jg.NewTracer(
					jaegercfg.Tag("namespace", conf.GetString("namespace")),
					jaegercfg.Tag("instance", conf.GetString("instance")),
					jaegercfg.Tag("sub_process", subProcess),
					jaegercfg.Logger(jaegerlog.NullLogger),
					jaegercfg.Metrics(metrics.NullFactory),
				); err != nil {
					return nil, err
				}

				opentracing.SetGlobalTracer(tracer)

				return tracer, nil
			},
			Close: func(obj interface{}) (err error) {
				if c, ok := obj.(io.Closer); ok {
					return c.Close()
				}
				return nil
			},
		})
	})
}
