package proxy_handler

import (
	"github.com/sid-sun/rptat/cmd/config"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/sid-sun/rptat/app/proxy/metrics"
	"go.uber.org/zap"
)

// Proxy defines and implements necessities for proxy
type Proxy struct {
	lgr   *zap.Logger
	tgt   *url.URL
	mt    *metrics.Metrics
	proxy *httputil.ReverseProxy
}

// NewProxyHandler creates and returns a new Proxy with requsites initialized
// an error is returned if config doesn't define a valid URL
func NewProxyHandler(cfg config.ProxyConfig, lgr *zap.Logger, mt *metrics.Metrics) (*Proxy, error) {
	serveURL, err := url.Parse(cfg.GetServeURL())
	if err != nil {
		return nil, err
	}

	p := httputil.NewSingleHostReverseProxy(serveURL)
	p.Transport = &http.Transport{
		MaxIdleConns:        50,
		MaxIdleConnsPerHost: 50,
		ForceAttemptHTTP2:   true,
	}

	return &Proxy{
		lgr:   lgr,
		mt:    mt,
		proxy: p,
		tgt:   serveURL,
	}, nil
}

// MetricsProxyHandler proxies requests between source and target resource while collecting metrics
func (p *Proxy) MetricsProxyHandler() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		go func() {
			err := p.mt.IncrementRequestCount(req.URL.Path)
			if err != nil {
				p.lgr.Sugar().Errorf("[Proxy] [MetricsProxyHandler] [Add] %s", err.Error())
			}
		}()

		//p.lgr.Info(req.Host)
		p.serveReverseProxy(res, req)
	}
}

func (p *Proxy) serveReverseProxy(res http.ResponseWriter, req *http.Request) {
	// Update the headers to allow for SSL redirection
	req.URL.Host = p.tgt.Host
	req.URL.Scheme = p.tgt.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	req.Host = p.tgt.Host

	p.proxy.ModifyResponse = p.readResCode((*req).URL.Path)
	p.proxy.ServeHTTP(res, req)
}

func (p *Proxy) readResCode(path string) func(res *http.Response) error {
	return func(res *http.Response) error {
		go func() {
			err := p.mt.IncrementResponseCount(path, (*res).StatusCode)
			if err != nil {
				p.lgr.Sugar().Errorf("[Proxy] [MetricsProxyHandler] [Add] %s", err.Error())
			}
		}()
		return nil
	}
}
