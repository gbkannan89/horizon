package metrics

import (
	"fmt"
	"net/http"
	"runtime"
	"sync"
	"time"
)

var (
	startTime     = time.Now()
	metricsMtx    sync.RWMutex
	reqCount      = map[string]int64{}
	errCount      = map[string]int64{}
	durationSum   = map[string]float64{}
)

func RecordRequest(method, path string, status int, duration time.Duration) {
	metricsMtx.Lock()
	defer metricsMtx.Unlock()
	key := method + " " + path
	reqCount[key]++
	durationSum[key] += duration.Seconds()
	if status >= 400 {
		errCount[key]++
	}
}

func Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metricsMtx.RLock()
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")

		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		uptime := time.Since(startTime).Seconds()

		fmt.Fprintf(w, "horizon_build_info{version=\"0.1.0\"} 1\n")
		fmt.Fprintf(w, "horizon_uptime_seconds %.2f\n", uptime)
		fmt.Fprintf(w, "horizon_go_mem_alloc_bytes %d\n", m.Alloc)
		fmt.Fprintf(w, "horizon_go_goroutines %d\n", runtime.NumGoroutine())

		for key, count := range reqCount {
			fmt.Fprintf(w, "horizon_requests_total{endpoint=\"%s\"} %d\n", key, count)
		}
		for key, count := range errCount {
			fmt.Fprintf(w, "horizon_request_errors_total{endpoint=\"%s\"} %d\n", key, count)
		}

		metricsMtx.RUnlock()
	}
}
