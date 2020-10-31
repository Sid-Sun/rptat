package service

import (
	"encoding/json"

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

// TODO: Implement Commit option and Threadsafety
type Service interface {
	RegisterRequests(reqs map[contract.Request]int) error
	RegisterResponses(res map[contract.Response]int) error
}

type metricsService struct {
	lgr *zap.Logger
	str *store.Store
}

func (m *metricsService) RegisterRequests(reqs map[contract.Request]int) error {
	currentMetrics, err := m.getCurrentMetrics()
	if err != nil {
		m.lgr.Sugar().Errorf("[Service] [RegisterRequests] [getCurrentMetrics] %v", err)
		return err
	}

	for reqElem, count := range reqs {
		foo := currentMetrics[reqElem.Date]
		if foo == nil {
			foo = make(dailyMetrics)
		}

		boo := foo[reqElem.Path].Response
		if boo == nil {
			boo = make(map[int]int)
		}

		foo[reqElem.Path] = routeMetrics{
			Requests: foo[reqElem.Path].Requests + count,
			Response: boo,
		}

		currentMetrics[reqElem.Date] = foo
	}

	raw, err := json.Marshal(currentMetrics)
	if err != nil {
		m.lgr.Sugar().Errorf("[Service] [RegisterRequests] [Marshal] %v", err)
		return err
	}

	err = (*m.str).Write(raw)
	if err != nil {
		m.lgr.Sugar().Errorf("[Service] [RegisterRequests] [Write] %v", err)
		return err
	}

	return nil
}

func (m *metricsService) RegisterResponses(res map[contract.Response]int) error {
	currentMetrics, err := m.getCurrentMetrics()
	if err != nil {
		m.lgr.Sugar().Errorf("[Service] [RegisterResponses] [getCurrentMetrics] %v", err)
		return err
	}

	// TODO: Refactor to be independent from relying on reqs to make struct
	for response, responseCount := range res {
		currentMetrics[response.Date][response.Path].Response[response.Code] += responseCount
	}

	raw, err := json.Marshal(currentMetrics)
	if err != nil {
		m.lgr.Sugar().Errorf("[Service] [RegisterRequests] [Marshal] %v", err)
		return err
	}

	err = (*m.str).Write(raw)
	if err != nil {
		m.lgr.Sugar().Errorf("[Service] [RegisterRequests] [Write] %v", err)
		return err
	}

	return nil
}

func (m *metricsService) getCurrentMetrics() (metrics, error) {
	raw, err := (*m.str).Read()
	if err != nil {
		m.lgr.Sugar().Errorf("[Service] [getCurrentMetrics] [Read] %v", err)
		return nil, err
	}

	var currentMetrics metrics
	err = json.Unmarshal(raw, &currentMetrics)
	if err != nil {
		m.lgr.Sugar().Errorf("[Service] [getCurrentMetrics] [Unmarshal] %v", err)
		return nil, err
	}

	return currentMetrics, nil
}

func NewService(str *store.Store, lgr *zap.Logger) Service {
	return &metricsService{
		lgr: lgr,
		str: str,
	}
}
