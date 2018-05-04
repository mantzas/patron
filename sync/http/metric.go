package http

import (
	"context"

	"github.com/mantzas/patron/log"
	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

var (
	defaultLatencyDistribution = view.Distribution(0, 1, 2, 3, 4, 5, 6, 8, 10, 13, 16, 20, 25, 30, 40, 50, 65, 80, 100, 130, 160, 200, 250, 300, 400, 500, 650, 800, 1000, 2000, 5000, 10000, 20000, 50000, 100000)
	serverRequestCount         = stats.Int64("http/server/request_count", "Number of HTTP requests started", stats.UnitDimensionless)
	serverLatency              = stats.Float64("http/server/latency_ms", "End-to-end latency", stats.UnitMilliseconds)
	statusKey                  tag.Key
	pathKey                    tag.Key
	methodKey                  tag.Key
	hostKey                    tag.Key

	serverRequestCountView = &view.View{
		Name:        "http/server/request_count",
		Description: "Count of HTTP requests started",
		Measure:     serverRequestCount,
		Aggregation: view.Count(),
	}

	serverLatencyView = &view.View{
		Name:        "http/server/latency_ms",
		Description: "Latency distribution of HTTP requests in ms",
		Measure:     serverLatency,
		TagKeys:     []tag.Key{hostKey, methodKey, pathKey, statusKey},
		Aggregation: defaultLatencyDistribution,
	}

	defaultServerViews = []*view.View{serverRequestCountView, serverLatencyView}
)

func init() {
	var err error
	statusKey, err = tag.NewKey("http.request.status")
	if err != nil {
		log.Fatal("failed to set tag status code")
	}
	pathKey, err = tag.NewKey("http.request.path")
	if err != nil {
		log.Fatal("failed to set tag path")
	}
	methodKey, err = tag.NewKey("http.request.method")
	if err != nil {
		log.Fatal("failed to set tag method")
	}
	hostKey, err = tag.NewKey("http.request.host")
	if err != nil {
		log.Fatal("failed to set tag host")
	}
}

func recordMetric(ctx context.Context, host, method, path, status string, latency float64) {
	ctx, err := tag.New(ctx, tag.Upsert(methodKey, method), tag.Upsert(pathKey, path),
		tag.Upsert(statusKey, status), tag.Upsert(hostKey, host))
	if err != nil {
		log.Errorf("failed to set tags")
	}
	stats.Record(ctx, serverLatency.M(latency), serverRequestCount.M(1))
}
