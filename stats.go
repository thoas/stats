package stats

import (
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

type StatsMiddleware struct {
	Lock              sync.RWMutex
	Start             time.Time
	Pid               int
	ResponseCounts    map[string]int
	TotalResponseTime time.Time
}

func New() *StatsMiddleware {
	return &StatsMiddleware{
		Start:             time.Now(),
		Pid:               os.Getpid(),
		ResponseCounts:    map[string]int{},
		TotalResponseTime: time.Time{},
	}
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

	writer := &recorderResponseWriter{w, 0}

	next(writer, r)

	mw.handleWriter(start, writer)
}

func (mw *StatsMiddleware) handleWriter(start time.Time, writer *recorderResponseWriter) {
	end := time.Now()

	responseTime := end.Sub(start)

	statusCode := writer.statusCode

	mw.Lock.Lock()

	defer mw.Lock.Unlock()

	mw.ResponseCounts[fmt.Sprintf("%d", statusCode)]++
	mw.TotalResponseTime = mw.TotalResponseTime.Add(responseTime)
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

	mw.Lock.RLock()

	now := time.Now()

	uptime := now.Sub(mw.Start)

	totalCount := 0
	for _, count := range mw.ResponseCounts {
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
		TotalCount:             totalCount,
		TotalResponseTime:      totalResponseTime.String(),
		TotalResponseTimeSec:   totalResponseTime.Seconds(),
		AverageResponseTime:    averageResponseTime.String(),
		AverageResponseTimeSec: averageResponseTime.Seconds(),
	}

	mw.Lock.RUnlock()

	return stats
}
