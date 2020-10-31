package metrics

import (
	"sync"
	"time"

	"github.com/nsnikhil/go-datastructures/queue"
	"github.com/sid-sun/rptat/app/contract"
	"github.com/sid-sun/rptat/app/service"
	"github.com/sid-sun/rptat/cmd/config"
	"go.uber.org/zap"
)

const (
	SyncShutdown = true
	SyncNow      = false
)

type (
	// Metrics defines and implements the requisites for metrics
	Metrics struct {
		svc *service.Service

		lock     sync.Mutex
		syncChan chan bool
		ackChan  chan bool
		lgr      *zap.Logger
		total    uint

		request  instance
		response instance

		maxPending uint
	}
	instance struct {
		lock  sync.Mutex
		queue *queue.LinkedQueue
	}
)

// NewMetrics creates and returns a new Metrics instance
// It returns the Metrics object for registering new requests
// And A sync method for syncing current metrics with service
// An error is returned if initialization of requirements fail
func NewMetrics(svc *service.Service, cfg config.MetricsConfig) (*Metrics, error) {
	reqQ, err := queue.NewLinkedQueue()
	if err != nil {
		return nil, err
	}

	resQ, err := queue.NewLinkedQueue()
	if err != nil {
		return nil, err
	}

	m := &Metrics{
		request: instance{
			queue: reqQ,
		},
		response: instance{
			queue: resQ,
		},
		syncChan:   make(chan bool),
		ackChan:    make(chan bool),
		svc:        svc,
		maxPending: cfg.GetMinForSync(),
	}
	go m.sync()
	return m, nil
}

// IncrementRequestCount registers a new request at the specified path
func (m *Metrics) IncrementRequestCount(path string) error {
	m.request.lock.Lock()

	err := m.request.queue.Add(contract.Request{Path: path, Date: time.Now().Format("01-02-2006")})
	if err != nil {
		m.lgr.Sugar().Errorf("[Metrics] [IncrementRequestCount] [Add] %v", err)
		m.request.lock.Unlock()
		return err
	}

	m.request.lock.Unlock()

	if m.maxPending != 0 {
		m.syncWithIncrement()
	}

	return nil
}

// IncrementResponseCount registers a response and its status code
func (m *Metrics) IncrementResponseCount(path string, code int) error {
	m.response.lock.Lock()

	err := m.response.queue.Add(contract.Response{Path: path, Code: code, Date: time.Now().Format("01-02-2006")})
	if err != nil {
		m.lgr.Sugar().Errorf("[Metrics] [IncrementResponseCount] [Add] %v", err)
		m.response.lock.Unlock()
		return err
	}

	m.response.lock.Unlock()

	if m.maxPending != 0 {
		m.syncWithIncrement()
	}

	return nil
}

// syncWithIncrement increments total count and syncs metrics
// syncWithIncrement is a BLOCKING function
func (m *Metrics) syncWithIncrement() {
	// Increment
	m.lock.Lock()
	m.total++
	if m.total >= m.maxPending {
		m.lock.Unlock()
		m.SyncNow()
		return
	}
	m.lock.Unlock()
}

// SyncNow performs a *BLOCKING* sync of metrics to service
func (m *Metrics) SyncNow() {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.syncChan <- SyncNow
	m.total = 0

	<-m.ackChan
}

// Sync is a blocking method (meant to be called in a routine)
// which listens for sync calls from either SyncPeriodically, SyncNow or (syncWithIncrement indirectly)
// and syncs current metrics with service then commits them to store via service
func (m *Metrics) sync() {
	for {
		cont := <-m.syncChan

		reqMetrics := make(map[contract.Request]int)
		var totalIncrement int
		m.request.lock.Lock()
		for !m.request.queue.Empty() {
			elem, err := m.request.queue.Remove()
			if err != nil {
				m.lgr.Sugar().Error("[Metrics] [Sync] [ReqRemove] %s", err.Error())
			}
			reqMetrics[elem.(contract.Request)]++
			totalIncrement++
		}
		m.request.lock.Unlock()

		(*m.svc).RegisterRequests(reqMetrics)

		resMetrics := make(map[contract.Response]int)
		m.response.lock.Lock()
		for !m.response.queue.Empty() {
			elem, err := m.response.queue.Remove()
			if err != nil {
				m.lgr.Sugar().Error("[Metrics] [Sync] [ResRemove] %s", err.Error())
			}
			resMetrics[elem.(contract.Response)]++
			totalIncrement++
		}
		m.response.lock.Unlock()

		(*m.svc).RegisterResponses(resMetrics)

		err := (*m.svc).Commit()
		if err != nil {
			m.lgr.Sugar().Errorf("[Metrics] [Sync] [Commit] %v", err)
		}

		// Acknowledge sync completion
		m.ackChan <- true
		// If SyncShutdown break
		if cont == SyncShutdown {
			// panic("Wow")
			break
		}
	}
}

//SyncPeriodically starts a routine which syncs metrics at given duration
func (m *Metrics) SyncPeriodically(d time.Duration) {
	go func() {
		ticker := time.NewTicker(d)
		for {
			select {
			case <-ticker.C:
				m.SyncNow()
			}
		}
	}()
}
