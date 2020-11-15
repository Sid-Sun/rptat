package proxy

import (
	auth "github.com/abbot/go-http-auth"
	"github.com/sid-sun/rptat/app/proxy/metrics"
	"github.com/sid-sun/rptat/app/proxy/proxy_handler"
	"github.com/sid-sun/rptat/app/proxy/service"
)

// Proxy bundles the requirements for a new proxy
type Proxy struct {
	Hostname      string
	Handler       *proxy_handler.Proxy
	Service       *service.Service
	Metrics       *metrics.Metrics
	Authenticator *auth.DigestAuth
}
