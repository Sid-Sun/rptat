package metrics

import (
	"sync"
	"time"

	"github.com/nsnikhil/go-datastructures/queue"
	"github.com/sid-sun/rptat/app/contract"
	"github.com/sid-sun/rptat/app/service"
	"go.uber.org/zap"
)

type (
	Metrics struct {
		svc *service.Service

		lock     sync.Mutex
		syncChan *chan bool
		lgr      *zap.Logger
		total    int

		request  instance
		response instance
	}
	instance struct {
		lock  sync.Mutex
		queue *queue.LinkedQueue
	}
)

func NewMetrics(svc *service.Service) (*Metrics, *chan bool, error) {
	reqQ, err := queue.NewLinkedQueue()
	if err != nil {
		return nil, nil, err
	}
	resQ, err := queue.NewLinkedQueue()
	if err != nil {
		return nil, nil, err
	}
	c := make(chan bool)
	return &Metrics{
		request: instance{
			queue: reqQ,
		},
		response: instance{
			queue: resQ,
		},
		syncChan: &c,
		svc:      svc,
	}, &c, nil
}

func (m *Metrics) IncrementRequestCount(path string) error {
	m.request.lock.Lock()

	err := m.request.queue.Add(contract.Request{Path: path, Date: time.Now().Format("01-02-2006")})
	if err != nil {
		m.lgr.Sugar().Errorf("[Metrics] [IncrementRequestCount] [Add] %v", err)
		m.request.lock.Unlock()
		return err
	}

	m.request.lock.Unlock()
	m.lock.Lock()

	m.total++
	if m.total >= 30 {
		*m.syncChan <- true
		m.total = 0
	}

	m.lock.Unlock()
	return nil
}

func (m *Metrics) IncrementResponseCount(path string, code int) error {
	m.response.lock.Lock()

	err := m.response.queue.Add(contract.Response{Path: path, Code: code, Date: time.Now().Format("01-02-2006")})
	if err != nil {
		m.lgr.Sugar().Errorf("[Metrics] [IncrementResponseCount] [Add] %v", err)
		m.request.lock.Unlock()
		return err
	}

	m.response.lock.Unlock()
	m.lock.Lock()

	m.total++
	if m.total >= 30 {
		*m.syncChan <- true
		m.total = 0
	}

	m.lock.Unlock()
	return nil
}

func (m *Metrics) Sync() {
	for {
		cont := <-*m.syncChan

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

		if !cont {
			// Send ack for sync ack
			*m.syncChan <- true
			break
		}
	}
}
