package metrics

import (
	"expvar"
	"net/http"
	"runtime"
	"time"
)

var (
	startTime    = time.Now()
	requestCount = expvar.NewInt("requests_total")
	errorCount   = expvar.NewInt("errors_total")
)

func init() {
	expvar.Publish("goroutines", expvar.Func(func() interface{} {
		return runtime.NumGoroutine()
	}))
	expvar.Publish("uptime_seconds", expvar.Func(func() interface{} {
		return int(time.Since(startTime).Seconds())
	}))
	expvar.Publish("memory_mb", expvar.Func(func() interface{} {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		return m.Alloc / 1024 / 1024
	}))
}

func IncrementRequests() {
	requestCount.Add(1)
}

func IncrementErrors() {
	errorCount.Add(1)
}

func Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		expvar.Handler().ServeHTTP(w, r)
	})
}
