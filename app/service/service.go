package service

import (
	"encoding/json"
	"sync"

	"github.com/sid-sun/rptat/app/contract"
	"github.com/sid-sun/rptat/app/store"
	"go.uber.org/zap"
)

// Metrics defines the structure for data
type Metrics map[string]DailyMetrics

// DailyMetrics defines structure for a day's metrics
type DailyMetrics map[string]RouteMetrics

// RouteMetrics defines structure for metrics for a path
type RouteMetrics struct {
	Requests int         `json:"requests"`
	Response map[int]int `json:"responses"`
}

// Service defines the interface for a service
type Service interface {
	RegisterRequests(reqs map[contract.Request]int)
	RegisterResponses(res map[contract.Response]int)
	GetCurrentMetrics() Metrics
	Commit() error
}

type metricsService struct {
	mtx            *sync.Mutex
	currentMetrics *Metrics
	lgr            *zap.Logger
	str            *store.Store
}

// RegisterRequests registers the requests in map with local metrics
func (m *metricsService) RegisterRequests(reqs map[contract.Request]int) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	for reqElem, count := range reqs {
		// metOn is the metrics on the date
		metOn := (*m.currentMetrics)[reqElem.Date]
		if metOn == nil {
			metOn = make(DailyMetrics)
		}

		// metAt is metrics at path on the day
		metAt := metOn[reqElem.Path].Response
		if metAt == nil {
			metAt = make(map[int]int)
		}

		metOn[reqElem.Path] = RouteMetrics{
			Requests: metOn[reqElem.Path].Requests + count,
			Response: metAt,
		}

		(*m.currentMetrics)[reqElem.Date] = metOn
	}
}

// RegisterResponses registers the requests in map with local metrics
func (m *metricsService) RegisterResponses(res map[contract.Response]int) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	for response, responseCount := range res {
		// metOn is the metrics on the date
		metOn := (*m.currentMetrics)[response.Date]
		if metOn == nil {
			metOn = make(DailyMetrics)
		}

		// metAt is metrics at path on the day
		metAt := metOn[response.Path].Response
		if metAt == nil {
			metAt = make(map[int]int)
		}

		metAt[response.Code] += responseCount

		(*m.currentMetrics)[response.Date] = metOn
	}
}

// GetCurrentMetrics returns the current metrics details
func (m *metricsService) GetCurrentMetrics() Metrics {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	return *m.currentMetrics
}

// Commit syncs the local registered metrics with store
func (m *metricsService) Commit() error {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	raw, err := json.Marshal(*m.currentMetrics)
	if err != nil {
		m.lgr.Sugar().Errorf("[Service] [Commit] [Marshal] %v", err)
		return err
	}

	err = (*m.str).Write(raw)
	if err != nil {
		m.lgr.Sugar().Errorf("[Service] [Commit] [Write] %v", err)
		return err
	}
	return nil
}

// loadMetricsFromStore loads current metrics info from store
func (m *metricsService) loadMetricsFromStore() (*Metrics, error) {
	raw, err := (*m.str).Read()
	if err != nil {
		m.lgr.Sugar().Errorf("[Service] [loadMetricsFromStore] [Read] %v", err)
		return nil, err
	}

	var currentMetrics Metrics
	err = json.Unmarshal(raw, &currentMetrics)
	if err != nil {
		m.lgr.Sugar().Errorf("[Service] [loadMetricsFromStore] [Unmarshal] %v", err)
		return nil, err
	}

	return &currentMetrics, nil
}

// NewService creates and returns a service implementation
func NewService(str *store.Store, lgr *zap.Logger) (Service, error) {
	ms := &metricsService{
		lgr: lgr,
		str: str,
		mtx: &sync.Mutex{},
	}

	var err error
	ms.currentMetrics, err = ms.loadMetricsFromStore()
	if err != nil {
		return nil, err
	}

	return ms, nil
}
