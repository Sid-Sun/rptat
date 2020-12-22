package app

import (
	"context"
	"fmt"
	"github.com/go-chi/hostrouter"
	"github.com/sid-sun/rptat/app/proxy"
	"github.com/sid-sun/rptat/app/router"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/sid-sun/rptat/cmd/config"
	"go.uber.org/zap"
)

var proxies map[string]*proxy.Proxy
var m sync.Mutex
var pidPath string

// StartServer starts the proxy, inits all the requited submodules and routine for shutdown
func StartServer(cfg config.Config, logger *zap.Logger) {
	proxies = make(map[string]*proxy.Proxy)

	proxyRouter := router.NewRouter()
	proxyRouter.Mount("/", initHostRouter(cfg.ProxyConfig, logger))
	proxyServer := &http.Server{Addr: cfg.API.Address(), Handler: proxyRouter}

	logger.Info(fmt.Sprintf("[StartServer] [Server] Listening on %s", cfg.API.Address()))
	go func() {
		if err := proxyServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(fmt.Sprintf("[StartServer] [Server] [ListenAndServe]: %s", err.Error()))
			panic(err)
		}
	}()

	c := make(chan os.Signal, 1)
	go liveReload(proxyServer, c, logger)
	gracefulShutdown(proxyServer, logger, c)
}

func initHostRouter(pxyCfg []config.ProxyConfig, logger *zap.Logger) hostrouter.Routes {
	hr := hostrouter.New()
	for _, pxy := range pxyCfg {
		prox := proxy.NewProxy(&pxy, logger)
		proxies[pxy.GetHostname()] = &prox
		hr.Map(pxy.GetHostname(), router.NewProxyRouter(prox, logger))
		logger.Sugar().Infof("[initHostRouter] [Map] Subscribed [%s] as [%s]", pxy.GetServeURL(), pxy.GetHostname())
	}
	return hr
}

func liveReload(httpServer *http.Server, shutdownChan chan os.Signal, logger *zap.Logger) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGUSR1)
	err := createPIDFile()
	if err != nil {
		logger.Sugar().Errorf("[liveReload] [createPIDFile] %v", err)
		self, _ := os.FindProcess(os.Getpid())
		_ = self.Signal(os.Interrupt)
	}
	for {
		select {
		case <-c:
			m.Lock()

			cfg, err := config.Load()
			if err != nil {
				logger.Sugar().Error("[liveReload] [Load] Cancelling reload")
				continue
			}
			pxyCfg := cfg.ProxyConfig

			// Save old proxies for shutdown
			// and create a new map for new proxies
			oldProx := proxies
			proxies = make(map[string]*proxy.Proxy)

			proxyRouter := router.NewRouter()
			proxyRouter.Mount("/", initHostRouter(pxyCfg, logger))

			httpServer.Handler = proxyRouter
			logger.Sugar().Infof("[liveReload] Mounted new router")

			for hostname, pxy := range oldProx {
				pxy.Metrics.SyncAndShutdown()
				logger.Sugar().Infof("[liveReload] [SyncAndShutdown] Shutdown [%s]", hostname)
			}

			m.Unlock()
		case <-shutdownChan:
			go func() {
				if pidPath != "" {
					err := os.Remove(pidPath)
					if err != nil {
						panic(err)
					}
				}
			}()
			shutdownChan <- os.Interrupt
			return
		}
	}
}

func gracefulShutdown(httpServer *http.Server, logger *zap.Logger, shutdownChan chan os.Signal) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)
	<-c
	shutdownChan <- os.Interrupt
	logger.Info("Attempting GracefulShutdown")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		if err := httpServer.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
			logger.Error(fmt.Sprintf("[GracefulShutdown] [Server] [Shutdown]: %s", err.Error()))
			panic(err)
		}
	}()

	m.Lock()
	for _, pxy := range proxies {
		(*pxy).Metrics.SyncAndShutdown()
	}
	<-shutdownChan
}

func createPIDFile() error {
	mode := 480 // 110-110-000 (rw-r-----)
	path := "/var/run/rptat"
	err := os.MkdirAll(path, os.FileMode(mode))
	if err != nil && !os.IsExist(err) {
		if !os.IsPermission(err) {
			pidPath = ""
			return err
		}
		// Fallback to user's run dir.
		path = "/var/run/user/" + strconv.Itoa(os.Getuid()) + "/rptat"
		err = os.Mkdir(path, os.FileMode(mode))
		if err != nil && !os.IsExist(err) {
			pidPath = ""
			return err
		}
	}
	pidPath = path + "/rptat.pid"

	info, err := os.Stat(pidPath)
	if err == nil || info != nil {
		pidPath = ""
		return err
	}

	err = ioutil.WriteFile(pidPath, []byte(strconv.Itoa(os.Getpid())), os.FileMode(mode))
	if err != nil {
		pidPath = ""
		return err
	}
	return nil
}
