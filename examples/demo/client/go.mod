module github.com/open-telemetry/opentelemetry-collector-contrib/examples/demo/client

go 1.19

require (
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.41.1
	go.opentelemetry.io/otel v1.15.1
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v0.38.1
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.15.1
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.15.1
	go.opentelemetry.io/otel/metric v0.38.1
	go.opentelemetry.io/otel/sdk v1.15.1
	go.opentelemetry.io/otel/sdk/metric v0.38.1
	google.golang.org/grpc v1.60.1
)

require (
	github.com/cenkalti/backoff/v4 v4.2.1 // indirect
	github.com/felixge/httpsnoop v1.0.3 // indirect
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.15.2 // indirect
	go.opentelemetry.io/otel/exporters/otlp/internal/retry v1.15.1 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric v0.38.1 // indirect
	go.opentelemetry.io/otel/trace v1.15.1 // indirect
	go.opentelemetry.io/proto/otlp v0.19.0 // indirect
	golang.org/x/net v0.16.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	google.golang.org/genproto v0.0.0-20231002182017-d307bd883b97 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20231002182017-d307bd883b97 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20231002182017-d307bd883b97 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
)

retract (
	v0.76.2
	v0.76.1
	v0.65.0
)
