package agent

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	reflect "reflect"
	"runtime"
	"sync"
	"time"
)

type Settings struct {
	PollInterval   time.Duration
	ReportInterval time.Duration
	Metrics        *[]string
	Host           string
	Port           string
	RequestTimeout time.Duration
}

type metrics struct {
	memstats    *runtime.MemStats
	randomValue float64
	pollCount   int64
}

type Agent struct {
	Settings Settings
	client   http.Client
	data     metrics
	SyncWG   sync.WaitGroup
}

func MakeAgent(settings Settings) *Agent {
	var A = new(Agent)
	A.Settings = settings
	A.client = http.Client{}
	transport := &http.Transport{}
	transport.MaxIdleConns = 20
	A.client.Transport = transport
	return A
}

func (A *Agent) Start(ctx context.Context) {

	pollInterval := time.NewTicker(A.Settings.PollInterval)
	defer pollInterval.Stop()

	sendInterval := time.NewTicker(A.Settings.ReportInterval)
	defer sendInterval.Stop()

	for {
		select {
		case <-pollInterval.C:
			A.parseMetrics()
		case <-sendInterval.C:
			A.sendMetrics(ctx)
		case <-ctx.Done():
			return
		}
	}

}

func (A *Agent) sendMetrics(ctx context.Context) {
	ctx2, cancel := context.WithTimeout(ctx, A.Settings.RequestTimeout)
	defer cancel()
	for _, name := range *A.Settings.Metrics {
		value, statType := GetMemStatByName(A.data.memstats, name)
		switch statType {
		case reflect.Uint64:
			A.sendRequest(ctx2, "gauge", name, value)
		case reflect.Float64:
			A.sendRequest(ctx2, "gauge", name, value)
		}
	}
	A.sendRequest(ctx, "counter", "PollCount", A.data.pollCount)
	A.sendRequest(ctx, "counter", "RandomValue", A.data.randomValue)
	A.SyncWG.Wait()
}

func (A *Agent) sendRequest(ctx context.Context, statType string, nameStat string, value interface{}) {
	var url string
	//x := reflect.TypeOf(value).Kind()
	url = fmt.Sprintf("http://%s:%s/update/%s/%s/%v",
		A.Settings.Host, A.Settings.Port,
		statType, nameStat, reflect.ValueOf(value))
	fmt.Println(url)
	A.SyncWG.Add(1)
	go func() {
		defer A.SyncWG.Done()
		req, err := http.NewRequestWithContext(
			ctx, "POST", url, nil)
		if err != nil {
			fmt.Println(err)
			return
		}
		resp, err := A.client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()
	}()
}

func (A *Agent) parseMetrics() {
	var a = runtime.MemStats{}
	runtime.ReadMemStats(&a)
	A.data.memstats = &a
	A.data.randomValue = rand.Float64()
	A.data.pollCount++
}
