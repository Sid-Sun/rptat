package service

import (
	"encoding/json"
	"sync"

	"github.com/sid-sun/rptat/app/contract"
	"github.com/sid-sun/rptat/app/store"
	"go.uber.org/zap"
)

type metrics map[string]dailyMetrics

type dailyMetrics map[string]routeMetrics

type routeMetrics struct {
	Requests int         `json:"requests"`
	Response map[int]int `json:"responses"`
}

type Service interface {
	RegisterRequests(reqs map[contract.Request]int)
	RegisterResponses(res map[contract.Response]int)
	GetCurrentMetrics() metrics
	Commit() error
}

type metricsService struct {
	mtx            *sync.Mutex
	currentMetrics *metrics
	lgr            *zap.Logger
	str            *store.Store
}

func (m *metricsService) RegisterRequests(reqs map[contract.Request]int) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	for reqElem, count := range reqs {
		// metOn is the metrics on the date
		metOn := (*m.currentMetrics)[reqElem.Date]
		if metOn == nil {
			metOn = make(dailyMetrics)
		}

		// metAt is metrics at path on the day
		metAt := metOn[reqElem.Path].Response
		if metAt == nil {
			metAt = make(map[int]int)
		}

		metOn[reqElem.Path] = routeMetrics{
			Requests: metOn[reqElem.Path].Requests + count,
			Response: metAt,
		}

		(*m.currentMetrics)[reqElem.Date] = metOn
	}
}

func (m *metricsService) RegisterResponses(res map[contract.Response]int) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	for response, responseCount := range res {
		// metOn is the metrics on the date
		metOn := (*m.currentMetrics)[response.Date]
		if metOn == nil {
			metOn = make(dailyMetrics)
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
func (m *metricsService) GetCurrentMetrics() metrics {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	return *m.currentMetrics
}

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

func (m *metricsService) loadMetricsFromStore() (*metrics, error) {
	raw, err := (*m.str).Read()
	if err != nil {
		m.lgr.Sugar().Errorf("[Service] [loadMetricsFromStore] [Read] %v", err)
		return nil, err
	}

	var currentMetrics metrics
	err = json.Unmarshal(raw, &currentMetrics)
	if err != nil {
		m.lgr.Sugar().Errorf("[Service] [loadMetricsFromStore] [Unmarshal] %v", err)
		return nil, err
	}

	return &currentMetrics, nil
}

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
