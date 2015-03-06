package stats

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

type StatsMiddleware struct {
	Lock                sync.RWMutex
	Start               time.Time
	Pid                 int
	ResponseCounts      map[string]int
	TotalResponseCounts map[string]int
	TotalResponseTime   time.Time
}

func New() *StatsMiddleware {
	stats := &StatsMiddleware{
		Start:               time.Now(),
		Pid:                 os.Getpid(),
		ResponseCounts:      map[string]int{},
		TotalResponseCounts: map[string]int{},
		TotalResponseTime:   time.Time{},
	}

	go func() {
		for {
			stats.ResetResponseCounts()

			time.Sleep(time.Second * 1)
		}
	}()

	return stats
}

func (mw *StatsMiddleware) ResetResponseCounts() {
	mw.Lock.Lock()
	defer mw.Lock.Unlock()
	mw.ResponseCounts = map[string]int{}
}

type recorderResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func (w *recorderResponseWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
	w.StatusCode = code
}

// MiddlewareFunc makes StatsMiddleware implement the Middleware interface.
func (mw *StatsMiddleware) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		writer := &recorderResponseWriter{w, 0}

		h.ServeHTTP(writer, r)

		mw.handleWriter(start, writer)
	})
}

// Negroni compatible interface
func (mw *StatsMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()

	writer := &recorderResponseWriter{w, 200}

	next(writer, r)

	mw.handleWriter(start, writer)
}

func (mw *StatsMiddleware) handleWriter(start time.Time, writer *recorderResponseWriter) {
	end := time.Now()

	responseTime := end.Sub(start)

	mw.Lock.Lock()

	defer mw.Lock.Unlock()

	statusCode := fmt.Sprintf("%d", writer.StatusCode)

	mw.ResponseCounts[statusCode]++
	mw.TotalResponseCounts[statusCode]++
	mw.TotalResponseTime = mw.TotalResponseTime.Add(responseTime)
}

type Stats struct {
	Pid                    int            `json: "pid"`
	UpTime                 string         `json: "uptime"`
	UpTimeSec              float64        `json: "uptime_sec"`
	Time                   string         `json: "time"`
	TimeUnix               int64          `json: "unixtime"`
	StatusCodeCount        map[string]int `json: "status_code_count"`
	TotalStatusCodeCount   map[string]int `json: "total_status_code_count"`
	Count                  int            `json: "count"`
	TotalCount             int            `json: "total_count"`
	TotalResponseTime      string         `json" "total_response_time`
	TotalResponseTimeSec   float64        `json: "total_response_time_sec"`
	AverageResponseTime    string         `json: "average_response_time"`
	AverageResponseTimeSec float64        `json: "average_response_time_sec"`
}

func (mw *StatsMiddleware) GetStats() *Stats {

	mw.Lock.RLock()

	now := time.Now()

	uptime := now.Sub(mw.Start)

	count := 0
	for _, current := range mw.ResponseCounts {
		count += current
	}

	totalCount := 0
	for _, count := range mw.TotalResponseCounts {
		totalCount += count
	}

	totalResponseTime := mw.TotalResponseTime.Sub(time.Time{})

	averageResponseTime := time.Duration(0)
	if totalCount > 0 {
		avgNs := int64(totalResponseTime) / int64(totalCount)
		averageResponseTime = time.Duration(avgNs)
	}

	stats := &Stats{
		Pid:                    mw.Pid,
		UpTime:                 uptime.String(),
		UpTimeSec:              uptime.Seconds(),
		Time:                   now.String(),
		TimeUnix:               now.Unix(),
		StatusCodeCount:        mw.ResponseCounts,
		TotalStatusCodeCount:   mw.TotalResponseCounts,
		Count:                  count,
		TotalCount:             totalCount,
		TotalResponseTime:      totalResponseTime.String(),
		TotalResponseTimeSec:   totalResponseTime.Seconds(),
		AverageResponseTime:    averageResponseTime.String(),
		AverageResponseTimeSec: averageResponseTime.Seconds(),
	}

	mw.Lock.RUnlock()

	return stats
}
