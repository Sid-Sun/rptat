package metrics

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/nsnikhil/go-datastructures/queue"
	"go.uber.org/zap"
)

type (
	Metrics struct {
		lock     sync.Mutex
		syncChan *chan bool
		lgr      *zap.Logger
		total    int

		request struct {
			lock  sync.Mutex
			queue *queue.LinkedQueue
		}

		response struct {
			lock  sync.Mutex
			queue *queue.LinkedQueue
		}
	}
	st struct {
		lock  sync.Mutex
		queue *queue.LinkedQueue
	}
)

type req struct {
	date string
	path string
}

type res struct {
	date string
	path string
	code int
}

// type n struct {

// Requests int `json:"requests"`
// // Responses []struct{
// // 	StatusCode int `json:"status_code"`
// // 	Count int `json:"count"`
// // } `json:"responses"`
// Res map[int]int `json:"responses"`
// }

type m map[string]n

type n map[string]Path

type Path struct {
	Requests int         `json:"requests"`
	Res      map[int]int `json:"responses"`
}

var allMetrics m

func NewMetrics() (*Metrics, *chan bool, error) {
	allMetrics = make(m)

	file, err := os.Open("data.json")
	defer func() {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}()
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, nil, err
	}

	json.Unmarshal(bytes, &allMetrics)

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
		request: st{
			queue: reqQ,
		},
		response: st{
			queue: resQ,
		},
		syncChan: &c,
	}, &c, nil
}

func (m *Metrics) IncrementRequestCount(path string) error {
	m.request.lock.Lock()
	err := m.request.queue.Add(req{path: path, date: time.Now().Format("01-02-2006")})
	m.request.lock.Unlock()
	m.lock.Lock()
	m.total++
	if m.total >= 30 {
		*m.syncChan <- true
		m.total = 0
	}
	m.lock.Unlock()
	return err
}

func (m *Metrics) IncrementResponseCount(path string, code int) error {
	m.response.lock.Lock()
	err := m.response.queue.Add(res{path: path, code: code, date: time.Now().Format("01-02-2006")})
	m.response.lock.Unlock()
	m.lock.Lock()
	m.total++
	if m.total >= 30 {
		*m.syncChan <- true
		m.total = 0
	}
	m.lock.Unlock()
	return err
}

func (m *Metrics) Sync(shutdownGracefully *bool) {
	for {
		cont := <-*m.syncChan

		reqMetrics := make(map[req]int)
		var totalIncrement int
		m.request.lock.Lock()
		for !m.request.queue.Empty() {
			elem, err := m.request.queue.Remove()
			if err != nil {
				m.lgr.Sugar().Error("[Metrics] [Sync] [ReqRemove] %s", err.Error())
			}
			reqMetrics[elem.(req)]++
			totalIncrement++
		}
		m.request.lock.Unlock()
		// fmt.Println("Requests:")
		for reqElem, count := range reqMetrics {
			foo := allMetrics[reqElem.date]
			if foo == nil {
				foo = make(n)
			}

			boo := foo[reqElem.path].Res
			if boo == nil {
				boo = make(map[int]int)
			}

			foo[reqElem.path] = Path{
				Requests: foo[reqElem.path].Requests + count,
				Res:      boo,
			}

			allMetrics[reqElem.date] = foo
		}

		resMetrics := make(map[res]int)
		m.response.lock.Lock()
		for !m.response.queue.Empty() {
			elem, err := m.response.queue.Remove()
			if err != nil {
				m.lgr.Sugar().Error("[Metrics] [Sync] [ResRemove] %s", err.Error())
			}
			resMetrics[elem.(res)]++
			totalIncrement++
		}
		m.response.lock.Unlock()

		// fmt.Println("Responses:")
		for req, count := range resMetrics {
			// fmt.Printf("  Path: [%s] Staus Code: [%d] Count: [%d]\n", req.path, req.code, count)
			allMetrics[req.date][req.path].Res[req.code] += count
		}

		for date, data := range allMetrics {
			fmt.Println(date)
			for key, val := range data {
				fmt.Printf("Path: %s Reqs: %d\nResponses: \n", key, val.Requests)
				for code, count := range val.Res {
					fmt.Printf("  Status Code: %d Count: %d\n", code, count)
				}
			}
		}

		data, err := json.Marshal(allMetrics)
		if err != nil {
			m.lgr.Sugar().Errorf("[Metrics] [Sync] [Marshal] %s", err.Error())
		}
		err = ioutil.WriteFile("data.json", data, 0644)
		if err != nil {
			m.lgr.Sugar().Errorf("[Metrics] [Sync] [WriteFile] %s", err.Error())
		}

		if !cont {
			// Send ack for sync ack
			*m.syncChan <- true
			break
		}
	}
}
