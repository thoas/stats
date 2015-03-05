package stats

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

type StatsMiddleware struct {
	lock              sync.RWMutex
	start             time.Time
	pid               int
	responseCounts    map[string]int
	totalResponseTime time.Time
}

type recorderResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *recorderResponseWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
	w.statusCode = code
}

// MiddlewareFunc makes StatsMiddleware implement the Middleware interface.
func (mw *StatsMiddleware) Handler(h http.Handler) http.Handler {

	mw.start = time.Now()
	mw.pid = os.Getpid()
	mw.responseCounts = map[string]int{}
	mw.totalResponseTime = time.Time{}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		writer := &recorderResponseWriter{w, 0}

		h.ServeHTTP(writer, r)

		end := time.Now()

		responseTime := end.Sub(start)

		statusCode := writer.statusCode

		mw.lock.Lock()
		mw.responseCounts[fmt.Sprintf("%d", statusCode)]++

		mw.totalResponseTime = mw.totalResponseTime.Add(responseTime)
		mw.lock.Unlock()
	})
}

type Stats struct {
	Pid                    int            `json: "pid"`
	UpTime                 string         `json: "uptime"`
	UpTimeSec              float64        `json: "uptime_sec"`
	Time                   string         `json: "time"`
	TimeUnix               int64          `json: "unixtime"`
	StatusCodeCount        map[string]int `json: "status_code_count"`
	TotalCount             int            `json: "total_count"`
	TotalResponseTime      string         `json" "total_response_time`
	TotalResponseTimeSec   float64        `json: "total_response_time_sec"`
	AverageResponseTime    string         `json: "average_response_time"`
	AverageResponseTimeSec float64        `json: "average_response_time_sec"`
}

func (mw *StatsMiddleware) GetStats() *Stats {

	mw.lock.RLock()

	now := time.Now()

	uptime := now.Sub(mw.start)

	totalCount := 0
	for _, count := range mw.responseCounts {
		totalCount += count
	}

	totalResponseTime := mw.totalResponseTime.Sub(time.Time{})

	averageResponseTime := time.Duration(0)
	if totalCount > 0 {
		avgNs := int64(totalResponseTime) / int64(totalCount)
		averageResponseTime = time.Duration(avgNs)
	}

	stats := &Stats{
		Pid:                    mw.pid,
		UpTime:                 uptime.String(),
		UpTimeSec:              uptime.Seconds(),
		Time:                   now.String(),
		TimeUnix:               now.Unix(),
		StatusCodeCount:        mw.responseCounts,
		TotalCount:             totalCount,
		TotalResponseTime:      totalResponseTime.String(),
		TotalResponseTimeSec:   totalResponseTime.Seconds(),
		AverageResponseTime:    averageResponseTime.String(),
		AverageResponseTimeSec: averageResponseTime.Seconds(),
	}

	mw.lock.RUnlock()

	return stats
}
