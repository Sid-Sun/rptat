package proxy

import (
	"github.com/sid-sun/rptat/app/proxy/metrics"
	"github.com/sid-sun/rptat/app/proxy/proxy_handler"
	"github.com/sid-sun/rptat/app/proxy/service"
)

type Proxy struct {
	Hostname string
	Handler  *proxy_handler.Proxy
	Service  *service.Service
	Metrics  *metrics.Metrics
}
