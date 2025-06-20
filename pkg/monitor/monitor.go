package monitor

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	HttpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total HTTP requests",
		},
		[]string{"path", "method"},
	)

	OrdersCreatedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "orders_created_total",
			Help: "Total number of orders created",
		},
	)

	StockTransfersTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "stock_transfers_total",
			Help: "Total number of stock transfers between warehouses",
		},
	)

	StockReservedTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "stock_reserved_total",
			Help: "Total number of items reserved for orders",
		},
	)

	StockReleaseTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "stock_release_total",
			Help: "Total number of items released due to order expiration",
		},
	)
)

func Init() {
	prometheus.MustRegister(
		HttpRequestsTotal,
		OrdersCreatedTotal,
		StockTransfersTotal,
		StockReservedTotal,
		StockReleaseTotal,
	)
}

func Handler() http.Handler {
	return promhttp.Handler()
}
