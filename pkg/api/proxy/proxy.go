package proxy

import (
	"github.com/sid-sun/rptat/cmd/config"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/sid-sun/rptat/pkg/api/metrics"
	"go.uber.org/zap"
)

type Proxy struct {
	lgr   *zap.Logger
	tgt   *url.URL
	mt    *metrics.Metrics
	proxy *httputil.ReverseProxy
}

func NewProxy(cfg config.ProxyConfig, lgr *zap.Logger, mt *metrics.Metrics) (*Proxy, error) {
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

func (p *Proxy) MetricsProxyHandler() func(res http.ResponseWriter, req *http.Request) {
	return func(res http.ResponseWriter, req *http.Request) {
		go func() {
			err := p.mt.IncrementRequestCount(req.URL.Path)
			if err != nil {
				p.lgr.Sugar().Errorf("[Proxy] [MetricsProxyHandler] [Add] %s", err.Error())
			}
		}()

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
